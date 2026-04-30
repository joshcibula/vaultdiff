package vault

import (
	"testing"
)

func TestClassifySecrets_Disabled(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"key": "val"},
	}
	opts := DefaultClassifyOptions()
	out := ClassifySecrets(secrets, opts)
	if _, ok := out["secret/app"]["_class"]; ok {
		t.Fatal("expected no _class key when disabled")
	}
}

func TestClassifySecrets_DefaultTag(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"key": "val"},
	}
	opts := DefaultClassifyOptions()
	opts.Enabled = true
	out := ClassifySecrets(secrets, opts)
	if got := out["secret/app"]["_class"]; got != "unclassified" {
		t.Fatalf("expected 'unclassified', got %q", got)
	}
}

func TestClassifySecrets_PathPrefixRule(t *testing.T) {
	secrets := map[string]map[string]string{
		"prod/db": {"password": "s3cr3t"},
		"dev/db":  {"password": "devpass"},
	}
	opts := ClassifyOptions{
		Enabled:    true,
		DefaultTag: "unclassified",
		Rules: []ClassifyRule{
			{PathPrefix: "prod/", Tag: "production"},
		},
	}
	out := ClassifySecrets(secrets, opts)
	if got := out["prod/db"]["_class"]; got != "production" {
		t.Fatalf("expected 'production', got %q", got)
	}
	if got := out["dev/db"]["_class"]; got != "unclassified" {
		t.Fatalf("expected 'unclassified', got %q", got)
	}
}

func TestClassifySecrets_KeyPrefixRule(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/svc": {"aws_access_key": "AKIA..."},
	}
	opts := ClassifyOptions{
		Enabled:    true,
		DefaultTag: "general",
		Rules: []ClassifyRule{
			{KeyPrefix: "aws_", Tag: "cloud-credentials"},
		},
	}
	out := ClassifySecrets(secrets, opts)
	if got := out["secret/svc"]["_class"]; got != "cloud-credentials" {
		t.Fatalf("expected 'cloud-credentials', got %q", got)
	}
}

func TestClassifySecrets_DoesNotMutateInput(t *testing.T) {
	original := map[string]map[string]string{
		"secret/app": {"key": "val"},
	}
	opts := ClassifyOptions{Enabled: true, DefaultTag: "x"}
	ClassifySecrets(original, opts)
	if _, ok := original["secret/app"]["_class"]; ok {
		t.Fatal("original map was mutated")
	}
}

func TestClassifySecrets_EmptySecrets(t *testing.T) {
	opts := ClassifyOptions{Enabled: true, DefaultTag: "x"}
	out := ClassifySecrets(map[string]map[string]string{}, opts)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d entries", len(out))
	}
}
