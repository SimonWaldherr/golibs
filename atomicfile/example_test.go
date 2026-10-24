package atomicfile_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"simonwaldherr.de/go/golibs/atomicfile"
)

func ExampleWriteFile() {
	directory, err := os.MkdirTemp("", "atomicfile-example-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(directory)
	path := filepath.Join(directory, "settings.json")
	if err := atomicfile.WriteFile(path, []byte(`{"enabled":true}`), 0600); err != nil {
		panic(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"enabled":true}
}

func ExampleWrite() {
	directory, err := os.MkdirTemp("", "atomicfile-example-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(directory)
	path := filepath.Join(directory, "settings.json")
	err = atomicfile.Write(path, 0600, func(writer io.Writer) error {
		return json.NewEncoder(writer).Encode(map[string]int{"count": 2})
	})
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(data))
	// Output: {"count":2}
}
