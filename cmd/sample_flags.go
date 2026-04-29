package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerSampleFlags attaches sampling-related flags to the given command.
func registerSampleFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("sample", false, "Enable random sampling of secret paths")
	cmd.Flags().Int("sample-max-paths", 100, "Maximum number of paths to include in the sample")
	cmd.Flags().Int64("sample-seed", 0, "Random seed for sampling (0 = non-deterministic)")
}

// resolveSampleOptions builds a SampleOptions from the command's flags.
func resolveSampleOptions(cmd *cobra.Command) (vault.SampleOptions, error) {
	enabled, err := cmd.Flags().GetBool("sample")
	if err != nil {
		return vault.DefaultSampleOptions(), err
	}

	maxPaths, err := cmd.Flags().GetInt("sample-max-paths")
	if err != nil {
		return vault.DefaultSampleOptions(), err
	}

	seed, err := cmd.Flags().GetInt64("sample-seed")
	if err != nil {
		return vault.DefaultSampleOptions(), err
	}

	return vault.SampleOptions{
		Enabled:  enabled,
		MaxPaths: maxPaths,
		Seed:     seed,
	}, nil
}
