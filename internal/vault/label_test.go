package vault

import (
	"testing"
)

func TestLabelSecrets_NoOp(t *testing.T) {
	input := map[string]map[string]string{
		"secret/foo": {"key": "val"},
	}
	opts := DefaultLabelOptions()
	out := LabelSecrets(input, opts)
	if _, ok := out["secret/foo"]; !ok {
		t.Fatal("expected path to be unchanged")
	}
}

func TestLabelSecrets_Prefix(t *testing.T) {
	input := map[string]map[string]string{
		"foo": {"a": "1"},
		"bar": {"b": "2"},
	}
	opts := DefaultLabelOptions()
	opts.Prefix = "prod/"
	out := LabelSecrets(input, opts)
	if _, ok := out["prod/foo"]; !ok {
		t.Error("expected prod/foo")
	}
	if _, ok := out["prod/bar"]; !ok {
		t.Error("expected prod/bar")
	}
}

func TestLabelSecrets_StripPrefix(t *testing.T) {
	input := map[string]map[string]string{
		"kv/data/myapp/db": {"pass": "secret"},
	}
	opts := DefaultLabelOptions()
	opts.StripPrefix = "kv/data/"
	out := LabelSecrets(input, opts)
	if _, ok := out["myapp/db"]; !ok {
		t.Error("expected stripped path myapp/db")
	}
}

func TestLabelSecrets_Alias(t *testing.T) {
	input := map[string]map[string]string{
		"secret/db": {"user": "admin"},
	}
	opts := DefaultLabelOptions()
	opts.Alias = map[string]string{"secret/db": "Database Credentials"}
	out := LabelSecrets(input, opts)
	if _, ok := out["Database Credentials"]; !ok {
		t.Error("expected alias 'Database Credentials'")
	}
}

func TestLabelSecrets_AliasBeatsPrefix(t *testing.T) {
	input := map[string]map[string]string{
		"secret/db": {"user": "admin"},
	}
	opts := DefaultLabelOptions()
	opts.Prefix = "prod/"
	opts.Alias = map[string]string{"secret/db": "DB"}
	out := LabelSecrets(input, opts)
	if _, ok := out["DB"]; !ok {
		t.Error("expected alias to take precedence over prefix")
	}
	if _, ok := out["prod/secret/db"]; ok {
		t.Error("prefix should not be applied when alias matches")
	}
}

func TestLabelSecrets_DoesNotMutateInput(t *testing.T) {
	input := map[string]map[string]string{
		"path/a": {"x": "1"},
	}
	opts := DefaultLabelOptions()
	opts.Prefix = "env/"
	LabelSecrets(input, opts)
	if _, ok := input["path/a"]; !ok {
		t.Error("original map should not be mutated")
	}
}
