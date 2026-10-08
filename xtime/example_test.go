package xtime_test

import (
	"fmt"
	"time"

	"simonwaldherr.de/go/golibs/xtime"
)

func ExampleFmt() {
	t := time.Unix(1234567890, 0)
	cet, _ := time.LoadLocation("CET")
	t = t.In(cet)
	fmt.Println(xtime.Fmt("%Y-%m-%d %H:%M:%S", t))
	// Also demonstrate parsing an ANSIC time string and formatting it.
	t2, _ := time.Parse(time.ANSIC, "Mon Jan 02 15:04:05 2006")
	fmt.Println(xtime.Fmt("%Y-%m-%d %H:%M:%S", t2))

	// Output:
	// 2009-02-14 00:31:30
	// 2006-01-02 15:04:05
}
