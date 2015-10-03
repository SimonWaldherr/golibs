package foreach

import (
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

/*
const xmlString = `
<?xml version="1.0" encoding="UTF-8"?>
<root>
  <group>
    <type>group</type>
    <value>
      <group>Lorem</group>
      <group>Ipsum</group>
      <group>dolor</group>
      <group>sit</group>
      <group>
        <group>A</group>
        <group>m</group>
        <group>e</group>
        <group>t</group>
      </group>
    </value>
  </group>
  <group>
    <type>value</type>
    <value>Hello World</value>
  </group>
  <group>
    <type>value</type>
    <value>Ipsum</value>
  </group>
</root>
`*/

func Test_JSON(t *testing.T) {
	i := 0
	err := JSON(jsonString, func(key *string, index *int, value *interface{}, depth int) {
		v := *value
		switch v.(type) {
		case string:
			if *value == "Ipsum" {
				i++
			}
		}
	})
	if i != 2 || err != nil {
		t.Fatalf("JSON Test failed: %v", err)
	}
}

/*
func Test_XML(t *testing.T) {
	i := 0
	err := XML(xmlString, func(key *string, index *int, value *interface{}, depth int) {
		v := *value
		switch v.(type) {
		case string:
  		t.Logf("%v\n", *value)
			if *value == "Ipsum" {
				i++
			}
		default:
  		t.Logf("%v\n", *value)
		}
	})
	if i != 2 || err != nil {
		t.Fatalf("XML Test failed: %v", err)
	}
}*/
