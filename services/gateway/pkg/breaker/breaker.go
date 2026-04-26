package breaker

import (
    "time"

    "github.com/sony/gobreaker"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// isInfraError returns true only for errors that indicate the service is actually down.
func isInfraError(err error) bool {
    s, ok := status.FromError(err)
    if !ok {
        return true // non-gRPC error, treat as infra failure
    }
    switch s.Code() {
    case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
        return true
    default:
        return false // business logic errors (InvalidArgument, AlreadyExists, etc.)
    }
}

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
            return !isInfraError(err)
        },
    })
}
