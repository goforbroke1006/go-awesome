package awesome

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker_Process(t *testing.T) {
	var (
		successfulJob = func() error { return nil }
		errFake       = errors.New("fake one")
		failedJob     = func() error { return errFake }
	)

	t.Run("always close", func(t *testing.T) {
		cb := NewCircuitBreaker(1, 3, time.Second)

		for idx := 0; idx < 10000; idx++ {
			assert.Nil(t, cb.Process(successfulJob))
		}
	})

	t.Run("open on failures limit", func(t *testing.T) {
		cb := NewCircuitBreaker(1, 3, time.Second)

		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, ErrCircuitBreakerOpened, cb.Process(failedJob))
		assert.Equal(t, ErrCircuitBreakerOpened, cb.Process(failedJob))
		assert.Equal(t, ErrCircuitBreakerOpened, cb.Process(failedJob))
		time.Sleep(time.Second + 100*time.Millisecond)
		assert.Equal(t, errFake, cb.Process(failedJob))
	})

	t.Run("open -> half-open -> open", func(t *testing.T) {
		cb := NewCircuitBreaker(1, 3, time.Second)

		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, ErrCircuitBreakerOpened, cb.Process(failedJob))
		assert.Equal(t, ErrCircuitBreakerOpened, cb.Process(failedJob))
		assert.Equal(t, ErrCircuitBreakerOpened, cb.Process(failedJob))
		time.Sleep(time.Second + 100*time.Millisecond)
		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, ErrCircuitBreakerOpened, cb.Process(failedJob))
	})

	t.Run("open -> half-open -> closed", func(t *testing.T) {
		cb := NewCircuitBreaker(1, 3, time.Second)

		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, ErrCircuitBreakerOpened, cb.Process(failedJob))
		assert.Equal(t, ErrCircuitBreakerOpened, cb.Process(failedJob))
		assert.Equal(t, ErrCircuitBreakerOpened, cb.Process(failedJob))
		time.Sleep(time.Second + 100*time.Millisecond)
		assert.Nil(t, cb.Process(successfulJob))
		assert.Equal(t, circuitBreakerStatusClosed, cb.status)
	})

	t.Run("success reset fails counter", func(t *testing.T) {
		cb := NewCircuitBreaker(1, 3, time.Second)

		for idx := 0; idx < 10000; idx++ {
			assert.Equal(t, errFake, cb.Process(failedJob))
			assert.Equal(t, errFake, cb.Process(failedJob))
			assert.Nil(t, cb.Process(successfulJob))
		}
		assert.Equal(t, circuitBreakerStatusClosed, cb.status)
	})

	t.Run("success dont accumulate", func(t *testing.T) {
		cb := NewCircuitBreaker(1, 1, time.Second)

		assert.Nil(t, cb.Process(successfulJob))
		assert.Nil(t, cb.Process(successfulJob))
		assert.Nil(t, cb.Process(successfulJob))
		assert.Nil(t, cb.Process(successfulJob))
		assert.Nil(t, cb.Process(successfulJob))
		assert.Equal(t, errFake, cb.Process(failedJob))
		assert.Equal(t, circuitBreakerStatusOpen, cb.status)
	})
}
