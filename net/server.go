package net

import (
	"context"
	"errors"
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
	defaultGracefulTimeout = time.Second
)

type ConnHandler interface {
	HandleConn(context.Context, net.Conn) error // It may need to recover the panic.
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

type Server struct {
	listener net.Listener

	started *atomic.Bool
	stopped *atomic.Bool
	stopc   chan struct{}
	donec   chan struct{}
}

func NewServer(l net.Listener) *Server {
	return &Server{
		listener: l,
		started:  atomic.NewBool(false),
		stopped:  atomic.NewBool(false),
		stopc:    make(chan struct{}),
		donec:    make(chan struct{}),
	}
}

func (s *Server) Serve(handler ConnHandler, opts ...WithServeOption) error {
	if s.started.Swap(true) {
		return ErrAlreadyStarted
	}
	serveOpts := makeServeOptions(opts...)

	ctx, cancel := context.WithCancel(serveOpts.context)
	wg := sync.WaitGroup{}
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
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			handler.HandleConn(ctx, conn)
		}()
	}
}

func (s *Server) Close() (err error) {
	if s.stopped.Swap(true) {
		err = ErrAlreadyStopped
	} else {
		err = s.listener.Close()
		close(s.stopc)
	}
	<-s.donec
	return
}
