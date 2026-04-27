package cmd

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func newWatchTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	registerWatchFlags(cmd)
	return cmd
}

func TestRegisterWatchFlags_FlagsPresent(t *testing.T) {
	cmd := newWatchTestCmd()
	for _, name := range []string{"watch", "watch-interval"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag %q to be registered", name)
		}
	}
}

func TestResolveWatchOptions_Defaults(t *testing.T) {
	cmd := newWatchTestCmd()
	_ = cmd.ParseFlags([]string{})
	opts := resolveWatchOptions(cmd)
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.Interval != 30*time.Second {
		t.Errorf("unexpected default interval: %v", opts.Interval)
	}
}

func TestResolveWatchOptions_Enabled(t *testing.T) {
	cmd := newWatchTestCmd()
	_ = cmd.ParseFlags([]string{"--watch"})
	opts := resolveWatchOptions(cmd)
	if !opts.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestResolveWatchOptions_CustomInterval(t *testing.T) {
	cmd := newWatchTestCmd()
	_ = cmd.ParseFlags([]string{"--watch-interval", "10s"})
	opts := resolveWatchOptions(cmd)
	if opts.Interval != 10*time.Second {
		t.Errorf("expected 10s, got %v", opts.Interval)
	}
}

func TestResolveWatchOptions_ZeroIntervalKeepsDefault(t *testing.T) {
	cmd := newWatchTestCmd()
	_ = cmd.ParseFlags([]string{"--watch-interval", "0s"})
	opts := resolveWatchOptions(cmd)
	if opts.Interval != 30*time.Second {
		t.Errorf("zero interval should fall back to default, got %v", opts.Interval)
	}
}
