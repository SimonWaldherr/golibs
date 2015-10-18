package log_test

import (
	"simonwaldherr.de/go/golibs/log"
)

func Example() {
	log.Info.Println("1")
	log.Warning.Println("2")
	log.Error.Println("3")
	log.Fatal.Println("4")
}
