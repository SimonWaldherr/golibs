package csv

import (
	"testing"
)

var userdata string = `id;name;email
0;John Doe;jDoe@example.org`

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
