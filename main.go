package main

import (
	"flag"
	"fmt"
	"time"
)

type limiter interface {
	Allow(n int) bool
}

//
type limiterFixWindow struct {
	lastW      int
	capacity   int
	winsizeSec int
	count      int
}

func newLimiterFixWindow(capacity int, winsizeSec int) *limiterFixWindow {
	return &limiterFixWindow{
		lastW:      int(time.Now().Unix() / int64(winsizeSec)),
		capacity:   capacity,
		winsizeSec: winsizeSec,
	}
}

func (l *limiterFixWindow) Allow(n int) bool {
	curW := int(time.Now().Unix() / int64(l.winsizeSec))
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

func doSomething() {
}

func main() {
	var dur time.Duration
	flag.DurationVar(&dur, "dur", 10*time.Second, "Run duration")
	flag.Parse()

	limiter := newLimiterFixWindow(100000, 5)

	t0 := time.Now()
	t1, lastT := time.Now(), time.Now()
	cntAllow := 0
	cntReject := 0
	for {
		if limiter.Allow(1) {
			cntAllow++
			doSomething()
		} else {
			cntReject++
		}

		t1 = time.Now()

		// Print & flush `cnt` every 100ms
		if t1.Sub(lastT).Milliseconds() >= 100 {
			fmt.Printf("[%0.2f] %-8v %v\n", t1.Sub(t0).Seconds(), cntAllow, cntReject)
			lastT = t1
			cntAllow = 0
			cntReject = 0
		}

		// Run duration
		if t1.Sub(t0).Seconds() >= dur.Seconds() {
			break
		}
	}
	fmt.Println("duration: ", time.Now().Sub(t0))
}
