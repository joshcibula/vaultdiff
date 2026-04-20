package vault

import (
	"sync"
	"time"
)

// RateLimiter implements a simple token-bucket rate limiter for Vault API calls.
type RateLimiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	refillPS float64 // tokens per second
	lastTime time.Time
}

// RateLimitOptions configures the rate limiter.
type RateLimitOptions struct {
	// RequestsPerSecond is the maximum sustained request rate.
	RequestsPerSecond float64
	// Burst is the maximum number of requests allowed in a burst.
	Burst float64
}

// DefaultRateLimitOptions returns sensible defaults.
func DefaultRateLimitOptions() RateLimitOptions {
	return RateLimitOptions{
		RequestsPerSecond: 10,
		Burst:             20,
	}
}

// NewRateLimiter creates a RateLimiter from the given options.
func NewRateLimiter(opts RateLimitOptions) *RateLimiter {
	return &RateLimiter{
		tokens:   opts.Burst,
		max:      opts.Burst,
		refillPS: opts.RequestsPerSecond,
		lastTime: time.Now(),
	}
}

// Wait blocks until a token is available, then consumes one.
func (r *RateLimiter) Wait() {
	for {
		r.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(r.lastTime).Seconds()
		r.lastTime = now
		r.tokens += elapsed * r.refillPS
		if r.tokens > r.max {
			r.tokens = r.max
		}
		if r.tokens >= 1 {
			r.tokens--
			r.mu.Unlock()
			return
		}
		// Calculate how long until the next token is available.
		waitDuration := time.Duration((1-r.tokens)/r.refillPS*1000) * time.Millisecond
		r.mu.Unlock()
		time.Sleep(waitDuration)
	}
}

// TryAcquire attempts to consume a token without blocking.
// Returns true if a token was acquired, false otherwise.
func (r *RateLimiter) TryAcquire() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(r.lastTime).Seconds()
	r.lastTime = now
	r.tokens += elapsed * r.refillPS
	if r.tokens > r.max {
		r.tokens = r.max
	}
	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}
