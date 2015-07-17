package re_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/re"
	"time"
)

func ExampleDo() {
	values := []string{"a", "b", "c", "d", "e"}
	data, stop := re.Do(time.Millisecond*1, func(data chan<- interface{}) {
		for i := 0; i < 5; i++ {
			data <- fmt.Sprintf("%v", values[i])
		}
	})
	j := 0
	for i := 0; i < 5; i++ {
		select {
		case x := <-data:
			fmt.Printf("%v:%v\n", x, i)
			j++
		}
	}
	stop <- true

	// Output:
	// a:0
	// b:1
	// c:2
	// d:3
	// e:4
}
