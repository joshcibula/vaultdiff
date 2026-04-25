package vault

import (
	"testing"
)

func TestAnnotateSecrets_Disabled(t *testing.T) {
	opts := DefaultAnnotateOptions()
	input := map[string]map[string]string{
		"secret/foo": {"key": "value"},
	}
	out := AnnotateSecrets(input, opts)
	if _, ok := out["secret/foo"]["_vaultdiff_source"]; ok {
		t.Error("expected no annotation when disabled")
	}
}

func TestAnnotateSecrets_InjectsTagKey(t *testing.T) {
	opts := DefaultAnnotateOptions()
	opts.Enabled = true
	opts.TagValue = "production"
	input := map[string]map[string]string{
		"secret/foo": {"key": "value"},
	}
	out := AnnotateSecrets(input, opts)
	if got := out["secret/foo"]["_vaultdiff_source"]; got != "production" {
		t.Errorf("expected 'production', got %q", got)
	}
}

func TestAnnotateSecrets_UsesPathWhenNoTagValue(t *testing.T) {
	opts := DefaultAnnotateOptions()
	opts.Enabled = true
	opts.PathPrefix = "secret/"
	input := map[string]map[string]string{
		"secret/myapp": {"db": "pass"},
	}
	out := AnnotateSecrets(input, opts)
	if got := out["secret/myapp"]["_vaultdiff_source"]; got != "myapp" {
		t.Errorf("expected 'myapp', got %q", got)
	}
}

func TestAnnotateSecrets_CustomTags(t *testing.T) {
	opts := DefaultAnnotateOptions()
	opts.Enabled = true
	opts.TagValue = "staging"
	opts.CustomTags = map[string]string{"env": "staging", "team": "platform"}
	input := map[string]map[string]string{
		"secret/svc": {"token": "abc"},
	}
	out := AnnotateSecrets(input, opts)
	if out["secret/svc"]["env"] != "staging" {
		t.Error("expected custom tag 'env'")
	}
	if out["secret/svc"]["team"] != "platform" {
		t.Error("expected custom tag 'team'")
	}
}

func TestAnnotateSecrets_DoesNotMutateInput(t *testing.T) {
	opts := DefaultAnnotateOptions()
	opts.Enabled = true
	opts.TagValue = "test"
	input := map[string]map[string]string{
		"secret/x": {"a": "1"},
	}
	_ = AnnotateSecrets(input, opts)
	if _, ok := input["secret/x"]["_vaultdiff_source"]; ok {
		t.Error("original input was mutated")
	}
}

func TestDefaultAnnotateOptions_Disabled(t *testing.T) {
	opts := DefaultAnnotateOptions()
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.TagKey != "_vaultdiff_source" {
		t.Errorf("unexpected default TagKey: %q", opts.TagKey)
	}
}
