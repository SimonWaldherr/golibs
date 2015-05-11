package ssl

import (
	"testing"
)

func Test_Main(t *testing.T) {
	err := Check("ssl_gwv.cert", "ssl_gwv.key")
	if err != nil {
		t.Log("Couldn't find certs, trying to create new ones\n")
		options := map[string]string{}
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
	}
}