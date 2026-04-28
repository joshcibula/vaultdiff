package vault

import (
	"testing"
)

func TestRedactSecrets_NoOptions(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"api_key": "abc123", "db_pass": "s3cr3t"},
	}
	opts := DefaultRedactOptions()
	result := RedactSecrets(secrets, opts)
	if result["secret/app"]["api_key"] != "abc123" {
		t.Errorf("expected value to be unchanged, got %s", result["secret/app"]["api_key"])
	}
}

func TestRedactSecrets_ByPath(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/sensitive/token": {"value": "topsecret"},
		"secret/public/config":   {"value": "visible"},
	}
	opts := RedactOptions{Paths: []string{"secret/sensitive"}}
	result := RedactSecrets(secrets, opts)

	if result["secret/sensitive/token"]["value"] != "[REDACTED]" {
		t.Errorf("expected redacted, got %s", result["secret/sensitive/token"]["value"])
	}
	if result["secret/public/config"]["value"] != "visible" {
		t.Errorf("expected visible, got %s", result["secret/public/config"]["value"])
	}
}

func TestRedactSecrets_ByKey(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"api_key": "abc123", "region": "us-east-1"},
	}
	opts := RedactOptions{Keys: []string{"api_key"}}
	result := RedactSecrets(secrets, opts)

	if result["secret/app"]["api_key"] != "[REDACTED]" {
		t.Errorf("expected api_key to be redacted, got %s", result["secret/app"]["api_key"])
	}
	if result["secret/app"]["region"] != "us-east-1" {
		t.Errorf("expected region to be unchanged, got %s", result["secret/app"]["region"])
	}
}

func TestRedactSecrets_PathTakesPrecedenceOverKey(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/vault": {"token": "xyz", "user": "admin"},
	}
	opts := RedactOptions{
		Paths: []string{"secret/vault"},
		Keys:  []string{"token"},
	}
	result := RedactSecrets(secrets, opts)
	for _, v := range result["secret/vault"] {
		if v != "[REDACTED]" {
			t.Errorf("expected all values redacted by path, got %s", v)
		}
	}
}

func TestRedactSecrets_EmptySecrets(t *testing.T) {
	secrets := map[string]map[string]string{}
	opts := RedactOptions{Paths: []string{"secret/sensitive"}, Keys: []string{"token"}}
	result := RedactSecrets(secrets, opts)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestMatchesAnyPrefix(t *testing.T) {
	if !matchesAnyPrefix("secret/sensitive/key", []string{"secret/sensitive"}) {
		t.Error("expected prefix match")
	}
	if matchesAnyPrefix("secret/public/key", []string{"secret/sensitive"}) {
		t.Error("expected no prefix match")
	}
	if matchesAnyPrefix("secret/app", nil) {
		t.Error("expected no match with empty prefixes")
	}
}
