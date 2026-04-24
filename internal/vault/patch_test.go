package vault

import (
	"testing"
)

func baseSecrets() map[string]map[string]string {
	return map[string]map[string]string{
		"secret/app": {"db_pass": "hunter2", "api_key": "abc123"},
	}
}

func TestPatchSecrets_Set(t *testing.T) {
	ops := PatchOptions{
		Operations: []PatchOperation{
			{Op: "set", Path: "secret/app", Key: "db_pass", Value: "newpass"},
		},
	}
	out, applied, err := PatchSecrets(baseSecrets(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["secret/app"]["db_pass"] != "newpass" {
		t.Errorf("expected newpass, got %s", out["secret/app"]["db_pass"])
	}
	if len(applied) != 1 {
		t.Errorf("expected 1 applied op, got %d", len(applied))
	}
}

func TestPatchSecrets_Delete(t *testing.T) {
	ops := PatchOptions{
		Operations: []PatchOperation{
			{Op: "delete", Path: "secret/app", Key: "api_key"},
		},
	}
	out, _, err := PatchSecrets(baseSecrets(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["secret/app"]["api_key"]; ok {
		t.Error("expected api_key to be deleted")
	}
}

func TestPatchSecrets_Rename(t *testing.T) {
	ops := PatchOptions{
		Operations: []PatchOperation{
			{Op: "rename", Path: "secret/app", Key: "db_pass", NewKey: "database_password"},
		},
	}
	out, _, err := PatchSecrets(baseSecrets(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["secret/app"]["db_pass"]; ok {
		t.Error("expected old key to be removed")
	}
	if out["secret/app"]["database_password"] != "hunter2" {
		t.Error("expected renamed key to carry original value")
	}
}

func TestPatchSecrets_DryRun(t *testing.T) {
	ops := PatchOptions{
		DryRun: true,
		Operations: []PatchOperation{
			{Op: "set", Path: "secret/app", Key: "db_pass", Value: "should-not-apply"},
		},
	}
	out, applied, err := PatchSecrets(baseSecrets(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["secret/app"]["db_pass"] == "should-not-apply" {
		t.Error("dry-run should not mutate secrets")
	}
	if len(applied) != 1 {
		t.Errorf("expected 1 applied op description, got %d", len(applied))
	}
}

func TestPatchSecrets_UnknownPath(t *testing.T) {
	ops := PatchOptions{
		Operations: []PatchOperation{
			{Op: "set", Path: "secret/missing", Key: "k", Value: "v"},
		},
	}
	_, _, err := PatchSecrets(baseSecrets(), ops)
	if err == nil {
		t.Error("expected error for unknown path")
	}
}

func TestPatchSecrets_UnknownOp(t *testing.T) {
	ops := PatchOptions{
		Operations: []PatchOperation{
			{Op: "upsert", Path: "secret/app", Key: "k"},
		},
	}
	_, _, err := PatchSecrets(baseSecrets(), ops)
	if err == nil {
		t.Error("expected error for unknown operation")
	}
}

func TestPatchSecrets_DoesNotMutateInput(t *testing.T) {
	input := baseSecrets()
	ops := PatchOptions{
		Operations: []PatchOperation{
			{Op: "set", Path: "secret/app", Key: "db_pass", Value: "changed"},
		},
	}
	_, _, _ = PatchSecrets(input, ops)
	if input["secret/app"]["db_pass"] != "hunter2" {
		t.Error("input map was mutated")
	}
}
