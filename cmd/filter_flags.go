package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerFilterFlags attaches secret filtering flags to the given command.
func registerFilterFlags(cmd *cobra.Command) {
	cmd.Flags().String("prefix", "", "Only compare secrets whose path starts with this prefix")
	cmd.Flags().StringSlice("exclude-keys", nil, "Comma-separated list of secret keys to exclude from comparison")
}

// resolveFilterOptions builds a FilterOptions from the command's parsed flags.
func resolveFilterOptions(cmd *cobra.Command) vault.FilterOptions {
	prefix, _ := cmd.Flags().GetString("prefix")
	excludeKeys, _ := cmd.Flags().GetStringSlice("exclude-keys")

	// Normalise: an empty slice from cobra is []string{""} when the flag is
	// not set; strip those out so FilterSecrets doesn't treat "" as a key.
	clean := make([]string, 0, len(excludeKeys))
	for _, k := range excludeKeys {
		if k != "" {
			clean = append(clean, k)
		}
	}

	return vault.FilterOptions{
		Prefix:      prefix,
		ExcludeKeys: clean,
	}
}
