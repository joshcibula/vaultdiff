package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newRateLimitTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerRateLimitFlags(cmd)
	return cmd
}

func TestRegisterRateLimitFlags_FlagsPresent(t *testing.T) {
	cmd := newRateLimitTestCmd()
	for _, name := range []string{"rate-limit", "rate-burst"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag %q to be registered", name)
		}
	}
}

func TestResolveRateLimitOptions_Defaults(t *testing.T) {
	cmd := newRateLimitTestCmd()
	// Parse with no args so defaults apply.
	_ = cmd.ParseFlags([]string{})
	opts := resolveRateLimitOptions(cmd)

	if opts.RequestsPerSecond != 10 {
		t.Errorf("expected RequestsPerSecond=10, got %f", opts.RequestsPerSecond)
	}
	if opts.Burst != 20 {
		t.Errorf("expected Burst=20, got %f", opts.Burst)
	}
}

func TestResolveRateLimitOptions_CustomValues(t *testing.T) {
	cmd := newRateLimitTestCmd()
	_ = cmd.ParseFlags([]string{"--rate-limit=50", "--rate-burst=100"})
	opts := resolveRateLimitOptions(cmd)

	if opts.RequestsPerSecond != 50 {
		t.Errorf("expected RequestsPerSecond=50, got %f", opts.RequestsPerSecond)
	}
	if opts.Burst != 100 {
		t.Errorf("expected Burst=100, got %f", opts.Burst)
	}
}

func TestResolveRateLimitOptions_ZeroRateIsUnlimited(t *testing.T) {
	cmd := newRateLimitTestCmd()
	_ = cmd.ParseFlags([]string{"--rate-limit=0"})
	opts := resolveRateLimitOptions(cmd)

	if opts.RequestsPerSecond < 1_000 {
		t.Errorf("expected very high RPS for unlimited mode, got %f", opts.RequestsPerSecond)
	}
}
