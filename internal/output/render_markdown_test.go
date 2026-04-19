package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestRenderMarkdown_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	err := renderMarkdown([]diff.Result{}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestRenderMarkdown_Changes(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/app", Key: "password", Status: diff.Modified, LeftValue: "old", RightValue: "new"},
		{Path: "secret/app", Key: "token", Status: diff.Added, LeftValue: "", RightValue: "abc"},
	}
	var buf bytes.Buffer
	err := renderMarkdown(results, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Path | Key |") {
		t.Errorf("expected table header, got: %s", out)
	}
	if !strings.Contains(out, "password") {
		t.Errorf("expected password key in output")
	}
	if !strings.Contains(out, "added") {
		t.Errorf("expected 'added' status in output")
	}
}

func TestRenderMarkdown_EscapesPipe(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/app", Key: "key", Status: diff.Modified, LeftValue: "a|b", RightValue: "c"},
	}
	var buf bytes.Buffer
	_ = renderMarkdown(results, &buf)
	if !strings.Contains(buf.String(), `\|`) {
		t.Errorf("expected escaped pipe in output, got: %s", buf.String())
	}
}

func TestRenderer_MarkdownFormat(t *testing.T) {
	r := NewRenderer("markdown")
	var buf bytes.Buffer
	r.Writer = &buf
	err := r.Render([]diff.Result{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Vault Diff") {
		t.Errorf("expected markdown heading")
	}
}
