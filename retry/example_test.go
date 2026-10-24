package retry_test

import (
	"context"
	"errors"
	"fmt"

	"simonwaldherr.de/go/golibs/retry"
)

func ExampleDo() {
	attempts := 0
	err := retry.Do(context.Background(), 3, 0, func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary failure")
		}
		return nil
	})
	fmt.Println(attempts, err)
	// Output: 3 <nil>
}
