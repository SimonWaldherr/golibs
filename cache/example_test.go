package cache_test

import (
	"bytes"
	"fmt"
	"time"

	"simonwaldherr.de/go/golibs/cache"
)

func ExampleCache() {
	c := cache.New(5*time.Second, 1*time.Second)
	defer c.Stop()

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

func ExampleCache_Export() {
	c := cache.New(5*time.Second, 1*time.Second)
	defer c.Stop()

	c.Add("0", 23)
	c.Add("2", "foobar")
	c.Add("4", false)
	c.Add("8", []byte("test"))

	var buf bytes.Buffer

	fmt.Println(c.String())

	if err := c.Export(&buf); err != nil {
		fmt.Println(err)
		return
	}

	c2 := cache.New(15*time.Second, 1*time.Second)
	defer c2.Stop()
	if err := c2.Import(&buf); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(c2.String())

	// Output:
	//0	23
	//2	foobar
	//4	false
	//8	[116 101 115 116]
	//
	//0	23
	//2	foobar
	//4	false
	//8	[116 101 115 116]
}
