package vault

import (
	"testing"
)

func TestPivotSecrets_Disabled(t *testing.T) {
	input := map[string]map[string]string{
		"secret/a": {"host": "localhost", "port": "5432"},
	}
	opts := DefaultPivotOptions()
	out := PivotSecrets(input, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 path, got %d", len(out))
	}
	if out["secret/a"]["host"] != "localhost" {
		t.Errorf("expected original map to be returned unchanged")
	}
}

func TestPivotSecrets_AllKeys(t *testing.T) {
	input := map[string]map[string]string{
		"secret/a": {"host": "host-a", "port": "5432"},
		"secret/b": {"host": "host-b", "port": "5433"},
	}
	opts := PivotOptions{Enabled: true}
	out := PivotSecrets(input, opts)

	if len(out) != 2 {
		t.Fatalf("expected 2 pivoted keys, got %d", len(out))
	}
	if out["host"]["secret/a"] != "host-a" {
		t.Errorf("expected host-a, got %q", out["host"]["secret/a"])
	}
	if out["host"]["secret/b"] != "host-b" {
		t.Errorf("expected host-b, got %q", out["host"]["secret/b"])
	}
	if out["port"]["secret/a"] != "5432" {
		t.Errorf("expected 5432, got %q", out["port"]["secret/a"])
	}
}

func TestPivotSecrets_KeyField(t *testing.T) {
	input := map[string]map[string]string{
		"secret/a": {"host": "host-a", "port": "5432"},
		"secret/b": {"host": "host-b", "port": "5433"},
	}
	opts := PivotOptions{Enabled: true, KeyField: "host"}
	out := PivotSecrets(input, opts)

	if len(out) != 1 {
		t.Fatalf("expected 1 pivoted key, got %d", len(out))
	}
	if _, ok := out["port"]; ok {
		t.Error("port key should have been filtered out")
	}
	if out["host"]["secret/a"] != "host-a" {
		t.Errorf("expected host-a, got %q", out["host"]["secret/a"])
	}
}

func TestPivotSecrets_PathPrefix(t *testing.T) {
	input := map[string]map[string]string{
		"secret/a": {"host": "host-a"},
	}
	opts := PivotOptions{Enabled: true, PathPrefix: "pivot"}
	out := PivotSecrets(input, opts)

	if _, ok := out["pivot/host"]; !ok {
		t.Fatalf("expected pivot/host key, got keys: %v", keys(out))
	}
	if out["pivot/host"]["secret/a"] != "host-a" {
		t.Errorf("expected host-a, got %q", out["pivot/host"]["secret/a"])
	}
}

func TestPivotSecrets_DoesNotMutateInput(t *testing.T) {
	input := map[string]map[string]string{
		"secret/a": {"host": "host-a"},
	}
	opts := PivotOptions{Enabled: true}
	PivotSecrets(input, opts)

	if _, ok := input["secret/a"]; !ok {
		t.Error("original input was mutated")
	}
}

func keys(m map[string]map[string]string) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
