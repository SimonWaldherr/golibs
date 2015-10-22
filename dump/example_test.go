package dump_test

import (
	"encoding/json"
	"simonwaldherr.de/go/golibs/dump"
)

func ExampleSprint_array() {
	var obj interface{}
	json.Unmarshal([]byte("[[[1,2,3]]]"), &obj)
	dump.Print(obj)

	// Output:
	// [
	//   0 => [
	//     0 => [
	//       0 => 1
	//       1 => 2
	//       2 => 3
	//     ]
	//   ]
	// ]
}

/*
func ExampleSprint_object() {
	var obj interface{}
	json.Unmarshal([]byte(`{"a": {"b": {"c": "d"}, "e": "f"}}`), &obj)
	dump.Print(obj)

	// Output:
	// {
	//   "a" => {
	//     "b" => {
	//       "c" => "d"
	//     }
	//   }
	//   "e" => "f"
	// }
}
*/