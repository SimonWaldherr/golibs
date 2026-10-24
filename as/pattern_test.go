package as

import (
	"sync"
	"testing"
)

func TestTypePatternPriority(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{" \tTRUE\r\n", "bool"},
		{"0", "bool"},
		{"14", "int"},
		{"-42", "float"},
		{"3.14159", "float"},
		{"€42", "price"},
		{"https://example.com", "url"},
		{"192.0.2.1", "ipv4"},
		{"2001:db8::1", "ipv6"},
		{"a:b:c:d:e:f", "mac"},
		{"person@example.com", "email"},
		{"MMXXIV", "date"},
		{"hello", "alpha"},
		{"ABC123", "alphanumeric"},
		{"AA==", "base64"},
		{"a\tb", "ascii"},
		{"", ""},
		{"😃", ""},
	} {
		got, err := Type(tc.input)
		if err != nil || got != tc.want {
			t.Errorf("Type(%q) = %q, %v; want %q", tc.input, got, err, tc.want)
		}
	}
}

func TestValidatorPatternBoundaries(t *testing.T) {
	for _, tc := range []struct {
		input string
		valid bool
	}{
		{"person+tag@example.com", true},
		{"a..b@example.com", true}, // Retain the existing pattern's acceptance of consecutive dots.
		{"PERSON@example.com", false},
		{"person@example.com\n", false},
		{"prefix person@example.com suffix", false},
		{"person@localhost", false},
		{"é@example.com", false},
		{"", false},
	} {
		if got := IsEmailValid(tc.input); got != tc.valid {
			t.Errorf("IsEmailValid(%q) = %v, want %v", tc.input, got, tc.valid)
		}
	}
	for _, tc := range []struct {
		input string
		valid bool
	}{
		{"example.com", true},
		{"http://example.com/docs/index.html", true},
		{"https://EXAMPLE.com", false},
		{"https://example.com:8443", false},
		{"https://example.com\n", false},
		{"ftp://example.com", false},
		{"", false},
	} {
		if got := IsURLValid(tc.input); got != tc.valid {
			t.Errorf("IsURLValid(%q) = %v, want %v", tc.input, got, tc.valid)
		}
	}
}

func TestPatternsConcurrent(t *testing.T) {
	const workers = 16
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				if got, err := Type("person@example.com"); err != nil || got != "email" {
					t.Errorf("Type = %q, %v", got, err)
				}
				if !IsEmailValid("person@example.com") || IsEmailValid("invalid") ||
					!IsURLValid("https://example.com") || IsURLValid("invalid") {
					t.Error("validator results changed during concurrent calls")
				}
			}
		}()
	}
	wg.Wait()
}
