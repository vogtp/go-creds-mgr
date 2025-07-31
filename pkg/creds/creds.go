package creds

import (
	"bytes"
	"context"
)

type Manager interface {
	List(context.Context) ([]string, error)
	Load(context.Context, string) ([]byte, error)
	Store(context.Context, string, []byte) error
}

func New(pass string, backend Manager) Manager {
	m := &manager{
		persistent:      backend,
		secretsPassword: []byte(pass),
	}
	return m
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
