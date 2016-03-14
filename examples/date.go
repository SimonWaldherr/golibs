package main

import (
	"os"
	"simonwaldherr.de/go/golibs/xtime"
)

func main() {
	print(xtime.FmtNow(os.Args[1]))
	print("\n")
}
