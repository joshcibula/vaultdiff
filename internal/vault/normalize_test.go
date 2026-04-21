package vault

import (
	"testing"
)

func TestNormalizeSecrets_NoOp(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {"KEY": "value"},
	}
	opts := DefaultNormalizeOptions()
	result := NormalizeSecrets(input, opts)

	if result["secret/app"]["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", result["secret/app"]["KEY"])
	}
}

func TestNormalizeSecrets_TrimKeyPrefix(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {"APP_DB_HOST": "localhost", "APP_PORT": "5432"},
	}
	opts := NormalizeOptions{TrimKeyPrefix: "APP_"}
	result := NormalizeSecrets(input, opts)

	kv := result["secret/app"]
	if _, ok := kv["DB_HOST"]; !ok {
		t.Error("expected key 'DB_HOST' after prefix trim")
	}
	if _, ok := kv["PORT"]; !ok {
		t.Error("expected key 'PORT' after prefix trim")
	}
	if _, ok := kv["APP_DB_HOST"]; ok {
		t.Error("original key should not exist after prefix trim")
	}
}

func TestNormalizeSecrets_StripTrailingSlash(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app/": {"key": "val"},
	}
	opts := NormalizeOptions{StripTrailingSlash: true}
	result := NormalizeSecrets(input, opts)

	if _, ok := result["secret/app"]; !ok {
		t.Error("expected path without trailing slash")
	}
	if _, ok := result["secret/app/"]; ok {
		t.Error("original path with trailing slash should not exist")
	}
}

func TestNormalizeSecrets_CollapseWhitespace(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {"desc": "  hello   world  "},
	}
	opts := NormalizeOptions{CollapseWhitespace: true}
	result := NormalizeSecrets(input, opts)

	got := result["secret/app"]["desc"]
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestNormalizeSecrets_DoesNotMutateInput(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {"APP_KEY": "  value  "},
	}
	opts := NormalizeOptions{TrimKeyPrefix: "APP_", CollapseWhitespace: true}
	NormalizeSecrets(input, opts)

	if _, ok := input["secret/app"]["APP_KEY"]; !ok {
		t.Error("original input should not be mutated")
	}
	if input["secret/app"]["APP_KEY"] != "  value  " {
		t.Error("original value should not be mutated")
	}
}

func TestCollapseWhitespace(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"hello world", "hello world"},
		{"  spaces  ", "spaces"},
		{"tab\there", "tab here"},
		{"", ""},
	}
	for _, tc := range cases {
		got := collapseWhitespace(tc.input)
		if got != tc.expected {
			t.Errorf("collapseWhitespace(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
