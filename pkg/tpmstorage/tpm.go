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

var tmpDevices = []string{"/dev/tpmrm0", "/dev/tpm0"}

func (ts *tpmStorage) initTPM(ctx context.Context) error {
	ts.tpm = tpmwrap.NewWrapper()
	if len(ts.tpmDevPaths) > 0 {
		tmpDevices = ts.tpmDevPaths
	}
	var allerr error
	for _, t := range tmpDevices {
		_, err := ts.tpm.SetConfig(ctx, wrapping.WithConfigMap(map[string]string{
			tpmwrap.TPM_PATH: t,
			// tpmwrap.PCR_VALUES: pcrValues,
			tpmwrap.USER_AUTH: string(ts.secretsPassword),
		}))
		if err != nil {
			allerr = fmt.Errorf("cannot initalise TPM: %w", err)
			continue
		}
		if err := ts.checkTPM(ctx); err != nil {
			allerr = fmt.Errorf("tpm %s not working: %w", t, err)
			continue
		}
		if t == Simulator {
			slog.Warn("Using SIMULATED TPM! This is INSECURE!")
		}
		return nil
	}
	return fmt.Errorf("no working tpm found in %v: %w", tmpDevices, allerr)
}
func (ts *tpmStorage) checkTPM(ctx context.Context) error {
	_, err := ts.tpm.Encrypt(ctx, []byte{0})
	return err
}
func (ts *tpmStorage) encrypt(ctx context.Context, name string, secret []byte) error {

	blobInfo, err := ts.tpm.Encrypt(ctx, secret)
	if err != nil {
		return fmt.Errorf("error encrypting %w", err)
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

	secret, err := ts.tpm.Decrypt(ctx, newBlobInfo)
	if err != nil {
		return nil, fmt.Errorf("error decrypting %w", err)
	}
	return secret, nil
}
