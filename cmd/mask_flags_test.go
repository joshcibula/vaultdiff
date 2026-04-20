package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newMaskTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	registerMaskFlags(cmd)
	return cmd
}

func TestResolveMaskOptions_Defaults(t *testing.T) {
	cmd := newMaskTestCmd()
	_ = cmd.ParseFlags([]string{})
	opts := resolveMaskOptions(cmd)
	if opts.Enabled {
		t.Error("expected mask disabled by default")
	}
	if opts.MaskString != "***" {
		t.Errorf("expected default mask string ***, got %s", opts.MaskString)
	}
	if len(opts.RevealKeys) != 0 {
		t.Errorf("expected no reveal keys by default, got %v", opts.RevealKeys)
	}
}

func TestResolveMaskOptions_Enabled(t *testing.T) {
	cmd := newMaskTestCmd()
	_ = cmd.ParseFlags([]string{"--mask", "--mask-string", "[redacted]", "--reveal-keys", "user,email"})
	opts := resolveMaskOptions(cmd)
	if !opts.Enabled {
		t.Error("expected mask enabled")
	}
	if opts.MaskString != "[redacted]" {
		t.Errorf("expected [redacted], got %s", opts.MaskString)
	}
	if len(opts.RevealKeys) != 2 {
		t.Errorf("expected 2 reveal keys, got %d", len(opts.RevealKeys))
	}
}

func TestRegisterMaskFlags_FlagsPresent(t *testing.T) {
	cmd := newMaskTestCmd()
	for _, name := range []string{"mask", "mask-string", "reveal-keys"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag %q to be registered", name)
		}
	}
}
