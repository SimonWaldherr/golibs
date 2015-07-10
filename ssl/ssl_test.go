package ssl

import (
	"testing"
)

func Test_Main(t *testing.T) {
	err := Check("ssl_gwv.cert", "ssl_gwv.key")
	options := map[string]string{}
	Generate(options)

	if err != nil {
		t.Log("Couldn't find certs, trying to create new ones\n")

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
		t.Fatalf("There should be an error about a missing file\n")
	}
	err = Check("x.cert", "ssl_gwv.key")
	if err == nil {
		t.Fatalf("There should be an error about a missing file\n")
	}
	options["keyPath"] = "/not/existing/path/x.key"
	err = Generate(options)
	if err == nil {
		t.Fatalf("There should be an error about a not writeable path\n", err)
	}
	options["certPath"] = "/not/existing/path/x.cert"
	err = Generate(options)
	if err == nil {
		t.Fatalf("There should be an error about a not writeable path\n", err)
	}
}
