package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newTransformTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	registerTransformFlags(cmd)
	return cmd
}

func TestRegisterTransformFlags_FlagsPresent(t *testing.T) {
	cmd := newTransformTestCmd()

	flags := []string{"trim-space", "lowercase-keys", "ignore-keys"}
	for _, name := range flags {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag %q to be registered", name)
		}
	}
}

func TestResolveTransformOptions_Defaults(t *testing.T) {
	cmd := newTransformTestCmd()
	_ = cmd.ParseFlags([]string{})

	opts, err := resolveTransformOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.TrimSpace {
		t.Error("expected TrimSpace to be false by default")
	}
	if opts.LowercaseKeys {
		t.Error("expected LowercaseKeys to be false by default")
	}
	if len(opts.IgnoreKeys) != 0 {
		t.Errorf("expected IgnoreKeys to be empty by default, got %v", opts.IgnoreKeys)
	}
}

func TestResolveTransformOptions_TrimSpace(t *testing.T) {
	cmd := newTransformTestCmd()
	_ = cmd.ParseFlags([]string{"--trim-space"})

	opts, err := resolveTransformOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !opts.TrimSpace {
		t.Error("expected TrimSpace to be true")
	}
}

func TestResolveTransformOptions_LowercaseKeys(t *testing.T) {
	cmd := newTransformTestCmd()
	_ = cmd.ParseFlags([]string{"--lowercase-keys"})

	opts, err := resolveTransformOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !opts.LowercaseKeys {
		t.Error("expected LowercaseKeys to be true")
	}
}

func TestResolveTransformOptions_IgnoreKeys(t *testing.T) {
	cmd := newTransformTestCmd()
	_ = cmd.ParseFlags([]string{"--ignore-keys", "foo,bar,baz"})

	opts, err := resolveTransformOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"foo", "bar", "baz"}
	if len(opts.IgnoreKeys) != len(expected) {
		t.Fatalf("expected %d ignore keys, got %d", len(expected), len(opts.IgnoreKeys))
	}
	for i, k := range expected {
		if opts.IgnoreKeys[i] != k {
			t.Errorf("expected IgnoreKeys[%d] = %q, got %q", i, k, opts.IgnoreKeys[i])
		}
	}
}
