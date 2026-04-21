package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerConcurrencyFlags adds concurrency-related flags to a command.
func registerConcurrencyFlags(cmd *cobra.Command) {
	cmd.Flags().Int(
		"workers",
		vault.DefaultConcurrencyOptions().Workers,
		"number of parallel workers for fetching secrets",
	)
}

// resolveConcurrencyOptions builds ConcurrencyOptions from parsed flags.
func resolveConcurrencyOptions(cmd *cobra.Command) vault.ConcurrencyOptions {
	workers, err := cmd.Flags().GetInt("workers")
	if err != nil || workers <= 0 {
		return vault.DefaultConcurrencyOptions()
	}
	return vault.ConcurrencyOptions{
		Workers: workers,
	}
}
