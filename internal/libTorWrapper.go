package internal

import (
	"context"

	"github.com/cretz/bine/process"
	libtor "github.com/ipsn/go-libtor"
)

var creator = libtor.Creator

type libTorWrapper struct{}

func (libTorWrapper) New(ctx context.Context, args ...string) (process.Process, error) {
	return creator.New(ctx, args...)
}
