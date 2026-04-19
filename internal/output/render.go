package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff"
)

type Renderer struct {
	Format string
	Writer io.Writer
}

func NewRenderer(format string) *Renderer {
	return &Renderer{
		Format: format,
		Writer: os.Stdout,
	}
}

func (r *Renderer) Render(results []diff.Result) error {
	switch strings.ToLower(r.Format) {
	case "json":
		return renderJSON(results, r.Writer)
	case "markdown", "md":
		return renderMarkdown(results, r.Writer)
	case "text", "":
		return renderText(results, r.Writer)
	default:
		return fmt.Errorf("unsupported output format: %s", r.Format)
	}
}

func renderText(results []diff.Result, w io.Writer) error {
	if len(results) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return nil
	}
	for _, r := range results {
		switch r.Status {
		case diff.Added:
			fmt.Fprintf(w, "+ [%s] %s = %s\n", r.Path, r.Key, r.RightValue)
		case diff.Removed:
			fmt.Fprintf(w, "- [%s] %s = %s\n", r.Path, r.Key, r.LeftValue)
		case diff.Modified:
			fmt.Fprintf(w, "~ [%s] %s: %s -> %s\n", r.Path, r.Key, r.LeftValue, r.RightValue)
		}
	}
	return nil
}
