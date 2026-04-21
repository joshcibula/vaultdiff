package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newAuditTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerAuditFlags(cmd)
	return cmd
}

func TestRegisterAuditFlags_FlagsPresent(t *testing.T) {
	cmd := newAuditTestCmd()
	for _, name := range []string{"audit", "audit-redact", "audit-file"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("expected flag --%s to be registered", name)
		}
	}
}

func TestResolveAuditOptions_Defaults(t *testing.T) {
	cmd := newAuditTestCmd()
	_ = cmd.ParseFlags([]string{})
	opts, err := resolveAuditOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected audit disabled by default")
	}
	if !opts.RedactValues {
		t.Error("expected redact values true by default")
	}
}

func TestResolveAuditOptions_Enabled(t *testing.T) {
	cmd := newAuditTestCmd()
	_ = cmd.ParseFlags([]string{"--audit", "--audit-redact=false"})
	opts, err := resolveAuditOptions(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected audit enabled")
	}
	if opts.RedactValues {
		t.Error("expected redact values false")
	}
}

func TestResolveAuditOptions_InvalidFile(t *testing.T) {
	cmd := newAuditTestCmd()
	_ = cmd.ParseFlags([]string{"--audit", "--audit-file=/nonexistent/dir/audit.log"})
	_, err := resolveAuditOptions(cmd)
	if err == nil {
		t.Error("expected error for unwritable audit file path")
	}
}
