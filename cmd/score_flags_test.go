package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newScoreTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	registerScoreFlags(cmd)
	return cmd
}

func TestRegisterScoreFlags_FlagsPresent(t *testing.T) {
	cmd := newScoreTestCmd()
	flags := []string{"score", "score-weight-modified", "score-weight-left", "score-weight-right"}
	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("expected flag %q to be registered", f)
		}
	}
}

func TestResolveScoreConfig_Defaults(t *testing.T) {
	cmd := newScoreTestCmd()
	_ = cmd.ParseFlags([]string{})

	cfg, err := resolveScoreConfig(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected score to be disabled by default")
	}
	if cfg.Options.WeightModified != 1.0 {
		t.Errorf("expected WeightModified=1.0, got %.2f", cfg.Options.WeightModified)
	}
	if cfg.Options.WeightOnlyInLeft != 1.0 {
		t.Errorf("expected WeightOnlyInLeft=1.0, got %.2f", cfg.Options.WeightOnlyInLeft)
	}
	if cfg.Options.WeightOnlyInRight != 1.0 {
		t.Errorf("expected WeightOnlyInRight=1.0, got %.2f", cfg.Options.WeightOnlyInRight)
	}
}

func TestResolveScoreConfig_Enabled(t *testing.T) {
	cmd := newScoreTestCmd()
	_ = cmd.ParseFlags([]string{"--score"})

	cfg, err := resolveScoreConfig(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Enabled {
		t.Error("expected score to be enabled")
	}
}

func TestResolveScoreConfig_CustomWeights(t *testing.T) {
	cmd := newScoreTestCmd()
	_ = cmd.ParseFlags([]string{
		"--score",
		"--score-weight-modified=2.0",
		"--score-weight-left=0.5",
		"--score-weight-right=3.0",
	})

	cfg, err := resolveScoreConfig(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Options.WeightModified != 2.0 {
		t.Errorf("expected WeightModified=2.0, got %.2f", cfg.Options.WeightModified)
	}
	if cfg.Options.WeightOnlyInLeft != 0.5 {
		t.Errorf("expected WeightOnlyInLeft=0.5, got %.2f", cfg.Options.WeightOnlyInLeft)
	}
	if cfg.Options.WeightOnlyInRight != 3.0 {
		t.Errorf("expected WeightOnlyInRight=3.0, got %.2f", cfg.Options.WeightOnlyInRight)
	}
}
