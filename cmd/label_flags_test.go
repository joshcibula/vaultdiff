package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newLabelTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	registerLabelFlags(cmd)
	return cmd
}

func TestRegisterLabelFlags_FlagsPresent(t *testing.T) {
	cmd := newLabelTestCmd()
	for _, name := range []string{"label-prefix", "label-strip-prefix", "label-alias"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag %q to be registered", name)
		}
	}
}

func TestResolveLabelOptions_Defaults(t *testing.T) {
	cmd := newLabelTestCmd()
	_ = cmd.Execute()
	opts, err := resolveLabelOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Prefix != "" {
		t.Errorf("expected empty prefix, got %q", opts.Prefix)
	}
	if opts.StripPrefix != "" {
		t.Errorf("expected empty strip-prefix, got %q", opts.StripPrefix)
	}
	if len(opts.Alias) != 0 {
		t.Errorf("expected empty alias map, got %v", opts.Alias)
	}
}

func TestResolveLabelOptions_Prefix(t *testing.T) {
	cmd := newLabelTestCmd()
	_ = cmd.Flags().Set("label-prefix", "staging/")
	_ = cmd.Execute()
	opts, err := resolveLabelOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Prefix != "staging/" {
		t.Errorf("expected prefix 'staging/', got %q", opts.Prefix)
	}
}

func TestResolveLabelOptions_StripPrefix(t *testing.T) {
	cmd := newLabelTestCmd()
	_ = cmd.Flags().Set("label-strip-prefix", "kv/data/")
	_ = cmd.Execute()
	opts, err := resolveLabelOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.StripPrefix != "kv/data/" {
		t.Errorf("expected strip-prefix 'kv/data/', got %q", opts.StripPrefix)
	}
}
