package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newConcurrencyTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	registerConcurrencyFlags(cmd)
	return cmd
}

func TestRegisterConcurrencyFlags_FlagsPresent(t *testing.T) {
	cmd := newConcurrencyTestCmd()
	if cmd.Flags().Lookup("workers") == nil {
		t.Error("expected --workers flag to be registered")
	}
}

func TestResolveConcurrencyOptions_Defaults(t *testing.T) {
	cmd := newConcurrencyTestCmd()
	_ = cmd.ParseFlags([]string{})

	opts := resolveConcurrencyOptions(cmd)
	if opts.Workers != 5 {
		t.Errorf("expected default workers=5, got %d", opts.Workers)
	}
}

func TestResolveConcurrencyOptions_CustomWorkers(t *testing.T) {
	cmd := newConcurrencyTestCmd()
	_ = cmd.ParseFlags([]string{"--workers", "10"})

	opts := resolveConcurrencyOptions(cmd)
	if opts.Workers != 10 {
		t.Errorf("expected workers=10, got %d", opts.Workers)
	}
}

func TestResolveConcurrencyOptions_ZeroFallsBackToDefault(t *testing.T) {
	cmd := newConcurrencyTestCmd()
	_ = cmd.ParseFlags([]string{"--workers", "0"})

	opts := resolveConcurrencyOptions(cmd)
	if opts.Workers <= 0 {
		t.Errorf("expected positive workers, got %d", opts.Workers)
	}
}
