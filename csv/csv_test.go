package csv

import (
	"sort"
	"testing"
)

var userdata string = `id;name;email
0;John Doe;jDoe@example.org
1;Jane Doe;jane.doe@example.com
2;Max Mustermann;m.mustermann@alpha.tld`

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

func Test_LoadCSVfromString(t *testing.T) {
	var str string
	var should string = "Max Mustermann, John Doe, Jane Doe"

	csvmap, _ := LoadCSVfromString(userdata)

	sortedmap := sorting(csvmap)

	for _, user := range sortedmap {
		if str != "" {
			str += ", "
		}
		str += string(user.Value[1])
	}
	if str != should {
		t.Fatalf("user is \"%s\", should be \"%s\"\n", str, should)
	}
}

func Test_LoadCSVfromFile(t *testing.T) {
	var str string
	var should string = "Max Mustermann, John Doe, Jane Doe"

	csvmap, _ := LoadCSVfromFile("./test.csv")

	sortedmap := sorting(csvmap)

	for _, user := range sortedmap {
		if str != "" {
			str += ", "
		}
		str += string(user.Value[1])
	}
	if str != should {
		t.Fatalf("user is \"%s\", should be \"%s\"\n", str, should)
	}
}
