package channel_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/channel"
	"sync"
)

func ExampleInit() {
	var w sync.WaitGroup
	var receiver = make(map[int]chan interface{})

	w.Add(5)

	con := channel.Init()

	for i := 0; i < 5; i++ {
		receiver[i] = con.AddReceiver()

		go func(j int) {
			fmt.Println(<-receiver[j])
			w.Done()
		}(i)
	}

	transmitter01 := con.AddTransmitter()

	transmitter01 <- "Hello World"

	w.Wait()

	for i := 0; i < 5; i++ {
		con.CloseReceiver(receiver[i])
	}
	// Output:
	// Hello World
	// Hello World
	// Hello World
	// Hello World
	// Hello World
}
