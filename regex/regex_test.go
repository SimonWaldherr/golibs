package regex

import (
	"fmt"
	"testing"
)

func Test_ReplaceAllString(t *testing.T) {
	ReplaceAllString("FooBaR LoReM IpSuM", "[a-z]", "")
	if str, _ := ReplaceAllString("FooBaR LoReM IpSuM", "[a-z]", ""); str != "FBR LRM ISM" {
		t.Fatalf("ReplaceAllString failed")
	}
}

func Test_ReplaceAllString_Not(t *testing.T) {
	if str, _ := ReplaceAllString("FooBaR", "[^a-z]", ""); str != "ooa" {
		t.Fatalf("ReplaceAllString_Not failed")
	}
}

func Test_ReplaceAllString_Match(t *testing.T) {
	if str, _ := ReplaceAllString("Ipsum Lorem", "([^ ]+) ([^ ]+)", "$2 $1"); str != "Lorem Ipsum" {
		t.Fatalf("ReplaceAllString_Match failed")
	}
}

func Test_ReplaceAllStringFunc(t *testing.T) {
	if str, _ := ReplaceAllStringFunc("abc def", "([^ ]+)", func(str string) string {
		if str == "def" {
			return "123"
		}
		return str
	}); str != "abc 123" {
		t.Fatalf("ReplaceAllStringFunc failed")
	}
	if str, _ := ReplaceAllStringFunc("abc def", "([^ ]+)", func(str string) string {
		if str == "def" {
			return "123"
		}
		return str
	}); str != "abc 123" {
		t.Fatalf("ReplaceAllStringFunc failed")
	}
}

func Test_FindAllString(t *testing.T) {
	FindAllString("a short text", "[a-z]+")
	ret, _ := FindAllString("a short text", "[a-z]+")
	if fmt.Sprintf("%v", ret) != "[a short text]" {
		t.Fatalf("FindAllString failed")
	}
}

func Test_FindAllStringSubmatch(t *testing.T) {
	ret, _ := FindAllStringSubmatch("a short text", "[a-z]+")
	if fmt.Sprintf("%v", ret) != "[[a] [short] [text]]" {
		t.Fatalf("FindAllStringSubmatch failed")
	}

	ret, _ = FindAllStringSubmatch("text", "[a-z]")
	if fmt.Sprintf("%v", ret) != "[[t] [e] [x] [t]]" {
		t.Fatalf("FindAllStringSubmatch failed")
	}

	ret, _ = FindAllStringSubmatch("1337", "\\d")
	if fmt.Sprintf("%v", ret) != "[[1] [3] [3] [7]]" {
		t.Fatalf("FindAllStringSubmatch failed")
	}
}

func Test_Count(t *testing.T) {
	if Count() != 6 {
		t.Fatalf("Count Test failed")
	}
	Cleanup()
	if Count() != 0 {
		t.Fatalf("Cleanup Count Test failed")
	}
}

func Test_CheckRegex(t *testing.T) {
	if CheckRegex(".+\\d?") != nil {
		t.Fatalf("CheckRegex Test failed")
	}
	if CheckRegex(".+(d?") == nil {
		t.Fatalf("CheckRegex Test failed")
	}
}

func Test_CacheRegex(t *testing.T) {
	if CacheRegex(".+\\d?") != nil {
		t.Fatalf("CacheRegex Test failed")
	}
	if CacheRegex(".+\\d?") != nil {
		t.Fatalf("CacheRegex Test failed")
	}
}

func Test_MatchString(t *testing.T) {
	if match, _ := MatchString("Fooobar", "o{3}"); match != true {
		t.Fatalf("MatchString Test failed")
	}
	if match, _ := MatchString("Fooobar", "o{3}"); match != true {
		t.Fatalf("MatchString Test failed")
	}
}
