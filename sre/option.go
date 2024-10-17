package sre

import (
	"time"
)

const (
	StateOpen int32 = iota
	StateClose
)

type BreakerOptions struct {
	successNum   float64
	requestNum   int64
	bucketSize   int
	windowPeriod time.Duration
}
type Option func(*BreakerOptions)

func WithSuccess(success float64) Option {
	return func(o *BreakerOptions) {
		o.successNum = success
	}
}
func WithRequest(request int64) Option {
	return func(o *BreakerOptions) {
		o.requestNum = request
	}
}
func WithBucket(bucket int) Option {
	return func(o *BreakerOptions) {
		o.bucketSize = bucket
	}
}
func WithWindowPeriod(windowPeriod time.Duration) Option {
	return func(o *BreakerOptions) {
		o.windowPeriod = windowPeriod
	}
}
