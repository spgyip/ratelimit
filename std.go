package main

import (
	"time"

	"golang.org/x/time/rate"
)

type stdLimiter struct {
	limiter *rate.Limiter
}

func (l *stdLimiter) Name() string {
	return "std"
}

func newLimiterStd(r int, burst int) *stdLimiter {
	return &stdLimiter{
		limiter: rate.NewLimiter(rate.Limit(r), burst),
	}
}

func (l *stdLimiter) Allow(n int) bool {
	return l.limiter.AllowN(time.Now(), n)
}
