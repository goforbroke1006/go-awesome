package awesome

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_semaphoreSimple_Acquire(t *testing.T) {
	t.Parallel()

	t.Run("positive 1 - can not acquire more than limit", func(t *testing.T) {
		t.Parallel()

		const limit = 5

		sem := NewSemaphoreSimple(limit)

		var jobsCount uint32

		for idx := 0; idx < 100; idx++ {
			go func() {
				_ = sem.Acquire()
				atomic.AddUint32(&jobsCount, 1)
			}()
		}

		time.Sleep(time.Second)

		assert.Equal(t, jobsCount, uint32(limit))
	})
}

func Test_semaphoreSimple_Release(t *testing.T) {
	t.Parallel()

	t.Run("positive 1 - can acquire after single release", func(t *testing.T) {
		t.Parallel()

		const limit = 5

		sem := NewSemaphoreSimple(limit)

		var jobsCount uint32

		for idx := 0; idx < limit; idx++ {
			go func() {
				_ = sem.Acquire()
				atomic.AddUint32(&jobsCount, 1)
			}()
		}

		time.Sleep(time.Second)

		// release one job
		_ = sem.Release()

		acquireCh := make(chan struct{})
		go func() {
			_ = sem.Acquire()
			acquireCh <- struct{}{}
		}()

		select {
		case <-acquireCh:
			// everything okay
		case <-time.After(time.Second):
			t.Error("should acquire semaphore if it has been released")
		}
	})

	t.Run("positive 2 - can acquire after release all", func(t *testing.T) {
		t.Parallel()

		const limit = 5

		sem := NewSemaphoreSimple(limit)

		var jobsCount uint32

		for idx := 0; idx < limit; idx++ {
			go func() {
				_ = sem.Acquire()
				atomic.AddUint32(&jobsCount, 1)
			}()
		}

		time.Sleep(time.Second)

		// release one job
		for idx := 0; idx < 100; idx++ {
			_ = sem.Release()
		}

		acquireCh := make(chan struct{})
		go func() {
			_ = sem.Acquire()
			acquireCh <- struct{}{}
		}()

		select {
		case <-acquireCh:
			// everything okay
		case <-time.After(time.Second):
			t.Error("should acquire semaphore if it has been released")
		}
	})

	t.Run("negative 1 - can not acquire on exceeded one without release", func(t *testing.T) {
		t.Parallel()

		const limit = 5

		sem := NewSemaphoreSimple(limit)

		var jobsCount uint32

		for idx := 0; idx < 100; idx++ {
			go func() {
				_ = sem.Acquire()
				atomic.AddUint32(&jobsCount, 1)
			}()
		}

		time.Sleep(2 * time.Second)

		// no release

		acquireCh := make(chan struct{})
		go func() {
			_ = sem.Acquire()
			acquireCh <- struct{}{}
		}()

		select {
		case <-acquireCh:
			t.Error("should not acquire exceeded semaphore")
		case <-time.After(time.Second):
			// everything okay
		}
	})
}
