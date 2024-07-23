package as_test

import (
	"encoding/json"
	"fmt"
	"time"

	"simonwaldherr.de/go/golibs/as"
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

func ExampleHexToRGB() {
	values := []string{"#ff5733", "#00ff00", "#0000ff", "123456"}

	fmt.Println("as.HexToRGB()")
	for _, v := range values {
		r, g, b, err := as.HexToRGB(v)
		if err != nil {
			fmt.Printf("%v => error: %v\n", v, err)
		} else {
			fmt.Printf("%v => R: %d, G: %d, B: %d\n", v, r, g, b)
		}
	}

	// Output:
	// as.HexToRGB()
	// #ff5733 => R: 255, G: 87, B: 51
	// #00ff00 => R: 0, G: 255, B: 0
	// #0000ff => R: 0, G: 0, B: 255
	// 123456 => R: 18, G: 52, B: 86
}

func ExampleRGBToHex() {
	values := [][]int{{255, 87, 51}, {0, 255, 0}, {0, 0, 255}}

	fmt.Println("as.RGBToHex()")
	for _, v := range values {
		hex := as.RGBToHex(v[0], v[1], v[2])
		fmt.Printf("R: %d, G: %d, B: %d => %v\n", v[0], v[1], v[2], hex)
	}

	// Output:
	// as.RGBToHex()
	// R: 255, G: 87, B: 51 => #ff5733
	// R: 0, G: 255, B: 0 => #00ff00
	// R: 0, G: 0, B: 255 => #0000ff
}

func ExampleMapToStruct() {
	type Person struct {
		Name string
		Age  int
	}

	values := []map[string]interface{}{
		{"Name": "Alice", "Age": 30},
		{"Name": "Bob", "Age": 25},
	}

	fmt.Println("as.MapToStruct()")
	for _, v := range values {
		var person Person
		err := as.MapToStruct(v, &person)
		if err != nil {
			fmt.Printf("%v => error: %v\n", v, err)
		} else {
			fmt.Printf("%v => %+v\n", v, person)
		}
	}

	// Output:
	// as.MapToStruct()
	// map[Age:30 Name:Alice] => {Name:Alice Age:30}
	// map[Age:25 Name:Bob] => {Name:Bob Age:25}
}

func ExampleStructToMap() {
	type Person struct {
		Name string
		Age  int
	}

	values := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	fmt.Println("as.StructToMap()")
	for _, v := range values {
		m := as.StructToMap(v)
		fmt.Printf("%+v => %v\n", v, m)
	}

	// Output:
	// as.StructToMap()
	// {Name:Alice Age:30} => map[Age:30 Name:Alice]
	// {Name:Bob Age:25} => map[Age:25 Name:Bob]
}

func ExampleIsEmailValid() {
	values := []string{"test@example.com", "invalid-email", "user@domain.co", "name@domain"}

	fmt.Println("as.IsEmailValid()")
	for _, v := range values {
		valid := as.IsEmailValid(v)
		fmt.Printf("%v => %v\n", v, valid)
	}

	// Output:
	// as.IsEmailValid()
	// test@example.com => true
	// invalid-email => false
	// user@domain.co => true
	// name@domain => false
}

func ExampleIsURLValid() {
	values := []string{"http://example.com", "https://www.google.com", "invalid-url", "ftp://example.com"}

	fmt.Println("as.IsURLValid()")
	for _, v := range values {
		valid := as.IsURLValid(v)
		fmt.Printf("%v => %v\n", v, valid)
	}

	// Output:
	// as.IsURLValid()
	// http://example.com => true
	// https://www.google.com => true
	// invalid-url => false
	// ftp://example.com => false
}

func ExampleIsCreditCardValid() {
	values := []string{"4111111111111111", "5500000000000004", "1234567890123456", "378282246310005"}

	fmt.Println("as.IsCreditCardValid()")
	for _, v := range values {
		valid := as.IsCreditCardValid(v)
		fmt.Printf("%v => %v\n", v, valid)
	}

	// Output:
	// as.IsCreditCardValid()
	// 4111111111111111 => true
	// 5500000000000004 => true
	// 1234567890123456 => false
	// 378282246310005 => true
}

func ExampleToJSON() {
	values := []interface{}{"foobar", 123, true, map[string]interface{}{"key": "value"}}

	fmt.Println("as.ToJSON()")
	for _, v := range values {
		jsonStr, err := as.ToJSON(v)
		if err != nil {
			fmt.Printf("%v => error: %v\n", v, err)
		} else {
			fmt.Printf("%v => %v\n", v, jsonStr)
		}
	}

	// Output:
	// as.ToJSON()
	// foobar => "foobar"
	// 123 => 123
	// true => true
	// map[key:value] => {"key":"value"}
}

func ExampleFromJSON() {
	type Person struct {
		Name string
		Age  int
	}
	values := []string{`{"Name": "Alice", "Age": 30}`, `{"Name": "Bob", "Age": 25}`}

	fmt.Println("as.FromJSON()")
	for _, v := range values {
		var person Person
		err := as.FromJSON(v, &person)
		if err != nil {
			fmt.Printf("%v => error: %v\n", v, err)
		} else {
			fmt.Printf("%v => %+v\n", v, person)
		}
	}

	// Output:
	// as.FromJSON()
	// {"Name": "Alice", "Age": 30} => {Name:Alice Age:30}
	// {"Name": "Bob", "Age": 25} => {Name:Bob Age:25}
}
