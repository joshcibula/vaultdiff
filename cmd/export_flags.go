package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

// registerExportFlags attaches export-related flags to the given command.
func registerExportFlags(cmd *cobra.Command) {
	cmd.Flags().String("export-format", "json", "Export format for secrets: json, csv, env")
	cmd.Flags().String("export-file", "", "Write exported secrets to this file (default: stdout)")
}

// resolveExportOptions reads export flags from the command and returns ExportOptions.
func resolveExportOptions(cmd *cobra.Command) (vault.ExportOptions, string, error) {
	opts := vault.DefaultExportOptions()

	fmt, err := cmd.Flags().GetString("export-format")
	if err != nil {
		return opts, "", fmt_err("export-format", err)
	}

	switch vault.ExportFormat(fmt) {
	case vault.ExportFormatJSON, vault.ExportFormatCSV, vault.ExportFormatEnv:
		opts.Format = vault.ExportFormat(fmt)
	default:
		return opts, "", fmt_errorf("unknown export format %q: must be one of json, csv, env", fmt)
	}

	filePath, err := cmd.Flags().GetString("export-file")
	if err != nil {
		return opts, "", fmt_err("export-file", err)
	}

	return opts, filePath, nil
}

func fmt_err(flag string, err error) error {
	return fmt.Errorf("reading flag %q: %w", flag, err)
}

func fmt_errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
