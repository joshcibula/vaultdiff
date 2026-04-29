package vault

import (
	"math/rand"
	"sort"
)

// DefaultSampleOptions returns a SampleOptions with sampling disabled.
func DefaultSampleOptions() SampleOptions {
	return SampleOptions{
		Enabled:  false,
		MaxPaths: 0,
		Seed:     0,
	}
}

// SampleOptions controls random sampling of secret paths.
type SampleOptions struct {
	// Enabled activates sampling.
	Enabled bool
	// MaxPaths is the maximum number of paths to retain. 0 means unlimited.
	MaxPaths int
	// Seed is the random seed. 0 means use a non-deterministic seed.
	Seed int64
}

// SampleSecrets randomly samples up to opts.MaxPaths paths from secrets.
// The input map is not mutated. If sampling is disabled or MaxPaths is 0,
// the original map is returned unchanged.
func SampleSecrets(secrets map[string]map[string]string, opts SampleOptions) map[string]map[string]string {
	if !opts.Enabled || opts.MaxPaths <= 0 || len(secrets) <= opts.MaxPaths {
		return secrets
	}

	// Collect and sort paths for deterministic ordering before sampling.
	paths := make([]string, 0, len(secrets))
	for p := range secrets {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	r := rand.New(rand.NewSource(opts.Seed)) //nolint:gosec
	r.Shuffle(len(paths), func(i, j int) {
		paths[i], paths[j] = paths[j], paths[i]
	})

	sampled := make(map[string]map[string]string, opts.MaxPaths)
	for _, p := range paths[:opts.MaxPaths] {
		kv := make(map[string]string, len(secrets[p]))
		for k, v := range secrets[p] {
			kv[k] = v
		}
		sampled[p] = kv
	}
	return sampled
}
