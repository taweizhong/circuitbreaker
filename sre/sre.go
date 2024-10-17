package sre

import (
	"errors"
	"maas-gateway/middleware/circuitbreaker/sre/window"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var ErrNotAllowed = errors.New("circuit breaker: not allowed for circuit open")

type CircuitBreaker interface {
	Allow() error
	MarkSuccess()
	MarkFailure()
}

type Breaker struct {
	stat     window.RollingCounter
	r        *rand.Rand
	randLock sync.Mutex
	k        float64
	request  int64
	state    int32
}

func NewBreaker(opts ...Option) *Breaker {
	opt := &BreakerOptions{
		successNum:   0.6,
		requestNum:   100,
		bucketSize:   10,
		windowPeriod: 3 * time.Second,
	}
	for _, o := range opts {
		o(opt)
	}
	counterOpts := window.RollingCounterOpts{
		Size:           opt.bucketSize,
		BucketDuration: time.Duration(int64(opt.windowPeriod) / int64(opt.bucketSize)),
	}
	stat := window.NewRollingCounter(counterOpts)
	return &Breaker{
		stat:    stat,
		r:       rand.New(rand.NewSource(int64(time.Now().UnixNano()))),
		request: opt.requestNum,
		k:       1 / opt.successNum,
		state:   StateClose,
	}
}
func (b *Breaker) Allow() error {
	accepts, total := b.Summery()
	requests := b.k * float64(accepts)
	if total < b.request || float64(total) < requests {
		atomic.CompareAndSwapInt32(&b.state, StateOpen, StateClose)
		return nil
	}
	atomic.CompareAndSwapInt32(&b.state, StateClose, StateOpen)
	dr := math.Max(0, (float64(total)-requests)/float64(total+1))
	drop := b.Judgment(dr)
	if drop {
		return ErrNotAllowed
	}
	return nil
}
func (b *Breaker) MarkSuccess() { b.stat.Add(1) }
func (b *Breaker) MarkFailure() { b.stat.Add(0) }
func (b *Breaker) Summery() (int64, int64) {
	var success, total = int64(0), int64(0)
	b.stat.Reduce(func(iterator window.Iterator) float64 {
		for iterator.Next() {
			bucket := iterator.Bucket()
			total += bucket.Count
			success += int64(bucket.Points[0])
		}
		// 为了满足函数的签名，不会影响实际的计算逻辑
		return 0
	})
	return success, total
}
func (b *Breaker) Judgment(p float64) bool {
	var flag bool
	b.randLock.Lock()
	flag = b.r.Float64() < p
	b.randLock.Unlock()
	return flag
}
