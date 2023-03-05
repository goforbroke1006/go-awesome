package awesome

import (
	"errors"
	"time"
)

type Semaphore interface {
	Acquire() error
	Release() error
}

func NewSemaphoreSimple(limit uint) Semaphore {
	return &semaphoreSimple{
		semaphore: make(chan struct{}, limit),
	}
}

var _ Semaphore = (*semaphoreSimple)(nil)

type semaphoreSimple struct {
	semaphore chan struct{}
}

func (s semaphoreSimple) Acquire() error {
	s.semaphore <- struct{}{}
	return nil
}

func (s semaphoreSimple) Release() error {
	select {
	case _ = <-s.semaphore:
		//
	default:
		//
	}

	return nil
}

var (
	ErrSemaphoreTimeout = errors.New("semaphore timeout")
)

func NewSemaphoreWithTimeout(limit uint, timeout time.Duration) Semaphore {
	return &semaphoreWithTimeout{
		semaphore: make(chan struct{}, limit),
		timeout:   timeout,
	}
}

var _ Semaphore = (*semaphoreWithTimeout)(nil)

type semaphoreWithTimeout struct {
	semaphore chan struct{}
	timeout   time.Duration
}

func (s semaphoreWithTimeout) Acquire() error {
	select {
	case s.semaphore <- struct{}{}:
		return nil
	case <-time.After(s.timeout):
		return ErrSemaphoreTimeout
	}
}

func (s semaphoreWithTimeout) Release() error {
	select {
	case _ = <-s.semaphore:
		//
	default:
		//
	}

	return nil
}
