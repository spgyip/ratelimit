package main

import "time"

// FixedWindow
type limiterFixedWindow struct {
	lastW       int
	capacity    int
	winSizeNSec int
	count       int
}

func (l *limiterFixedWindow) Name() string {
	return "FixedWindow"
}

// Allow `capacity` in each `winSizeNSec` seconds time window
func newLimiterFixedWindow(capacity int, winSizeNSec int) *limiterFixedWindow {
	return &limiterFixedWindow{
		lastW:       int(time.Now().Unix() / int64(winSizeNSec)),
		capacity:    capacity,
		winSizeNSec: winSizeNSec,
		count:       0,
	}
}

func (l *limiterFixedWindow) Allow(n int) bool {
	curW := int(time.Now().Unix() / int64(l.winSizeNSec))
	if curW != l.lastW {
		l.count = 0
		l.lastW = curW
	}

	nc := l.count + n
	if nc > l.capacity {
		return false
	}

	l.count += n
	return true
}
