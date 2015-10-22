package pseudodb

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const jsonString = `
  [
    {
      "type": "group",
      "value": [
        "Lorem",
        "Ipsum",
        "dolor",
        "sit",
        ["A", "m", "e", "t"]
      ]
    },
    {
      "type": "value",
      "value": "Hello World"
    },
    {
      "type": "value",
      "value": "Ipsum"
    }
  ]
`

func Test_DB(t *testing.T) {
	var s string
	var j interface{}
	err := json.Unmarshal([]byte(jsonString), &j)

	if err != nil {
		t.Fatalf("JSON Test failed: %v", err)
	}

	db := New()
	db.Insert(j)

	db.Each(func(i *int, obj *interface{}) {
		s = fmt.Sprintf("%+v", *obj)
	})

	if !strings.Contains(s, "[Lorem Ipsum dolor sit [A m e t]]") {
		t.Fatalf("DB Test 1 failed")
	}
}
