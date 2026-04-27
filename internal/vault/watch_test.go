package vault

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestDefaultWatchOptions(t *testing.T) {
	opts := DefaultWatchOptions()
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.Interval != 30*time.Second {
		t.Errorf("unexpected default interval: %v", opts.Interval)
	}
}

func TestNewWatcher_FieldsSet(t *testing.T) {
	opts := DefaultWatchOptions()
	w := NewWatcher(opts, nil)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
	if w.prev == nil {
		t.Error("expected prev map to be initialised")
	}
}

func TestWatcher_CallsOnChange(t *testing.T) {
	var mu sync.Mutex
	changes := map[string]int{}

	call := 0
	fetch := func(_ context.Context, path string) (map[string]string, error) {
		call++
		if call == 1 {
			return map[string]string{"key": "v1"}, nil
		}
		return map[string]string{"key": "v2"}, nil
	}

	opts := WatchOptions{
		Enabled:  true,
		Interval: 10 * time.Millisecond,
		OnChange: func(path string, _, _ map[string]map[string]string) {
			mu.Lock()
			changes[path]++
			mu.Unlock()
		},
	}
	w := NewWatcher(opts, fetch)
	// seed initial state
	w.prev["secret/app"] = map[string]string{"key": "v1"}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	w.Watch(ctx, []string{"secret/app"})

	mu.Lock()
	defer mu.Unlock()
	if changes["secret/app"] == 0 {
		t.Error("expected at least one change notification")
	}
}

func TestWatcher_NoChangeNoCallback(t *testing.T) {
	called := false
	fetch := func(_ context.Context, _ string) (map[string]string, error) {
		return map[string]string{"key": "same"}, nil
	}
	opts := WatchOptions{
		Enabled:  true,
		Interval: 10 * time.Millisecond,
		OnChange: func(_ string, _, _ map[string]map[string]string) { called = true },
	}
	w := NewWatcher(opts, fetch)
	w.prev["secret/app"] = map[string]string{"key": "same"}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	w.Watch(ctx, []string{"secret/app"})

	if called {
		t.Error("onChange should not be called when secrets are unchanged")
	}
}

func TestMapsEqual(t *testing.T) {
	if !mapsEqual(map[string]string{"a": "1"}, map[string]string{"a": "1"}) {
		t.Error("expected equal")
	}
	if mapsEqual(map[string]string{"a": "1"}, map[string]string{"a": "2"}) {
		t.Error("expected not equal")
	}
	if mapsEqual(map[string]string{"a": "1"}, map[string]string{}) {
		t.Error("expected not equal on length mismatch")
	}
}
