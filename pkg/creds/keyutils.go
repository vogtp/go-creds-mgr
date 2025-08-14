//go:build linux

package creds

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dpeckett/keyutils"
)

func getCache() (Manager, error) {
	keyring, err := keyutils.UserKeyring()
	if err != nil {
		return nil, fmt.Errorf("cannot open keyring: %w", err)
	}
	return &keyutilsManager{keyring: keyring}, nil
}

type keyutilsManager struct {
	keyring keyutils.Keyring
}

func (ku keyutilsManager) List(_ context.Context) ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ku *keyutilsManager) Load(_ context.Context, name string) ([]byte, error) {
	key, err := ku.keyring.Search(name)
	if err != nil {
		return nil, fmt.Errorf("cannot search key: %w", err)
	}

	data, err := key.Get()
	if err != nil {
		return nil, fmt.Errorf("cannot read key: %w", err)
	}
	return data, nil
}

func (ku *keyutilsManager) Store(_ context.Context, name string, secret []byte) error {
	id, err := ku.keyring.Add(name, secret)
	if err != nil {
		return fmt.Errorf("cannot store key: %w", err)
	}
	slog.Info("Saved secret to keyutils", "name", name, "keyid", id)
	return nil
}

func (keyutilsManager) ValidatePass([]byte) bool {
	return false
}
