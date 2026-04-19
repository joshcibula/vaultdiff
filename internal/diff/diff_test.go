package diff

import (
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	left := SecretMap{"key1": "val1", "key2": "val2"}
	right := SecretMap{"key1": "val1", "key2": "val2"}

	res := Compare(left, right)

	if res.HasDifferences() {
		t.Errorf("expected no differences, got: %s", res.Summary())
	}
	if len(res.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(res.Unchanged))
	}
}

func TestCompare_Modified(t *testing.T) {
	left := SecretMap{"key1": "old"}
	right := SecretMap{"key1": "new"}

	res := Compare(left, right)

	if !res.HasDifferences() {
		t.Error("expected differences")
	}
	pair, ok := res.Modified["key1"]
	if !ok {
		t.Fatal("expected key1 in Modified")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected values: %v", pair)
	}
}

func TestCompare_OnlyInLeft(t *testing.T) {
	left := SecretMap{"gone": "val"}
	right := SecretMap{}

	res := Compare(left, right)

	if _, ok := res.OnlyInLeft["gone"]; !ok {
		t.Error("expected 'gone' in OnlyInLeft")
	}
}

func TestCompare_OnlyInRight(t *testing.T) {
	left := SecretMap{}
	right := SecretMap{"new": "val"}

	res := Compare(left, right)

	if _, ok := res.OnlyInRight["new"]; !ok {
		t.Error("expected 'new' in OnlyInRight")
	}
}

func TestResult_Summary(t *testing.T) {
	left := SecretMap{"a": "1", "b": "2"}
	right := SecretMap{"b": "changed", "c": "3"}

	res := Compare(left, right)
	summary := res.Summary()

	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
