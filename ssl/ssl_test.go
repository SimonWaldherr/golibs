package ssl

import (
	"testing"
)

func Test_Main(t *testing.T) {
	err := Check("ssl_gwv.cert", "ssl_gwv.key")
	options := map[string]string{}
	Generate(options)

	if err != nil {
		options["certPath"] = "x.cert"
		options["keyPath"] = "x.key"
		options["host"] = "*"
		options["countryName"] = "DE"
		options["provinceName"] = "Bavaria"
		options["organizationName"] = "Lorem Ipsum Ltd"
		options["commonName"] = "*"

		err = Generate(options)
		if err != nil {
			t.Fatalf("Couldn't create certs: %v\n", err)
		}
	} else {
		t.Fatalf("There should be an error about a missing file: %v\n", err)
	}
	err = Check("x.cert", "ssl_gwv.key")
	if err == nil {
		t.Fatalf("There should be an error about a missing file: %v\n", err)
	}
	err = Check("x.cert", "x.key")
	if err != nil {
		t.Fatalf("There should be no error: %v\n", err)
	}
	options["keyPath"] = "/not/existing/path/x.key"
	err = Generate(options)
	if err == nil {
		t.Fatalf("There should be an error about a not writeable path: %v\n", err)
	}
	options["certPath"] = "/not/existing/path/x.cert"
	err = Generate(options)
	if err == nil {
		t.Fatalf("There should be an error about a not writeable path: %v\n", err)
	}
}
