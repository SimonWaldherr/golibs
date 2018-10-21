package csv

import (
	"testing"
)

var userdata string = `id;name;email
0;John Doe;jDoe@example.org
1;Jane Doe;jane.doe@example.com
2;Max Mustermann;m.mustermann@alpha.tld`

func Test_LoadCSVfromString(t *testing.T) {
	var str string
	var should string = "John Doe, Jane Doe, Max Mustermann"

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
