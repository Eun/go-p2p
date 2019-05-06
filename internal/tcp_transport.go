package internal

import (
	"context"
	"net"

	"go.uber.org/zap"
)

func (t *tcpTransport) Close() error {
	return nil
}

func TCPTransport(ctx context.Context, logger *zap.Logger, _ func(info string, percentage int)) (Transport, error) {
	return &tcpTransport{
		Context: ctx,
		Logger:  logger,
	}, nil
}

type tcpTransport struct {
	Context context.Context
	Logger  *zap.Logger
}

func (t *tcpTransport) Listen(addr string) (listener Listener, err error) {
	listener.Listener, err = net.Listen("tcp", addr)
	if err != nil {
		return listener, err
	}

	listener.ID, _, err = net.SplitHostPort(addr)
	if err != nil {
		return listener, err
	}
	var ips []net.IP
	switch listener.ID {
	case "":
		listener.ID = "127.0.0.1"
		ips, err = t.getIPs(0)
	case "0.0.0.0":
		listener.ID = "127.0.0.1"
		ips, err = t.getIPs(4)
	case "::":
		listener.ID = "[::1]"
		ips, err = t.getIPs(6)
	}

	if err != nil {
		return listener, err
	}

	if len(ips) > 0 {
		listener.ID = ips[0].String()
	}
	return listener, err
}

func (t *tcpTransport) Dial(addr string) (net.Conn, error) {
	return t.DialContext(t.Context, "", addr)
}

func (t *tcpTransport) DialContext(ctx context.Context, _, addr string) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(ctx, "tcp", addr)
}

func (t *tcpTransport) getIPs(version int) ([]net.IP, error) {
	var ips []net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}

			if ip.IsLoopback() {
				continue
			}

			switch version {
			case 4:
				if ip.To4() != nil {
					ips = append(ips, ip)
				}
			case 6:
				if ip.To4() == nil {
					ips = append(ips, ip)
				}
			default:
				ips = append(ips, ip)
			}
		}
	}
	return ips, nil
}
