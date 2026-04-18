package env_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"simonwaldherr.de/go/golibs/env"
)

func ExampleString() {
	os.Setenv("EXAMPLE_HOST", "localhost")
	fmt.Println(env.String("EXAMPLE_HOST", "127.0.0.1"))
	fmt.Println(env.String("EXAMPLE_MISSING", "127.0.0.1"))
	// Output:
	// localhost
	// 127.0.0.1
}

func ExampleInt() {
	os.Setenv("EXAMPLE_PORT", "8080")
	fmt.Println(env.Int("EXAMPLE_PORT", 80))
	fmt.Println(env.Int("EXAMPLE_MISSING_PORT", 80))
	// Output:
	// 8080
	// 80
}

func ExampleBool() {
	os.Setenv("EXAMPLE_DEBUG", "true")
	fmt.Println(env.Bool("EXAMPLE_DEBUG", false))
	fmt.Println(env.Bool("EXAMPLE_MISSING_BOOL", false))
	// Output:
	// true
	// false
}

func ExampleDuration() {
	os.Setenv("EXAMPLE_TIMEOUT", "30s")
	fmt.Println(env.Duration("EXAMPLE_TIMEOUT", 5*time.Second))
	fmt.Println(env.Duration("EXAMPLE_MISSING_TIMEOUT", 5*time.Second))
	// Output:
	// 30s
	// 5s
}

func ExampleSlice() {
	os.Setenv("EXAMPLE_HOSTS", "host1, host2, host3")
	fmt.Println(env.Slice("EXAMPLE_HOSTS", ",", nil))
	fmt.Println(env.Slice("EXAMPLE_MISSING_HOSTS", ",", []string{"default"}))
	// Output:
	// [host1 host2 host3]
	// [default]
}

func ExampleIsSet() {
	os.Setenv("EXAMPLE_SET", "value")
	os.Unsetenv("EXAMPLE_UNSET")
	fmt.Println(env.IsSet("EXAMPLE_SET"))
	fmt.Println(env.IsSet("EXAMPLE_UNSET"))
	// Output:
	// true
	// false
}

func ExampleFloat64() {
	os.Setenv("EXAMPLE_RATIO", "3.14")
	fmt.Printf("%.2f\n", env.Float64("EXAMPLE_RATIO", 1.0))
	fmt.Printf("%.2f\n", env.Float64("EXAMPLE_MISSING_RATIO", 1.0))
	// Output:
	// 3.14
	// 1.00
}

func TestString(t *testing.T) {
	os.Setenv("TEST_STR", "hello")
	defer os.Unsetenv("TEST_STR")
	if got := env.String("TEST_STR", "world"); got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
	if got := env.String("TEST_STR_MISSING", "world"); got != "world" {
		t.Errorf("expected %q, got %q", "world", got)
	}
}

func TestInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")
	if got := env.Int("TEST_INT", 0); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
	if got := env.Int("TEST_INT_INVALID", 7); got != 7 {
		t.Errorf("expected 7, got %d", got)
	}
}

func TestInt64(t *testing.T) {
	os.Setenv("TEST_INT64", "9999999999")
	defer os.Unsetenv("TEST_INT64")
	if got := env.Int64("TEST_INT64", 0); got != 9999999999 {
		t.Errorf("expected 9999999999, got %d", got)
	}
}

func TestFloat64(t *testing.T) {
	os.Setenv("TEST_F64", "2.718")
	defer os.Unsetenv("TEST_F64")
	got := env.Float64("TEST_F64", 0)
	if got < 2.717 || got > 2.719 {
		t.Errorf("expected ~2.718, got %f", got)
	}
}

func TestBool(t *testing.T) {
	cases := []struct {
		val  string
		want bool
	}{
		{"true", true}, {"1", true}, {"yes", true}, {"on", true},
		{"false", false}, {"0", false}, {"no", false}, {"off", false},
	}
	for _, c := range cases {
		os.Setenv("TEST_BOOL", c.val)
		if got := env.Bool("TEST_BOOL", !c.want); got != c.want {
			t.Errorf("Bool(%q): expected %v, got %v", c.val, c.want, got)
		}
	}
	os.Unsetenv("TEST_BOOL")
}

func TestDuration(t *testing.T) {
	os.Setenv("TEST_DUR", "2m30s")
	defer os.Unsetenv("TEST_DUR")
	want := 2*time.Minute + 30*time.Second
	if got := env.Duration("TEST_DUR", 0); got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestSlice(t *testing.T) {
	os.Setenv("TEST_SLICE", "a,b,c")
	defer os.Unsetenv("TEST_SLICE")
	got := env.Slice("TEST_SLICE", ",", nil)
	if len(got) != 3 || got[0] != "a" || got[2] != "c" {
		t.Errorf("unexpected slice: %v", got)
	}
}

func TestMustString(t *testing.T) {
	os.Setenv("TEST_MUST", "present")
	defer os.Unsetenv("TEST_MUST")
	if got := env.MustString("TEST_MUST"); got != "present" {
		t.Errorf("expected %q, got %q", "present", got)
	}
}

func TestMustString_Panic(t *testing.T) {
	os.Unsetenv("TEST_MUST_MISSING")
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing required env var")
		}
	}()
	env.MustString("TEST_MUST_MISSING")
}

func TestIsSet(t *testing.T) {
	os.Setenv("TEST_IS_SET", "")
	defer os.Unsetenv("TEST_IS_SET")
	if !env.IsSet("TEST_IS_SET") {
		t.Error("IsSet should return true even for empty variable")
	}
	if env.IsSet("TEST_IS_SET_MISSING") {
		t.Error("IsSet should return false for unset variable")
	}
}

func TestMap(t *testing.T) {
	os.Setenv("MYAPP_HOST", "localhost")
	os.Setenv("MYAPP_PORT", "8080")
	os.Setenv("OTHER_VAR", "ignored")
	defer func() {
		os.Unsetenv("MYAPP_HOST")
		os.Unsetenv("MYAPP_PORT")
		os.Unsetenv("OTHER_VAR")
	}()

	m := env.Map("MYAPP_")
	if m["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", m["HOST"])
	}
	if m["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", m["PORT"])
	}
	if _, ok := m["OTHER_VAR"]; ok {
		t.Error("OTHER_VAR should not be in map with prefix MYAPP_")
	}
}
