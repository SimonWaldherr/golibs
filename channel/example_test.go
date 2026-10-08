package channel_test

import (
	"fmt"
	"sync"

	"simonwaldherr.de/go/golibs/channel"
)

func ExampleInit() {
	var w sync.WaitGroup
	var receiver = make(map[int]chan interface{})

	w.Add(5)

	con := channel.Init()
	defer con.Close()

	for i := 0; i < 5; i++ {
		receiver[i] = con.AddReceiver()

		go func(ch chan interface{}) {
			fmt.Println(<-ch)
			w.Done()
		}(receiver[i])
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
