package vault

import (
	"testing"
)

func TestFlattenSecrets_Disabled(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {"db.host": "localhost", "db.port": "5432"},
	}
	opts := DefaultFlattenOptions()
	opts.Enabled = false

	result := FlattenSecrets(input, opts)
	if result["secret/app"]["db.host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %q", result["secret/app"]["db.host"])
	}
}

func TestFlattenSecrets_PassthroughStringValues(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {
			"host": "localhost",
			"port": "5432",
		},
	}
	opts := DefaultFlattenOptions()
	opts.Enabled = true

	result := FlattenSecrets(input, opts)
	if result["secret/app"]["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %q", result["secret/app"]["host"])
	}
	if result["secret/app"]["port"] != "5432" {
		t.Errorf("expected port=5432, got %q", result["secret/app"]["port"])
	}
}

func TestFlattenSecrets_DoesNotMutateInput(t *testing.T) {
	original := map[string]map[string]string{
		"secret/app": {"key": "value"},
	}
	opts := DefaultFlattenOptions()
	opts.Enabled = true

	_ = FlattenSecrets(original, opts)
	if original["secret/app"]["key"] != "value" {
		t.Error("input was mutated")
	}
}

func TestFlattenSecrets_MultiplePathsPreserved(t *testing.T) {
	input := map[string]map[string]string{
		"secret/a": {"x": "1"},
		"secret/b": {"y": "2"},
	}
	opts := DefaultFlattenOptions()
	opts.Enabled = true

	result := FlattenSecrets(input, opts)
	if len(result) != 2 {
		t.Errorf("expected 2 paths, got %d", len(result))
	}
	if result["secret/a"]["x"] != "1" {
		t.Errorf("expected x=1, got %q", result["secret/a"]["x"])
	}
	if result["secret/b"]["y"] != "2" {
		t.Errorf("expected y=2, got %q", result["secret/b"]["y"])
	}
}

func TestDefaultFlattenOptions(t *testing.T) {
	opts := DefaultFlattenOptions()
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.Separator != "." {
		t.Errorf("expected separator='.', got %q", opts.Separator)
	}
	if opts.MaxDepth != 10 {
		t.Errorf("expected MaxDepth=10, got %d", opts.MaxDepth)
	}
}

func TestFlattenMap_SingleLevel(t *testing.T) {
	m := map[string]interface{}{
		"host": "localhost",
		"port": "5432",
	}
	out := make(map[string]string)
	flattenMap("", m, ".", 0, 10, out)
	if out["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %q", out["host"])
	}
	if out["port"] != "5432" {
		t.Errorf("expected port=5432, got %q", out["port"])
	}
}

func TestFlattenMap_Nested(t *testing.T) {
	m := map[string]interface{}{
		"db": map[string]interface{}{
			"host": "localhost",
			"port": "5432",
		},
	}
	out := make(map[string]string)
	flattenMap("", m, ".", 0, 10, out)
	if out["db.host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %q", out["db.host"])
	}
	if out["db.port"] != "5432" {
		t.Errorf("expected db.port=5432, got %q", out["db.port"])
	}
}
