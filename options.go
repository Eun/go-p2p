package p2p

import (
	"github.com/Eun/go-p2p/internal"
	"go.uber.org/zap"
)

type OnProgressFunc func(info string, percentage int)

type ListenerOption func(*listener) error

var ListenerOptions = struct {
	Logger     func(logger *zap.Logger) ListenerOption
	OnProgress func(onProgress OnProgressFunc) ListenerOption
	UseTCP     func() ListenerOption
	SecretSize func(size int) ListenerOption
}{
	Logger: func(logger *zap.Logger) ListenerOption {
		return func(l *listener) error {
			l.Logger = logger
			return nil
		}
	},
	OnProgress: func(onProgress OnProgressFunc) ListenerOption {
		return func(l *listener) error {
			l.OnProgress = onProgress
			return nil
		}
	},
	UseTCP: func() ListenerOption {
		return func(l *listener) error {
			l.TransportFunc = internal.TCPTransport
			return nil
		}
	},
	SecretSize: func(size int) ListenerOption {
		return func(l *listener) error {
			l.SecretSize = size
			return nil
		}
	},
}

type DialerOption func(*dialer) error

var DialerOptions = struct {
	Logger     func(logger *zap.Logger) DialerOption
	OnProgress func(onProgress OnProgressFunc) DialerOption
	UseTCP     func() DialerOption
}{
	Logger: func(logger *zap.Logger) DialerOption {
		return func(d *dialer) error {
			d.Logger = logger
			return nil
		}
	},
	OnProgress: func(onProgress OnProgressFunc) DialerOption {
		return func(d *dialer) error {
			d.OnProgress = onProgress
			return nil
		}
	},
	UseTCP: func() DialerOption {
		return func(d *dialer) error {
			d.TransportFunc = internal.TCPTransport
			d.IdentifierSuffix = ""
			return nil
		}
	},
}
