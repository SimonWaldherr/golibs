package dump

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
        ["A", "m", "e", "t", [[5,23,42]]]
      ]
    },
    {
      "value": "Hello World"
    },
    {
      "value": "foobar"
    }
  ]
`

func Test_Sprint_JSON(t *testing.T) {
	var obj interface{}
	json.Unmarshal([]byte(jsonString), &obj)
	t.Logf("%#v \n%#v", obj, &obj)
	if Sprint(obj) != `[
  {
    "type" => "group"
    "value" => [
      0 => "Lorem"
      1 => "Ipsum"
      2 => "dolor"
      3 => "sit"
      4 => [
        0 => "A"
        1 => "m"
        2 => "e"
        3 => "t"
        4 => [
          0 => [
            0 => 5
            1 => 23
            2 => 42
          ]
        ]
      ]
    ]
  }
  {
    "value" => "Hello World"
  }
  {
    "value" => "foobar"
  }
]
` {
		t.Fatalf("Dump_Sprint JSON Test failed: %v", Sprint(obj))
	}
}
