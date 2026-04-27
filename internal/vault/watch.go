package vault

import (
	"context"
	"time"
)

// WatchOptions configures the secret watcher behaviour.
type WatchOptions struct {
	Enabled  bool
	Interval time.Duration
	OnChange func(path string, prev, curr map[string]map[string]string)
}

// DefaultWatchOptions returns sensible defaults for watching.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Enabled:  false,
		Interval: 30 * time.Second,
	}
}

// Watcher polls two sets of secret paths at a fixed interval and calls
// OnChange whenever the resolved secrets differ from the previous snapshot.
type Watcher struct {
	opts    WatchOptions
	fetch   func(ctx context.Context, path string) (map[string]string, error)
	prev    map[string]map[string]string
}

// NewWatcher creates a Watcher with the supplied fetch function.
func NewWatcher(opts WatchOptions, fetch func(ctx context.Context, path string) (map[string]string, error)) *Watcher {
	return &Watcher{opts: opts, fetch: fetch, prev: make(map[string]map[string]string)}
}

// Watch starts polling the given paths until ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context, paths []string) {
	ticker := time.NewTicker(w.opts.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.poll(ctx, paths)
		}
	}
}

func (w *Watcher) poll(ctx context.Context, paths []string) {
	for _, p := range paths {
		curr, err := w.fetch(ctx, p)
		if err != nil {
			continue
		}
		prev := w.prev[p]
		if !mapsEqual(prev, curr) {
			if w.opts.OnChange != nil {
				w.opts.OnChange(p, prev, map[string]map[string]string{p: curr})
			}
			w.prev[p] = curr
		}
	}
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
