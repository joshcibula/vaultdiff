package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/your-org/vaultdiff/internal/vault"
)

// renderPlan writes a human-readable plan to w.
func renderPlan(w io.Writer, plan vault.Plan) error {
	if len(plan.Entries) == 0 {
		_, err := fmt.Fprintln(w, "No changes planned.")
		return err
	}
	_, err := fmt.Fprintf(w, "Plan: %s\n\n", plan.Summary())
	if err != nil {
		return err
	}
	currentPath := ""
	for _, e := range plan.Entries {
		if e.Path != currentPath {
			currentPath = e.Path
			_, err = fmt.Fprintf(w, "  [%s]\n", currentPath)
			if err != nil {
				return err
			}
		}
		switch e.Action {
		case vault.PlanActionAdd:
			_, err = fmt.Fprintf(w, "    + %s = %s\n", e.Key, e.NewVal)
		case vault.PlanActionRemove:
			_, err = fmt.Fprintf(w, "    - %s (was: %s)\n", e.Key, e.OldVal)
		case vault.PlanActionUpdate:
			_, err = fmt.Fprintf(w, "    ~ %s: %s -> %s\n", e.Key, e.OldVal, e.NewVal)
		case vault.PlanActionNoop:
			_, err = fmt.Fprintf(w, "    = %s\n", e.Key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// planActionSymbol returns a short symbol for use in compact output.
func planActionSymbol(a vault.PlanAction) string {
	symbols := map[vault.PlanAction]string{
		vault.PlanActionAdd:    "+",
		vault.PlanActionRemove: "-",
		vault.PlanActionUpdate: "~",
		vault.PlanActionNoop:   "=",
	}
	if s, ok := symbols[a]; ok {
		return s
	}
	return "?"
}

// PlanToLines converts a plan to a slice of formatted strings.
func PlanToLines(plan vault.Plan) []string {
	lines := make([]string, 0, len(plan.Entries))
	for _, e := range plan.Entries {
		lines = append(lines, fmt.Sprintf("%s %s/%s", planActionSymbol(e.Action), e.Path, e.Key))
	}
	return lines
}

// PlanToCompact returns a compact single-line summary of the plan.
func PlanToCompact(plan vault.Plan) string {
	lines := PlanToLines(plan)
	if len(lines) == 0 {
		return "no changes"
	}
	return strings.Join(lines, "; ")
}
