package vault

import (
	"bytes"
	"strings"
	"testing"
)

func TestExportSecrets_JSON(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"DB_PASS": "s3cr3t", "API_KEY": "abc123"},
	}
	var buf bytes.Buffer
	opts := DefaultExportOptions()
	if err := ExportSecrets(&buf, secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "secret/app") {
		t.Errorf("expected path in JSON output, got: %s", out)
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected key in JSON output, got: %s", out)
	}
}

func TestExportSecrets_CSV(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"KEY": "value"},
	}
	var buf bytes.Buffer
	opts := ExportOptions{Format: ExportFormatCSV, PathLabel: "path"}
	if err := ExportSecrets(&buf, secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 row, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "secret/app") {
		t.Errorf("expected path in CSV row: %s", lines[0])
	}
	if !strings.Contains(lines[0], "KEY") {
		t.Errorf("expected key in CSV row: %s", lines[0])
	}
}

func TestExportSecrets_Env(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"TOKEN": "xyz"},
	}
	var buf bytes.Buffer
	opts := ExportOptions{Format: ExportFormatEnv}
	if err := ExportSecrets(&buf, secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "TOKEN=xyz") {
		t.Errorf("expected TOKEN=xyz in env output, got: %s", out)
	}
}

func TestExportSecrets_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := ExportOptions{Format: "xml"}
	err := ExportSecrets(&buf, map[string]map[string]string{}, opts)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestDefaultExportOptions(t *testing.T) {
	opts := DefaultExportOptions()
	if opts.Format != ExportFormatJSON {
		t.Errorf("expected JSON format, got %q", opts.Format)
	}
	if opts.PathLabel != "path" {
		t.Errorf("expected path label 'path', got %q", opts.PathLabel)
	}
}

func TestExportSecrets_MultiplePathsCSV(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/a": {"K1": "v1"},
		"secret/b": {"K2": "v2"},
	}
	var buf bytes.Buffer
	opts := ExportOptions{Format: ExportFormatCSV, PathLabel: "path"}
	if err := ExportSecrets(&buf, secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 rows, got %d: %v", len(lines), lines)
	}
}
