package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newSnapshotTestCmd() *cobra.Command {
	c := &cobra.Command{Use: "test"}
	registerSnapshotFlags(c)
	return c
}

func TestRegisterSnapshotFlags_FlagsPresent(t *testing.T) {
	cmd := newSnapshotTestCmd()

	for _, name := range []string{"snapshot-save", "snapshot-load"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag %q to be registered", name)
		}
	}
}

func TestResolveSnapshotOptions_Defaults(t *testing.T) {
	cmd := newSnapshotTestCmd()
	// parse with no args so defaults apply
	_ = cmd.ParseFlags([]string{})

	opts := resolveSnapshotOptions(cmd)
	if opts.SavePath != "" {
		t.Errorf("expected empty SavePath by default, got %q", opts.SavePath)
	}
	if opts.LoadPath != "" {
		t.Errorf("expected empty LoadPath by default, got %q", opts.LoadPath)
	}
}

func TestResolveSnapshotOptions_SavePath(t *testing.T) {
	cmd := newSnapshotTestCmd()
	_ = cmd.ParseFlags([]string{"--snapshot-save", "/tmp/snap.json"})

	opts := resolveSnapshotOptions(cmd)
	if opts.SavePath != "/tmp/snap.json" {
		t.Errorf("expected SavePath %q, got %q", "/tmp/snap.json", opts.SavePath)
	}
	if opts.LoadPath != "" {
		t.Errorf("expected empty LoadPath, got %q", opts.LoadPath)
	}
}

func TestResolveSnapshotOptions_LoadPath(t *testing.T) {
	cmd := newSnapshotTestCmd()
	_ = cmd.ParseFlags([]string{"--snapshot-load", "/tmp/prev.json"})

	opts := resolveSnapshotOptions(cmd)
	if opts.LoadPath != "/tmp/prev.json" {
		t.Errorf("expected LoadPath %q, got %q", "/tmp/prev.json", opts.LoadPath)
	}
	if opts.SavePath != "" {
		t.Errorf("expected empty SavePath, got %q", opts.SavePath)
	}
}
