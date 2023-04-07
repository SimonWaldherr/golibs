// run a function for each ...
package foreach

import (
	"encoding/json"
	"os"
	"reflect"

	"simonwaldherr.de/go/golibs/file"
	"simonwaldherr.de/go/golibs/node"
	"simonwaldherr.de/go/golibs/structs"
)

// Array runs a function for each element of an array
func Array(arr interface{}, handler func(int, interface{}, int)) {
	v := reflect.ValueOf(arr)
	t := reflect.TypeOf(arr)

	structs.ReflectHelper(v, t, 0, func(name string, vtype string, value interface{}, depth int) {
		handler(depth, value, depth)
	})
}

// Dir runs a function for each file in a directory
func Dir(dirname string, recursive bool, fnc func(string, string, string, bool, os.FileInfo)) error {
	return file.Each(dirname, recursive, fnc)
}

// File runs a function for each file in a directory
func File(dirname string, recursive bool, fnc func(string, string, string, bool, os.FileInfo)) error {
	return file.Each(dirname, recursive, fnc)
}

// JSON runs a function for each element of a JSON string
func JSON(str string, handler func(*string, *int, *interface{}, int)) error {
	var j interface{}
	err := json.Unmarshal([]byte(str), &j)
	if err == nil {
		node.Node(&j, handler)
		return nil
	}
	return err
}

// Struct runs a function for each element of a struct
func Struct(sstruct interface{}, handler func(string, string, interface{}, int)) {
	v := reflect.ValueOf(sstruct)
	t := reflect.TypeOf(sstruct)

	structs.ReflectHelper(v, t, 0, handler)
}

// XML runs a function for each element of a XML string
func XML(str string, handler func(*string, *int, *interface{}, int)) error {
	var x interface{}
	err := json.Unmarshal([]byte(str), &x)
	if err == nil {
		node.Node(&x, handler)
		return nil
	}
	return err
}
