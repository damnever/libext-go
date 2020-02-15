package net

import (
	"context"
	"errors"
	"math"
	"net"
	"sync"
	"time"

	"go.uber.org/atomic"
)

var (
	ErrAlreadyStarted = errors.New("libext-go/net: server already started")
	ErrAlreadyStopped = errors.New("libext-go/net: server already stopped")
)

const (
	defaultGracefulTimeout = 100 * time.Millisecond
)

type (
	ConnHandleFunc func(context.Context, net.Conn)

	Server struct {
		*GenericServer

		listener net.Listener
	}
)

func NewTCPServer(laddr string, handleConn ConnHandleFunc) (*Server, error) {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		return nil, err
	}
	return NewServerFromListener(l, handleConn), nil
}

func NewUnixServer(sockpath string, handleConn ConnHandleFunc) (*Server, error) {
	l, err := net.Listen("unix", sockpath)
	if err != nil {
		return nil, err
	}
	return NewServerFromListener(l, handleConn), nil
}

func NewServerFromListener(l net.Listener, handleConn ConnHandleFunc) *Server {
	return &Server{
		GenericServer: NewGenericServer(func(ctx context.Context) (func(), error) {
			conn, err := l.Accept()
			if err != nil {
				return nil, err
			}
			return func() { handleConn(ctx, conn) }, nil
		}, l.Close),
		listener: l,
	}
}

func (s *Server) ListenAddr() net.Addr {
	return s.listener.Addr()
}

type (
	PacketHandleFunc func(context.Context, net.PacketConn, net.Addr, []byte)

	// PacketServer processing one packet per goroutine, use with caution.
	PacketServer struct {
		*GenericServer

		conn net.PacketConn
	}
)

func NewUDPServer(laddr string, handleConn PacketHandleFunc) (*PacketServer, error) {
	conn, err := net.ListenPacket("udp", laddr)
	if err != nil {
		return nil, err
	}
	return NewPacketServerFromConn(conn, handleConn), nil
}

func NewPacketServerFromConn(conn net.PacketConn, handleConn PacketHandleFunc) *PacketServer {
	// Here we can't do PMTUD and assume Ethernet frames are large,
	// so take the IP packet maximum size as the buffer size.
	buf := make([]byte, math.MaxUint16, math.MaxUint16)
	return &PacketServer{
		GenericServer: NewGenericServer(func(ctx context.Context) (func(), error) {
			n, addr, err := conn.ReadFrom(buf[:])
			if err != nil {
				return nil, err
			}
			data := make([]byte, n, n)
			copy(data, buf[:n])

			return func() {
				// Multiple goroutines may invoke methods on a PacketConn simultaneously.
				handleConn(ctx, conn, addr, data)
			}, nil
		}, conn.Close),
		conn: conn,
	}
}

func (s *PacketServer) ListenAddr() net.Addr {
	return s.conn.LocalAddr()
}

type (
	ServeOptions struct {
		context         context.Context
		gracefulTimeout time.Duration
	}
	WithServeOption func(opts *ServeOptions)
)

func WithContext(ctx context.Context) WithServeOption {
	return func(opts *ServeOptions) {
		opts.context = ctx
	}
}

func WithGracefulTimeout(gracefulTimeout time.Duration) WithServeOption {
	return func(opts *ServeOptions) {
		opts.gracefulTimeout = gracefulTimeout
	}
}

var (
	defaultServeOptions = []WithServeOption{
		WithContext(context.Background()),
		WithGracefulTimeout(defaultGracefulTimeout),
	}
)

func makeServeOptions(opts ...WithServeOption) ServeOptions {
	var serveOpts ServeOptions
	opts = append(defaultServeOptions, opts...)
	for _, opt := range opts {
		opt(&serveOpts)
	}
	return serveOpts
}

type GenericServer struct {
	poll  func(context.Context) (func(), error)
	close func() error

	started *atomic.Bool
	stopped *atomic.Bool
	stopc   chan struct{}
	donec   chan struct{}
}

func NewGenericServer(poller func(context.Context) (func(), error), closer func() error) *GenericServer {
	return &GenericServer{
		poll:    poller,
		close:   closer,
		started: atomic.NewBool(false),
		stopped: atomic.NewBool(false),
		stopc:   make(chan struct{}),
		donec:   make(chan struct{}),
	}
}

func (s *GenericServer) Serve(opts ...WithServeOption) error {
	// FIXME(damnever):
	// - https://github.com/golang/go/issues/5045
	// - https://golang.org/src/sync/once.go?s=1473:1500#L30
	if s.started.Swap(true) {
		return ErrAlreadyStarted
	}
	serveOpts := makeServeOptions(opts...)

	ctx, cancel := context.WithCancel(serveOpts.context)
	wg := &sync.WaitGroup{}
	defer func() {
		cancel() // Cancel sub-contexts.

		donec := make(chan struct{})
		go func() {
			wg.Wait()
			close(donec)
		}()

		select {
		case <-time.After(serveOpts.gracefulTimeout):
		case <-donec:
		}
		close(s.donec)
	}()

	for {
		select {
		case <-s.stopc:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		handle, err := s.poll(ctx)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			// Panic handling is up to the caller.
			handle()
		}()
	}
}

func (s *GenericServer) Close() (err error) {
	if s.stopped.Swap(true) {
		err = ErrAlreadyStopped
	} else {
		err = s.close()
		close(s.stopc)
	}
	<-s.donec
	return
}
