package cachedfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"simonwaldherr.de/go/golibs/file"
)

func Test_Init(t *testing.T) {
	t.Cleanup(Stop)
	Init(15*time.Minute, 1*time.Minute)
	Init(15*time.Minute, 1*time.Minute)
}

func Test_Cache1(t *testing.T) {
	var fs, ca string
	original, err := file.Read("./test.txt")
	if err != nil {
		t.Fatal(err)
	}
	fn := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(fn, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	Init(15*time.Minute, 1*time.Minute)
	t.Cleanup(Stop)

	if err := Write(fn, "\nFoobar\n", true); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	Reset()
	if err := Write(fn, "\nFoobar\n", true); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if err := Write(fn, "\nFoobar\n", true); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	fs, _ = file.Read(fn)
	ca, _ = Read(fn)

	if fs == ca {
		t.Fatalf("CachedFile Test 1 failed")
	}

	Stop()

	fs, _ = file.Read(fn)
	ca, _ = Read(fn)

	if fs != ca {
		t.Fatalf("CachedFile Test 2 failed")
	}

	if s, _ := Size(fn); s != int64(len(original)+3*len("\nFoobar\n")) {
		t.Fatalf("CachedFile Test 3 failed")
	}
	if Clean(fn) != nil {
		t.Fatalf("CachedFile Test 4 failed")
	}
	if s, _ := Size(fn); s != 0 {
		t.Fatalf("CachedFile Test 5 failed")
	}
	if err := Write(fn, original, false); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	Stop()
}

func Test_Cache2(t *testing.T) {
	fn := "./test.txt"

	str, _ := Read(fn)
	if str != "Lorem ipsum dolor sit amet, consectetur adipisici elit, sed eiusmod tempor incidunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquid ex ea commodi consequat. Quis aute iure reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint obcaecat cupiditat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." {
		t.Fatalf("CachedFile Test 6 failed")
	}
}

func Test_Cache3(t *testing.T) {
	var err error
	var fn string
	fn = ""

	_, err = Read(fn)
	if err == nil {
		t.Fatalf("CachedFile Test 7 failed")
	}

	err = Write(fn, fn, false)
	if err == nil {
		t.Fatalf("CachedFile Test 8 failed")
	}
}

func Test_Cache4(t *testing.T) {
	var err error
	var fn string
	fn = ""

	err = Write(fn, fn, true)
	if err == nil {
		t.Fatalf("CachedFile Test 9 failed")
	}
}

func Test_Cache5(t *testing.T) {
	var err error
	var fn string

	Reset()
	fn = "😃"

	_, err = Size(fn)
	if err == nil {
		t.Fatalf("CachedFile Test 10 failed")
	}
}
