package dump

import (
	"fmt"
	"simonwaldherr.de/go/golibs/node"
	"simonwaldherr.de/go/golibs/stack"
)

func nodeWalker(obj *interface{}) string {
	var str string
	var ldepth int
	var ltype string

	array := stack.Lifo()

	node.Node(obj, func(key *string, index *int, value *interface{}, depth int) {
		if ldepth > depth {
			for i := 1; i <= ldepth; i++ {
				for j := 0; j < ldepth+1-i; j++ {
					str += fmt.Sprint("  ")
				}
				str += fmt.Sprintln(array.Pop().(string))
			}
		} else if depth > ldepth {
			if key != nil {
				if ltype == "[" {
					str += fmt.Sprintf("%q {\n", *key)
				} else {
					str += fmt.Sprintf("{\n")
				}
				array.Push("}")
			} else if index != nil {
				str += fmt.Sprintf("[\n")
			}
		}
		for i := 0; i <= depth; i++ {
			str += fmt.Sprint("  ")
		}

		v := *value
		switch v.(type) {
		case map[string]interface{}:
			ltype = "{"
		case []interface{}:
			if key != nil {
				str += fmt.Sprintf("%q => ", *key)
				array.Push("]")
			} else {
				str += fmt.Sprintf("%v => ", *index)
				array.Push("]")
			}
			ltype = "["
		case string:
			if key != nil {
				str += fmt.Sprintf("%q => %#v\n", *key, *value)
			} else {
				str += fmt.Sprintf("%d => %#v\n", *index, *value)
			}
		default:
			if key != nil {
				str += fmt.Sprintf("%q => ", *key)
			} else if value != nil {
				str += fmt.Sprintf("%d => %v\n", *index, *value)
			}
		}
		ldepth = depth
	})

	for array.Len() > 0 {
		for i := 0; i < array.Len(); i++ {
			str += "  "
		}
		str += array.Pop().(string) + "\n"
	}

	o := *obj

	switch o.(type) {
	case []interface{}:
		str = "[\n" + str + "]\n"
	case map[string]interface{}:
		str = "{\n" + str + "}\n"
	}
	return str
}

func Sprint(obj interface{}) string {
	return nodeWalker(&obj)
}

func Print(obj interface{}) {
	fmt.Print(Sprint(obj))
}

func Println(obj interface{}) {
	fmt.Println(Sprint(obj))
}
