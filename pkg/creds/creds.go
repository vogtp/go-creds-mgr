//go:build linux

package creds

import (
	"bytes"
	"context"
	"log/slog"
)

func New(pass string, backend Manager) (Manager, error) {
	m := &manager{
		persistent:      backend,
		secretsPassword: []byte(pass),
	}
	cache, err := getCache()
	if err != nil {
		slog.Warn("Cannot initalise cache", "err", err)
	}
	m.cache = cache
	return m, nil
}

type manager struct {
	persistent      Manager
	cache           Manager
	secretsPassword []byte
}

func (m *manager) List(ctx context.Context) ([]string, error) {
	return m.persistent.List(ctx)
}

func (m *manager) Load(ctx context.Context, name string) ([]byte, error) {
	s, err := m.cache.Load(ctx, name)
	if err == nil && len(s) > 0 {
		return s, nil
	}
	slog.Info("Keyutils cache miss", "secret-name", name, "err", err)
	s, err = m.persistent.Load(ctx, name)
	if err != nil {
		return nil, err
	}
	m.cache.Store(ctx, name, s)
	return s, nil
}

func (m *manager) Store(ctx context.Context, name string, secret []byte) error {
	secret = bytes.TrimSpace(secret)
	if err := m.cache.Store(ctx, name, secret); err != nil {
		slog.Warn("Could not save secret in keyutils", "name", name, "err", err)
	}
	return m.persistent.Store(ctx, name, secret)
}

func (m manager) ValidatePass(s []byte) bool {
	return bytes.Equal(m.secretsPassword, s)
}
