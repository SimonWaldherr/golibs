package ssl_test

import (
	"simonwaldherr.de/go/golibs/ssl"
)

func ExampleGenerate() {
	options := map[string]string{}
	options["certPath"] = "x.cert"
	options["keyPath"] = "x.key"
	options["host"] = "*"
	options["countryName"] = "DE"
	options["provinceName"] = "Bavaria"
	options["organizationName"] = "Lorem Ipsum Ltd"
	options["commonName"] = "*"
	ssl.Generate(options)

	// Output:
}
