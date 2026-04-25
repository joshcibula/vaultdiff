package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newAnnotateTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	registerAnnotateFlags(cmd)
	return cmd
}

func TestRegisterAnnotateFlags_FlagsPresent(t *testing.T) {
	cmd := newAnnotateTestCmd()
	for _, name := range []string{"annotate", "annotate-tag-key", "annotate-tag-value", "annotate-path-prefix", "annotate-custom-tags"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag %q to be registered", name)
		}
	}
}

func TestResolveAnnotateOptions_Defaults(t *testing.T) {
	cmd := newAnnotateTestCmd()
	_ = cmd.ParseFlags([]string{})
	opts, err := resolveAnnotateOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.TagKey != "_vaultdiff_source" {
		t.Errorf("unexpected default TagKey: %q", opts.TagKey)
	}
}

func TestResolveAnnotateOptions_Enabled(t *testing.T) {
	cmd := newAnnotateTestCmd()
	_ = cmd.ParseFlags([]string{"--annotate", "--annotate-tag-value", "production"})
	opts, err := resolveAnnotateOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected Enabled=true")
	}
	if opts.TagValue != "production" {
		t.Errorf("expected TagValue='production', got %q", opts.TagValue)
	}
}

func TestResolveAnnotateOptions_CustomTags(t *testing.T) {
	cmd := newAnnotateTestCmd()
	_ = cmd.ParseFlags([]string{"--annotate", "--annotate-custom-tags", "env=staging,team=platform"})
	opts, err := resolveAnnotateOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.CustomTags["env"] != "staging" {
		t.Errorf("expected custom tag env=staging, got %q", opts.CustomTags["env"])
	}
	if opts.CustomTags["team"] != "platform" {
		t.Errorf("expected custom tag team=platform, got %q", opts.CustomTags["team"])
	}
}
