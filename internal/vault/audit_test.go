package vault

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestAuditLogger_Disabled(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultAuditOptions()
	opts.Writer = &buf
	// Disabled by default — nothing should be written.
	logger := NewAuditLogger(opts)
	logger.Log(AuditEntry{
		Timestamp: time.Now(),
		Path:      "secret/data/foo",
		Keys:      []string{"api_key"},
		Source:    "left",
	})
	if buf.Len() != 0 {
		t.Errorf("expected no output when disabled, got: %s", buf.String())
	}
}

func TestAuditLogger_Enabled(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultAuditOptions()
	opts.Enabled = true
	opts.Writer = &buf
	logger := NewAuditLogger(opts)
	logger.Log(AuditEntry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Path:      "secret/data/myapp",
		Keys:      []string{"db_pass", "api_key"},
		Source:    "right",
	})
	out := buf.String()
	if !strings.Contains(out, "[audit]") {
		t.Errorf("expected [audit] prefix, got: %s", out)
	}
	if !strings.Contains(out, "secret/data/myapp") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "source=right") {
		t.Errorf("expected source in output, got: %s", out)
	}
}

func TestAuditLogger_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultAuditOptions()
	opts.Enabled = true
	opts.Writer = &buf
	logger := NewAuditLogger(opts)
	for i := 0; i < 3; i++ {
		logger.Log(AuditEntry{
			Timestamp: time.Now(),
			Path:      "secret/data/item",
			Keys:      []string{"key"},
			Source:    "left",
		})
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 audit lines, got %d", len(lines))
	}
}

func TestDefaultAuditOptions(t *testing.T) {
	opts := DefaultAuditOptions()
	if opts.Enabled {
		t.Error("expected audit to be disabled by default")
	}
	if !opts.RedactValues {
		t.Error("expected redact values to be true by default")
	}
	if opts.Writer == nil {
		t.Error("expected a non-nil default writer")
	}
}
