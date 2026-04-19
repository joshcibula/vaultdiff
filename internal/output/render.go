package output

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"

	"vaultdiff/internal/diff"
)

// Format controls the output format of the diff.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Renderer writes diff results to an output destination.
type Renderer struct {
	Writer io.Writer
	Format Format
	NoColor bool
}

// NewRenderer creates a Renderer writing to stdout.
func NewRenderer(format Format, noColor bool) *Renderer {
	if noColor {
		color.NoColor = true
	}
	return &Renderer{Writer: os.Stdout, Format: format, NoColor: noColor}
}

// Render writes the diff results using the configured format.
func (r *Renderer) Render(results []diff.Result) error {
	switch r.Format {
	case FormatJSON:
		return renderJSON(r.Writer, results)
	default:
		return renderText(r.Writer, results)
	}
}

func renderText(w io.Writer, results []diff.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return nil
	}
	add := color.New(color.FgGreen)
	remove := color.New(color.FgRed)
	modify := color.New(color.FgYellow)
	for _, r := range results {
		switch r.Change {
		case diff.Added:
			add.Fprintf(w, "+ %s: %v\n", r.Key, r.RightValue)
		case diff.Removed:
			remove.Fprintf(w, "- %s: %v\n", r.Key, r.LeftValue)
		case diff.Modified:
			modify.Fprintf(w, "~ %s: %v -> %v\n", r.Key, r.LeftValue, r.RightValue)
		}
	}
	return nil
}
