// Package atomicfile writes complete replacement files through temporary files.
// Renaming a file in the same directory provides atomic replacement on systems
// that support it, including Unix. It is not a cross-platform atomicity promise.
package atomicfile

import (
	"crypto/rand"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// WriteFile replaces path with data using Write. Requested POSIX permission bits
// replace the old permissions; other metadata, such as ownership and extended
// attributes, is not preserved. Other platforms apply their OS permission rules.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	return Write(path, perm, func(writer io.Writer) error {
		_, err := writer.Write(data)
		return err
	})
}

// Write calls fn with a temporary file created with mode 0600 in path's directory.
// After fn succeeds, it applies perm.Perm(), syncs and closes the file, then
// renames it over path. Only the final rename changes the destination. Errors
// before that point leave an existing destination unchanged. Panics propagate
// after temporary-file cleanup.
//
// An existing target must be a regular file or a symlink. A final symlink is
// replaced, never followed. Parent paths are caller-owned and may be resolved
// normally by the operating system; this is not a filesystem sandbox.
// The parent directory must already exist and be accessible and writable; it
// is not created. Its handle is captured before fn, preserving the operating
// system's symlink/.. resolution and staying anchored if fn changes the working
// directory or renames a parent directory.
//
// fn must finish its writes before returning and must not retain the writer.
// Only permission bits are used. On POSIX systems they are applied exactly,
// independently of the umask; these bits do not establish Windows ACLs.
// Concurrent writes are not serialized; the last successful rename wins.
// The file is synced, but its parent directory is not, so success does not
// promise survival of the directory update across a crash or power loss.
func Write(path string, perm os.FileMode, fn func(io.Writer) error) error {
	if fn == nil {
		return errors.New("atomicfile: nil write function")
	}
	directory, base := filepath.Split(path)
	if base == "" || base == "." || base == ".." {
		return &os.PathError{Op: "atomicfile", Path: path, Err: os.ErrInvalid}
	}
	if directory == "" {
		directory = "."
	}
	root, err := os.OpenRoot(directory)
	if err != nil {
		return err
	}
	defer root.Close()
	if info, err := root.Lstat(base); err == nil {
		if !info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0 {
			return &os.PathError{Op: "atomicfile", Path: path, Err: os.ErrInvalid}
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	file, name, err := createTemp(root)
	if err != nil {
		return err
	}
	renamed := false
	defer func() {
		_ = file.Close()
		if !renamed {
			_ = root.Remove(name)
		}
	}()
	if err := fn(writeOnly{file}); err != nil {
		return err
	}
	if err := file.Chmod(perm.Perm()); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := root.Rename(name, base); err != nil {
		return err
	}
	renamed = true
	return nil
}

func createTemp(root *os.Root) (*os.File, string, error) {
	for attempt := 0; attempt < 100; attempt++ {
		name := ".atomicfile-" + rand.Text() + ".tmp"
		file, err := root.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
		if err == nil {
			return file, name, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, "", err
		}
	}
	return nil, "", errors.New("atomicfile: could not create a unique temporary file")
}

// Expose writing without giving callbacks ownership of Close, Chmod or Name.
type writeOnly struct{ file *os.File }

func (writer writeOnly) Write(data []byte) (int, error) { return writer.file.Write(data) }

func (writer writeOnly) WriteString(data string) (int, error) { return writer.file.WriteString(data) }

func (writer writeOnly) ReadFrom(reader io.Reader) (int64, error) {
	return writer.file.ReadFrom(reader)
}
