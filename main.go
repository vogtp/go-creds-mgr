package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/vogtp/go-creds-mgr/pkg/creds"
	"github.com/vogtp/go-creds-mgr/pkg/tpmstorage"
)

var secrets_pass = "SECRETS PASSWORD"

type commander interface {
	CobraCommand() *cobra.Command
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	tpm, err := tpmstorage.New(ctx,
		tpmstorage.SecretsPassword(secrets_pass),
		tpmstorage.StorePath("./tmp"),
		tpmstorage.TPMDevice("/dev/tpmrm0", tpmstorage.Simulator),
	)
	if err != nil {
		log.Fatalf("Could not open tpm persisten storage: %s", err)
	}
	credsManager, err := creds.New(secrets_pass, tpm)
	if err != nil {
		log.Fatalf("Cannot create credential manager: %s", err)
	}
	if cmder, ok := credsManager.(commander); ok {
		rootCtl := cmder.CobraCommand()
		if err := rootCtl.ExecuteContext(ctx); err != nil {
			log.Fatal(err)
		}
	}
}
