package main

import "time"

// SlidingWindow
// Improval for `FixedWindow`
// Taking account for `lastCount` into current count
type limiterSlidingWindow struct {
	lastW       int
	capacity    int
	winSizeNSec int
	count       int

	lastCount int
}

func (l *limiterSlidingWindow) Name() string {
	return "SlidingWindow"
}

// Allow `capacity` in each `winSizeNSec` seconds time window
func newLimiterSlidingWindow(capacity int, winSizeNSec int) *limiterSlidingWindow {
	return &limiterSlidingWindow{
		lastW:       int(time.Now().Unix() / int64(winSizeNSec)),
		capacity:    capacity,
		winSizeNSec: winSizeNSec,
		count:       0,

		// For smoothly startup, give an initial `lastCount `
		lastCount: capacity,
	}
}

func (l *limiterSlidingWindow) Allow(n int) bool {
	now := time.Now()

	curW := int(now.Unix() / int64(l.winSizeNSec))
	if curW != l.lastW {
		l.lastCount = l.count
		l.count = 0
		l.lastW = curW
	}

	winSizeMSec := int64(l.winSizeNSec * 1000)
	percent := float64((now.UnixNano()/1000000)%winSizeMSec) / float64(winSizeMSec)

	nc := int(float64(l.lastCount)*(1-percent)) + l.count + n
	if nc > l.capacity {
		return false
	}

	l.count += n
	return true
}
