package vault

import (
	"testing"
)

func baseSecretsForSample() map[string]map[string]string {
	return map[string]map[string]string{
		"secret/a": {"key": "val-a"},
		"secret/b": {"key": "val-b"},
		"secret/c": {"key": "val-c"},
		"secret/d": {"key": "val-d"},
		"secret/e": {"key": "val-e"},
	}
}

func TestSampleSecrets_Disabled(t *testing.T) {
	secrets := baseSecretsForSample()
	opts := DefaultSampleOptions()
	out := SampleSecrets(secrets, opts)
	if len(out) != len(secrets) {
		t.Errorf("expected %d paths, got %d", len(secrets), len(out))
	}
}

func TestSampleSecrets_MaxPathsZero(t *testing.T) {
	secrets := baseSecretsForSample()
	opts := SampleOptions{Enabled: true, MaxPaths: 0, Seed: 42}
	out := SampleSecrets(secrets, opts)
	if len(out) != len(secrets) {
		t.Errorf("expected all %d paths when MaxPaths=0, got %d", len(secrets), len(out))
	}
}

func TestSampleSecrets_ReducesPaths(t *testing.T) {
	secrets := baseSecretsForSample()
	opts := SampleOptions{Enabled: true, MaxPaths: 3, Seed: 42}
	out := SampleSecrets(secrets, opts)
	if len(out) != 3 {
		t.Errorf("expected 3 sampled paths, got %d", len(out))
	}
}

func TestSampleSecrets_DoesNotMutateInput(t *testing.T) {
	secrets := baseSecretsForSample()
	orig := make(map[string]map[string]string, len(secrets))
	for k, v := range secrets {
		orig[k] = v
	}
	opts := SampleOptions{Enabled: true, MaxPaths: 2, Seed: 7}
	SampleSecrets(secrets, opts)
	if len(secrets) != len(orig) {
		t.Error("input map was mutated")
	}
}

func TestSampleSecrets_DeterministicWithSameSeed(t *testing.T) {
	secrets := baseSecretsForSample()
	opts := SampleOptions{Enabled: true, MaxPaths: 2, Seed: 99}
	out1 := SampleSecrets(secrets, opts)
	out2 := SampleSecrets(secrets, opts)
	if len(out1) != len(out2) {
		t.Fatal("expected same length for same seed")
	}
	for k := range out1 {
		if _, ok := out2[k]; !ok {
			t.Errorf("key %q present in first run but not second", k)
		}
	}
}

func TestSampleSecrets_AllPathsRetainedWhenUnderMax(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/x": {"a": "1"},
	}
	opts := SampleOptions{Enabled: true, MaxPaths: 10, Seed: 1}
	out := SampleSecrets(secrets, opts)
	if len(out) != 1 {
		t.Errorf("expected 1 path, got %d", len(out))
	}
}
