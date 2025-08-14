package credsctl

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vogtp/go-creds-mgr/pkg/creds"
	"golang.org/x/term"
)

// Command adds a cobra.Command to manage credentials
func Command(getManager func() creds.Manager) *cobra.Command {
	credsCmd := &cobra.Command{
		Use:          "credentials",
		Short:        "Manage credentials",
		Long:         ``,
		Aliases:      []string{"creds", "cred"},
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	credsListCmd := &cobra.Command{
		Use:     "list",
		Short:   "List credentials",
		Long:    ``,
		Aliases: []string{"ls", "show"},
		RunE: func(cmd *cobra.Command, args []string) error {
			crds, err := getManager().List(cmd.Context())
			if err != nil {
				return err
			}
			for _, c := range crds {
				fmt.Println(c)
			}
			return nil
		},
	}
	credsCmd.AddCommand(credsListCmd)

	credsStoreCmd := &cobra.Command{
		Use:     "store <name> <secret>",
		Short:   "store a credentials",
		Long:    ``,
		Aliases: []string{"add", "save"},

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return cmd.Usage()
			}
			name := args[0]
			sec := []byte(args[1])
			return getManager().Store(cmd.Context(), name, sec)
		},
	}
	credsCmd.AddCommand(credsStoreCmd)

	credsLoadCmd := &cobra.Command{
		Use:     "load <name> [<password>]",
		Short:   "load a credential",
		Long:    ``,
		Aliases: []string{"get"},

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Usage()
			}
			name := args[0]
			s, err := getManager().Load(cmd.Context(), name)
			if err != nil {
				return err
			}
			var bytePassword []byte
			if len(args) > 1 {
				bytePassword = []byte(args[1])
			} else {
				fmt.Print("Enter Password: ")
				bytePassword, err = term.ReadPassword(0)
				if err != nil {
					return err
				}
			}
			if !getManager().ValidatePass(bytePassword) {
				s = []byte("Secret not shown, password not valid")
			}
			fmt.Printf("%s: %q\n", name, string(s))
			return nil
		},
	}
	credsCmd.AddCommand(credsLoadCmd)
	return credsCmd
}
