//go:build linux

package creds

import (
	"bytes"
	"context"
)

func New(pass string, backend Manager) (Manager, error) {
	m := &manager{
		persistent:      backend,
		secretsPassword: []byte(pass),
	}
	return m, nil
}

type manager struct {
	persistent      Manager
	secretsPassword []byte
}

func (m *manager) List(ctx context.Context) ([]string, error) {
	return m.persistent.List(ctx)
}
func (m *manager) Load(ctx context.Context, name string) ([]byte, error) {
	return m.persistent.Load(ctx, name)
}
func (m *manager) Store(ctx context.Context, name string, secret []byte) error {
	return m.persistent.Store(ctx, name, secret)
}

func (m manager) isSecretsPassword(s []byte) bool {
	return bytes.Equal(m.secretsPassword, s)
}
