package csv_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/csv"
	"sort"
)

var userdata string = `id;name;email
0;John Doe;jDoe@example.org
1;Jane Doe;jane.doe@example.com
2;Max Mustermann;m.mustermann@alpha.tld`

var userdata2 string = `id;name;email
0;John Doe;jDoe@example.org`

type kv struct {
	Key   int
	Value []string
}

func sorting(m map[int][]string) []kv {
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value[1] > ss[j].Value[1]
	})
	return ss
}

func Example() {
	csvmap, _ := csv.LoadCSVfromString(userdata)
	sortedmap := sorting(csvmap)

	for _, user := range sortedmap {
		fmt.Println(string(user.Value[1]))
	}

	// Output: Max Mustermann
	// John Doe
	// Jane Doe
}

func Example_second() {
	csvmap, _ := csv.LoadCSVfromFile("./test.csv")
	sortedmap := sorting(csvmap)

	for _, user := range sortedmap {
		fmt.Println(string(user.Value[1]))
	}

	// Output: Max Mustermann
	// John Doe
	// Jane Doe
}

func Example_third() {
	csvmap, k := csv.LoadCSVfromString(userdata2)
	for _, user := range csvmap {
		fmt.Println(user[k["email"]])
	}

	// Output: jDoe@example.org
}
