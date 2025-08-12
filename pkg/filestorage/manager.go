package filestorage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vogtp/go-creds-mgr/pkg/creds"
)

func New(ctx context.Context, options ...Option) (creds.Manager, error) {

	ts := &fileStorage{
		storagePath:   "/var/lib/go-creds-mgr/",
		fileExtention: "srt",
	}
	for _, o := range options {
		o(ts)
	}

	if _, err := os.ReadDir(ts.storagePath); err != nil {
		return nil, err
	}
	return ts, nil
}

type fileStorage struct {
	storagePath     string
	secretsPassword []byte
	fileExtention   string
}

func (s fileStorage) List(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(s.storagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open secrect store: %w", err)
	}
	res := make([]string, 0, len(entries))
	trimExt := len(s.fileExtention) + 1
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := e.Name()
		if !strings.HasSuffix(n, s.fileExtention) {
			continue
		}
		res = append(res, n[:len(n)-trimExt])
	}
	return res, nil
}
func (s *fileStorage) Load(ctx context.Context, name string) ([]byte, error) {
	b, err := os.ReadFile(s.getStorageFilename(name))
	if err != nil {
		return nil, fmt.Errorf("reading secret file %s: %w", s.getStorageFilename(name), err)
	}
	return decrypt(b, s.secretsPassword), nil
}
func (s *fileStorage) Store(ctx context.Context, name string, secret []byte) error {
	b := encrypt(secret, s.secretsPassword)
	return os.WriteFile(s.getStorageFilename(name), b, 0640)
}

func (s fileStorage) getStorageFilename(name string) string {
	return filepath.Clean(fmt.Sprintf("%s/%s.%s", s.storagePath, name, s.fileExtention))
}

func (ts fileStorage) isSecretsPassword(s []byte) bool {
	return bytes.Equal(ts.secretsPassword, s)
}
