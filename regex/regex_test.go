package regex

import (
	"fmt"
	"testing"
)

func Test_ReplaceAllString(t *testing.T) {
	if ReplaceAllString("FooBaR LoReM IpSuM", "[a-z]", "") != "FBR LRM ISM" {
		t.Fatalf("ReplaceAllString failed")
	}
}

func Test_ReplaceAllString_Not(t *testing.T) {
	if ReplaceAllString("FooBaR", "[^a-z]", "") != "ooa" {
		t.Fatalf("ReplaceAllString_Not failed")
	}
}

func Test_ReplaceAllString_Match(t *testing.T) {
	if ReplaceAllString("Ipsum Lorem", "([^ ]+) ([^ ]+)", "$2 $1") != "Lorem Ipsum" {
		t.Fatalf("ReplaceAllString_Match failed")
	}
}

func Test_ReplaceAllStringFunc(t *testing.T) {
	if ReplaceAllStringFunc("abc def", "([^ ]+)", func(str string) string {
		if str == "def" {
			return "123"
		}
		return str
	}) != "abc 123" {
		t.Fatalf("ReplaceAllStringFunc failed")
	}
}

func Test_FindAllString(t *testing.T) {
	if fmt.Sprintf("%v", FindAllString("a short text", "[a-z]+")) != "[a short text]" {
		t.Fatalf("FindAllString failed")
	}
}

func Test_FindAllStringSubmatch(t *testing.T) {
	if fmt.Sprintf("%v", FindAllStringSubmatch("a short text", "[a-z]+")) != "[[a] [short] [text]]" {
		t.Fatalf("FindAllStringSubmatch failed")
	}

	if fmt.Sprintf("%v", FindAllStringSubmatch("text", "[a-z]")) != "[[t] [e] [x] [t]]" {
		t.Fatalf("FindAllStringSubmatch failed")
	}

	if fmt.Sprintf("%v", FindAllStringSubmatch("1337", "\\d")) != "[[1] [3] [3] [7]]" {
		t.Fatalf("FindAllStringSubmatch failed")
	}
}
