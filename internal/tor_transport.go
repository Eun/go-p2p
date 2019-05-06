package internal

import (
	"context"
	"net"

	"fmt"

	"os"

	"io/ioutil"

	"io"

	"github.com/cretz/bine/tor"
	"go.uber.org/zap"
)

func TORTransport(ctx context.Context, logger *zap.Logger, progressFunc func(info string, percentage int)) (Transport, error) {
	var err error
	t := torTransport{
		Context: ctx,
		Logger:  logger,
	}

	t.TempDir, err = ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	// Start tor with some defaults + elevated verbosity
	logger.Debug("Starting and registering onion service, please wait a bit...")
	startConf := &tor.StartConf{
		ProcessCreator:  libTorWrapper{},
		TempDataDirBase: t.TempDir,
	}

	var writers []io.Writer

	if progressFunc != nil {
		writers = append(writers, &TorProgress{OnProgress: progressFunc})
	}

	if logger.Core().Enabled(zap.DebugLevel) {
		writers = append(writers, os.Stdout)
	}

	if size := len(writers); size == 1 {
		startConf.DebugWriter = writers[0]
	} else if size > 1 {
		startConf.DebugWriter = io.MultiWriter(writers...)
	}

	t.Tor, err = tor.Start(ctx, startConf)
	if err != nil {
		_ = os.RemoveAll(t.TempDir)
		return nil, err
	}
	return &t, nil
}

type torTransport struct {
	Context context.Context
	Logger  *zap.Logger
	Tor     *tor.Tor
	TempDir string
}

func (t *torTransport) Listen(addr string) (listener Listener, err error) {
	var port string
	var portNumber int
	_, port, err = net.SplitHostPort(addr)
	if err != nil {
		return listener, fmt.Errorf("unable to get port from `%s'", addr)
	}

	portNumber, err = net.DefaultResolver.LookupPort(t.Context, "tcp", port)
	if err != nil {
		return listener, fmt.Errorf("unable to lookupPort port `%s'", port)
	}

	var torListener *tor.OnionService
	torListener, err = t.Tor.Listen(t.Context, &tor.ListenConf{
		RemotePorts: []int{portNumber},
	})
	if err != nil {
		return listener, err
	}

	listener.Listener = torListener
	listener.ID = torListener.ID
	return listener, err
}

func (t *torTransport) Dial(addr string) (net.Conn, error) {
	return t.DialContext(t.Context, "", addr)
}

func (t *torTransport) DialContext(ctx context.Context, _, addr string) (net.Conn, error) {
	d, err := t.Tor.Dialer(ctx, nil)
	if err != nil {
		return nil, err
	}
	return d.DialContext(ctx, "tcp", addr)
}

func (t *torTransport) Close() error {
	_ = t.Tor.Close()
	_ = os.RemoveAll(t.TempDir)
	return nil
}
