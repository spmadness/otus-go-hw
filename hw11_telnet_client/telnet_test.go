package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("connection success", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		client := NewTelnetClient(l.Addr().String(), time.Second*10, nil, nil)

		stderr := os.Stderr
		r, w, err := os.Pipe()
		require.NoError(t, err)
		defer r.Close()

		os.Stderr = w

		defer func() { os.Stderr = stderr }()

		// Connect
		err = client.Connect()
		require.NoError(t, err)

		scanner := bufio.NewScanner(r)
		require.True(t, scanner.Scan(), "Failed to scan stderr: %v", scanner.Err())
		actualMsg := scanner.Text()

		expectedMsg := fmt.Sprintf("...Connected to %s", l.Addr().String())
		require.Truef(t, actualMsg == expectedMsg,
			"not equal, expectedMsg: %s, actualMsg: %s", expectedMsg, actualMsg)
	})

	t.Run("connection fail: non-existent socket address", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)

		address := l.Addr().String()
		l.Close()

		client := NewTelnetClient(address, time.Second*10, nil, nil)

		err = client.Connect()

		opError := &net.OpError{}

		require.ErrorAsf(t, err, &opError, "expected net.OpError actual error: %s", err)
	})

	t.Run("connection fail: timeout reached", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		client := NewTelnetClient(l.Addr().String(), time.Microsecond*1, nil, nil)

		err = client.Connect()
		require.Error(t, err, "expected an error")

		var netError net.Error
		require.True(t, errors.As(err, &netError), "expected err to be of type net.Error")
		require.True(t, netError.Timeout(), "expected err to be timeout error")
	})
}
