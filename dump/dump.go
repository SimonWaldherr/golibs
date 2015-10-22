package dump

import (
	"fmt"
	"simonwaldherr.de/go/golibs/node"
	"simonwaldherr.de/go/golibs/stack"
)

func nodeWalker(obj *interface{}) string {
	var str string
	var ldepth int
	
	array := stack.Lifo()

	//iobj =
	node.Node(obj, func(key *string, index *int, value *interface{}, depth int) {
  	if ldepth > depth {
    	for i := 1; i < ldepth; i++ {
    	  str += fmt.Sprintln("  " + array.Pop().(string))
    	}
  	} else if depth > ldepth {
    	for i := 1; i < ldepth; i++ {
    	  str += fmt.Sprintln("  " + array.Pop().(string))
    	}
  	}
		for i := 0; i < depth; i++ {
			str += fmt.Sprint("  ")
		}
		v := *value
		switch v.(type) {
		case string:
			if key != nil {
				str += fmt.Sprintf("OBJECT: key=%q, value=%#v\n", *key, *value)
			} else if index != nil {
				str += fmt.Sprintf("ARRAY: index=%d, value=%#v\n", *index, *value)
			} else {
  			str += fmt.Sprintf("VALUE: index=%d, value=%#v\n", *index, *value)
			}
		default:
			if key != nil {
				str += fmt.Sprintf("OBJECT key=%v {\n", *key)
				array.Push("}")
			} else {
				str += fmt.Sprint("[\n")
				array.Push("]")
			}
		}
		ldepth = depth
	})
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
