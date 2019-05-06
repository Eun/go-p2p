package p2p

import (
	"fmt"
	"net"

	"context"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha512"
	"crypto/x509"
	"math/big"
	"time"

	"github.com/Eun/go-p2p/internal"
	onetimetlsserver "github.com/Eun/onetimetls/server"
	"go.uber.org/zap"
)

const ProtoVersion = 1

type listener struct {
	TransportFunc     func(context.Context, *zap.Logger, func(info string, percentage int)) (internal.Transport, error)
	Transport         internal.Transport
	TransportListener internal.Listener
	TLSListener       net.Listener
	Logger            *zap.Logger
	OnProgress        OnProgressFunc
	SecretSize        int
	addr              Addr
}

// Accept waits for and returns the next connection to the listener.
func (l *listener) Accept() (net.Conn, error) {
	return l.TLSListener.Accept()
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *listener) Close() error {
	if l.Transport != nil {
		l.Transport.Close()
	}
	if l.TransportListener.Listener != nil {
		l.TransportListener.Listener.Close()
	}
	if l.TLSListener != nil {
		l.TLSListener.Close()
	}
	return nil
}

// Addr returns the listener's network address.
func (l *listener) Addr() net.Addr {
	return &l.addr
}

func Listen(opts ...ListenerOption) (net.Listener, error) {
	l := listener{
		TransportFunc: internal.TORTransport,
		Logger:        zap.NewNop(),
	}

	for _, o := range opts {
		if err := o(&l); err != nil {
			return nil, err
		}
	}

	ctx := context.Background()

	var err error
	// setup transport
	l.Transport, err = l.TransportFunc(ctx, l.Logger, l.OnProgress)
	if err != nil {
		return nil, err
	}

	port := randomPort()
	l.TransportListener, err = l.Transport.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		l.Close()
		return nil, fmt.Errorf("failed to listen on %d: %v", port, err)
	}

	l.addr.Version = ProtoVersion
	l.addr.Identifier = l.TransportListener.ID
	l.addr.Port = port

	l.addr.Secret, err = generateSecret(l.SecretSize)
	if err != nil {
		l.Close()
		return nil, err
	}

	// setup onetimelts
	// create a server cert
	cert, err := onetimetlsserver.MakeCert(time.Duration(time.Hour * 24))
	if err != nil {
		l.Close()
		return nil, err
	}

	// onetimetls server
	l.TLSListener = &onetimetlsserver.Server{
		Timeout:     time.Second * 60,
		Certificate: &cert,
		Listener:    l.TransportListener,
		EncryptKey: func(conn net.Conn) (secret []byte, cipher x509.PEMCipher, err error) {
			// the secret for the current key is sha512(ID as string representation)
			hash := sha512.Sum512([]byte(l.addr.String()))
			return hash[:], x509.PEMCipherAES256, nil
		},
	}

	fingerPrint := sha1.Sum(cert.Certificate[0])
	l.addr.FingerPrint = fingerPrint[:10]
	return &l, nil
}

func randomPort() uint16 {
	max := big.NewInt(65535 - 1024)
	i, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	return uint16(i.Uint64() + 1024)
}
