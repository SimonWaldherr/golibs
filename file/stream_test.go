package file

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func fixture(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDelimitedReaders(t *testing.T) {
	for _, read := range []struct {
		name string
		fn   func(string, []string) (string, string, int, error)
	}{{"ReadUntil", ReadUntil}, {"ReadBytes", ReadBytes}} {
		t.Run(read.name, func(t *testing.T) {
			for _, tc := range []struct {
				input, text, delimiter string
				delimiters             []string
				position               int
			}{
				{"Hello\nworld", "Hello", "\n", []string{"\r", "\n"}, 5},
				{"漢字🙂stop", "漢字", "🙂", []string{"🙂"}, 2},
				{"\n", "", "\n", []string{"\n"}, 0},
				{"", "", "", []string{"\n"}, 0},
				{"end", "end", "", nil, 3},
				{"x\xffz", "x�z", "", []string{""}, 3},
				{"abc", "abc", "", []string{"ab"}, 3},
			} {
				text, delimiter, position, err := read.fn(fixture(t, tc.input), tc.delimiters)
				if err != nil || text != tc.text || delimiter != tc.delimiter || position != tc.position {
					t.Errorf("input %q: (%q, %q, %d, %v), want (%q, %q, %d, nil)", tc.input,
						text, delimiter, position, err, tc.text, tc.delimiter, tc.position)
				}
			}
			if _, _, _, err := read.fn(filepath.Join(t.TempDir(), "missing"), nil); err == nil {
				t.Error("missing input accepted")
			}
			if _, _, _, err := read.fn(t.TempDir(), nil); err == nil {
				t.Error("directory read failure was swallowed")
			}
		})
	}
}

func TestReadBlocksBoundariesAndErrors(t *testing.T) {
	var blocks []string
	output, err := ReadBlocks(fixture(t, "a,b,,tail"), []string{","}, func(block string) (string, error) {
		blocks = append(blocks, block)
		return "[" + block + "]", nil
	})
	if err != nil || output != "[a][,b][,]" || !reflect.DeepEqual(blocks, []string{"a", ",b", ","}) {
		t.Fatalf("ReadBlocks: %q, blocks %q, error %v", output, blocks, err)
	}
	// Preserve the legacy contract: only delimiter-terminated blocks are emitted,
	// and each delimiter starts the following block.
	for _, input := range []string{"", "unterminated"} {
		output, err := ReadBlocks(fixture(t, input), []string{","}, func(string) (string, error) {
			t.Fatal("callback invoked for an unterminated block")
			return "", nil
		})
		if err != nil || output != "" {
			t.Fatalf("unterminated %q: %q, %v", input, output, err)
		}
	}
	wantErr := errors.New("callback failed")
	if output, err := ReadBlocks(fixture(t, "a,b"), []string{","}, func(string) (string, error) {
		return "ignored", wantErr
	}); output != "" || !errors.Is(err, wantErr) {
		t.Fatalf("callback failure: %q, %v", output, err)
	}
	if _, err := ReadBlocks(t.TempDir(), nil, func(string) (string, error) { return "", nil }); err == nil {
		t.Fatal("directory read failure was swallowed")
	}
}

func TestReadLinesReportsScannerErrors(t *testing.T) {
	lines, err := ReadLines(fixture(t, "ok\n"+strings.Repeat("x", bufio.MaxScanTokenSize+1)))
	if !errors.Is(err, bufio.ErrTooLong) || !reflect.DeepEqual(lines, []string{"ok"}) {
		t.Fatalf("oversized line: %q, %v", lines, err)
	}
	if _, err := ReadLines(t.TempDir()); err == nil {
		t.Fatal("directory read failure was swallowed")
	}
}

func TestCopyRejectsSameFileBeforeTruncation(t *testing.T) {
	source := fixture(t, "keep this content")
	for name, destination := range map[string]string{
		"same path": source,
		"hard link": filepath.Join(t.TempDir(), "hardlink"),
		"symlink":   filepath.Join(t.TempDir(), "symlink"),
	} {
		t.Run(name, func(t *testing.T) {
			if name == "hard link" {
				if err := os.Link(source, destination); err != nil {
					t.Skipf("hard links unavailable: %v", err)
				}
			} else if name == "symlink" {
				if err := os.Symlink(source, destination); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
			}
			if err := Copy(source, destination); err == nil {
				t.Fatal("copying onto the same file succeeded")
			}
			data, err := os.ReadFile(source)
			if err != nil || string(data) != "keep this content" {
				t.Fatalf("source was changed: %q, %v", data, err)
			}
		})
	}
}

func TestCopyReplacesLongerDestination(t *testing.T) {
	source := fixture(t, "short")
	destination := fixture(t, "a much longer original destination")
	if err := Copy(source, destination); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(destination)
	if err != nil || string(data) != "short" {
		t.Fatalf("copy replacement: %q, %v", data, err)
	}
}

func TestCopyDirectoryLeavesDestinationUnchanged(t *testing.T) {
	destination := fixture(t, "keep destination")
	if err := Copy(t.TempDir(), destination); err == nil {
		t.Fatal("directory source accepted")
	}
	data, err := os.ReadFile(destination)
	if err != nil || string(data) != "keep destination" {
		t.Fatalf("invalid source changed destination: %q, %v", data, err)
	}
}

func TestWriteLinesPreservesConcatenationAndAppend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "output.txt")
	if err := WriteLines(path, []string{"hello", "", " world\n"}, false); err != nil {
		t.Fatal(err)
	}
	if err := AppendLines(path, []string{"next", " line"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "hello world\nnext line" {
		t.Fatalf("line write: %q, %v", data, err)
	}
	if err := WriteLines(t.TempDir(), []string{"x"}, false); err == nil {
		t.Fatal("writing to a directory succeeded")
	}
}

func TestWriteErrorsAreReturned(t *testing.T) {
	// On systems exposing /dev/full, real write and buffered flush failures
	// must reach the caller instead of being mistaken for successful output.
	if _, err := os.Stat("/dev/full"); err != nil {
		t.Skip("no full-device fixture on this platform")
	}
	if err := Write("/dev/full", "x", false); err == nil {
		t.Fatal("write failure was ignored")
	}
	for _, data := range []string{"x", strings.Repeat("x", 8192)} {
		if err := WriteLines("/dev/full", []string{data}, false); err == nil {
			t.Fatal("buffered write/flush failure was ignored")
		}
	}
}

func TestEachPathsAndDirectorySymlinks(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "nested"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "nested", "entry.txt"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "nested", "back")
	if err := os.Symlink(root, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	seen := make(map[string]bool)
	err := Each(root, true, func(name, extension, path string, directory bool, info os.FileInfo) {
		if !filepath.IsAbs(path) || filepath.Base(path) != name {
			t.Errorf("invalid callback path: %q, name %q", path, name)
		}
		if seen[path] {
			t.Errorf("duplicate traversal: %q", path)
		}
		seen[path] = true
		if name == "entry.txt" && extension != "txt" {
			t.Errorf("extension = %q", extension)
		}
		if name == "back" && (directory || info.Mode()&os.ModeSymlink == 0) {
			t.Error("directory symlink was followed")
		}
	})
	if err != nil || len(seen) != 3 {
		t.Fatalf("Each: %d entries, %v", len(seen), err)
	}
	if err := Each(filepath.Join(root, "nested", "entry.txt"), false,
		func(string, string, string, bool, os.FileInfo) { t.Fatal("iterated a regular file") }); err == nil {
		t.Fatal("non-directory input accepted")
	}
}

func TestEachPreservesSymlinkParentSemantics(t *testing.T) {
	directory := t.TempDir()
	lexical := filepath.Join(directory, "lexical")
	actual := filepath.Join(directory, "actual")
	child := filepath.Join(actual, "child")
	for _, path := range []string{lexical, child} {
		if err := os.MkdirAll(path, 0700); err != nil {
			t.Fatal(err)
		}
	}
	link := filepath.Join(lexical, "link")
	if err := os.Symlink(child, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := os.WriteFile(filepath.Join(child, "entry.txt"), []byte("actual content"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(lexical, "unrelated.txt"), []byte("lexical content"), 0600); err != nil {
		t.Fatal(err)
	}
	// Join would clean away link/.. and incorrectly select lexical.
	path := link + string(os.PathSeparator) + ".."
	seen := make(map[string]bool)
	err := Each(path, true, func(name, extension, path string, directory bool, info os.FileInfo) {
		seen[name] = true
		if directory {
			return
		}
		content, err := os.ReadFile(path)
		if err != nil || string(content) != "actual content" {
			t.Errorf("callback path selects the wrong file: %q, %q, %v", path, content, err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(seen) != 2 || !seen["child"] || !seen["entry.txt"] {
		t.Fatalf("traversed the wrong directory: %v", seen)
	}
}
