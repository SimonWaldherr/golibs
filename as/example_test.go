package as_test

import (
	"encoding/json"
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
	values := []string{
		"false",
		"0",
		"3",
		"5€",
		"23$",
		"42¢",
		"1.337£",
		"020161586X",
		"13,37",
		"https://simonwaldherr.de/",
		"contact@simonwaldherr.de",
		"127.0.0.1",
		"rgb(255, 231, 200)",
		"#ffffff",
		"#03a",
		"06.06.1989",
		"2010/07/29",
		"2006-01-02",
		"¢«∑€®†Ω¨⁄øπ•±‘æœ@∆ºª©ƒ∂‚å≤¥≈ç√∫~µ∞…–¡“≠¿'¬”ﬁ",
	}

	fmt.Println("as.Type()")
	for _, v := range values {
		str, _ := as.Type(v)
		fmt.Printf("%v => %v\n", v, str)
	}

	// Output:
	// as.Type()
	// false => bool
	// 0 => bool
	// 3 => int
	// 5€ => price
	// 23$ => price
	// 42¢ => price
	// 1.337£ => price
	// 020161586X => isbn
	// 13,37 => float
	// https://simonwaldherr.de/ => url
	// contact@simonwaldherr.de => email
	// 127.0.0.1 => ipv4
	// rgb(255, 231, 200) => color
	// #ffffff => color
	// #03a => color
	// 06.06.1989 => date
	// 2010/07/29 => date
	// 2006-01-02 => date
	// ¢«∑€®†Ω¨⁄øπ•±‘æœ@∆ºª©ƒ∂‚å≤¥≈ç√∫~µ∞…–¡“≠¿'¬”ﬁ =>
}

type dynStruct struct {
	First  as.Dynamic `json:"first"`
	Second as.Dynamic `json:"second"`
	Third  as.Dynamic `json:"third"`
}

func ExampleType_json() {
	jsonStr := `
{
	"first": "2019-01-04 00:00:00",
	"second": "foobar@example.tld",
	"third": "15€"
}`

	x := dynStruct{}

	if err := json.Unmarshal([]byte(jsonStr), &x); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v: %v\n", x.First.Type, x.First.Data)
	fmt.Printf("%v: %v\n", x.Second.Type, x.Second.Data)
	fmt.Printf("%v: %v\n", x.Third.Type, x.Third.Data)

	// Output:
	// date: 2019-01-04 00:00:00 +0000 UTC
	// email: foobar@example.tld
	// price: 15€
}
