package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerMergeFlags attaches merge-related flags to the given command.
func registerMergeFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(
		"merge-prefer-left",
		false,
		"When merging secrets, prefer left-side values on key conflicts (default: right wins)",
	)
	cmd.Flags().Bool(
		"merge-skip-empty",
		false,
		"Omit keys with empty string values during merge",
	)
}

// resolveMergeOptions builds a MergeOptions from the command's flags.
func resolveMergeOptions(cmd *cobra.Command) (vault.MergeOptions, error) {
	opts := vault.DefaultMergeOptions()

	preferLeft, err := cmd.Flags().GetBool("merge-prefer-left")
	if err != nil {
		return opts, err
	}
	opts.PreferLeft = preferLeft

	skipEmpty, err := cmd.Flags().GetBool("merge-skip-empty")
	if err != nil {
		return opts, err
	}
	opts.SkipEmpty = skipEmpty

	return opts, nil
}
