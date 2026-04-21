package vault

import (
	"testing"
)

func TestMergeSecrets_NoOverlap(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"key1": "val1"},
	}
	right := map[string]map[string]string{
		"secret/b": {"key2": "val2"},
	}
	result := MergeSecrets(left, right, DefaultMergeOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(result))
	}
	if result["secret/a"]["key1"] != "val1" {
		t.Errorf("expected val1, got %s", result["secret/a"]["key1"])
	}
	if result["secret/b"]["key2"] != "val2" {
		t.Errorf("expected val2, got %s", result["secret/b"]["key2"])
	}
}

func TestMergeSecrets_RightWinsOnConflict(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"key": "from-left"},
	}
	right := map[string]map[string]string{
		"secret/a": {"key": "from-right"},
	}
	opts := DefaultMergeOptions() // PreferLeft = false
	result := MergeSecrets(left, right, opts)
	if result["secret/a"]["key"] != "from-right" {
		t.Errorf("expected from-right, got %s", result["secret/a"]["key"])
	}
}

func TestMergeSecrets_LeftWinsOnConflict(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"key": "from-left"},
	}
	right := map[string]map[string]string{
		"secret/a": {"key": "from-right"},
	}
	opts := MergeOptions{PreferLeft: true}
	result := MergeSecrets(left, right, opts)
	if result["secret/a"]["key"] != "from-left" {
		t.Errorf("expected from-left, got %s", result["secret/a"]["key"])
	}
}

func TestMergeSecrets_SkipEmpty(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"present": "yes", "empty": ""},
	}
	right := map[string]map[string]string{}
	opts := MergeOptions{SkipEmpty: true}
	result := MergeSecrets(left, right, opts)
	if _, ok := result["secret/a"]["empty"]; ok {
		t.Error("expected empty key to be skipped")
	}
	if result["secret/a"]["present"] != "yes" {
		t.Errorf("expected present key to survive, got %s", result["secret/a"]["present"])
	}
}

func TestMergeSecrets_DoesNotMutateInput(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"key": "original"},
	}
	right := map[string]map[string]string{
		"secret/a": {"key": "override"},
	}
	_ = MergeSecrets(left, right, DefaultMergeOptions())
	if left["secret/a"]["key"] != "original" {
		t.Error("MergeSecrets mutated the left input map")
	}
}
