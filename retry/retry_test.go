package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestInvalidAttempts(t *testing.T) {
	for _, attempts := range []int{0, -1} {
		called := false
		err := Do(context.Background(), attempts, 0, func(context.Context) error {
			called = true
			return nil
		})
		if !errors.Is(err, ErrAttempts) || called {
			t.Fatalf("attempts %d: error = %v, called = %v", attempts, err, called)
		}
	}
}

func TestSuccessDoesNotWait(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	calls := 0
	err := Do(ctx, 3, time.Hour, func(received context.Context) error {
		if received != ctx {
			t.Fatal("callback did not receive the supplied context")
		}
		calls++
		return nil
	})
	if err != nil || calls != 1 {
		t.Fatalf("Do = %v, calls = %d; want immediate success", err, calls)
	}
}

func TestRetryThenSuccess(t *testing.T) {
	for _, delay := range []time.Duration{-time.Second, 0, time.Nanosecond} {
		calls := 0
		err := Do(context.Background(), 4, delay, func(context.Context) error {
			calls++
			if calls < 3 {
				return errors.New("try again")
			}
			return nil
		})
		if err != nil || calls != 3 {
			t.Fatalf("delay %v: Do = %v, calls = %d", delay, err, calls)
		}
	}
}

func TestExhaustedReturnsFinalError(t *testing.T) {
	want := errors.New("final error")
	calls := 0
	err := Do(context.Background(), 3, 0, func(context.Context) error {
		calls++
		if calls == 3 {
			return want
		}
		return errors.New("earlier error")
	})
	if err != want || calls != 3 {
		t.Fatalf("Do = %v, calls = %d; want final error and 3 calls", err, calls)
	}
	// There must be no delay after the final failed attempt either.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = Do(ctx, 1, time.Hour, func(context.Context) error { return want })
	if err != want {
		t.Fatalf("single failed attempt returned %v", err)
	}
}

func TestCancellationBeforeAttempt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Do(ctx, 3, 0, func(context.Context) error {
		t.Fatal("called an operation with an already canceled context")
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Do = %v, want context.Canceled", err)
	}
}

func TestCancellationBetweenAttempts(t *testing.T) {
	for _, delay := range []time.Duration{0, time.Hour} {
		ctx, cancel := context.WithCancel(context.Background())
		calls := 0
		err := Do(ctx, 3, delay, func(context.Context) error {
			calls++
			cancel()
			return errors.New("retryable failure")
		})
		cancel()
		if !errors.Is(err, context.Canceled) || calls != 1 {
			t.Fatalf("delay %v: Do = %v, calls = %d", delay, err, calls)
		}
	}
}

func TestSuccessTakesPrecedenceOverCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := Do(ctx, 3, time.Hour, func(context.Context) error {
		cancel()
		return nil
	}); err != nil {
		t.Fatalf("successful callback returned %v", err)
	}
}

func TestPanicsPropagate(t *testing.T) {
	defer func() {
		if got := recover(); got != "callback panic" {
			t.Errorf("panic = %v, want callback panic", got)
		}
	}()
	_ = Do(context.Background(), 3, 0, func(context.Context) error { panic("callback panic") })
	t.Fatal("Do swallowed the callback panic")
}

func TestDelayedRetriesReturnFinalError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	want := errors.New("last attempt")
	calls := 0
	err := Do(ctx, 5, time.Nanosecond, func(context.Context) error {
		calls++
		if calls == 5 {
			return want
		}
		return errors.New("earlier attempt")
	})
	if err != want || calls != 5 {
		t.Fatalf("Do = %v, calls = %d; want final error and five attempts", err, calls)
	}
}

func TestCancellationAfterDelayedRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls := 0
	err := Do(ctx, 5, time.Nanosecond, func(context.Context) error {
		calls++
		if calls == 2 {
			cancel()
		}
		return errors.New("retryable failure")
	})
	if !errors.Is(err, context.Canceled) || calls != 2 {
		t.Fatalf("Do = %v, calls = %d; want cancellation after second attempt", err, calls)
	}
}
