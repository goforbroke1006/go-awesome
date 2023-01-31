package awesome

import (
	"errors"
	"sync"
	"time"
)

var ErrThrottlerCancellation = errors.New("throttler cancellation")

func NewThrottler() *throttler {
	return &throttler{
		invocations: map[interface{}]time.Time{},
	}
}

type throttler struct {
	mu sync.Mutex

	invocations map[interface{}]time.Time
}

func (t *throttler) Throttle(key interface{}, period time.Duration) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if last, ok := t.invocations[key]; ok && time.Since(last) < period {
		return ErrThrottlerCancellation
	}

	t.invocations[key] = time.Now()

	go t.awaitRemoving(key, period)

	return nil
}

func (t *throttler) awaitRemoving(key interface{}, period time.Duration) {
	time.Sleep(period)

	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.invocations, key)
}
