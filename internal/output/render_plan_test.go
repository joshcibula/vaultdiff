package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/vault"
)

func buildTestPlan(includeNoops bool) vault.Plan {
	results := []diff.Result{
		{Path: "secret/app", Key: "API_KEY", Change: diff.Added, RightVal: "newkey"},
		{Path: "secret/app", Key: "DB_PASS", Change: diff.Modified, LeftVal: "old", RightVal: "new"},
		{Path: "secret/app", Key: "LEGACY", Change: diff.Removed, LeftVal: "gone"},
	}
	if includeNoops {
		results = append(results, diff.Result{Path: "secret/app", Key: "HOST", Change: diff.Unchanged})
	}
	return vault.BuildPlan(results, vault.PlanOptions{IncludeNoops: includeNoops})
}

func TestRenderPlan_NoChanges(t *testing.T) {
	plan := vault.BuildPlan(nil, vault.DefaultPlanOptions())
	var buf bytes.Buffer
	if err := renderPlan(&buf, plan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes', got: %s", buf.String())
	}
}

func TestRenderPlan_ContainsSummaryLine(t *testing.T) {
	plan := buildTestPlan(false)
	var buf bytes.Buffer
	if err := renderPlan(&buf, plan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Plan:") {
		t.Errorf("expected 'Plan:' header, got: %s", out)
	}
}

func TestRenderPlan_ContainsActionSymbols(t *testing.T) {
	plan := buildTestPlan(false)
	var buf bytes.Buffer
	_ = renderPlan(&buf, plan)
	out := buf.String()
	for _, sym := range []string{"+", "-", "~"} {
		if !strings.Contains(out, sym) {
			t.Errorf("expected symbol %q in output:\n%s", sym, out)
		}
	}
}

func TestPlanToLines_ReturnsOnePerEntry(t *testing.T) {
	plan := buildTestPlan(false)
	lines := PlanToLines(plan)
	if len(lines) != len(plan.Entries) {
		t.Errorf("expected %d lines, got %d", len(plan.Entries), len(lines))
	}
}

func TestPlanToCompact_Empty(t *testing.T) {
	plan := vault.BuildPlan(nil, vault.DefaultPlanOptions())
	out := PlanToCompact(plan)
	if out != "no changes" {
		t.Errorf("expected 'no changes', got %q", out)
	}
}

func TestPlanToCompact_NonEmpty(t *testing.T) {
	plan := buildTestPlan(false)
	out := PlanToCompact(plan)
	if out == "" || out == "no changes" {
		t.Errorf("expected non-empty compact output, got %q", out)
	}
	if !strings.Contains(out, "/") {
		t.Errorf("expected path/key format in compact output, got %q", out)
	}
}
