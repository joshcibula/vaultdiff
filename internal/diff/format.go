package diff

import (
	"fmt"
	"io"
	"sort"
)

// Format writes a human-readable diff to w.
func Format(w io.Writer, result Result) {
	keys := func(m map[string]string) []string {
		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		return ks
	}

	for _, k := range keys(result.OnlyInLeft) {
		fmt.Fprintf(w, "- %s = %s\n", k, result.OnlyInLeft[k])
	}

	for _, k := range keys(result.OnlyInRight) {
		fmt.Fprintf(w, "+ %s = %s\n", k, result.OnlyInRight[k])
	}

	modKeys := make([]string, 0, len(result.Modified))
	for k := range result.Modified {
		modKeys = append(modKeys, k)
	}
	sort.Strings(modKeys)

	for _, k := range modKeys {
		pair := result.Modified[k]
		fmt.Fprintf(w, "~ %s: %s -> %s\n", k, pair[0], pair[1])
	}

	if !result.HasDifferences() {
		fmt.Fprintln(w, "No differences found.")
	}
}
