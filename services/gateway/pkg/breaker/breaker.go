package breaker

import (
    "time"
    "github.com/sony/gobreaker"
)

func New(name string) *gobreaker.CircuitBreaker {
    return gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        name,
        MaxRequests: 3,                // allow a few request in half-open
        Interval:    10 * time.Second, // reset counter
        Timeout:     30 * time.Second, // open state duration
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            return counts.ConsecutiveFailures >= 3
        },
    })
}
