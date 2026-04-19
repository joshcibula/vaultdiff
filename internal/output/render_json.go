package output

import (
	"encoding/json"
	"fmt"
	"io"

	"vaultdiff/internal/diff"
)

type jsonEntry struct {
	Key        string      `json:"key"`
	Change     string      `json:"change"`
	LeftValue  interface{} `json:"left_value,omitempty"`
	RightValue interface{} `json:"right_value,omitempty"`
}

func renderJSON(w io.Writer, results []diff.Result) error {
	entries := make([]jsonEntry, 0, len(results))
	for _, r := range results {
		e := jsonEntry{
			Key:    r.Key,
			Change: string(r.Change),
		}
		switch r.Change {
		case diff.Added:
			e.RightValue = r.RightValue
		case diff.Removed:
			e.LeftValue = r.LeftValue
		case diff.Modified:
			e.LeftValue = r.LeftValue
			e.RightValue = r.RightValue
		}
		entries = append(entries, e)
	}
	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}
