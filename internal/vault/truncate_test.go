package vault

import (
	"strings"
	"testing"
)

func TestTruncateSecrets_Disabled(t *testing.T) {
	opts := DefaultTruncateOptions() // Enabled == false
	input := map[string]map[string]string{
		"secret/app": {"key": strings.Repeat("x", 200)},
	}
	out := TruncateSecrets(input, opts)
	if out["secret/app"]["key"] != strings.Repeat("x", 200) {
		t.Error("expected value to be unchanged when disabled")
	}
}

func TestTruncateSecrets_ShortValueUnchanged(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.Enabled = true
	input := map[string]map[string]string{
		"secret/app": {"key": "short"},
	}
	out := TruncateSecrets(input, opts)
	if out["secret/app"]["key"] != "short" {
		t.Errorf("expected 'short', got %q", out["secret/app"]["key"])
	}
}

func TestTruncateSecrets_LongValueTruncated(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.Enabled = true
	opts.MaxLength = 10
	input := map[string]map[string]string{
		"secret/app": {"token": "abcdefghijklmnopqrstuvwxyz"},
	}
	out := TruncateSecrets(input, opts)
	val := out["secret/app"]["token"]
	if !strings.HasSuffix(val, "...") {
		t.Errorf("expected ellipsis suffix, got %q", val)
	}
	if len([]rune(val)) > opts.MaxLength+len([]rune(opts.Ellipsis)) {
		t.Errorf("truncated value too long: %q", val)
	}
}

func TestTruncateSecrets_SkipKeys(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.Enabled = true
	opts.MaxLength = 5
	opts.SkipKeys = []string{"exempt"}
	long := strings.Repeat("z", 100)
	input := map[string]map[string]string{
		"secret/app": {
			"exempt":    long,
			"truncated": long,
		},
	}
	out := TruncateSecrets(input, opts)
	if out["secret/app"]["exempt"] != long {
		t.Error("exempt key should not be truncated")
	}
	if out["secret/app"]["truncated"] == long {
		t.Error("non-exempt key should be truncated")
	}
}

func TestTruncateSecrets_DoesNotMutateInput(t *testing.T) {
	opts := DefaultTruncateOptions()
	opts.Enabled = true
	opts.MaxLength = 3
	original := "hello world"
	input := map[string]map[string]string{
		"secret/app": {"k": original},
	}
	TruncateSecrets(input, opts)
	if input["secret/app"]["k"] != original {
		t.Error("input map was mutated")
	}
}

func TestDefaultTruncateOptions(t *testing.T) {
	opts := DefaultTruncateOptions()
	if opts.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if opts.MaxLength <= 0 {
		t.Errorf("expected positive MaxLength, got %d", opts.MaxLength)
	}
	if opts.Ellipsis == "" {
		t.Error("expected non-empty default Ellipsis")
	}
}
