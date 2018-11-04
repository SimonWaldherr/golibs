package csv

import (
	"sort"
	"testing"
)

var userdata string = `id;name;email
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

func Test_LoadCSVfromString(t *testing.T) {
	var str string
	var should string = "John Doe"

	csvmap, k := LoadCSVfromString(userdata)

	for _, user := range csvmap {
		if str != "" {
			str += ", "
		}
		str += user[k["name"]]
	}
	if str != should {
		t.Fatalf("user is \"%s\", should be \"%s\"\n", str, should)
	}
}

func Test_LoadCSVfromFile(t *testing.T) {
	var str string
	var should string = "John Doe"

	csvmap, k := LoadCSVfromFile("./test.csv")

	for _, user := range csvmap {
		if str != "" {
			str += ", "
		}
		str += user[k["name"]]
	}
	if str != should {
		t.Fatalf("user is \"%s\", should be \"%s\"\n", str, should)
	}
}
