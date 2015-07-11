package cachedfile

import (
	"../file"
	"testing"
)

func Test_Cache(t *testing.T) {
	var fn, fs, ca, original string

	fn = "./test.txt"

	CachedWrite(fn, "\nFoobar\n", true)
	CachedWrite(fn, "\nFoobar\n", true)

	fs, _ = file.Read(fn)
	ca, _ = CachedRead(fn)
	original = fs

	if fs == ca {
		t.Fatalf("CachedFile Test 1 failed")
	}

	StopCache()

	fs, _ = file.Read(fn)
	ca, _ = CachedRead(fn)

	if fs != ca {
		t.Fatalf("CachedFile Test 2 failed")
	}

	CachedWrite(fn, original, false)
}
