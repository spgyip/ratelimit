package main

import (
	"flag"
	"fmt"
	"time"
)

func doSomething() {
}

func main() {
	var dur time.Duration
	var typ string
	flag.DurationVar(&dur, "dur", 10*time.Second, "Run duration")
	flag.StringVar(&typ, "type", "fw", "Type of limiter(fw|sl|tb|std)")
	flag.Parse()

	var lmt limiter
	switch typ {
	case "fw":
		lmt = newLimiterFixedWindow(10000)
	case "sl":
		lmt = newLimiterSlidingWindow(10000)
	case "tb":
		lmt = newLimiterTokenBucket(10000, 100)
	case "std":
		lmt = newLimiterStd(10000, 100)
	default:
		fmt.Printf("unknown type: %v\n", typ)
		return
	}

	fmt.Printf("Using limiter %v, run duration %v.\n", lmt.Name(), dur)

	t0 := time.Now()
	t1, lastT := time.Now(), time.Now()
	cntAllow := 0
	cntReject := 0
	for t1.Sub(t0).Seconds() < dur.Seconds() {
		// Call func
		if lmt.Allow(1) {
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

	}
	fmt.Println("duration: ", time.Now().Sub(t0))
}
