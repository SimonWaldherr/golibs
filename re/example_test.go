package re_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"simonwaldherr.de/go/golibs/re"
	"time"
)

func ExampleDo() {
	values := []string{"a", "b", "c", "d", "e"}
	data, stop := re.Do(time.Millisecond*100, func(data chan<- interface{}) {
		for i := 0; i < 5; i++ {
			data <- fmt.Sprintf("%v", values[i])
		}
	})
	for i := 0; i < 5; i++ {
		x := <-data
		fmt.Printf("%v:%v\n", x, i)
	}
	stop <- true
	time.Sleep(time.Millisecond * 500)

	// Output:
	// a:0
	// b:1
	// c:2
	// d:3
	// e:4
}

func ExampleTry() {
	_ = re.Try(5, func() error {
		var err error
		response, err := http.Get("http://golang.org/")
		if err != nil {
			fmt.Printf("%s", err)
			return err
		} else {
			defer func() { _ = response.Body.Close() }()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("%s", err)
				return err
			}
			fmt.Printf("%s\n", string(contents))
		}
		return err
	})
}
