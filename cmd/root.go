package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultdiff/internal/diff"
	"vaultdiff/internal/vault"
)

var (
	address   string
	token     string
	namespace string
)

var rootCmd = &cobra.Command{
	Use:   "vaultdiff <path1> <path2>",
	Short: "Diff secrets between two HashiCorp Vault paths",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(address, token, namespace)
		if err != nil {
			return fmt.Errorf("failed to create vault client: %w", err)
		}

		left, err := client.ReadSecrets(args[0])
		if err != nil {
			return fmt.Errorf("failed to read secrets from %s: %w", args[0], err)
		}

		right, err := client.ReadSecrets(args[1])
		if err != nil {
			return fmt.Errorf("failed to read secrets from %s: %w", args[1], err)
		}

		results := diff.Compare(left, right)
		output := diff.Format(results, args[0], args[1])
		fmt.Print(output)

		if len(results) > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&address, "address", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Vault token (overrides VAULT_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Vault namespace")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
