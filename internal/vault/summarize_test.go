package vault

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestSummarizeDiff_Empty(t *testing.T) {
	results := []diff.Result{}
	summary := SummarizeDiff(results, DefaultSummaryOptions())

	if summary.TotalAdded != 0 || summary.TotalRemoved != 0 || summary.TotalModified != 0 {
		t.Errorf("expected all zeros for empty results, got %s", summary.String())
	}
}

func TestSummarizeDiff_Counts(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/foo", Change: diff.Added},
		{Path: "secret/bar", Change: diff.Removed},
		{Path: "secret/baz", Change: diff.Modified},
		{Path: "secret/qux", Change: diff.Unchanged},
	}

	summary := SummarizeDiff(results, DefaultSummaryOptions())

	if summary.TotalAdded != 1 {
		t.Errorf("expected 1 added, got %d", summary.TotalAdded)
	}
	if summary.TotalRemoved != 1 {
		t.Errorf("expected 1 removed, got %d", summary.TotalRemoved)
	}
	if summary.TotalModified != 1 {
		t.Errorf("expected 1 modified, got %d", summary.TotalModified)
	}
	if summary.TotalUnchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", summary.TotalUnchanged)
	}
}

func TestSummarizeDiff_ByMount(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/foo", Change: diff.Added},
		{Path: "secret/bar", Change: diff.Modified},
		{Path: "kv/alpha", Change: diff.Removed},
	}

	summary := SummarizeDiff(results, DefaultSummaryOptions())

	if len(summary.ByMount) != 2 {
		t.Fatalf("expected 2 mounts, got %d", len(summary.ByMount))
	}

	secret := summary.ByMount["secret"]
	if secret == nil {
		t.Fatal("expected 'secret' mount in summary")
	}
	if secret.Added != 1 || secret.Modified != 1 || secret.Total != 2 {
		t.Errorf("unexpected secret mount stats: %+v", secret)
	}

	kv := summary.ByMount["kv"]
	if kv == nil {
		t.Fatal("expected 'kv' mount in summary")
	}
	if kv.Removed != 1 || kv.Total != 1 {
		t.Errorf("unexpected kv mount stats: %+v", kv)
	}
}

func TestSummarizeDiff_String(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/a", Change: diff.Added},
		{Path: "secret/b", Change: diff.Removed},
	}

	summary := SummarizeDiff(results, DefaultSummaryOptions())
	str := summary.String()

	if str == "" {
		t.Error("expected non-empty summary string")
	}
}

func TestSummarizeDiff_SortedMounts(t *testing.T) {
	results := []diff.Result{
		{Path: "z-mount/x", Change: diff.Added},
		{Path: "a-mount/y", Change: diff.Added},
		{Path: "m-mount/z", Change: diff.Added},
	}

	summary := SummarizeDiff(results, DefaultSummaryOptions())
	mounts := summary.SortedMounts()

	if len(mounts) != 3 {
		t.Fatalf("expected 3 mounts, got %d", len(mounts))
	}
	if mounts[0] != "a-mount" || mounts[1] != "m-mount" || mounts[2] != "z-mount" {
		t.Errorf("mounts not sorted: %v", mounts)
	}
}
