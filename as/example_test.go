package as_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/as"
	"time"
)

func ExampleBool() {
	values := []interface{}{"foobar", "true", 1, true}

	fmt.Println("as.Bool()")
	for _, v := range values {
		fmt.Printf("%v => %v\n", v, as.Bool(v))
	}

	// Output:
	// as.Bool()
	// foobar => false
	// true => true
	// 1 => true
	// true => true
}

func ExampleBytes() {
	values := []interface{}{"foobar", "true", 1, true}

	fmt.Println("as.Bytes()")
	for _, v := range values {
		fmt.Printf("%v => %v\n", v, as.Bytes(v))
	}

	// Output:
	// as.Bytes()
	// foobar => [102 111 111 98 97 114]
	// true => [116 114 117 101]
	// 1 => [49]
	// true => [116 114 117 101]
}

func ExampleString() {
	values := []interface{}{[]byte("foobar"), 1, true, 5 * time.Second, 0xfefe}

	fmt.Println("as.String()")
	for _, v := range values {
		fmt.Printf("%v => %v\n", v, as.String(v))
	}

	// Output:
	// as.String()
	// [102 111 111 98 97 114] => foobar
	// 1 => 1
	// true => true
	// 5s => 5s
	// 65278 => 65278
}

func ExampleTime() {
	values := []interface{}{"01.01.1970", "01.01.1999", "2000/10/30"}

	fmt.Println("as.Time()")
	for _, v := range values {
		fmt.Printf("%v => %v\n", v, as.Time(v))
	}

	// Output:
	// as.Time()
	// 01.01.1970 => 1970-01-01 00:00:00 +0000 UTC
	// 01.01.1999 => 1999-01-01 00:00:00 +0000 UTC
	// 2000/10/30 => 2000-10-30 00:00:00 +0000 UTC
}

func ExampleType() {
	values := []string{"5€", "13,37", "https://simonwaldherr.de/", "127.0.0.1", "#ffffff"}

	fmt.Println("as.Type()")
	for _, v := range values {
		str, _ := as.Type(v)
		fmt.Printf("%v => %v\n", v, str)
	}

	// Output:
	// as.Type()
	// 5€ => price
	// 13,37 => float
	// https://simonwaldherr.de/ => url
	// 127.0.0.1 => ipv4
	// #ffffff => color
}
