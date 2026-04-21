package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerTransformFlags adds secret-transformation flags to the given command.
func registerTransformFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("trim-space", false, "Trim leading/trailing whitespace from secret values before diffing")
	cmd.Flags().Bool("lowercase-keys", false, "Normalize secret keys to lowercase before diffing")
	cmd.Flags().StringSlice("ignore-keys", nil, "Comma-separated list of secret keys to zero out before diffing")
}

// resolveTransformOptions builds a vault.TransformOptions from the command's flags.
func resolveTransformOptions(cmd *cobra.Command) vault.TransformOptions {
	trimSpace, _ := cmd.Flags().GetBool("trim-space")
	lowercaseKeys, _ := cmd.Flags().GetBool("lowercase-keys")
	ignoreRaw, _ := cmd.Flags().GetStringSlice("ignore-keys")

	var ignoreKeys []string
	for _, k := range ignoreRaw {
		k = strings.TrimSpace(k)
		if k != "" {
			ignoreKeys = append(ignoreKeys, k)
		}
	}

	return vault.TransformOptions{
		TrimSpace:     trimSpace,
		LowercaseKeys: lowercaseKeys,
		IgnoreKeys:    ignoreKeys,
	}
}
