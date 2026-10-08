// wanna do something several times?
// Then this is the correct package for you.
package re

import (
	"fmt"
	"time"
)

var MaxAttempts = 10

// Try tries to execute a function several times
func Try(retrys int, fn func() (err error)) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	var err error

	for attempt := 1; ; attempt++ {
		err = fn()
		if err == nil {
			break
		}

		if attempt >= MaxAttempts || attempt >= retrys {
			return fmt.Errorf("Reached number of attempts (%v)\n%v", (attempt), err)
		}
	}
	return err
}

// Do executes fn sequentially, waiting wait after each callback returns.
// Close stop or send one value on it to request shutdown; Do does not close it.
// An active callback finishes before the output channel closes. Callers must
// continue consuming the output until it closes so callback sends can finish.
// The callback must return and must not retain or close the output channel.
func Do(wait time.Duration, fn func(chan<- interface{})) (<-chan interface{}, chan<- bool) {
	stop := make(chan bool, 1)
	ret := make(chan interface{})

	go func() {
		defer close(ret)
		for {
			select {
			case <-stop:
				return
			default:
			}
			fn(ret)

			timer := time.NewTimer(wait)
			select {
			case <-stop:
				timer.Stop()
				return
			case <-timer.C:
			}
		}
	}()
	return ret, stop
}
