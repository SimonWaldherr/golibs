package cache_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/cache"
	"time"
)

func ExampleCache() {
	c := cache.New(5*time.Second, 1*time.Second)

	c.Add("0", 23)
	c.Add("1", 1*time.Second)
	c.Add("2", "foobar")
	c.Add("3", false)
	c.Add("4", []byte("test"))

	fmt.Println(c.String())

	// Output:
	//0	23
	//1	1s
	//2	foobar
	//3	false
	//4	[116 101 115 116]
}
