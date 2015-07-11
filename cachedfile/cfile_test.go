package cachedfile

import (
	"simonwaldherr.de/go/golibs/file"
	"testing"
)

func Test_Cache(t *testing.T) {
	var fn, fs, ca, original string

	fn = "./test.txt"

	Write(fn, "\nFoobar\n", true)
	Write(fn, "\nFoobar\n", true)

	fs, _ = file.Read(fn)
	ca, _ = Read(fn)
	original = fs

	if fs == ca {
		t.Fatalf("CachedFile Test 1 failed")
	}

	Stop()

	fs, _ = file.Read(fn)
	ca, _ = Read(fn)

	if fs != ca {
		t.Fatalf("CachedFile Test 2 failed")
	}

	Write(fn, original, false)
}
