package csv_test

import (
	"fmt"

	"simonwaldherr.de/go/golibs/csv"
)

var userdata string = `id;name;email
0;John Doe;jDoe@example.org
1;Jane Doe;jane.doe@example.com
2;Max Mustermann;m.mustermann@alpha.tld`

var userdata2 string = `id;name;email
0;John Doe;jDoe@example.org`

// helper types removed; not used in examples

func Example() {
	csvmap, k := csv.LoadCSVfromString(userdata)

	for _, user := range csvmap {
		fmt.Println(user[k["email"]])
	}
}

func Example_second() {
	csvmap, k := csv.LoadCSVfromFile("./test.csv")

	for _, user := range csvmap {
		fmt.Println(user[k["email"]])
	}
}

func Example_third() {
	csvmap, k := csv.LoadCSVfromString(userdata2)
	for _, user := range csvmap {
		fmt.Println(user[k["email"]])
	}

	// Output: jDoe@example.org
}

func Example_fourth() {
	csvmap, k := csv.LoadCSVfromFile("./test.csv")
	for _, user := range csvmap {
		fmt.Println(user[k["email"]])
	}

	// Output: jDoe@example.org
}
