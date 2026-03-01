// Package ratelimit provides in-memory per-key rate limiting for the HTTP server.
// When SECURITY_RATE_LIMIT_RPM > 0, requests are limited per client key (IP or identity).
package ratelimit

import (
	"sync"
	"time"
)

// Limiter limits the number of requests per key within a time window.
// Keys are typically client IP or authenticated identity.
type Limiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	rpm    int
	window time.Duration
}

// NewLimiter returns a limiter that allows rpm requests per key per minute.
// Burst is ignored in this implementation; the window is one minute.
func NewLimiter(rpm int, burst int) *Limiter {
	if rpm <= 0 {
		return nil
	}
	if burst <= 0 {
		burst = rpm * 2
	}
	_ = burst // reserve for future use
	return &Limiter{
		hits:   make(map[string][]time.Time),
		rpm:    rpm,
		window: time.Minute,
	}
}

// Allow reports whether the request for key is allowed.
// It records the request and prunes old entries for key. Thread-safe.
func (l *Limiter) Allow(key string) bool {
	if l == nil || l.rpm <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-l.window)
	times := l.hits[key]
	// Prune entries outside the window
	i := 0
	for _, t := range times {
		if t.After(cutoff) {
			times[i] = t
			i++
		}
	}
	times = times[:i]
	if len(times) >= l.rpm {
		l.hits[key] = times
		return false
	}
	l.hits[key] = append(times, now)
	return true
}
