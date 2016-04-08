// run a function for each ...
package foreach

import (
	"encoding/json"
	"os"
	"reflect"
	"simonwaldherr.de/go/golibs/file"
	"simonwaldherr.de/go/golibs/node"
	"simonwaldherr.de/go/golibs/structs"
	//"encoding/xml"
)

func File(dirname string, recursive bool, fnc func(string, string, string, bool, os.FileInfo)) error {
	return file.Each(dirname, recursive, fnc)
}

func JSON(str string, handler func(*string, *int, *interface{}, int)) error {
	var j interface{}
	err := json.Unmarshal([]byte(str), &j)
	if err == nil {
		node.Node(&j, handler)
		return nil
	}
	return err
}

func Struct(sstruct interface{}, handler func(string, string, interface{}, int)) {
	v := reflect.ValueOf(sstruct)
	t := reflect.TypeOf(sstruct)

	structs.ReflectHelper(v, t, 0, handler)
}

/*
func XML(str string, handler func(*string, *int, *interface{}, int)) error {
	var x interface{}
	err := xml.Unmarshal([]byte(str), &x)
	if err == nil {
		Node(&x, handler)
		return nil
	}
	return err
}*/
