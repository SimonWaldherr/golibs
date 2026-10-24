package semaphore_test

import (
	"context"
	"errors"
	"fmt"

	"simonwaldherr.de/go/golibs/semaphore"
)

func ExampleSemaphore_TryAcquire() {
	limit, err := semaphore.New(2)
	if err != nil {
		panic(err)
	}
	fmt.Println(limit.TryAcquire(), limit.TryAcquire(), limit.TryAcquire())
	fmt.Println(limit.InUse(), limit.Capacity())
	if err := limit.Release(); err != nil {
		panic(err)
	}
	if err := limit.Release(); err != nil {
		panic(err)
	}
	// Output:
	// true true false
	// 2 2
}

func ExampleSemaphore_Acquire() {
	limit, err := semaphore.New(1)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = limit.Acquire(ctx)
	fmt.Println(errors.Is(err, context.Canceled), limit.InUse())
	// Output: true 0
}
