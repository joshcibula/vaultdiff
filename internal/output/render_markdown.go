package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff"
)

func renderMarkdown(results []diff.Result, w io.Writer) error {
	if len(results) == 0 {
		fmt.Fprintln(w, "## Vault Diff\n\nNo differences found.")
		return nil
	}

	fmt.Fprintln(w, "## Vault Diff")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| Path | Key | Status | Left | Right |")
	fmt.Fprintln(w, "|------|-----|--------|------|-------|")

	for _, r := range results {
		status := strings.ToLower(string(r.Status))
		left := escapeMarkdown(r.LeftValue)
		right := escapeMarkdown(r.RightValue)
		fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n",
			r.Path, r.Key, status, left, right)
	}

	return nil
}

func escapeMarkdown(s string) string {
	if s == "" {
		return "_empty_"
	}
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}
