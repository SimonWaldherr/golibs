package regex_test

import (
	"fmt"
	"simonwaldherr.de/go/golibs/regex"
)

func ExampleReplaceAllString() {
	str, _ := regex.ReplaceAllString("FooBaR LoReM IpSuM", "[a-z]", "")
	fmt.Print(str)

	// Output: FBR LRM ISM
}
