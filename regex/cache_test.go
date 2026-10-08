package regex

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestInvalidExpressionsReturnErrors(t *testing.T) {
	Cleanup()
	const invalid = "["
	checks := []struct {
		name string
		call func() error
	}{
		{"cache", func() error { return CacheRegex(invalid) }},
		{"match", func() error { _, err := MatchString("text", invalid); return err }},
		{"replace", func() error { _, err := ReplaceAllString("text", invalid, ""); return err }},
		{"replace callback", func() error {
			_, err := ReplaceAllStringFunc("text", invalid, func(string) string {
				t.Fatal("callback called for invalid expression")
				return ""
			})
			return err
		}},
		{"find", func() error { _, err := FindAllString("text", invalid); return err }},
		{"submatches", func() error { _, err := FindAllStringSubmatch("text", invalid); return err }},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			if err := check.call(); err == nil {
				t.Fatal("expected expression compilation error")
			}
		})
	}
	if Count() != 0 {
		t.Fatal("invalid expression was cached")
	}
}

func TestReplacementCallbackCanUseCache(t *testing.T) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		result, err := ReplaceAllStringFunc("foo bar", "[a-z]+", func(s string) string {
			Cleanup()
			matched, err := MatchString(s, "^[a-z]+$")
			if err != nil || !matched {
				t.Errorf("nested match: %v, %v", matched, err)
			}
			return "word"
		})
		if err != nil || result != "word word" {
			t.Errorf("replacement: %q, %v", result, err)
		}
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("replacement callback blocked on cache")
	}
}

func TestConcurrentCacheAccess(t *testing.T) {
	Cleanup()
	var wg sync.WaitGroup
	for worker := 0; worker < 16; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			pattern := fmt.Sprintf("^value-%d$", worker)
			value := fmt.Sprintf("value-%d", worker)
			for i := 0; i < 100; i++ {
				if err := CacheRegex(pattern); err != nil {
					t.Error(err)
				}
				if matched, err := MatchString(value, pattern); err != nil || !matched {
					t.Errorf("match = %v, %v", matched, err)
				}
				_ = Count()
				if i%10 == 0 {
					Cleanup()
				}
			}
		}(worker)
	}
	wg.Wait()
}
