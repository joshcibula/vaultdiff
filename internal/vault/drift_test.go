package vault

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestDefaultDriftOptions(t *testing.T) {
	opts := DefaultDriftOptions()
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.Threshold != 10.0 {
		t.Errorf("expected Threshold=10.0, got %v", opts.Threshold)
	}
}

func TestDetectDrift_DisabledReturnsEmpty(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/a", Key: "foo", Change: diff.Modified},
	}
	opts := DefaultDriftOptions() // Enabled=false
	report := DetectDrift(results, opts)
	if report.TotalKeys != 0 || report.ChangedKeys != 0 {
		t.Error("expected empty report when disabled")
	}
}

func TestDetectDrift_NoChanges(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/a", Key: "x", Change: diff.NoChange},
		{Path: "secret/b", Key: "y", Change: diff.NoChange},
	}
	opts := DriftOptions{Enabled: true, Threshold: 10.0}
	report := DetectDrift(results, opts)
	if report.ChangedKeys != 0 {
		t.Errorf("expected 0 changed keys, got %d", report.ChangedKeys)
	}
	if report.DriftPercent != 0 {
		t.Errorf("expected 0%% drift, got %v", report.DriftPercent)
	}
	if report.Significant {
		t.Error("expected Significant=false")
	}
}

func TestDetectDrift_SignificantDrift(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/a", Key: "k1", Change: diff.Modified},
		{Path: "secret/a", Key: "k2", Change: diff.Modified},
		{Path: "secret/b", Key: "k3", Change: diff.NoChange},
	}
	opts := DriftOptions{Enabled: true, Threshold: 50.0}
	report := DetectDrift(results, opts)
	if report.TotalKeys != 3 {
		t.Errorf("expected TotalKeys=3, got %d", report.TotalKeys)
	}
	if report.ChangedKeys != 2 {
		t.Errorf("expected ChangedKeys=2, got %d", report.ChangedKeys)
	}
	expected := 2.0 / 3.0 * 100
	if report.DriftPercent != expected {
		t.Errorf("expected DriftPercent=%v, got %v", expected, report.DriftPercent)
	}
	if !report.Significant {
		t.Error("expected Significant=true")
	}
}

func TestDetectDrift_IgnorePaths(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/ignored/a", Key: "k", Change: diff.Modified},
		{Path: "secret/keep/b", Key: "k", Change: diff.NoChange},
	}
	opts := DriftOptions{Enabled: true, Threshold: 10.0, IgnorePaths: []string{"secret/ignored"}}
	report := DetectDrift(results, opts)
	if report.TotalKeys != 1 {
		t.Errorf("expected TotalKeys=1 after ignore, got %d", report.TotalKeys)
	}
	if report.ChangedKeys != 0 {
		t.Errorf("expected ChangedKeys=0, got %d", report.ChangedKeys)
	}
}

func TestDriftReport_String(t *testing.T) {
	results := []diff.Result{
		{Path: "secret/x", Key: "a", Change: diff.OnlyInLeft},
	}
	opts := DriftOptions{Enabled: true, Threshold: 5.0}
	report := DetectDrift(results, opts)
	s := report.String()
	if len(s) == 0 {
		t.Error("expected non-empty String()")
	}
}
