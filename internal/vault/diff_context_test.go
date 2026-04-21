package vault

import (
	"testing"
)

func makeVaultPath(ns, mount, secret string) VaultPath {
	return VaultPath{
		Namespace:  ns,
		Mount:      mount,
		SecretPath: secret,
	}
}

func TestNewDiffContext_FieldsPopulated(t *testing.T) {
	left := makeVaultPath("ns1", "secret", "app/config")
	right := makeVaultPath("ns2", "secret", "app/config")
	ctx := NewDiffContext(left, right)

	if ctx.LeftNamespace != "ns1" {
		t.Errorf("expected LeftNamespace ns1, got %s", ctx.LeftNamespace)
	}
	if ctx.RightNamespace != "ns2" {
		t.Errorf("expected RightNamespace ns2, got %s", ctx.RightNamespace)
	}
	if ctx.LeftMount != "secret" || ctx.RightMount != "secret" {
		t.Errorf("unexpected mount values: %s / %s", ctx.LeftMount, ctx.RightMount)
	}
}

func TestDiffContext_SameNamespace(t *testing.T) {
	ctx := NewDiffContext(
		makeVaultPath("shared", "kv", "a"),
		makeVaultPath("shared", "kv", "b"),
	)
	if !ctx.SameNamespace() {
		t.Error("expected SameNamespace to return true")
	}

	ctx2 := NewDiffContext(
		makeVaultPath("ns1", "kv", "a"),
		makeVaultPath("ns2", "kv", "b"),
	)
	if ctx2.SameNamespace() {
		t.Error("expected SameNamespace to return false")
	}
}

func TestDiffContext_SameMount(t *testing.T) {
	ctx := NewDiffContext(
		makeVaultPath("", "kv", "a"),
		makeVaultPath("", "kv", "b"),
	)
	if !ctx.SameMount() {
		t.Error("expected SameMount to return true")
	}

	ctx2 := NewDiffContext(
		makeVaultPath("", "kv", "a"),
		makeVaultPath("", "secret", "b"),
	)
	if ctx2.SameMount() {
		t.Error("expected SameMount to return false")
	}
}

func TestDiffContext_Summary_WithNamespace(t *testing.T) {
	ctx := NewDiffContext(
		makeVaultPath("teamA", "kv", "app/prod"),
		makeVaultPath("teamB", "kv", "app/staging"),
	)
	got := ctx.Summary()
	want := "teamA/app/prod → teamB/app/staging"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}

func TestDiffContext_Summary_NoNamespace(t *testing.T) {
	ctx := NewDiffContext(
		makeVaultPath("", "kv", "app/prod"),
		makeVaultPath("", "kv", "app/staging"),
	)
	got := ctx.Summary()
	want := "app/prod → app/staging"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}
