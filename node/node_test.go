package node

import (
	"encoding/json"
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

func Test_Node(t *testing.T) {
	var i int = 0
	var obj interface{}
	json.Unmarshal([]byte(jsonString), &obj)
	Node(&obj, func(key *string, index *int, value *interface{}, depth int) {
		v := *value
		switch v.(type) {
		case string:
			if *value == "Ipsum" {
				i++
			}
		}
	})
	if i != 2 {
		t.Fatalf("Node Test failed")
	}
}
