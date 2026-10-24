package semaphore

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func mustNew(t *testing.T, capacity int) *Semaphore {
	t.Helper()
	s, err := New(capacity)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestCapacityAndRelease(t *testing.T) {
	for _, capacity := range []int{-1, 0} {
		if s, err := New(capacity); s != nil || !errors.Is(err, ErrInvalidCapacity) {
			t.Fatalf("New(%d) = %v, %v", capacity, s, err)
		}
	}
	s := mustNew(t, 3)
	if s.Capacity() != 3 || s.InUse() != 0 {
		t.Fatal("unexpected initial statistics")
	}
	for i := 0; i < 3; i++ {
		if !s.TryAcquire() || s.InUse() != i+1 {
			t.Fatalf("failed to acquire permit %d", i+1)
		}
	}
	if s.TryAcquire() {
		t.Fatal("acquired beyond capacity")
	}
	for i := 0; i < 3; i++ {
		if err := s.Release(); err != nil || s.InUse() != 2-i {
			t.Fatalf("release %d: %v, in use %d", i, err, s.InUse())
		}
	}
	if err := s.Release(); !errors.Is(err, ErrNotAcquired) {
		t.Fatalf("unbalanced release = %v", err)
	}
}

func TestZeroAndNilAreNonblocking(t *testing.T) {
	var zero Semaphore
	for _, s := range []*Semaphore{nil, &zero} {
		if err := s.Acquire(context.Background()); !errors.Is(err, ErrUninitialized) {
			t.Fatalf("uninitialized Acquire = %v", err)
		}
		if s.TryAcquire() || s.Capacity() != 0 || s.InUse() != 0 {
			t.Fatal("uninitialized semaphore reported available permits")
		}
		if err := s.Release(); !errors.Is(err, ErrUninitialized) {
			t.Fatalf("uninitialized Release = %v", err)
		}
	}
}

func TestCopiesSharePermits(t *testing.T) {
	s := mustNew(t, 1)
	alias := *s
	if !s.TryAcquire() || alias.TryAcquire() || alias.InUse() != 1 {
		t.Fatal("copy did not share the same pool")
	}
	if err := alias.Release(); err != nil || s.InUse() != 0 {
		t.Fatalf("release through copy: %v, %d", err, s.InUse())
	}
	independent := mustNew(t, 1)
	if !independent.TryAcquire() {
		t.Fatal("independent pool shared capacity")
	}
	if err := independent.Release(); err != nil {
		t.Fatal(err)
	}
}

func TestAcquireAlreadyCanceledOrExpired(t *testing.T) {
	s := mustNew(t, 1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := s.Acquire(ctx); !errors.Is(err, context.Canceled) || s.InUse() != 0 {
		t.Fatalf("already canceled Acquire = %v, %d", err, s.InUse())
	}
	ctx, cancel = context.WithDeadline(context.Background(), time.Unix(0, 0))
	defer cancel()
	if err := s.Acquire(ctx); !errors.Is(err, context.DeadlineExceeded) || s.InUse() != 0 {
		t.Fatalf("expired Acquire = %v, %d", err, s.InUse())
	}
}

// Observing Done proves that Acquire passed the initial Err check and evaluated
// its blocking select. Tests can then release or cancel without arbitrary sleeps.
type observedContext struct {
	context.Context
	entered chan struct{}
	once    sync.Once
}

func (ctx *observedContext) Done() <-chan struct{} {
	ctx.once.Do(func() { close(ctx.entered) })
	return ctx.Context.Done()
}

func waitEntered(t *testing.T, entered <-chan struct{}) {
	t.Helper()
	select {
	case <-entered:
	case <-time.After(3 * time.Second):
		t.Fatal("Acquire did not start waiting")
	}
}

func awaitResult(t *testing.T, result <-chan error) error {
	t.Helper()
	select {
	case err := <-result:
		return err
	case <-time.After(3 * time.Second):
		t.Fatal("Acquire did not finish")
		return nil
	}
}

func TestAcquireWaitsUntilRelease(t *testing.T) {
	s := mustNew(t, 1)
	if !s.TryAcquire() {
		t.Fatal("initial acquire failed")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	observed := &observedContext{Context: ctx, entered: make(chan struct{})}
	result := make(chan error, 1)
	go func() { result <- s.Acquire(observed) }()
	waitEntered(t, observed.entered)
	select {
	case err := <-result:
		t.Fatalf("Acquire finished before a permit was released: %v", err)
	default:
	}
	if err := s.Release(); err != nil {
		t.Fatal(err)
	}
	if err := awaitResult(t, result); err != nil || s.InUse() != 1 {
		t.Fatalf("waiting Acquire = %v, %d", err, s.InUse())
	}
	if err := s.Release(); err != nil {
		t.Fatal(err)
	}
}

func TestCancellationWhileWaitingDoesNotConsumePermit(t *testing.T) {
	s := mustNew(t, 1)
	if !s.TryAcquire() {
		t.Fatal("initial acquire failed")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	observed := &observedContext{Context: ctx, entered: make(chan struct{})}
	result := make(chan error, 1)
	go func() { result <- s.Acquire(observed) }()
	waitEntered(t, observed.entered)
	cancel()
	if err := awaitResult(t, result); !errors.Is(err, context.Canceled) || s.InUse() != 1 {
		t.Fatalf("canceled Acquire = %v, %d", err, s.InUse())
	}
	if err := s.Release(); err != nil || !s.TryAcquire() {
		t.Fatalf("cancellation lost pool capacity: %v", err)
	}
	if err := s.Release(); err != nil {
		t.Fatal(err)
	}
}

func TestConcurrentCapacityBound(t *testing.T) {
	const capacity, workers = 4, 48
	s := mustNew(t, capacity)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	start := make(chan struct{})
	var active int32
	var group sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			for i := 0; i < 50; i++ {
				if err := s.Acquire(ctx); err != nil {
					t.Error(err)
					return
				}
				if got := atomic.AddInt32(&active, 1); got > capacity {
					t.Errorf("active work %d exceeded capacity %d", got, capacity)
				}
				if got := s.InUse(); got < 1 || got > capacity {
					t.Errorf("in-use snapshot outside bounds: %d", got)
				}
				atomic.AddInt32(&active, -1)
				if err := s.Release(); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	close(start)
	group.Wait()
	if active != 0 || s.InUse() != 0 {
		t.Fatalf("permits leaked: active %d, in use %d", active, s.InUse())
	}
}
