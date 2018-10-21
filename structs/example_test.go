package structs_test

import (
	"fmt"
	"reflect"
	"simonwaldherr.de/go/golibs/structs"
	"sort"
)

type LIPS struct {
	VBELN string
	POSNR string
	MATNR string
	MATKL string
	ARKTX string
	EANNR string
	LGORT string
	LFIMG string
	VRKME string
	VKBUR string
}

func ExampleReflect() {
	var lips LIPS
	x := structs.Reflect(lips)

	keys := []string{}
	for k, _ := range x {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("key: %v, value: %v\n", k, x[k])
	}

	// Output:
	// key: ARKTX, value: string
	// key: EANNR, value: string
	// key: LFIMG, value: string
	// key: LGORT, value: string
	// key: MATKL, value: string
	// key: MATNR, value: string
	// key: POSNR, value: string
	// key: VBELN, value: string
	// key: VKBUR, value: string
	// key: VRKME, value: string
}

func ExampleReflectHelper() {
	var lips LIPS
	v := reflect.ValueOf(lips)
	t := reflect.TypeOf(lips)

	structs.ReflectHelper(v, t, 0, func(name string, vtype string, value interface{}, depth int) {
		fmt.Printf("%v - %v - %v - %v\n", name, vtype, value, depth)
	})

	// Output:
	// VBELN - string -  - 0
	// POSNR - string -  - 0
	// MATNR - string -  - 0
	// MATKL - string -  - 0
	// ARKTX - string -  - 0
	// EANNR - string -  - 0
	// LGORT - string -  - 0
	// LFIMG - string -  - 0
	// VRKME - string -  - 0
	// VKBUR - string -  - 0
}
