package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerAuditFlags attaches audit-related flags to the given command.
func registerAuditFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("audit", false, "enable audit logging of accessed secret paths")
	cmd.Flags().Bool("audit-redact", true, "redact secret values in audit log (show keys only)")
	cmd.Flags().String("audit-file", "", "write audit log to file instead of stderr")
}

// resolveAuditOptions builds AuditOptions from parsed command flags.
func resolveAuditOptions(cmd *cobra.Command) (vault.AuditOptions, error) {
	opts := vault.DefaultAuditOptions()

	enabled, err := cmd.Flags().GetBool("audit")
	if err != nil {
		return opts, err
	}
	opts.Enabled = enabled

	redact, err := cmd.Flags().GetBool("audit-redact")
	if err != nil {
		return opts, err
	}
	opts.RedactValues = redact

	auditFile, err := cmd.Flags().GetString("audit-file")
	if err != nil {
		return opts, err
	}
	if auditFile != "" {
		f, err := os.OpenFile(auditFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			return opts, err
		}
		opts.Writer = f
	}

	return opts, nil
}
