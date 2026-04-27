package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerGroupFlags attaches grouping flags to the given command.
func registerGroupFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("group", false, "Enable grouping of secrets")
	cmd.Flags().String("group-by", "mount", "Group secrets by: mount, prefix, or depth")
	cmd.Flags().Int("group-depth", 1, "Path depth used when --group-by=depth")
}

// resolveGroupOptions builds a GroupOptions from the command's flags.
func resolveGroupOptions(cmd *cobra.Command) vault.GroupOptions {
	opts := vault.DefaultGroupOptions()

	if v, err := cmd.Flags().GetBool("group"); err == nil {
		opts.Enabled = v
	}
	if v, err := cmd.Flags().GetString("group-by"); err == nil && v != "" {
		opts.GroupBy = v
	}
	if v, err := cmd.Flags().GetInt("group-depth"); err == nil && v > 0 {
		opts.Depth = v
	}

	return opts
}
