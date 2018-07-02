package re

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func Test_Try(t *testing.T) {
	MaxAttempts = 20
	SomeFunction := func() (string, error) {
		return "", nil
	}
	var value string
	err := Try(5, func() error {
		var err error
		value, err = SomeFunction()
		_ = value
		return err
	})
	if err != nil {
		t.Fatalf("re.Try Test failed: %v\n", err)
	}
}

func Test_Try_Failing(t *testing.T) {
	MaxAttempts = 10
	callCount := 0
	err := Try(15, func() error {
		callCount++
		return errors.New("something went wrong")
	})
	if callCount != 10 {
		t.Fatalf("re.Try Failing failed: should: 10 is: %v\n", callCount)
	}
	if err == nil {
		t.Fatalf("re.Try Failing failed: err should not be nil\n")
	}
}

func Test_Try_Panic(t *testing.T) {
	panicFunction := func() (string, error) {
		panic("something went badly wrong")
	}
	err := Try(5, func() error {
		_, err := panicFunction()
		return err
	})
	if err != nil {
		t.Fatalf("re.Try panic Test failed: %v\n", err)
	}
}

func Test_Do(t *testing.T) {
	data, stop := Do(time.Millisecond*100, func(data chan<- interface{}) {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Microsecond * 600)
			data <- fmt.Sprintf("\n%v: %v", time.Now().Format("02.01.2006 15:04:05"), i)
		}
	})
	j := 0
	for i := 0; i < 5; i++ {
		select {
		case x := <-data:
			fmt.Printf("%v:1:%v", x, i)
			j++
		}
	}
	stop <- true
	for i := 0; i < 5; i++ {
		select {
		case x := <-data:
			fmt.Printf("%v:2:%v", x, i)
			j++
		}
	}
	fmt.Println()
	if j != 10 {
		t.Fatalf("re.Do Test failed: should: 10 is: %v\n", j)
	}
}
