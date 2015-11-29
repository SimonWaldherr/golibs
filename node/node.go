package node

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
