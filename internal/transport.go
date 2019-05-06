package internal

import (
	"context"
	"net"
)

type Transport interface {
	Close() error
	Listen(addr string) (Listener, error)
	Dial(addr string) (net.Conn, error)
	DialContext(ctx context.Context, _, addr string) (net.Conn, error)
}

type Listener struct {
	net.Listener
	ID string
}
