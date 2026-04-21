package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerValidateFlags attaches validation-related flags to cmd.
func registerValidateFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("validate", false, "Enable secret validation before diffing")
	cmd.Flags().Bool("validate-require-nonempty", false, "Fail if any secret value is empty")
	cmd.Flags().Int("validate-max-value-length", 0, "Maximum allowed length for secret values (0 = unlimited)")
	cmd.Flags().StringSlice("validate-forbidden-keys", nil, "Comma-separated list of forbidden key names")
}

// resolveValidateOptions builds a ValidationOptions from the parsed flags on cmd.
func resolveValidateOptions(cmd *cobra.Command) (vault.ValidationOptions, error) {
	opts := vault.DefaultValidationOptions()

	enabled, err := cmd.Flags().GetBool("validate")
	if err != nil {
		return opts, err
	}
	opts.Enabled = enabled

	requireNonEmpty, err := cmd.Flags().GetBool("validate-require-nonempty")
	if err != nil {
		return opts, err
	}
	opts.RequireNonEmpty = requireNonEmpty

	maxLen, err := cmd.Flags().GetInt("validate-max-value-length")
	if err != nil {
		return opts, err
	}
	opts.MaxValueLength = maxLen

	forbidden, err := cmd.Flags().GetStringSlice("validate-forbidden-keys")
	if err != nil {
		return opts, err
	}
	opts.ForbiddenKeys = forbidden

	return opts, nil
}
