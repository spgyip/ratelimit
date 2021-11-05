package main

import "time"

// FixedWindow
type limiterFixedWindow struct {
	lastSec int
	r       int
	count   int
}

func (l *limiterFixedWindow) Name() string {
	return "FixedWindow"
}

// Allow `r` request every 1 second window size.
func newLimiterFixedWindow(r int) *limiterFixedWindow {
	return &limiterFixedWindow{
		lastSec: int(time.Now().Unix()),
		r:       r,
		count:   0,
	}
}

func (l *limiterFixedWindow) Allow(n int) bool {
	curSec := int(time.Now().Unix())
	if curSec != l.lastSec {
		l.count = 0
		l.lastSec = curSec
	}

	nc := l.count + n
	if nc > l.r {
		return false
	}

	l.count += n
	return true
}
