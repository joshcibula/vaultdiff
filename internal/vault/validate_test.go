package vault

import (
	"strings"
	"testing"
)

func TestValidateSecrets_Disabled(t *testing.T) {
	opts := DefaultValidationOptions()
	secrets := map[string]map[string]string{
		"secret/app": {"password": ""},
	}
	if err := ValidateSecrets(secrets, opts); err != nil {
		t.Fatalf("expected no error when disabled, got: %v", err)
	}
}

func TestValidateSecrets_RequireNonEmpty(t *testing.T) {
	opts := DefaultValidationOptions()
	opts.Enabled = true
	opts.RequireNonEmpty = true

	secrets := map[string]map[string]string{
		"secret/app": {"token": "", "key": "value"},
	}
	err := ValidateSecrets(secrets, opts)
	if err == nil {
		t.Fatal("expected validation error for empty value")
	}
	if !strings.Contains(err.Error(), "must not be empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidateSecrets_ForbiddenKeys(t *testing.T) {
	opts := DefaultValidationOptions()
	opts.Enabled = true
	opts.ForbiddenKeys = []string{"root_password", "master_key"}

	secrets := map[string]map[string]string{
		"secret/db": {"ROOT_PASSWORD": "secret123"},
	}
	err := ValidateSecrets(secrets, opts)
	if err == nil {
		t.Fatal("expected validation error for forbidden key")
	}
	if !strings.Contains(err.Error(), "forbidden") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidateSecrets_MaxValueLength(t *testing.T) {
	opts := DefaultValidationOptions()
	opts.Enabled = true
	opts.MaxValueLength = 10

	secrets := map[string]map[string]string{
		"secret/app": {"cert": strings.Repeat("x", 20)},
	}
	err := ValidateSecrets(secrets, opts)
	if err == nil {
		t.Fatal("expected validation error for long value")
	}
	if !strings.Contains(err.Error(), "exceeds max") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidateSecrets_NoViolations(t *testing.T) {
	opts := DefaultValidationOptions()
	opts.Enabled = true
	opts.RequireNonEmpty = true
	opts.MaxValueLength = 50
	opts.ForbiddenKeys = []string{"banned"}

	secrets := map[string]map[string]string{
		"secret/app": {"api_key": "abc123"},
	}
	if err := ValidateSecrets(secrets, opts); err != nil {
		t.Fatalf("expected no violations, got: %v", err)
	}
}

func TestIsValidationError(t *testing.T) {
	ve := &ValidationError{Violations: []string{"oops"}}
	if !IsValidationError(ve) {
		t.Error("expected IsValidationError to return true for *ValidationError")
	}
	if IsValidationError(nil) {
		t.Error("expected IsValidationError to return false for nil")
	}
}
