package vault

import (
	"context"
	"sync"
)

// ConcurrencyOptions controls parallel secret fetching.
type ConcurrencyOptions struct {
	Workers int
}

// DefaultConcurrencyOptions returns sensible defaults.
func DefaultConcurrencyOptions() ConcurrencyOptions {
	return ConcurrencyOptions{
		Workers: 5,
	}
}

// FetchResult holds the result of fetching a single path.
type FetchResult struct {
	Path    string
	Secrets map[string]string
	Err     error
}

// FetchAllConcurrent fetches secrets for all paths in parallel using a worker pool.
func FetchAllConcurrent(ctx context.Context, paths []string, opts ConcurrencyOptions, fetch func(ctx context.Context, path string) (map[string]string, error)) []FetchResult {
	if opts.Workers <= 0 {
		opts.Workers = 1
	}

	pathCh := make(chan string, len(paths))
	for _, p := range paths {
		pathCh <- p
	}
	close(pathCh)

	results := make([]FetchResult, 0, len(paths))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < opts.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathCh {
				secrets, err := fetch(ctx, path)
				mu.Lock()
				results = append(results, FetchResult{Path: path, Secrets: secrets, Err: err})
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return results
}
