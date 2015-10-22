package dump

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
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
        ["A", "m", "e", "t", [["golibs"]]]
      ]
    },
    {
      "type": "value",
      "value": "Hello World"
    },
    {
      "type": "value",
      "value": "foobar"
    }
  ]
`

func Test_Sprint(t *testing.T) {
	var obj interface{}
	json.Unmarshal([]byte(jsonString), &obj)
	fmt.Println(Sprint(obj))
	spew.Dump(obj)
}
