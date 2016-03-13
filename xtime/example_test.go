package xtime_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/xtime"
	"time"
)

func ExampleFmt() {
	t := time.Unix(1234567890, 0)
	cet, _ := time.LoadLocation("CET")
	t = t.In(cet)
	fmt.Println(xtime.Fmt("%Y-%m-%d %H:%M:%S", t))
	// Output:
	// 2009-02-14 00:31:30
}

func ExampleFmt2() {
	t, _ := time.Parse(time.ANSIC, "Mon Jan 02 15:04:05 2006")
	fmt.Println(xtime.Fmt("%Y-%m-%d %H:%M:%S", t))
	// Output:
	// 2006-01-02 15:04:05
}
