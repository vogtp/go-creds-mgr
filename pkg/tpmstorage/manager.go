//go:build linux

package tpmstorage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tpmwrap "github.com/salrashid123/go-tpm-wrapping"
	"github.com/vogtp/go-creds-mgr/pkg/creds"
)

func New(ctx context.Context, options ...Option) (creds.Manager, error) {

	ts := &tpmStorage{
		storagePath:   "/var/lib/go-creds-mgr/",
		fileExtention: "scrt",
	}
	for _, o := range options {
		o(ts)
	}

	if _, err := os.ReadDir(ts.storagePath); err != nil {
		return nil, err
	}
	return ts, nil
}

type tpmStorage struct {
	tpmDevPaths     []string
	tpm             *tpmwrap.TPMWrapper
	storagePath     string
	secretsPassword []byte
	fileExtention   string
}

func (ts tpmStorage) List(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(ts.storagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open secrect store: %w", err)
	}
	res := make([]string, 0, len(entries))
	trimExt := len(ts.fileExtention) + 1
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := e.Name()
		if !strings.HasSuffix(n, ts.fileExtention) {
			continue
		}
		res = append(res, n[:len(n)-trimExt])
	}
	return res, nil
}
func (ts *tpmStorage) Load(ctx context.Context, name string) ([]byte, error) {
	return ts.decrypt(ctx, name)
}
func (ts *tpmStorage) Store(ctx context.Context, name string, secret []byte) error {
	return ts.encrypt(ctx, name, secret)
}

func (ts tpmStorage) getStorageFilename(name string) string {
	return filepath.Clean(fmt.Sprintf("%s/%s.%s", ts.storagePath, name, ts.fileExtention))
}

func (ts tpmStorage) isSecretsPassword(s []byte) bool {
	return bytes.Equal(ts.secretsPassword, s)
}
