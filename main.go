package main

import (
	"flag"
	"fmt"
	"time"
)

type limiter interface {
	Allow(n int) bool
}

// TokenBucket
type limiterTokenBucket struct {
	N        int
	burst    int
	tokens   int
	everySec int
	lastT    time.Time
}

// Fill `N` of tokens every `everySec` seconds.
func newLimiterTokenBucket(N int, everySec, burst int) *limiterTokenBucket {
	return &limiterTokenBucket{
		N:        N,
		burst:    burst,
		tokens:   burst,
		everySec: everySec,
		lastT:    time.Now(),
	}
}

func (l *limiterTokenBucket) Allow(n int) bool {
	now := time.Now()

	// Fill tokens
	elapse := now.Sub(l.lastT).Milliseconds()
	addN := int(float64(l.N) * (float64(elapse) / float64(l.everySec*1000)))
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

// FixedWindow
type limiterFixedWindow struct {
	lastW    int
	capacity int
	everySec int
	count    int
}

func newLimiterFixedWindow(capacity int, everySec int) *limiterFixedWindow {
	return &limiterFixedWindow{
		lastW:    int(time.Now().Unix() / int64(everySec)),
		capacity: capacity,
		everySec: everySec,
	}
}

func (l *limiterFixedWindow) Allow(n int) bool {
	curW := int(time.Now().Unix() / int64(l.everySec))
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

	//limiter := newLimiterFixedWindow(10000, 1)
	limiter := newLimiterTokenBucket(10000, 1, 100)

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

		if int(t1.Sub(t0).Seconds()) == 2 {
			fmt.Println("sleep 2 seconds now")
			time.Sleep(2 * time.Second)
		}

		// Run duration
		if t1.Sub(t0).Seconds() >= dur.Seconds() {
			break
		}
	}
	fmt.Println("duration: ", time.Now().Sub(t0))
}
