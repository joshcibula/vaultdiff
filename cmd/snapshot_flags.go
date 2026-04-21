package cmd

import (
	"github.com/spf13/cobra"
)

// SnapshotOptions holds CLI flag values related to snapshot save/load behaviour.
type SnapshotOptions struct {
	// SavePath, if non-empty, writes the left-side secrets as a snapshot JSON file.
	SavePath string
	// LoadPath, if non-empty, loads secrets for the left side from a snapshot file
	// instead of reading live from Vault.
	LoadPath string
}

// registerSnapshotFlags attaches snapshot-related flags to cmd.
func registerSnapshotFlags(cmd *cobra.Command) {
	cmd.Flags().String(
		"snapshot-save", "",
		"write a snapshot of the left-side secrets to this JSON file",
	)
	cmd.Flags().String(
		"snapshot-load", "",
		"load left-side secrets from a previously saved snapshot file instead of Vault",
	)
}

// resolveSnapshotOptions reads snapshot flag values from the given command.
func resolveSnapshotOptions(cmd *cobra.Command) SnapshotOptions {
	save, _ := cmd.Flags().GetString("snapshot-save")
	load, _ := cmd.Flags().GetString("snapshot-load")
	return SnapshotOptions{
		SavePath: save,
		LoadPath: load,
	}
}
