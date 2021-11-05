package main

import "time"

// TokenBucket
type limiterTokenBucket struct {
	N         int
	burst     int
	tokens    int
	everyNSec int
	lastT     time.Time
}

// Fill `N` tokens every `everyNSec` seconds.
func newLimiterTokenBucket(N int, everyNSec, burst int) *limiterTokenBucket {
	return &limiterTokenBucket{
		N:         N,
		burst:     burst,
		tokens:    burst,
		everyNSec: everyNSec,
		lastT:     time.Now(),
	}
}

func (l *limiterTokenBucket) Name() string {
	return "TokenBucket"
}

func (l *limiterTokenBucket) Allow(n int) bool {
	now := time.Now()

	// Fill tokens
	elapse := now.Sub(l.lastT).Milliseconds()
	addN := int(float64(l.N) * (float64(elapse) / float64(l.everyNSec*1000)))
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
