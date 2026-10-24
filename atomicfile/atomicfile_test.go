package atomicfile

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func assertContents(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil || string(data) != want {
		t.Fatalf("ReadFile(%s) = %q, %v; want %q", path, data, err, want)
	}
}

func assertNoTemporaryFiles(t *testing.T, directory string) {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(directory, ".atomicfile-*.tmp"))
	if err != nil || len(files) != 0 {
		t.Fatalf("temporary files remain: %v, %v", files, err)
	}
}

func TestCreateReplaceAndMode(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "settings.json")
	if err := WriteFile(path, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := WriteFile(path, []byte("new"), 0600); err != nil {
		t.Fatal(err)
	}
	assertContents(t, path, "new")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Fatalf("replacement mode = %v, want 0600", info.Mode())
	}
	assertNoTemporaryFiles(t, directory)
}

func TestCallbackFailurePreservesDestination(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "data")
	if err := os.WriteFile(path, []byte("original"), 0600); err != nil {
		t.Fatal(err)
	}
	wantErr := errors.New("encoding failed")
	err := Write(path, 0600, func(writer io.Writer) error {
		if _, err := io.WriteString(writer, "partial"); err != nil {
			t.Fatal(err)
		}
		assertContents(t, path, "original")
		if _, ownsClose := writer.(io.Closer); ownsClose {
			t.Fatal("callback received ownership of Close")
		}
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Write = %v, want callback error", err)
	}
	assertContents(t, path, "original")
	assertNoTemporaryFiles(t, directory)
}

func TestStreamingAndEmptyWrites(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "data")
	err := Write(path, 0600, func(writer io.Writer) error {
		_, err := io.Copy(writer, strings.NewReader("streamed content"))
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	assertContents(t, path, "streamed content")
	err = Write(path, 0600, func(writer io.Writer) error {
		// Hide WriterTo to exercise the destination's ReaderFrom fast path.
		reader := struct{ io.Reader }{strings.NewReader("copied content")}
		_, err := io.Copy(writer, reader)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	assertContents(t, path, "copied content")
	if err := WriteFile(path, nil, 0600); err != nil {
		t.Fatal(err)
	}
	assertContents(t, path, "")
	assertNoTemporaryFiles(t, directory)
}

func TestInvalidTargetsDoNotCallCallback(t *testing.T) {
	directory := t.TempDir()
	regular := filepath.Join(directory, "regular")
	if err := os.WriteFile(regular, []byte("untouched"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"", directory, filepath.Join(directory, "missing", "data"), filepath.Join(regular, "child")} {
		called := false
		err := Write(path, 0600, func(io.Writer) error { called = true; return nil })
		if err == nil || called {
			t.Errorf("Write(%q) = %v, called = %v", path, err, called)
		}
	}
	if err := Write(filepath.Join(directory, "data"), 0600, nil); err == nil {
		t.Fatal("nil callback was accepted")
	}
	assertNoTemporaryFiles(t, directory)
}

func TestCallbackTemporaryFileIsPrivate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not implement POSIX permission bits")
	}
	directory := t.TempDir()
	path := filepath.Join(directory, "data")
	err := Write(path, 0644, func(writer io.Writer) error {
		files, err := filepath.Glob(filepath.Join(directory, ".atomicfile-*.tmp"))
		if err != nil || len(files) != 1 {
			t.Fatalf("temporary files = %v, %v", files, err)
		}
		info, err := os.Stat(files[0])
		if err != nil || info.Mode().Perm()&0077 != 0 {
			t.Fatalf("temporary file isn't private: %v, %v", info, err)
		}
		_, err = io.WriteString(writer, "complete")
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	assertContents(t, path, "complete")
	assertNoTemporaryFiles(t, directory)
}

func TestWriterIsClosedOnCallbackError(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "data")
	var retained io.Writer
	wantErr := errors.New("callback error")
	err := Write(path, 0600, func(writer io.Writer) error {
		retained = writer
		return wantErr
	})
	if err != wantErr {
		t.Fatalf("callback error = %v", err)
	}
	if _, err := retained.Write([]byte("late write")); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("temporary writer remained open: %v", err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("failed callback created its destination: %v", err)
	}
	assertNoTemporaryFiles(t, directory)
}

func TestPanicCleansTemporaryFile(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "data")
	if err := os.WriteFile(path, []byte("original"), 0600); err != nil {
		t.Fatal(err)
	}
	func() {
		defer func() {
			if got := recover(); got != "failure" {
				t.Errorf("panic = %v", got)
			}
		}()
		_ = Write(path, 0600, func(io.Writer) error { panic("failure") })
	}()
	assertContents(t, path, "original")
	assertNoTemporaryFiles(t, directory)
}

func TestSymlinkAndHardLinkAreReplacedNotFollowed(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symbolic links may require elevated privileges")
	}
	for _, makeLink := range []struct {
		name string
		link func(string, string) error
	}{{"symbolic", os.Symlink}, {"hard", os.Link}} {
		t.Run(makeLink.name, func(t *testing.T) {
			directory := t.TempDir()
			target := filepath.Join(directory, "target")
			link := filepath.Join(directory, "link")
			if err := os.WriteFile(target, []byte("original"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := makeLink.link(target, link); err != nil {
				t.Fatal(err)
			}
			if err := WriteFile(link, []byte("replacement"), 0600); err != nil {
				t.Fatal(err)
			}
			assertContents(t, target, "original")
			assertContents(t, link, "replacement")
			info, err := os.Lstat(link)
			if err != nil || !info.Mode().IsRegular() {
				t.Fatalf("replacement isn't a regular file: %v, %v", info, err)
			}
			assertNoTemporaryFiles(t, directory)
		})
	}
}

func TestRenameFailureCleansTemporaryFile(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "occupied")
	err := Write(path, 0600, func(writer io.Writer) error {
		if _, err := io.WriteString(writer, "replacement"); err != nil {
			return err
		}
		return os.Mkdir(path, 0700)
	})
	if err == nil {
		t.Fatal("rename over a directory was accepted")
	}
	info, statErr := os.Stat(path)
	if statErr != nil || !info.IsDir() {
		t.Fatalf("concurrently created directory was replaced: %v, %v", info, statErr)
	}
	assertNoTemporaryFiles(t, directory)
}

func TestZeroPermissionBits(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not implement POSIX permission bits")
	}
	directory := t.TempDir()
	path := filepath.Join(directory, "data")
	if err := WriteFile(path, []byte("private"), 0); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0 {
		t.Fatalf("zero mode = %v, %v", info, err)
	}
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	assertContents(t, path, "private")
}

func TestLongDestinationBasename(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, strings.Repeat("a", 255))
	if err := os.WriteFile(path, []byte("original"), 0600); err != nil {
		t.Skipf("filesystem does not permit this basename length: %v", err)
	}
	if err := WriteFile(path, []byte("replacement"), 0600); err != nil {
		t.Fatal(err)
	}
	assertContents(t, path, "replacement")
	assertNoTemporaryFiles(t, directory)
}

func TestSymlinkParentDotDotUsesRealDestinationDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symbolic links may require elevated privileges")
	}
	directory := t.TempDir()
	lexical := filepath.Join(directory, "lexical")
	realParent := filepath.Join(directory, "actual")
	child := filepath.Join(realParent, "child")
	for _, path := range []string{lexical, child} {
		if err := os.MkdirAll(path, 0700); err != nil {
			t.Fatal(err)
		}
	}
	link := filepath.Join(lexical, "link")
	if err := os.Symlink(child, link); err != nil {
		t.Fatal(err)
	}
	actualTarget := filepath.Join(realParent, "data")
	lexicalTarget := filepath.Join(lexical, "data")
	for _, path := range []string{actualTarget, lexicalTarget} {
		if err := os.WriteFile(path, []byte("original"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	// filepath.Join/Clean would discard the symlink-sensitive link/.. pair.
	path := link + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "data"
	err := Write(path, 0600, func(writer io.Writer) error {
		files, err := filepath.Glob(filepath.Join(realParent, ".atomicfile-*.tmp"))
		if err != nil || len(files) != 1 {
			t.Fatalf("temporary file was not in the real target directory: %v, %v", files, err)
		}
		assertNoTemporaryFiles(t, lexical)
		_, err = io.WriteString(writer, "replacement")
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	assertContents(t, actualTarget, "replacement")
	assertContents(t, lexicalTarget, "original")
	assertNoTemporaryFiles(t, realParent)
}

func TestCallbackChdirCannotRetargetRelativeDestination(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "failure"}[fail], func(t *testing.T) {
			directory := t.TempDir()
			first := filepath.Join(directory, "first")
			second := filepath.Join(directory, "second")
			for _, path := range []string{first, second} {
				if err := os.Mkdir(path, 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(path, "data"), []byte("original"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			t.Chdir(first)
			wantErr := errors.New("callback failure")
			err := Write("data", 0600, func(writer io.Writer) error {
				if _, err := io.WriteString(writer, "replacement"); err != nil {
					return err
				}
				if err := os.Chdir(second); err != nil {
					return err
				}
				if fail {
					return wantErr
				}
				return nil
			})
			want := "replacement"
			if fail {
				want = "original"
				if err != wantErr {
					t.Fatalf("callback error = %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			assertContents(t, filepath.Join(first, "data"), want)
			assertContents(t, filepath.Join(second, "data"), "original")
			assertNoTemporaryFiles(t, first)
			assertNoTemporaryFiles(t, second)
		})
	}
}

func TestCallbackParentRenameKeepsDestinationAnchored(t *testing.T) {
	directory := t.TempDir()
	parent := filepath.Join(directory, "parent")
	moved := filepath.Join(directory, "moved")
	if err := os.Mkdir(parent, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(parent, "data")
	if err := os.WriteFile(path, []byte("original"), 0600); err != nil {
		t.Fatal(err)
	}
	err := Write(path, 0600, func(writer io.Writer) error {
		if err := os.Rename(parent, moved); err != nil {
			return err
		}
		_, err := io.WriteString(writer, "replacement")
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	assertContents(t, filepath.Join(moved, "data"), "replacement")
	assertNoTemporaryFiles(t, moved)
}
