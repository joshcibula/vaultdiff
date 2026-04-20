package cmd

import (
	"github.com/spf13/cobra"
	"strings"

	"vaultdiff/internal/vault"
)

// registerMaskFlags adds masking-related flags to a command.
func registerMaskFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("mask", false, "Mask secret values in output")
	cmd.Flags().String("mask-string", "***", "String used to replace masked values")
	cmd.Flags().StringSlice("reveal-keys", []string{}, "Comma-separated list of keys whose values are not masked")
}

// resolveMaskOptions builds a MaskOptions from parsed command flags.
func resolveMaskOptions(cmd *cobra.Command) vault.MaskOptions {
	enabled, _ := cmd.Flags().GetBool("mask")
	maskStr, _ := cmd.Flags().GetString("mask-string")
	revealRaw, _ := cmd.Flags().GetStringSlice("reveal-keys")

	var reveal []string
	for _, r := range revealRaw {
		r = strings.TrimSpace(r)
		if r != "" {
			reveal = append(reveal, r)
		}
	}

	return vault.MaskOptions{
		Enabled:    enabled,
		MaskString: maskStr,
		RevealKeys: reveal,
	}
}
