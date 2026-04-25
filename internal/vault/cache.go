package vault

import (
	"sync"
	"time"
)

// CacheEntry holds a cached secret map with an expiry timestamp.
type CacheEntry struct {
	Secrets   map[string]string
	FetchedAt time.Time
	TTL       time.Duration
}

// IsExpired returns true if the cache entry has exceeded its TTL.
func (e *CacheEntry) IsExpired() bool {
	if e.TTL == 0 {
		return false
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// SecretCache is a thread-safe in-memory cache for secret maps keyed by Vault path.
type SecretCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewSecretCache creates a new SecretCache with the given TTL.
// A TTL of zero means entries never expire.
func NewSecretCache(ttl time.Duration) *SecretCache {
	return &SecretCache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

// Get retrieves secrets for the given path. Returns nil, false if not found or expired.
func (c *SecretCache) Get(path string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[path]
	if !ok || entry.IsExpired() {
		return nil, false
	}
	return entry.Secrets, true
}

// Set stores secrets for the given path.
func (c *SecretCache) Set(path string, secrets map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[path] = &CacheEntry{
		Secrets:   secrets,
		FetchedAt: time.Now(),
		TTL:       c.ttl,
	}
}

// Invalidate removes a single path from the cache.
func (c *SecretCache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// Purge removes all expired entries from the cache and returns the number of entries removed.
func (c *SecretCache) Purge() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for path, entry := range c.entries {
		if entry.IsExpired() {
			delete(c.entries, path)
			count++
		}
	}
	return count
}

// Len returns the number of entries currently in the cache.
func (c *SecretCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
