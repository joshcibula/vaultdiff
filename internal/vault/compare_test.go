package vault

import (
	"testing"
)

func TestCompareSecrets_EqualMaps(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"foo": "bar"},
	}
	right := map[string]map[string]string{
		"secret/a": {"foo": "bar"},
	}
	results := CompareSecrets(left, right, DefaultCompareOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Equal {
		t.Error("expected Equal=true")
	}
}

func TestCompareSecrets_DifferentValues(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"key": "v1"},
	}
	right := map[string]map[string]string{
		"secret/a": {"key": "v2"},
	}
	results := CompareSecrets(left, right, DefaultCompareOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Equal {
		t.Error("expected Equal=false")
	}
}

func TestCompareSecrets_IgnoreCase(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"key": "Hello"},
	}
	right := map[string]map[string]string{
		"secret/a": {"key": "hello"},
	}
	opts := DefaultCompareOptions()
	opts.IgnoreCase = true
	results := CompareSecrets(left, right, opts)
	if !results[0].Equal {
		t.Error("expected Equal=true with IgnoreCase")
	}
}

func TestCompareSecrets_IgnoreWhitespace(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"key": "  value  "},
	}
	right := map[string]map[string]string{
		"secret/a": {"key": "value"},
	}
	opts := DefaultCompareOptions()
	opts.IgnoreWhitespace = true
	results := CompareSecrets(left, right, opts)
	if !results[0].Equal {
		t.Error("expected Equal=true with IgnoreWhitespace")
	}
}

func TestCompareSecrets_IgnoreKeys(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"skip": "x", "keep": "same"},
	}
	right := map[string]map[string]string{
		"secret/a": {"skip": "y", "keep": "same"},
	}
	opts := DefaultCompareOptions()
	opts.IgnoreKeys = []string{"skip"}
	results := CompareSecrets(left, right, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result after ignore, got %d", len(results))
	}
	if results[0].Key != "keep" {
		t.Errorf("expected key 'keep', got %q", results[0].Key)
	}
}

func TestCompareSecrets_MissingOnOneSide(t *testing.T) {
	left := map[string]map[string]string{
		"secret/a": {"only_left": "val"},
	}
	right := map[string]map[string]string{}
	results := CompareSecrets(left, right, DefaultCompareOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Equal {
		t.Error("expected Equal=false for missing key on right")
	}
}
