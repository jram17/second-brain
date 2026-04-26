package breaker

import (
    "net"
    "time"

    "github.com/sony/gobreaker"
)

func New(name string) *gobreaker.CircuitBreaker {
    return gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        name,
        MaxRequests: 3,
        Interval:    10 * time.Second,
        Timeout:     30 * time.Second,
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            return counts.ConsecutiveFailures >= 3
        },
        IsSuccessful: func(err error) bool {
            if err == nil {
                return true
            }
            // only count network-level failures
            _, isNet := err.(*net.OpError)
            return !isNet
        },
    })
}
