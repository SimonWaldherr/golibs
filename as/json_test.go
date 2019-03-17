package as

import (
	"encoding/json"
	"testing"
)

type intTest struct {
	I Int64 `json:"x"`
}

var testintdata = []struct {
	str string
	i   int64
}{
	{"{\"x\":1}", 1},
	{"{\"x\":\"05\"}", 5},
	{"{\"x\":\"42\"}", 42},
	{"{\"x\":\"42.03\"}", 42},
}

func TestJSONInt(t *testing.T) {
	for _, x := range testintdata {
		ts := &intTest{}
		if err := json.Unmarshal([]byte(x.str), ts); err != nil {
			t.Fatal(err)
		}
		if int64(ts.I) != x.i {
			t.Errorf("Unexpected result. Int not %v: %#v", x.i, *ts)
		}
	}
}

type dynTest struct {
	D Dynamic `json:"x"`
}

var testdyndata = []struct {
	str   string
	dtype string
}{
	{"{\"x\":\"2018-06-13\"}", "date"},
}

func TestJSONDynamic(t *testing.T) {
	for _, x := range testdyndata {
		ts := &dynTest{}
		if err := json.Unmarshal([]byte(x.str), ts); err != nil {
			t.Fatal(err)
		}
		if string(ts.D.Type) != x.dtype {
			t.Errorf("Unexpected result. Type not %v: %#v", x.dtype, *ts)
		}
	}
}
