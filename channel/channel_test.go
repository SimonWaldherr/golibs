// Package channel adds more channel functionality.
// In Golang it is easy to send from multiple transmitters to one receiver
// but it is complicated to receive the same data an multiple endpoints from one transmitter
package channel

import (
	"fmt"
	"sync"
	"testing"
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

	con.CloseReceiver(receiver01)
	con.CloseReceiver(receiver02)
	con.CloseReceiver(receiver03)
}
