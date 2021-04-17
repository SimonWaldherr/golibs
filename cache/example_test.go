package cache_test

import (
	"bytes"
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

func ExampleExport() {
	c := cache.New(5*time.Second, 1*time.Second)

	c.Add("0", 23)
	c.Add("2", "foobar")
	c.Add("4", false)
	c.Add("8", []byte("test"))

	var buf bytes.Buffer

	fmt.Println(c.String())

	c.Export(&buf)

	c2 := cache.New(15*time.Second, 1*time.Second)
	c2.Import(&buf)

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
