package tokenbucket

import (
	"errors"
	"math"
	"sync"
	"time"
)

var ErrInvalidConfig = errors.New("invalid token bucket configuration")

// Limiter defines the behavior expected from a rate limiter.
type Limiter interface {
	Allow() bool
	AllowN(n int) bool
}

// Clock abstracts time to make tests deterministic.
type Clock interface {
	Now() time.Time
}

type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now()
}

// Config defines the token bucket parameters.
type Config struct {
	// Capacity is the maximum number of tokens.
	Capacity int

	// RefillRate is the number of tokens added per second.
	RefillRate float64
}

// TokenBucket is a concurrency-safe token bucket rate limiter.
type TokenBucket struct {
	mu sync.Mutex

	capacity   float64
	refillRate float64
	tokens     float64
	lastRefill time.Time
	clock      Clock
}

// New creates a token bucket that starts full.
// Pass nil for clock to use the system clock.
func New(cfg Config, clock Clock) (*TokenBucket, error) {
	if cfg.Capacity <= 0 ||
		cfg.RefillRate <= 0 ||
		math.IsNaN(cfg.RefillRate) ||
		math.IsInf(cfg.RefillRate, 0) {
		return nil, ErrInvalidConfig
	}

	if clock == nil {
		clock = systemClock{}
	}

	now := clock.Now()

	return &TokenBucket{
		capacity:   float64(cfg.Capacity),
		refillRate: cfg.RefillRate,
		tokens:     float64(cfg.Capacity),
		lastRefill: now,
		clock:      clock,
	}, nil
}

// Allow attempts to consume one token.
func (b *TokenBucket) Allow() bool {
	return b.AllowN(1)
}

// AllowN atomically consumes n tokens if enough are available.
// Invalid requests and requests larger than the capacity are rejected.
func (b *TokenBucket) AllowN(n int) bool {
	if n <= 0 {
		return false
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.refill()

	if float64(n) > b.capacity {
		return false
	}

	if b.tokens < float64(n) {
		return false
	}

	b.tokens -= float64(n)
	return true
}

// refill updates the bucket based on elapsed time.
// The caller must hold b.mu.
func (b *TokenBucket) refill() {
	now := b.clock.Now()
	elapsed := now.Sub(b.lastRefill)

	// A clock that moves backwards must not create tokens.
	if elapsed <= 0 {
		return
	}

	b.tokens += elapsed.Seconds() * b.refillRate

	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}

	b.lastRefill = now
}
