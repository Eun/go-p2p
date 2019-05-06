package p2p

import (
	"io"
	"testing"

	"io/ioutil"

	"github.com/stretchr/testify/require"
)

func TestTCP(t *testing.T) {
	listener, err := Listen(ListenerOptions.UseTCP())
	require.NoError(t, err)
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		require.NoError(t, err)
		io.WriteString(conn, "Hello World")
		conn.Close()
	}()

	conn, err := Dial(listener.Addr().String(), DialerOptions.UseTCP())
	require.NoError(t, err)
	defer conn.Close()

	bytes, err := ioutil.ReadAll(conn)
	require.NoError(t, err)
	require.EqualValues(t, []byte("Hello World"), bytes)
}

func TestTOR(t *testing.T) {
	listener, err := Listen()
	require.NoError(t, err)
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		require.NoError(t, err)
		io.WriteString(conn, "Hello World")
		conn.Close()
	}()

	conn, err := Dial(listener.Addr().String())
	require.NoError(t, err)
	defer conn.Close()

	bytes, err := ioutil.ReadAll(conn)
	require.NoError(t, err)
	require.EqualValues(t, []byte("Hello World"), bytes)
}
