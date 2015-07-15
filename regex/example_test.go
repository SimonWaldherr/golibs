package regex_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/regex"
)

func ExampleReplaceAllString() {
	fmt.Print(regex.ReplaceAllString("FooBaR LoReM IpSuM", "[a-z]", ""))

	// Output: FBR LRM ISM
}
