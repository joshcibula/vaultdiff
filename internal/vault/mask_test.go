package vault

import (
	"testing"
)

func TestMaskSecrets_Disabled(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"password": "s3cr3t", "user": "admin"},
	}
	opts := DefaultMaskOptions()
	result := MaskSecrets(secrets, opts)
	if result["secret/app"]["password"] != "s3cr3t" {
		t.Errorf("expected unmasked value, got %s", result["secret/app"]["password"])
	}
}

func TestMaskSecrets_Enabled(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"password": "s3cr3t", "user": "admin"},
	}
	opts := MaskOptions{Enabled: true, MaskString: "***", RevealKeys: []string{}}
	result := MaskSecrets(secrets, opts)
	for _, v := range result["secret/app"] {
		if v != "***" {
			t.Errorf("expected masked value, got %s", v)
		}
	}
}

func TestMaskSecrets_RevealKeys(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"password": "s3cr3t", "user": "admin"},
	}
	opts := MaskOptions{Enabled: true, MaskString: "***", RevealKeys: []string{"user"}}
	result := MaskSecrets(secrets, opts)
	if result["secret/app"]["user"] != "admin" {
		t.Errorf("expected revealed value for 'user', got %s", result["secret/app"]["user"])
	}
	if result["secret/app"]["password"] != "***" {
		t.Errorf("expected masked password, got %s", result["secret/app"]["password"])
	}
}

func TestMaskValue_Disabled(t *testing.T) {
	opts := DefaultMaskOptions()
	if got := MaskValue("key", "val", opts); got != "val" {
		t.Errorf("expected val, got %s", got)
	}
}

func TestMaskValue_Enabled(t *testing.T) {
	opts := MaskOptions{Enabled: true, MaskString: "[hidden]", RevealKeys: []string{}}
	if got := MaskValue("key", "val", opts); got != "[hidden]" {
		t.Errorf("expected [hidden], got %s", got)
	}
}

func TestMaskValue_RevealCaseInsensitive(t *testing.T) {
	opts := MaskOptions{Enabled: true, MaskString: "***", RevealKeys: []string{"USER"}}
	if got := MaskValue("user", "admin", opts); got != "admin" {
		t.Errorf("expected admin, got %s", got)
	}
}
