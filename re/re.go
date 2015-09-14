// wanna do something several times?
// Then this is the correct package for you.
package re

import (
	"fmt"
	"time"
)

var MaxAttempts = 10

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

func Do(wait time.Duration, fn func(chan<- interface{})) (<-chan interface{}, chan<- bool) {
	stop := make(chan bool, 1)
	ret := make(chan interface{})

	go func() {
		defer func() {
			close(stop)
			time.Sleep(wait)
			close(ret)
		}()

		var stopbool = false
		for stopbool == false {
			select {
			case <-stop:
				stopbool = true
			default:
				go fn(ret)
				time.Sleep(wait)
			}
		}
	}()
	return ret, stop
}
