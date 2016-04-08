package structs

import (
	"reflect"
)

func Reflect(sstruct interface{}) map[string]interface{} {
	attrs := make(map[string]interface{})
	v := reflect.ValueOf(sstruct)
	t := reflect.TypeOf(sstruct)

	attrs = ReflectHelper(v, t, 0, func(name string, vtype string, value interface{}, depth int) {

	})

	return attrs
}

func ReflectHelper(v reflect.Value, t reflect.Type, depth int, handler func(string, string, interface{}, int)) map[string]interface{} {
	attrs := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		e := v.Field(i)
		f := t.Field(i)
		handler(f.Name, e.Type().String(), e.Interface(), depth)
		if e.Kind().String() == "struct" {
			attrs[f.Name] = ReflectHelper(e, e.Type(), depth+1, handler)
		} else {
			attrs[f.Name] = e.Type().String()
		}
	}
	return attrs
}
