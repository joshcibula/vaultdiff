package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerCompareFlags attaches field-level comparison flags to a command.
func registerCompareFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("compare-ignore-case", false, "Ignore case when comparing secret values")
	cmd.Flags().Bool("compare-ignore-whitespace", false, "Ignore leading/trailing whitespace in values")
	cmd.Flags().StringSlice("compare-ignore-keys", nil, "Keys to exclude from field-level comparison")
}

// resolveCompareOptions builds a CompareOptions from parsed command flags.
func resolveCompareOptions(cmd *cobra.Command) (vault.CompareOptions, error) {
	opts := vault.DefaultCompareOptions()

	ignoreCase, err := cmd.Flags().GetBool("compare-ignore-case")
	if err != nil {
		return opts, err
	}
	opts.IgnoreCase = ignoreCase

	ignoreWS, err := cmd.Flags().GetBool("compare-ignore-whitespace")
	if err != nil {
		return opts, err
	}
	opts.IgnoreWhitespace = ignoreWS

	ignoreKeys, err := cmd.Flags().GetStringSlice("compare-ignore-keys")
	if err != nil {
		return opts, err
	}
	if len(ignoreKeys) > 0 {
		opts.IgnoreKeys = ignoreKeys
	}

	return opts, nil
}
