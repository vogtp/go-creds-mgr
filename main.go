package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/vogtp/go-creds-mgr/pkg/creds"
	"github.com/vogtp/go-creds-mgr/pkg/credsctl"
	"github.com/vogtp/go-creds-mgr/pkg/filestorage"
	"github.com/vogtp/go-creds-mgr/pkg/tpmstorage"
)

var secrets_pass = "SECRETS PASSWORD"

func main() {
	defer func(t time.Time) { fmt.Printf("Duration %v\n", time.Since(t)) }(time.Now())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	s := getFileStorage(ctx)
	// s:=getTpmStorage(ctx)
	credsManager, err := creds.New(secrets_pass, s)
	if err != nil {
		log.Fatalf("Cannot create credential manager: %s", err)
	}

	rootCtl := credsctl.Command(func() creds.Manager { return credsManager })
	if err := rootCtl.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}

func getFileStorage(ctx context.Context) creds.Manager {
	tpm, err := filestorage.New(ctx,
		filestorage.SecretsPassword(secrets_pass),
		filestorage.StorePath("./tmp"),
	)
	if err != nil {
		log.Fatalf("Could not open tpm persisten storage: %s", err)
	}
	return tpm
}
func getTpmStorage(ctx context.Context) creds.Manager {
	tpm, err := tpmstorage.New(ctx,
		tpmstorage.SecretsPassword(secrets_pass),
		tpmstorage.StorePath("./tmp"),
		tpmstorage.TPMDevice("/dev/tpmrm0", tpmstorage.Simulator),
	)
	if err != nil {
		log.Fatalf("Could not open tpm persisten storage: %s", err)
	}
	return tpm
}
