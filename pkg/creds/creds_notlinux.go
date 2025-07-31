//go:build !linux

package creds

import "fmt"

func New(pass string, backend Manager) (Manager, error) {
	return nil, fmt.Errorf("Credentials manager is only implemented for linux")
}
