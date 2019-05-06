package p2p

import (
	"context"
	"crypto/sha512"
	"errors"
	"net"

	onetimetlsclient "github.com/Eun/onetimetls/client"

	"bytes"
	"crypto/sha1"
	"crypto/x509"
	"fmt"

	"strconv"

	"github.com/Eun/go-p2p/internal"
	"go.uber.org/zap"
)

type dialer struct {
	TransportFunc    func(context.Context, *zap.Logger, func(info string, percentage int)) (internal.Transport, error)
	Transport        internal.Transport
	IdentifierSuffix string
	Logger           *zap.Logger
	OnProgress       OnProgressFunc
}

func Dial(addr string, opts ...DialerOption) (net.Conn, error) {
	return DialContext(context.Background(), "", addr, opts...)
}

func DialContext(ctx context.Context, _, addr string, opts ...DialerOption) (net.Conn, error) {
	var a Addr
	if err := a.ParseString(addr); err != nil {
		return nil, err
	}

	d := dialer{
		TransportFunc:    internal.TORTransport,
		IdentifierSuffix: ".onion",
		Logger:           zap.NewNop(),
	}

	for _, o := range opts {
		if err := o(&d); err != nil {
			return nil, err
		}
	}
	var err error
	// setup transport
	d.Transport, err = d.TransportFunc(ctx, d.Logger, d.OnProgress)
	if err != nil {
		return nil, err
	}

	hash := sha512.Sum512([]byte(addr))

	tlsClient := onetimetlsclient.Client{
		Dialer:                d.Transport,
		VerifyPeerCertificate: verifyServerCert(&a),
		Secret:                hash[:],
	}

	return tlsClient.DialContext(ctx, "tcp", net.JoinHostPort(a.Identifier+d.IdentifierSuffix, strconv.FormatUint(uint64(a.Port), 10)))
}

func verifyServerCert(a *Addr) func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if len(rawCerts) <= 0 {
			return fmt.Errorf("no certs present")
		}
		switch a.Version {
		case 1:
			fp := sha1.Sum(rawCerts[0])
			if !bytes.Equal(a.FingerPrint, fp[:10]) {
				return errors.New("invalid fingerprint")
			}
		default:
			return fmt.Errorf("invalid version %d", a.Version)
		}
		return nil
	}
}
