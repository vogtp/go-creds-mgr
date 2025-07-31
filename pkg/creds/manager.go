package creds

import "context"

type Manager interface {
	List(context.Context) ([]string, error)
	Load(context.Context, string) ([]byte, error)
	Store(context.Context, string, []byte) error
}
