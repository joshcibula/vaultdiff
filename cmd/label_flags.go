package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerLabelFlags attaches label-related flags to cmd.
func registerLabelFlags(cmd *cobra.Command) {
	cmd.Flags().String("label-prefix", "", "Prepend a string to every secret path label in output")
	cmd.Flags().String("label-strip-prefix", "", "Strip a leading string from every secret path label in output")
	cmd.Flags().StringToString("label-alias", map[string]string{}, "Map exact paths to display aliases (e.g. secret/db=Database)")
}

// resolveLabelOptions builds a vault.LabelOptions from cmd flags.
func resolveLabelOptions(cmd *cobra.Command) (vault.LabelOptions, error) {
	opts := vault.DefaultLabelOptions()

	if v, err := cmd.Flags().GetString("label-prefix"); err == nil {
		opts.Prefix = v
	}
	if v, err := cmd.Flags().GetString("label-strip-prefix"); err == nil {
		opts.StripPrefix = v
	}
	if v, err := cmd.Flags().GetStringToString("label-alias"); err == nil && len(v) > 0 {
		opts.Alias = v
	}

	return opts, nil
}
