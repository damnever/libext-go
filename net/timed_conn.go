package net

import (
	"io"
	"net"
	"time"
)

type timedConn struct {
	io.Reader
	io.Writer

	conn net.Conn
}

func NewTimedConn(conn net.Conn, readTimeout, writeTimeout time.Duration) io.ReadWriteCloser {
	return timedConn{
		conn:   conn,
		Reader: NewTimedConnReader(conn, readTimeout),
		Writer: NewTimedConnWriter(conn, writeTimeout),
	}
}

func (c timedConn) Close() error {
	return c.conn.Close()
}

type timedConnReader struct {
	conn    net.Conn
	timeout time.Duration
}

func NewTimedConnReader(conn net.Conn, timeout time.Duration) io.Reader {
	return timedConnReader{
		conn:    conn,
		timeout: timeout,
	}
}

func (r timedConnReader) Read(p []byte) (int, error) {
	if r.timeout > 0 {
		if err := r.conn.SetReadDeadline(time.Now().Add(r.timeout)); err != nil {
			return 0, err
		}
	}
	return r.conn.Read(p)
}

type timedConnWriter struct {
	conn    net.Conn
	timeout time.Duration
}

func NewTimedConnWriter(conn net.Conn, timeout time.Duration) io.Writer {
	return timedConnWriter{
		conn:    conn,
		timeout: timeout,
	}
}

func (w timedConnWriter) Write(b []byte) (int, error) {
	if w.timeout > 0 {
		if err := w.conn.SetWriteDeadline(time.Now().Add(w.timeout)); err != nil {
			return 0, err
		}
	}
	return w.conn.Write(b)
}
