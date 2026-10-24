// Package semaphore limits concurrent work using context-aware permit acquisition.
package semaphore

import (
	"context"
	"errors"
)

var (
	// ErrInvalidCapacity reports a nonpositive capacity supplied to New.
	ErrInvalidCapacity = errors.New("semaphore: capacity must be positive")
	// ErrUninitialized reports use of a zero or nil Semaphore.
	ErrUninitialized = errors.New("semaphore: use New to initialize the semaphore")
	// ErrNotAcquired reports a release when the permit pool is empty.
	ErrNotAcquired = errors.New("semaphore: no acquired permit to release")
)

// Semaphore limits the number of simultaneously acquired permits.
// Its methods are safe for concurrent use. Initialize it with New; its zero value
// cannot acquire permits. Copying an initialized Semaphore shares its permit pool.
// Call New to create an independent pool. Acquisition order is unspecified.
//
// Permits have no goroutine ownership: callers must release exactly once for each
// successful acquisition. Release detects an empty pool, but cannot detect a
// mistaken release while another caller holds a permit.
type Semaphore struct {
	permits chan struct{}
}

// New constructs a semaphore with a positive capacity and no acquired permits.
func New(capacity int) (*Semaphore, error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}
	return &Semaphore{permits: make(chan struct{}, capacity)}, nil
}

// Acquire waits for a permit or context cancellation. A context that is already
// canceled cannot acquire a permit. If cancellation and acquisition become ready
// concurrently, either may win. Every nil result requires a matching Release,
// even if ctx is canceled afterward. Like other context APIs, ctx must be nonnil.
func (s *Semaphore) Acquire(ctx context.Context) error {
	if s == nil || s.permits == nil {
		return ErrUninitialized
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case s.permits <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire acquires a permit without waiting. It returns false when all permits
// are held, or when s is a zero or nil Semaphore.
func (s *Semaphore) TryAcquire() bool {
	if s == nil || s.permits == nil {
		return false
	}
	select {
	case s.permits <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release returns one acquired permit without waiting. It returns ErrNotAcquired
// if no permit is held, or ErrUninitialized for a zero or nil Semaphore.
func (s *Semaphore) Release() error {
	if s == nil || s.permits == nil {
		return ErrUninitialized
	}
	select {
	case <-s.permits:
		return nil
	default:
		return ErrNotAcquired
	}
}

// Capacity returns the permit limit, or zero for a zero or nil Semaphore.
func (s *Semaphore) Capacity() int {
	if s == nil {
		return 0
	}
	return cap(s.permits)
}

// InUse returns a momentary count of acquired permits, or zero for a zero or nil
// Semaphore. Concurrent operations may change the count immediately afterward.
func (s *Semaphore) InUse() int {
	if s == nil {
		return 0
	}
	return len(s.permits)
}
