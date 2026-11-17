package helper

import (
	"errors"
	"log"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type CircuitState string

const (
	StateClosed   CircuitState = "CLOSED"
	StateOpen     CircuitState = "OPEN"
	StateHalfOpen CircuitState = "HALF_OPEN"
)

// CircuitBreaker protects downstream calls from retry storms during partner outages.
type CircuitBreaker struct {
	mu                  sync.Mutex
	state               CircuitState
	consecutiveFailures int
	failureThreshold    int
	cooldown            time.Duration
	openedAt            time.Time
	halfOpenInFlight    bool
}

func NewCircuitBreaker(failureThreshold int, cooldown time.Duration) *CircuitBreaker {
	if failureThreshold < 1 {
		failureThreshold = 3
	}
	if cooldown <= 0 {
		cooldown = 5 * time.Second
	}
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		cooldown:         cooldown,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error, isSystemFailure func(error) bool) error {
	cb.mu.Lock()
	switch cb.state {
	case StateOpen:
		if time.Since(cb.openedAt) < cb.cooldown {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
		cb.state = StateHalfOpen
		cb.halfOpenInFlight = true
		log.Printf("[circuit-breaker] state=HALF_OPEN (probing partner)")
	case StateHalfOpen:
		if cb.halfOpenInFlight {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
		cb.halfOpenInFlight = true
	}
	cb.mu.Unlock()

	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.halfOpenInFlight = false

	if err == nil {
		if cb.state != StateClosed {
			log.Printf("[circuit-breaker] state=CLOSED (partner recovered)")
		}
		cb.state = StateClosed
		cb.consecutiveFailures = 0
		return nil
	}

	if isSystemFailure != nil && !isSystemFailure(err) {
		// Business/non-retryable errors do not indicate partner outage.
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
		}
		cb.consecutiveFailures = 0
		return err
	}

	cb.consecutiveFailures++
	if cb.state == StateHalfOpen || cb.consecutiveFailures >= cb.failureThreshold {
		if cb.state != StateOpen {
			log.Printf("[circuit-breaker] state=OPEN failures=%d cooldown=%s", cb.consecutiveFailures, cb.cooldown)
		}
		cb.state = StateOpen
		cb.openedAt = time.Now()
	}

	return err
}
