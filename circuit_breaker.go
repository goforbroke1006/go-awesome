package awesome

import (
	"errors"
	"sync"
	"time"
)

func NewCircuitBreaker(success, failure uint, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		success: success,
		failure: failure,
		timeout: timeout,

		status:         circuitBreakerStatusClosed,
		counterSuccess: 0,
		counterFailure: 0,
	}
}

var ErrCircuitBreakerOpened = errors.New("circuit breaker opened")

type CircuitBreaker struct {
	mu sync.Mutex

	success uint
	failure uint

	timeout time.Duration
	status  circuitBreakerStatus

	counterSuccess uint
	counterFailure uint
}

func (cb *CircuitBreaker) Process(f func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.status == circuitBreakerStatusOpen {
		return ErrCircuitBreakerOpened
	}

	err := f()

	if cb.status == circuitBreakerStatusHalfOpen {
		if err == nil {
			cb.counterSuccess++
			if cb.counterSuccess >= cb.success {
				cb.close()
			}
		} else {
			cb.open()
			go cb.awaitHalfOpen()
		}
		return err
	}

	// closed
	if err != nil {
		cb.counterFailure++
		if cb.counterFailure >= cb.failure {
			cb.open()
			go cb.awaitHalfOpen()
		}
	} else {
		cb.counterFailure = 0
	}
	return err
}

func (cb *CircuitBreaker) close() {
	cb.status = circuitBreakerStatusClosed
	cb.counterSuccess = 0
	cb.counterFailure = 0
}

func (cb *CircuitBreaker) open() {
	cb.status = circuitBreakerStatusOpen
	cb.counterSuccess = 0
	cb.counterFailure = 0
}

func (cb *CircuitBreaker) awaitHalfOpen() {
	timer := time.NewTimer(cb.timeout)
	<-timer.C

	cb.mu.Lock()
	cb.status = circuitBreakerStatusHalfOpen
	cb.mu.Unlock()
}

type circuitBreakerStatus uint16

const (
	circuitBreakerStatusClosed   = circuitBreakerStatus(1)
	circuitBreakerStatusHalfOpen = circuitBreakerStatus(2)
	circuitBreakerStatusOpen     = circuitBreakerStatus(3)
)
