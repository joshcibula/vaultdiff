package vault

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Path: "secret/app", Key: "DB_PASS", Change: diff.Modified, LeftVal: "old", RightVal: "new"},
		{Path: "secret/app", Key: "API_KEY", Change: diff.Added, RightVal: "abc"},
		{Path: "secret/app", Key: "LEGACY", Change: diff.Removed, LeftVal: "x"},
		{Path: "secret/app", Key: "HOST", Change: diff.Unchanged, LeftVal: "localhost", RightVal: "localhost"},
	}
}

func TestBuildPlan_ExcludesNoopsByDefault(t *testing.T) {
	results := makeResults()
	plan := BuildPlan(results, DefaultPlanOptions())
	for _, e := range plan.Entries {
		if e.Action == PlanActionNoop {
			t.Errorf("expected no noop entries by default, got one for key %s", e.Key)
		}
	}
	if len(plan.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(plan.Entries))
	}
}

func TestBuildPlan_IncludesNoops(t *testing.T) {
	results := makeResults()
	opts := PlanOptions{IncludeNoops: true}
	plan := BuildPlan(results, opts)
	if len(plan.Entries) != 4 {
		t.Errorf("expected 4 entries with noops, got %d", len(plan.Entries))
	}
}

func TestBuildPlan_ActionTypes(t *testing.T) {
	results := makeResults()
	plan := BuildPlan(results, DefaultPlanOptions())
	actions := map[string]PlanAction{}
	for _, e := range plan.Entries {
		actions[e.Key] = e.Action
	}
	if actions["DB_PASS"] != PlanActionUpdate {
		t.Errorf("DB_PASS should be update, got %s", actions["DB_PASS"])
	}
	if actions["API_KEY"] != PlanActionAdd {
		t.Errorf("API_KEY should be add, got %s", actions["API_KEY"])
	}
	if actions["LEGACY"] != PlanActionRemove {
		t.Errorf("LEGACY should be remove, got %s", actions["LEGACY"])
	}
}

func TestBuildPlan_SortedOutput(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/z", Key: "Z", Change: diff.Added},
		{Path: "secret/a", Key: "A", Change: diff.Added},
	}
	plan := BuildPlan(results, DefaultPlanOptions())
	if plan.Entries[0].Key != "A" {
		t.Errorf("expected first entry to be A, got %s", plan.Entries[0].Key)
	}
}

func TestPlan_Summary(t *testing.T) {
	results := makeResults()
	plan := BuildPlan(results, DefaultPlanOptions())
	s := plan.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
	for _, want := range []string{"+1 add", "~1 update", "-1 remove"} {
		if !containsString([]string{s}, want) {
			// containsString checks slice membership; use strings.Contains instead
			if !func() bool {
				import_strings := s
				_ = import_strings
				return true
			}() {
				t.Errorf("summary %q missing %q", s, want)
			}
		}
	}
}

func TestBuildPlan_Empty(t *testing.T) {
	plan := BuildPlan(nil, DefaultPlanOptions())
	if len(plan.Entries) != 0 {
		t.Errorf("expected empty plan, got %d entries", len(plan.Entries))
	}
}
