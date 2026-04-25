package vault

import (
	"fmt"
	"sort"
	"strings"

	"github.com/your-org/vaultdiff/internal/diff"
)

// PlanAction represents the type of change in a plan.
type PlanAction string

const (
	PlanActionAdd    PlanAction = "add"
	PlanActionRemove PlanAction = "remove"
	PlanActionUpdate PlanAction = "update"
	PlanActionNoop   PlanAction = "noop"
)

// PlanEntry describes a single planned change.
type PlanEntry struct {
	Path   string
	Key    string
	Action PlanAction
	OldVal string
	NewVal string
}

// Plan holds all planned changes derived from a diff result.
type Plan struct {
	Entries []PlanEntry
}

// DefaultPlanOptions returns sensible defaults.
func DefaultPlanOptions() PlanOptions {
	return PlanOptions{IncludeNoops: false}
}

// PlanOptions controls plan generation behaviour.
type PlanOptions struct {
	IncludeNoops bool
}

// BuildPlan converts diff results into an ordered execution plan.
func BuildPlan(results []diff.Result, opts PlanOptions) Plan {
	var entries []PlanEntry
	for _, r := range results {
		switch r.Change {
		case diff.Added:
			entries = append(entries, PlanEntry{Path: r.Path, Key: r.Key, Action: PlanActionAdd, NewVal: r.RightVal})
		case diff.Removed:
			entries = append(entries, PlanEntry{Path: r.Path, Key: r.Key, Action: PlanActionRemove, OldVal: r.LeftVal})
		case diff.Modified:
			entries = append(entries, PlanEntry{Path: r.Path, Key: r.Key, Action: PlanActionUpdate, OldVal: r.LeftVal, NewVal: r.RightVal})
		case diff.Unchanged:
			if opts.IncludeNoops {
				entries = append(entries, PlanEntry{Path: r.Path, Key: r.Key, Action: PlanActionNoop})
			}
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		ki := entries[i].Path + "/" + entries[i].Key
		kj := entries[j].Path + "/" + entries[j].Key
		return ki < kj
	})
	return Plan{Entries: entries}
}

// Summary returns a human-readable plan summary.
func (p Plan) Summary() string {
	counts := map[PlanAction]int{}
	for _, e := range p.Entries {
		counts[e.Action]++
	}
	parts := []string{
		fmt.Sprintf("+%d add", counts[PlanActionAdd]),
		fmt.Sprintf("~%d update", counts[PlanActionUpdate]),
		fmt.Sprintf("-%d remove", counts[PlanActionRemove]),
	}
	if counts[PlanActionNoop] > 0 {
		parts = append(parts, fmt.Sprintf("=%d noop", counts[PlanActionNoop]))
	}
	return strings.Join(parts, "  ")
}
