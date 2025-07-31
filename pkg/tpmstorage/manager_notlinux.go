//go:build !linux

package tpmstorage

import (
	"context"
	"fmt"

	"github.com/vogtp/go-creds-mgr/pkg/creds"
)

func New(ctx context.Context, options ...Option) (creds.Manager, error) {
	return nil, fmt.Errorf("Credentials manager is only implemented for linux")
}

type tpmStorage struct {
	tpmDevPaths     []string
	storagePath     string
	secretsPassword []byte
	fileExtention   string
}
