package vault

import (
	"testing"
)

func TestTransformSecrets_NoOp(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {"Key": "  value  "},
	}
	opts := DefaultTransformOptions()
	out := TransformSecrets(input, opts)
	if out["secret/app"]["Key"] != "  value  " {
		t.Errorf("expected value unchanged, got %q", out["secret/app"]["Key"])
	}
}

func TestTransformSecrets_TrimSpace(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {"db_pass": "  s3cr3t  "},
	}
	opts := TransformOptions{TrimSpace: true}
	out := TransformSecrets(input, opts)
	if out["secret/app"]["db_pass"] != "s3cr3t" {
		t.Errorf("expected trimmed value, got %q", out["secret/app"]["db_pass"])
	}
}

func TestTransformSecrets_LowercaseKeys(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {"DB_PASS": "hunter2", "ApiKey": "abc"},
	}
	opts := TransformOptions{LowercaseKeys: true}
	out := TransformSecrets(input, opts)
	if _, ok := out["secret/app"]["db_pass"]; !ok {
		t.Error("expected key 'db_pass' after lowercasing")
	}
	if _, ok := out["secret/app"]["apikey"]; !ok {
		t.Error("expected key 'apikey' after lowercasing")
	}
	if _, ok := out["secret/app"]["DB_PASS"]; ok {
		t.Error("expected original key 'DB_PASS' to be absent")
	}
}

func TestTransformSecrets_IgnoreKeys(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {"password": "s3cr3t", "username": "admin"},
	}
	opts := TransformOptions{IgnoreKeys: []string{"password"}}
	out := TransformSecrets(input, opts)
	if out["secret/app"]["password"] != "" {
		t.Errorf("expected ignored key to be empty, got %q", out["secret/app"]["password"])
	}
	if out["secret/app"]["username"] != "admin" {
		t.Errorf("expected username unchanged, got %q", out["secret/app"]["username"])
	}
}

func TestTransformSecrets_DoesNotMutateInput(t *testing.T) {
	input := map[string]map[string]string{
		"secret/app": {"key": "  val  "},
	}
	opts := TransformOptions{TrimSpace: true}
	TransformSecrets(input, opts)
	if input["secret/app"]["key"] != "  val  " {
		t.Error("original input was mutated")
	}
}

func TestTransformSecrets_CombinedOptions(t *testing.T) {
	input := map[string]map[string]string{
		"secret/svc": {"SECRET_KEY": "  topsecret  ", "HOST": "  localhost  "},
	}
	opts := TransformOptions{TrimSpace: true, LowercaseKeys: true, IgnoreKeys: []string{"secret_key"}}
	out := TransformSecrets(input, opts)
	if out["secret/svc"]["secret_key"] != "" {
		t.Errorf("expected ignored key to be empty, got %q", out["secret/svc"]["secret_key"])
	}
	if out["secret/svc"]["host"] != "localhost" {
		t.Errorf("expected trimmed lowercase host, got %q", out["secret/svc"]["host"])
	}
}
