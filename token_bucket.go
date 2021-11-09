package main

import "time"

// TokenBucket
type limiterTokenBucket struct {
	N      int
	burst  int
	tokens int
	lastT  time.Time
}

// Fill `N` tokens every 1 second.
func newLimiterTokenBucket(N int, burst int) *limiterTokenBucket {
	return &limiterTokenBucket{
		N:     N,
		burst: burst,

		// Leave tokens=0, lastT=0.
		// The initial `tokens` will be filled when the first `Allow()`.
	}
}

func (l *limiterTokenBucket) Name() string {
	return "TokenBucket"
}

func (l *limiterTokenBucket) Allow(n int) bool {
	now := time.Now()

	// Fill tokens
	elapse := now.Sub(l.lastT)
	addN := int(float64(l.N) * elapse.Seconds())
	if addN > 0 {
		l.tokens += addN
		if l.tokens >= l.burst {
			l.tokens = l.burst
		}
		l.lastT = now
	}

	if l.tokens < n {
		return false
	}
	l.tokens -= n
	return true
}
