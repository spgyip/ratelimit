package main

import (
	"go.uber.org/ratelimit"
)

//rl := ratelimit.New(100) // per second
//    prev := time.Now()
//    for i := 0; i < 10; i++ {
//        now := rl.Take()
//        fmt.Println(i, now.Sub(prev))
//        prev = now
//    }

type uberLimiter struct {
	limiter ratelimit.Limiter
}

func (l *uberLimiter) Name() string {
	return "uber"
}

func newLimiterUber(r int) *uberLimiter {
	return &uberLimiter{
		limiter: ratelimit.New(r),
	}
}

func (l *uberLimiter) Allow(n int) bool {
	l.limiter.Take()
	return true
}
