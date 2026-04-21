package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// registerScoreFlags attaches score-related flags to the given command.
func registerScoreFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("score", false, "print a similarity score after the diff")
	cmd.Flags().Float64("score-weight-modified", 1.0, "penalty weight for modified keys")
	cmd.Flags().Float64("score-weight-left", 1.0, "penalty weight for keys only in left path")
	cmd.Flags().Float64("score-weight-right", 1.0, "penalty weight for keys only in right path")
}

// ScoreConfig bundles the score flag values resolved from a command.
type ScoreConfig struct {
	Enabled bool
	Options vault.ScoreOptions
}

// resolveScoreConfig reads score flags from the command and returns a ScoreConfig.
func resolveScoreConfig(cmd *cobra.Command) (ScoreConfig, error) {
	enabled, err := cmd.Flags().GetBool("score")
	if err != nil {
		return ScoreConfig{}, err
	}

	weightMod, err := cmd.Flags().GetFloat64("score-weight-modified")
	if err != nil {
		return ScoreConfig{}, err
	}

	weightLeft, err := cmd.Flags().GetFloat64("score-weight-left")
	if err != nil {
		return ScoreConfig{}, err
	}

	weightRight, err := cmd.Flags().GetFloat64("score-weight-right")
	if err != nil {
		return ScoreConfig{}, err
	}

	return ScoreConfig{
		Enabled: enabled,
		Options: vault.ScoreOptions{
			WeightModified:    weightMod,
			WeightOnlyInLeft:  weightLeft,
			WeightOnlyInRight: weightRight,
		},
	}, nil
}
