package net

import (
	"context"
	"io"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTCPServer(t *testing.T) {
	addr := randAddr(t)
	ts, err := NewTCPServer(addr, func(_ context.Context, conn net.Conn) {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err == io.EOF {
				break
			}
			require.Nil(t, err)
			_, err = conn.Write(buf[:n])
			require.Nil(t, err)
		}
	})
	require.Nil(t, err)
	go ts.Serve()

	time.Sleep(30 * time.Millisecond)
	conn, err := net.Dial("tcp", addr)
	require.Nil(t, err)
	buf := make([]byte, 1024)
	for i := 1; i <= 23; i++ {
		data := strconv.Itoa(i)
		_, err := conn.Write([]byte(data))
		require.Nil(t, err)
		n, err := conn.Read(buf)
		require.Nil(t, err)
		require.Equal(t, data, string(buf[:n]))
	}
	require.Nil(t, ts.Close())
}

func TestUDPServer(t *testing.T) {
	addr := randAddr(t)
	us, err := NewUDPServer(addr, func(_ context.Context, conn net.PacketConn, addr net.Addr, data []byte) {
		_, err := conn.WriteTo(data, addr)
		require.Nil(t, err)
	})
	require.Nil(t, err)
	go us.Serve()

	uaddr, err := net.ResolveUDPAddr("udp", addr)
	require.Nil(t, err)
	conn, err := net.DialUDP("udp", nil, uaddr)
	buf := make([]byte, 1024)
	for i := 1; i <= 23; i++ {
		data := strconv.Itoa(i)
		_, err := conn.Write([]byte(data))
		require.Nil(t, err)
		n, err := conn.Read(buf)
		require.Nil(t, err)
		require.Equal(t, data, string(buf[:n]))
	}
	require.Nil(t, us.Close())

}

func randAddr(t *testing.T) string {
	l, err := net.Listen("tcp", ":0")
	require.Nil(t, err)
	defer l.Close()
	return l.Addr().String()
}
