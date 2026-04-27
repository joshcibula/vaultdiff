package vault

import (
	"testing"
)

func TestDedupeSecrets_Disabled(t *testing.T) {
	secrets := map[string]map[string]string{
		"a": {"key": "1"},
		"b": {"key": "2"},
	}
	opts := DefaultDedupeOptions()
	out := DedupeSecrets(secrets, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(out))
	}
}

func TestDedupeSecrets_RemovesDuplicates(t *testing.T) {
	secrets := map[string]map[string]string{
		"path/a": {"shared": "from-a", "unique-a": "x"},
		"path/b": {"shared": "from-b", "unique-b": "y"},
	}
	opts := DedupeOptions{Enabled: true, CaseSensitive: true, PreferLongerPath: false}
	out := DedupeSecrets(secrets, opts)

	// Exactly one path should own "shared".
	owners := 0
	for _, kv := range out {
		if _, ok := kv["shared"]; ok {
			owners++
		}
	}
	if owners != 1 {
		t.Errorf("expected exactly 1 owner of 'shared', got %d", owners)
	}
}

func TestDedupeSecrets_PreferLongerPath(t *testing.T) {
	secrets := map[string]map[string]string{
		"short":       {"key": "short-val"},
		"much/longer/path": {"key": "long-val"},
	}
	opts := DedupeOptions{Enabled: true, CaseSensitive: true, PreferLongerPath: true}
	out := DedupeSecrets(secrets, opts)

	longKV, ok := out["much/longer/path"]
	if !ok {
		t.Fatal("expected longer path to be retained")
	}
	if longKV["key"] != "long-val" {
		t.Errorf("expected 'long-val', got %q", longKV["key"])
	}
	if _, ok := out["short"]; ok {
		if _, has := out["short"]["key"]; has {
			t.Error("shorter path should not own 'key' when PreferLongerPath is set")
		}
	}
}

func TestDedupeSecrets_CaseInsensitive(t *testing.T) {
	secrets := map[string]map[string]string{
		"path/a": {"Key": "val-a"},
		"path/b": {"key": "val-b"},
	}
	opts := DedupeOptions{Enabled: true, CaseSensitive: false, PreferLongerPath: false}
	out := DedupeSecrets(secrets, opts)

	owners := 0
	for _, kv := range out {
		for k := range kv {
			_ = k
			owners++
		}
	}
	if owners != 1 {
		t.Errorf("case-insensitive dedup: expected 1 key total, got %d", owners)
	}
}

func TestDedupeSecrets_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]map[string]string{
		"p": {"k": "v"},
	}
	opts := DedupeOptions{Enabled: true, CaseSensitive: true}
	_ = DedupeSecrets(secrets, opts)
	if _, ok := secrets["p"]; !ok {
		t.Error("input map was mutated")
	}
}
