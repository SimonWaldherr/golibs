// Package retry retries failing operations with a fixed delay and cancellation.
package retry

import (
	"context"
	"errors"
	"time"
)

// ErrAttempts reports a nonpositive attempt limit.
var ErrAttempts = errors.New("retry: attempts must be positive")

// Do calls fn at most attempts times, returning immediately on the first nil
// result. After exhausting attempts it returns the final function error unchanged.
// It checks ctx before each attempt and during delays; fn must respect ctx to
// cancel an active call. Nonpositive delays retry immediately. Do does not recover
// panics. A successful call takes precedence over concurrent context cancellation.
func Do(ctx context.Context, attempts int, delay time.Duration, fn func(context.Context) error) error {
	if attempts < 1 {
		return ErrAttempts
	}
	var timer *time.Timer
	for remaining := attempts; ; remaining-- {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := fn(ctx)
		if err == nil || remaining == 1 {
			return err
		}
		if delay > 0 {
			if timer == nil {
				timer = time.NewTimer(delay)
			} else {
				timer.Reset(delay)
			}
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}
}
