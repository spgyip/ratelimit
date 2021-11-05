package main

type limiter interface {
	Name() string
	Allow(n int) bool
}
