package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func registerWatchFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("watch", false, "enable continuous watch mode")
	cmd.Flags().Duration("watch-interval", 30*time.Second, "polling interval for watch mode")
}

func resolveWatchOptions(cmd *cobra.Command) vault.WatchOptions {
	opts := vault.DefaultWatchOptions()

	if enabled, err := cmd.Flags().GetBool("watch"); err == nil {
		opts.Enabled = enabled
	}
	if interval, err := cmd.Flags().GetDuration("watch-interval"); err == nil && interval > 0 {
		opts.Interval = interval
	}
	return opts
}
