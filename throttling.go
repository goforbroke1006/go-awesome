package awesome

import (
	"errors"
	"sync"
	"time"
)

var ErrThrottlerCancellation = errors.New("throttler cancellation")

type Throttler interface {
	Throttle(key interface{}) error
}

func NewThrottlerFirstEntry(period time.Duration) Throttler {
	return &throttlerFirstEntry{
		period:      period,
		invocations: map[interface{}]time.Time{},
	}
}

var _ Throttler = (*throttlerFirstEntry)(nil)

type throttlerFirstEntry struct {
	period      time.Duration
	mu          sync.Mutex
	invocations map[interface{}]time.Time
}

func (t *throttlerFirstEntry) Throttle(key interface{}) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if last, ok := t.invocations[key]; ok && time.Since(last) < t.period {
		return ErrThrottlerCancellation
	}

	t.invocations[key] = time.Now()

	go t.awaitRemoving(key)

	return nil
}

func (t *throttlerFirstEntry) awaitRemoving(key interface{}) {
	time.Sleep(t.period)

	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.invocations, key)
}
