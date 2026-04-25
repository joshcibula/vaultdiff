package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerPromoteFlags attaches promotion-related flags to the given command.
func registerPromoteFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("promote", false, "promote secrets from source path to destination path")
	cmd.Flags().Bool("promote-dry-run", true, "preview promotion without writing to Vault (default true)")
	cmd.Flags().Bool("promote-overwrite", false, "overwrite existing keys in the destination")
	cmd.Flags().String("promote-prefix", "", "path prefix to prepend to promoted secret paths")
}

// resolvePromoteOptions builds a PromoteOptions from parsed CLI flags.
func resolvePromoteOptions(cmd *cobra.Command) (vault.PromoteOptions, error) {
	opts := vault.DefaultPromoteOptions()

	if v, err := cmd.Flags().GetBool("promote"); err == nil {
		opts.Enabled = v
	}
	if v, err := cmd.Flags().GetBool("promote-dry-run"); err == nil {
		opts.DryRun = v
	}
	if v, err := cmd.Flags().GetBool("promote-overwrite"); err == nil {
		opts.Overwrite = v
	}
	if v, err := cmd.Flags().GetString("promote-prefix"); err == nil {
		opts.PathPrefix = v
	}

	return opts, nil
}
