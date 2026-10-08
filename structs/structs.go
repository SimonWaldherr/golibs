package structs

import (
	"reflect"
)

// Package structs contains small helpers to inspect, validate and
// modify Go struct values using reflection and struct tag rules.

// Reflect returns a map describing the fields of the provided struct
// value. The returned map's keys are the field names and values are
// either the field's type string or a nested map for nested structs.
// Reflect expects a concrete struct value (not a pointer).
func Reflect(sstruct interface{}) map[string]interface{} {
	v := reflect.ValueOf(sstruct)
	t := reflect.TypeOf(sstruct)
	return ReflectHelper(v, t, 0, func(name string, vtype string, value interface{}, depth int) {})
}

// ReflectHelper is a lower-level helper used by Reflect. It walks the
// provided value `v` using type `t` and returns a description map. The
// `handler` callback is invoked for every field and can be used to
// collect additional metadata.
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
