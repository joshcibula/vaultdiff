package vault

import (
	"testing"
)

func TestPromoteSecrets_Disabled(t *testing.T) {
	src := map[string]map[string]string{"a": {"x": "1"}}
	dst := map[string]map[string]string{}
	opts := DefaultPromoteOptions() // Enabled=false

	out, results, err := PromoteSecrets(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
	if len(out) != 0 {
		t.Errorf("expected unchanged dst, got %v", out)
	}
}

func TestPromoteSecrets_DryRun(t *testing.T) {
	src := map[string]map[string]string{"secret/app": {"db_pass": "s3cr3t"}}
	dst := map[string]map[string]string{}
	opts := PromoteOptions{Enabled: true, DryRun: true}

	out, results, err := PromoteSecrets(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Skipped {
		t.Errorf("expected 1 non-skipped result, got %+v", results)
	}
	// dry run: dst should remain empty
	if _, ok := out["secret/app"]; ok {
		t.Error("dry run should not mutate destination")
	}
}

func TestPromoteSecrets_AppliesChanges(t *testing.T) {
	src := map[string]map[string]string{"secret/app": {"key": "val"}}
	dst := map[string]map[string]string{}
	opts := PromoteOptions{Enabled: true, DryRun: false, Overwrite: true}

	out, results, err := PromoteSecrets(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if out["secret/app"]["key"] != "val" {
		t.Errorf("expected promoted value, got %q", out["secret/app"]["key"])
	}
}

func TestPromoteSecrets_SkipsConflictWithoutOverwrite(t *testing.T) {
	src := map[string]map[string]string{"secret/app": {"key": "new"}}
	dst := map[string]map[string]string{"secret/app": {"key": "existing"}}
	opts := PromoteOptions{Enabled: true, DryRun: false, Overwrite: false}

	out, results, err := PromoteSecrets(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected skipped result, got %+v", results)
	}
	if out["secret/app"]["key"] != "existing" {
		t.Errorf("expected original value preserved, got %q", out["secret/app"]["key"])
	}
}

func TestPromoteSecrets_WithPathPrefix(t *testing.T) {
	src := map[string]map[string]string{"app": {"token": "abc"}}
	dst := map[string]map[string]string{}
	opts := PromoteOptions{Enabled: true, DryRun: false, Overwrite: false, PathPrefix: "prod"}

	out, _, err := PromoteSecrets(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["prod/app"]["token"] != "abc" {
		t.Errorf("expected key under prefixed path, got %v", out)
	}
}

func TestPromoteSecrets_DoesNotMutateInput(t *testing.T) {
	src := map[string]map[string]string{"s": {"k": "v"}}
	dst := map[string]map[string]string{"s": {"other": "x"}}
	opts := PromoteOptions{Enabled: true, DryRun: false, Overwrite: true}

	_, _, _ = PromoteSecrets(src, dst, opts)

	if _, ok := dst["s"]["k"]; ok {
		t.Error("original dst map must not be mutated")
	}
}
