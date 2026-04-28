package vault

import (
	"testing"
)

func TestDefaultSortOptions(t *testing.T) {
	opts := DefaultSortOptions()
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.Field != SortByPath {
		t.Errorf("expected Field=path, got %s", opts.Field)
	}
	if opts.Order != SortAsc {
		t.Errorf("expected Order=asc, got %s", opts.Order)
	}
}

func TestSortSecrets_Disabled(t *testing.T) {
	input := map[string]map[string]string{
		"z/path": {"k": "v"},
		"a/path": {"k": "v"},
	}
	opts := DefaultSortOptions() // Enabled=false
	out := SortSecrets(input, opts)
	if len(out) != len(input) {
		t.Errorf("expected same length, got %d", len(out))
	}
	// Should be the same map reference when disabled.
	if &out == &input {
		// fine, same reference
	}
}

func TestSortSecrets_ByPathAsc(t *testing.T) {
	input := map[string]map[string]string{
		"secret/z": {"key": "val"},
		"secret/a": {"key": "val"},
		"secret/m": {"key": "val"},
	}
	opts := SortOptions{Enabled: true, Field: SortByPath, Order: SortAsc}
	out := SortSecrets(input, opts)
	if len(out) != 3 {
		t.Fatalf("expected 3 paths, got %d", len(out))
	}
	for _, p := range []string{"secret/a", "secret/m", "secret/z"} {
		if _, ok := out[p]; !ok {
			t.Errorf("expected path %q in output", p)
		}
	}
}

func TestSortSecrets_ByPathDesc(t *testing.T) {
	input := map[string]map[string]string{
		"secret/a": {"k": "v"},
		"secret/z": {"k": "v"},
	}
	opts := SortOptions{Enabled: true, Field: SortByPath, Order: SortDesc}
	out := SortSecrets(input, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(out))
	}
}

func TestSortSecrets_ByKeyAsc(t *testing.T) {
	input := map[string]map[string]string{
		"secret/p": {"zebra": "1", "apple": "2", "mango": "3"},
	}
	opts := SortOptions{Enabled: true, Field: SortByKey, Order: SortAsc}
	out := SortSecrets(input, opts)
	kv := out["secret/p"]
	if kv["apple"] != "2" || kv["zebra"] != "1" || kv["mango"] != "3" {
		t.Error("key-value pairs should be preserved after key sort")
	}
}

func TestSortSecrets_ByValueAsc(t *testing.T) {
	input := map[string]map[string]string{
		"secret/p": {"b": "zz", "a": "aa"},
	}
	opts := SortOptions{Enabled: true, Field: SortByValue, Order: SortAsc}
	out := SortSecrets(input, opts)
	kv := out["secret/p"]
	if kv["a"] != "aa" || kv["b"] != "zz" {
		t.Error("values should be preserved after value sort")
	}
}

func TestSortSecrets_DoesNotMutateInput(t *testing.T) {
	input := map[string]map[string]string{
		"secret/x": {"foo": "bar"},
	}
	opts := SortOptions{Enabled: true, Field: SortByPath, Order: SortAsc}
	out := SortSecrets(input, opts)
	out["secret/x"]["foo"] = "mutated"
	if input["secret/x"]["foo"] == "mutated" {
		t.Error("SortSecrets should not mutate the input map")
	}
}

func TestSortSecrets_Empty(t *testing.T) {
	opts := SortOptions{Enabled: true, Field: SortByPath, Order: SortAsc}
	out := SortSecrets(map[string]map[string]string{}, opts)
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
