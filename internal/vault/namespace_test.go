package vault

import (
	"testing"
)

func TestParseVaultPath_SimpleMount(t *testing.T) {
	p, err := ParseVaultPath("secret/myapp", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", p.Mount)
	}
	if p.SecretPath != "myapp" {
		t.Errorf("expected secret path 'myapp', got %q", p.SecretPath)
	}
	if p.Namespace != "" {
		t.Errorf("expected empty namespace, got %q", p.Namespace)
	}
}

func TestParseVaultPath_WithDefaultNamespace(t *testing.T) {
	p, err := ParseVaultPath("secret/myapp", "team-ns")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Namespace != "team-ns" {
		t.Errorf("expected namespace 'team-ns', got %q", p.Namespace)
	}
}

func TestParseVaultPath_NamespaceInPath(t *testing.T) {
	p, err := ParseVaultPath("ns1/secret/myapp", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Namespace != "ns1" {
		t.Errorf("expected namespace 'ns1', got %q", p.Namespace)
	}
	if p.Mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", p.Mount)
	}
	if p.SecretPath != "myapp" {
		t.Errorf("expected secret path 'myapp', got %q", p.SecretPath)
	}
}

func TestParseVaultPath_Empty(t *testing.T) {
	_, err := ParseVaultPath("", "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestParseVaultPath_TooShort(t *testing.T) {
	_, err := ParseVaultPath("onlymount", "")
	if err == nil {
		t.Fatal("expected error for path with only one segment")
	}
}

func TestFullKVv2Path_AddsDataPrefix(t *testing.T) {
	p := &ParsedPath{Mount: "secret", SecretPath: "myapp"}
	got := p.FullKVv2Path()
	expected := "secret/data/myapp"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFullKVv2Path_NoDoubleData(t *testing.T) {
	p := &ParsedPath{Mount: "secret", SecretPath: "data/myapp"}
	got := p.FullKVv2Path()
	expected := "secret/data/myapp"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
