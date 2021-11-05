package main

import "time"

// SlidingWindow
// Improval for `FixedWindow`
// Taking account for `lastCount` into current count
type limiterSlidingWindow struct {
	lastSec int
	r       int
	count   int

	lastCount int
}

func (l *limiterSlidingWindow) Name() string {
	return "SlidingWindow"
}

// Allow `r` request every 1 second window size.
func newLimiterSlidingWindow(r int) *limiterSlidingWindow {
	return &limiterSlidingWindow{
		lastSec: int(time.Now().Unix()),
		r:       r,
		count:   0,

		// For smoothly startup, give an initial `lastCount `
		lastCount: r,
	}
}

func (l *limiterSlidingWindow) Allow(n int) bool {
	now := time.Now()

	curSec := int(now.Unix())
	if curSec != l.lastSec {
		l.lastCount = l.count
		l.count = 0
		l.lastSec = curSec
	}

	percent := float64((now.UnixNano()/1000000)%1000) / float64(1000)
	nc := int(float64(l.lastCount)*(1-percent)) + l.count + n
	if nc > l.r {
		return false
	}

	l.count += n
	return true
}
