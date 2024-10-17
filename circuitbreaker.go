package main

import (
	"google.golang.org/protobuf/types/known/durationpb"
	v1 "maas-gateway/api/gateway/middleware/circuitbreaker/v1"
	"maas-gateway/middleware/circuitbreaker/sre"
	"time"
)

var Breaker sre.CircuitBreaker

func init() {
	trigger := v1.CircuitBreaker_SuccessRatio{SuccessRatio: &v1.SuccessRatio{
		Success: 0.6,
		Request: 50,
		Bucket:  8,
		Window:  durationpb.New(time.Second * 3),
	}}
	Breaker = NewCircuitBreaker(trigger)

}
func NewCircuitBreaker(trigger v1.CircuitBreaker_SuccessRatio) sre.CircuitBreaker {
	var opts []sre.Option
	opts = append(opts, sre.WithSuccess(trigger.SuccessRatio.Success))
	opts = append(opts, sre.WithBucket(int(trigger.SuccessRatio.Bucket)))
	opts = append(opts, sre.WithRequest(int64(trigger.SuccessRatio.Request)))
	opts = append(opts, sre.WithWindowPeriod(trigger.SuccessRatio.Window.AsDuration()))
	return sre.NewBreaker(opts...)
}
func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := Breaker.Allow(); err != nil {
			Breaker.MarkFailure()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"circuit-breaker": err.Error()})
		}
		err := c.Next()
		if err != nil || c.Response().StatusCode() != fiber.StatusOK {
			Breaker.MarkFailure()
		}
		Breaker.MarkSuccess()
		return nil
	}
}
