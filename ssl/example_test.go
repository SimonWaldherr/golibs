package ssl_test

import (
	"os"
	"path/filepath"

	"simonwaldherr.de/go/golibs/ssl"
)

func ExampleGenerate() {
	directory, err := os.MkdirTemp("", "golibs-certificates-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(directory)
	options := map[string]string{
		"certPath":    filepath.Join(directory, "localhost.cert"),
		"keyPath":     filepath.Join(directory, "localhost.key"),
		"countryName": "DE", "provinceName": "Bavaria",
		"organizationName": "Example", "commonName": "localhost",
	}
	if err := ssl.Generate(options); err != nil {
		panic(err)
	}
	// Output:
}
