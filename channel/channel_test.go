// Package channel adds more channel functionality.
// In Golang it is easy to send from multiple transmitters to one receiver
// but it is complicated to receive the same data an multiple endpoints from one transmitter
package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_Init(t *testing.T) {
	var w sync.WaitGroup
	w.Add(3)

	con := Init()

	receiver01 := con.AddReceiver()
	receiver02 := con.AddReceiver()
	receiver03 := con.AddReceiver()

	go func() {
		fmt.Println(<-receiver01)
		w.Done()
	}()

	go func() {
		fmt.Println(<-receiver02)
		w.Done()
	}()

	go func() {
		fmt.Println(<-receiver03)
		w.Done()
	}()

	transmitter01 := con.AddTransmitter()

	transmitter01 <- "Hello World"

	w.Wait()

	//closer count is not synchronous
	fmt.Println("receiver remaining:", con.CloseReceiver(receiver01))
	fmt.Println("receiver remaining:", con.CloseReceiver(receiver02))
	fmt.Println("receiver remaining:", con.CloseReceiver(receiver03))
	time.Sleep(5 * time.Millisecond)

	fmt.Println("receiver:", con.CountReceiver())
}

func Test_WritingToDeadChannel(t *testing.T) {
	var w sync.WaitGroup
	var receiver = make(map[int]chan interface{})

	w.Add(4)

	con := Init()

	for i := 0; i < 5; i++ {
		receiver[i] = con.AddReceiver()

		go func(j int) {
			fmt.Println(<-receiver[j])
			w.Done()
		}(i)
	}

	con.CloseReceiver(receiver[2])

	transmitter01 := con.AddTransmitter()

	transmitter01 <- "Lorem"

	w.Wait()

	for i := 0; i < 5; i++ {
		con.CloseReceiver(receiver[i])
	}
}

func Test_CloseFromTheInside(t *testing.T) {
	var w sync.WaitGroup
	con := Init()
	w.Add(1)

	receiver := con.AddReceiver()

	go func() {
		fmt.Println(<-receiver)
		con.CloseReceiver(receiver)
		w.Done()
	}()

	transmitter := con.AddTransmitter()
	transmitter <- "Ipsum"

	w.Wait()
	time.Sleep(5 * time.Millisecond)

	if count := con.CountReceiver(); count != 0 {
		t.Fatalf("channel count is %d, should be 0\n", count)
	}
}
