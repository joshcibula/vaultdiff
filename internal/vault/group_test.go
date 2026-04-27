package vault

import (
	"testing"
)

func TestGroupSecrets_Disabled(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app/db": {"pass": "x"},
		"secret/app/api": {"key": "y"},
	}
	opts := DefaultGroupOptions()
	opts.Enabled = false

	groups := GroupSecrets(secrets, opts)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "" {
		t.Errorf("expected empty group name, got %q", groups[0].Name)
	}
	if len(groups[0].Secrets) != 2 {
		t.Errorf("expected 2 secrets in group, got %d", len(groups[0].Secrets))
	}
}

func TestGroupSecrets_ByMount(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app/db":  {"pass": "x"},
		"secret/app/api": {"key": "y"},
		"kv/infra/redis": {"url": "z"},
	}
	opts := DefaultGroupOptions()
	opts.Enabled = true
	opts.GroupBy = "mount"

	groups := GroupSecrets(secrets, opts)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "kv/infra" {
		t.Errorf("expected kv/infra, got %q", groups[0].Name)
	}
	if groups[1].Name != "secret/app" {
		t.Errorf("expected secret/app, got %q", groups[1].Name)
	}
}

func TestGroupSecrets_ByPrefix(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/a/x": {"k": "v"},
		"secret/b/y": {"k": "v"},
		"kv/a/z":     {"k": "v"},
	}
	opts := DefaultGroupOptions()
	opts.Enabled = true
	opts.GroupBy = "prefix"

	groups := GroupSecrets(secrets, opts)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestGroupSecrets_ByDepth(t *testing.T) {
	secrets := map[string]map[string]string{
		"a/b/c/d": {"k": "1"},
		"a/b/e/f": {"k": "2"},
		"a/x/y/z": {"k": "3"},
	}
	opts := DefaultGroupOptions()
	opts.Enabled = true
	opts.GroupBy = "depth"
	opts.Depth = 2

	groups := GroupSecrets(secrets, opts)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestGroupSecrets_Empty(t *testing.T) {
	opts := DefaultGroupOptions()
	opts.Enabled = true

	groups := GroupSecrets(map[string]map[string]string{}, opts)
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}

func TestDefaultGroupOptions(t *testing.T) {
	opts := DefaultGroupOptions()
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.GroupBy != "mount" {
		t.Errorf("expected GroupBy=mount, got %q", opts.GroupBy)
	}
	if opts.Depth != 1 {
		t.Errorf("expected Depth=1, got %d", opts.Depth)
	}
}
