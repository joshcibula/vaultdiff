package vault

import (
	"context"
	"errors"
	"sort"
	"sync/atomic"
	"testing"
)

func TestFetchAllConcurrent_AllSucceed(t *testing.T) {
	paths := []string{"secret/a", "secret/b", "secret/c"}
	fetch := func(ctx context.Context, path string) (map[string]string, error) {
		return map[string]string{"key": path}, nil
	}

	results := FetchAllConcurrent(context.Background(), paths, DefaultConcurrencyOptions(), fetch)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %s: %v", r.Path, r.Err)
		}
	}
}

func TestFetchAllConcurrent_PartialError(t *testing.T) {
	paths := []string{"secret/ok", "secret/fail"}
	errFetch := errors.New("fetch failed")
	fetch := func(ctx context.Context, path string) (map[string]string, error) {
		if path == "secret/fail" {
			return nil, errFetch
		}
		return map[string]string{"k": "v"}, nil
	}

	results := FetchAllConcurrent(context.Background(), paths, DefaultConcurrencyOptions(), fetch)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Path < results[j].Path })
	if results[1].Err == nil {
		t.Error("expected error for secret/fail")
	}
}

func TestFetchAllConcurrent_WorkerCount(t *testing.T) {
	var concurrent int64
	var peak int64

	paths := make([]string, 20)
	for i := range paths {
		paths[i] = "secret/path"
	}

	fetch := func(ctx context.Context, path string) (map[string]string, error) {
		cur := atomic.AddInt64(&concurrent, 1)
		for {
			p := atomic.LoadInt64(&peak)
			if cur <= p || atomic.CompareAndSwapInt64(&peak, p, cur) {
				break
			}
		}
		atomic.AddInt64(&concurrent, -1)
		return map[string]string{}, nil
	}

	FetchAllConcurrent(context.Background(), paths, ConcurrencyOptions{Workers: 3}, fetch)
	if peak > 3 {
		t.Errorf("peak concurrency %d exceeded worker limit 3", peak)
	}
}

func TestFetchAllConcurrent_EmptyPaths(t *testing.T) {
	results := FetchAllConcurrent(context.Background(), []string{}, DefaultConcurrencyOptions(), nil)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
