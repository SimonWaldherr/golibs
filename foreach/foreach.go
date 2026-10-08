// Package foreach provides helpers to iterate over different
// container-like values (arrays, JSON, XML, directories, structs)
// and invoke a user-supplied handler for each element.
package foreach

import (
	"encoding/json"
	"os"
	"reflect"

	"simonwaldherr.de/go/golibs/file"
	"simonwaldherr.de/go/golibs/node"
	"simonwaldherr.de/go/golibs/structs"
)

// Array runs handler for every element in `arr`. The handler
// signature is `func(index int, value interface{}, depth int)` and the
// function will be called for each discovered field. This function uses
// reflection and the package's `ReflectHelper` under the hood.
func Array(arr interface{}, handler func(int, interface{}, int)) {
	v := reflect.ValueOf(arr)
	t := reflect.TypeOf(arr)

	structs.ReflectHelper(v, t, 0, func(name string, vtype string, value interface{}, depth int) {
		handler(depth, value, depth)
	})
}

// Dir calls `fnc` for each file found in `dirname`. If `recursive` is
// true sub-directories are traversed.
func Dir(dirname string, recursive bool, fnc func(string, string, string, bool, os.FileInfo)) error {
	return file.Each(dirname, recursive, fnc)
}

// File is an alias of Dir kept for backward compatibility.
func File(dirname string, recursive bool, fnc func(string, string, string, bool, os.FileInfo)) error {
	return file.Each(dirname, recursive, fnc)
}

// JSON parses `str` and calls `handler` for each element. The handler
// has signature `func(*string, *int, *interface{}, int)` and is used by
// `node.Node` to walk generic JSON structures.
func JSON(str string, handler func(*string, *int, *interface{}, int)) error {
	var j interface{}
	err := json.Unmarshal([]byte(str), &j)
	if err == nil {
		node.Node(&j, handler)
		return nil
	}
	return err
}

// Struct walks the provided struct value and calls `handler` for each
// field using the same callback signature as ReflectHelper.
func Struct(sstruct interface{}, handler func(string, string, interface{}, int)) {
	v := reflect.ValueOf(sstruct)
	t := reflect.TypeOf(sstruct)

	structs.ReflectHelper(v, t, 0, handler)
}

// XML parses `str` (XML already converted to a generic structure) and
// invokes `handler` for each element. XML input is parsed by the same
// path as JSON (uses `node.Node`) so callers must ensure the content
// is converted to a compatible intermediate representation.
func XML(str string, handler func(*string, *int, *interface{}, int)) error {
	var x interface{}
	err := json.Unmarshal([]byte(str), &x)
	if err == nil {
		node.Node(&x, handler)
		return nil
	}
	return err
}
