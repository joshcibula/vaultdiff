package vault

import (
	"fmt"
	"sort"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// SummaryOptions controls how a diff summary is generated.
type SummaryOptions struct {
	Enabled       bool
	GroupByMount  bool
	ShowCounts    bool
}

// DefaultSummaryOptions returns sensible defaults.
func DefaultSummaryOptions() SummaryOptions {
	return SummaryOptions{
		Enabled:      true,
		GroupByMount: false,
		ShowCounts:   true,
	}
}

// MountSummary holds per-mount diff statistics.
type MountSummary struct {
	Mount    string
	Added    int
	Removed  int
	Modified int
	Total    int
}

// DiffSummary holds the overall summary of a diff operation.
type DiffSummary struct {
	TotalAdded    int
	TotalRemoved  int
	TotalModified int
	TotalUnchanged int
	ByMount       map[string]*MountSummary
}

// String returns a human-readable summary line.
func (s *DiffSummary) String() string {
	return fmt.Sprintf(
		"added=%d removed=%d modified=%d unchanged=%d",
		s.TotalAdded, s.TotalRemoved, s.TotalModified, s.TotalUnchanged,
	)
}

// SortedMounts returns mount names in sorted order.
func (s *DiffSummary) SortedMounts() []string {
	keys := make([]string, 0, len(s.ByMount))
	for k := range s.ByMount {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SummarizeDiff computes a DiffSummary from a slice of diff results.
func SummarizeDiff(results []diff.Result, opts SummaryOptions) *DiffSummary {
	summary := &DiffSummary{
		ByMount: make(map[string]*MountSummary),
	}

	for _, r := range results {
		mount := mountFromPath(r.Path)
		ms := ensureMount(summary.ByMount, mount)
		ms.Total++

		switch r.Change {
		case diff.Added:
			summary.TotalAdded++
			ms.Added++
		case diff.Removed:
			summary.TotalRemoved++
			ms.Removed++
		case diff.Modified:
			summary.TotalModified++
			ms.Modified++
		default:
			summary.TotalUnchanged++
		}
	}

	return summary
}

func mountFromPath(path string) string {
	for i, c := range path {
		if c == '/' && i > 0 {
			return path[:i]
		}
	}
	return path
}

func ensureMount(m map[string]*MountSummary, mount string) *MountSummary {
	if ms, ok := m[mount]; ok {
		return ms
	}
	ms := &MountSummary{Mount: mount}
	m[mount] = ms
	return ms
}
