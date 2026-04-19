package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"vaultdiff/internal/diff"
)

func TestRenderText_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := &Renderer{Writer: &buf, Format: FormatText, NoColor: true}
	if err := r.Render(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %q", buf.String())
	}
}

func TestRenderText_Changes(t *testing.T) {
	results := []diff.Result{
		{Key: "foo", Change: diff.Added, RightValue: "bar"},
		{Key: "baz", Change: diff.Removed, LeftValue: "old"},
		{Key: "qux", Change: diff.Modified, LeftValue: "v1", RightValue: "v2"},
	}
	var buf bytes.Buffer
	r := &Renderer{Writer: &buf, Format: FormatText, NoColor: true}
	if err := r.Render(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"+", "-", "~", "foo", "baz", "qux"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output: %q", want, out)
		}
	}
}

func TestRenderJSON_Changes(t *testing.T) {
	results := []diff.Result{
		{Key: "secret", Change: diff.Modified, LeftValue: "a", RightValue: "b"},
	}
	var buf bytes.Buffer
	r := &Renderer{Writer: &buf, Format: FormatJSON, NoColor: true}
	if err := r.Render(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entries []jsonEntry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(entries) != 1 || entries[0].Key != "secret" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}
