package foreach

import (
	"encoding/json"
	"os"
	"simonwaldherr.de/go/golibs/file"
	//"encoding/xml"
)

func File(dirname string, recursive bool, fnc func(string, string, string, bool, os.FileInfo)) error {
	return file.Each(dirname, recursive, fnc)
}

func JSON(str string, handler func(*string, *int, *interface{}, int)) error {
	var j interface{}
	err := json.Unmarshal([]byte(str), &j)
	if err == nil {
		Node(&j, handler)
		return nil
	}
	return err
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

func Node(in *interface{}, handler func(*string, *int, *interface{}, int)) {
	nodeHelper(in, handler, 0)
}

func nodeHelper(node *interface{}, handler func(*string, *int, *interface{}, int), depth int) {
	if node == nil {
		return
	}
	o, isObject := (*node).(map[string]interface{})
	if isObject {
		for k, v := range o {
			handler(&k, nil, &v, depth)
			nodeHelper(&v, handler, depth+1)
		}
	}
	a, isArray := (*node).([]interface{})
	if isArray {
		for i, x := range a {
			handler(nil, &i, &x, depth)
			nodeHelper(&x, handler, depth+1)
		}
	}
}
