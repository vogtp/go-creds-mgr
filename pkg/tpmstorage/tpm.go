//go:build linux

package tpmstorage

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	tpmwrap "github.com/salrashid123/go-tpm-wrapping"
	"google.golang.org/protobuf/encoding/protojson"
)

const Simulator = "simulator"

var tpmDevices = []string{"/dev/tpmrm0", "/dev/tpm0"}

func (ts *tpmStorage) initTPM(ctx context.Context, check func(context.Context, *tpmwrap.TPMWrapper) error) error {
	if ts.tpm != nil {
		if err := check(ctx, ts.tpm); err == nil {
			return nil
		}
	}
	ts.tpm = tpmwrap.NewWrapper()
	if len(ts.tpmDevPaths) > 0 {
		tpmDevices = ts.tpmDevPaths
	}
	var allerr error
	for _, t := range tpmDevices {
		_, err := ts.tpm.SetConfig(ctx, wrapping.WithConfigMap(map[string]string{
			tpmwrap.TPM_PATH: t,
			// tpmwrap.PCR_VALUES: pcrValues,
			tpmwrap.USER_AUTH: string(ts.secretsPassword),
		}))
		if err != nil {
			allerr = fmt.Errorf("cannot initalise TPM: %w", err)
			continue
		}
		if err := check(ctx, ts.tpm); err != nil {
			allerr = fmt.Errorf("tpm %s not working: %w", t, err)
			continue
		}
		if t == Simulator {
			slog.Warn("Using SIMULATED TPM! This is INSECURE!")
		}
		return nil
	}
	return fmt.Errorf("no working tpm found in %v: %w", tpmDevices, allerr)
}

func (ts *tpmStorage) encrypt(ctx context.Context, name string, secret []byte) error {
	var blobInfo *wrapping.BlobInfo
	var err error

	encrypt := func(ctx context.Context, tpm *tpmwrap.TPMWrapper) error {
		blobInfo, err = tpm.Encrypt(ctx, secret)
		if err != nil {
			return fmt.Errorf("error encrypting %w", err)
		}
		return nil
	}

	if err := ts.initTPM(ctx, encrypt); err != nil {
		return err
	}

	slog.Debug("Encrypted secret", "name", name, "encrypted", hex.EncodeToString(blobInfo.Ciphertext))

	b, err := protojson.Marshal(blobInfo)
	if err != nil {
		return fmt.Errorf("error marshalling bytes %w", err)
	}

	// var prettyJSON bytes.Buffer
	// err = json.Indent(&prettyJSON, b, "", "\t")
	// if err != nil {
	// 	return fmt.Errorf("error marshalling json %w", err)
	// }

	// slog.Debug("Marshalled encryptedBlob", "blob", prettyJSON.String())

	err = os.WriteFile(ts.getStorageFilename(name), b, 0666)
	if err != nil {
		return fmt.Errorf("error writing encrypted blob %w", err)
	}
	return nil
}

func (ts *tpmStorage) decrypt(ctx context.Context, name string) ([]byte, error) {

	b, err := os.ReadFile(ts.getStorageFilename(name))
	if err != nil {
		return nil, fmt.Errorf("error reading encrypted file %w", err)
	}

	newBlobInfo := &wrapping.BlobInfo{}
	err = protojson.Unmarshal(b, newBlobInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling %w", err)
	}
	var secret []byte

	decrypt := func(ctx context.Context, tpm *tpmwrap.TPMWrapper) error {
		secret, err = tpm.Decrypt(ctx, newBlobInfo)
		if err != nil {
			return fmt.Errorf("error decrypting %w", err)
		}
		return nil
	}

	if err := ts.initTPM(ctx, decrypt); err != nil {
		return nil, err
	}

	return secret, nil
}
