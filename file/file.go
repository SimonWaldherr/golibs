// file wraps around the standard functions to simplify reading and writing on disk
package file

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"simonwaldherr.de/go/golibs/gopath"
)

func delimiterSet(delimiters []string) map[string]struct{} {
	set := make(map[string]struct{}, len(delimiters))
	for _, delimiter := range delimiters {
		set[delimiter] = struct{}{}
	}
	return set
}

// Exists checks if a file exists
func Exists(fn string) bool {
	if _, err := os.Stat(fn); err == nil {
		return true
	}
	return false
}

// IsDir checks if a file is a directory
func IsDir(fn string) bool {
	file, err := os.Stat(fn)
	if err != nil {
		return false
	}
	if file.IsDir() {
		return true
	}
	return false
}

// IsFile checks if a file is a regular file
func IsFile(fn string) bool {
	file, err := os.Stat(fn)
	if err != nil {
		return false
	}
	if file.Mode().IsRegular() {
		return true
	}
	return false
}

func isSymlink(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeSymlink == os.ModeSymlink
}

// IsSymlink checks if a file is a symlink
func IsSymlink(fn string) bool {
	file, err := os.Lstat(fn)
	if err != nil {
		return false
	}
	return isSymlink(file)
}

// Read reads a file and returns the content as a string
func Read(fn string) (string, error) {
	var file *os.File
	var err error

	buf := bytes.NewBuffer(nil)

	file, err = os.Open(fn)

	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	if _, err = io.Copy(buf, file); err != nil {
		return "", err
	}

	s := buf.String()
	return s, nil
}

// ReadUntil reads through the first matching rune delimiter. The returned text
// excludes the delimiter and position counts preceding runes, not bytes.
// Each delimiter must be one rune; longer strings never match. Invalid UTF-8 is
// replaced with U+FFFD, matching bufio.ScanRunes. Read errors return partial text.
func ReadUntil(fn string, delim []string) (string, string, int, error) {
	file, err := os.Open(fn)
	if err != nil {
		return "", "", 0, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	delimiters := delimiterSet(delim)
	var content strings.Builder
	position := 0
	for scanner.Scan() {
		char := scanner.Bytes()
		if _, found := delimiters[string(char)]; found {
			return content.String(), string(char), position, nil
		}
		content.Write(char)
		position++
	}
	return content.String(), "", position, scanner.Err()
}

// ReadBytes is a legacy alias for ReadUntil. Despite the name, delimiters and
// positions use runes, and invalid UTF-8 is replaced with U+FFFD.
func ReadBytes(fn string, delim []string) (string, string, int, error) {
	return ReadUntil(fn, delim)
}

// ReadBlocks transforms delimiter-terminated blocks and concatenates their
// results. For compatibility, each delimiter starts the following block and an
// unterminated final block is not emitted. Read or callback errors return no output.
func ReadBlocks(fn string, delim []string, fnc func(string) (string, error)) (string, error) {
	file, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	delimiters := delimiterSet(delim)
	var block, result strings.Builder
	for scanner.Scan() {
		char := scanner.Bytes()
		if _, found := delimiters[string(char)]; found {
			transformed, err := fnc(block.String())
			if err != nil {
				return "", err
			}
			result.WriteString(transformed)
			block.Reset()
		}
		block.Write(char)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return result.String(), nil
}

// ReadLines reads newline-separated lines. It returns successfully read lines
// together with any read or scanner error, including lines exceeding
// bufio.MaxScanTokenSize. Newline characters are excluded from the results.
func ReadLines(fn string) ([]string, error) {
	var file *os.File
	var err error

	file, err = os.Open(fn)

	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Write writes a string to a file, syncing its contents before returning.
// Existing permissions are retained; new files use 0666 filtered by the umask.
// Write, sync, and close errors are returned to the caller.
func Write(fn, str string, append bool) (err error) {
	flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	if append {
		flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	}
	file, err := os.OpenFile(fn, flags, 0666)
	if err != nil {
		return err
	}
	defer closeWriter(file, &err)
	if _, err = file.WriteString(str); err != nil {
		return err
	}
	return file.Sync()
}

// WriteLines concatenates strings without adding separators and syncs the file.
// Writes are buffered to avoid a system call for every small string. Flush,
// sync, and close failures are returned to the caller.
func WriteLines(fn string, lines []string, append bool) (err error) {
	flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	if append {
		flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	}
	file, err := os.OpenFile(fn, flags, 0666)
	if err != nil {
		return err
	}
	defer closeWriter(file, &err)
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err = writer.WriteString(line); err != nil {
			return err
		}
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	return file.Sync()
}

func closeWriter(file *os.File, result *error) {
	if err := file.Close(); *result == nil {
		*result = err
	}
}

// Append appends a string to a file
func Append(fn, str string) error {
	return Write(fn, str, true)
}

// AppendLines appends a string slice to a file
func AppendLines(fn string, lines []string) error {
	return WriteLines(fn, lines, true)
}

// Size returns the size of a file
func Size(fn string) (int64, error) {
	file, err := os.Open(fn)

	if err != nil {
		return -1, err
	}
	defer func() { _ = file.Close() }()

	fi, err := file.Stat()

	if err != nil {
		return -1, err
	}

	return fi.Size(), nil
}

// Clean removes a file
func Clean(fn string) error {
	return Write(fn, "", false)
}

// Remove removes a file
func Remove(fn string) error {
	return os.Remove(fn)
}

// Rename renames a file
func Rename(from, to string) error {
	return os.Rename(from, to)
}

// Copy replaces the destination with the source's contents. It rejects source
// directories and identical files (including hard links and symlinks) before
// truncating the destination. Copy is not an atomic replacement operation.
func Copy(from, to string) (err error) {
	r, err := os.Open(from)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()
	source, err := r.Stat()
	if err != nil {
		return err
	}
	if source.IsDir() {
		return fmt.Errorf("file: cannot copy a directory: %w", os.ErrInvalid)
	}

	// Open without O_TRUNC so identity checks cannot destroy the source first.
	w, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer closeWriter(w, &err)
	destination, err := w.Stat()
	if err != nil {
		return err
	}
	if os.SameFile(source, destination) {
		return fmt.Errorf("file: source and destination are the same file: %w", os.ErrInvalid)
	}
	if destination.Mode().IsRegular() {
		if err := w.Truncate(0); err != nil {
			return err
		}
	}
	_, err = io.Copy(w, r)
	return err
}

// Delete deletes a file
func Delete(fn string) error {
	return os.Remove(fn)
}

// ReadDir reads a directory and returns the content as a string slice
func ReadDir(dn string) ([]string, error) {
	var flist []string
	dn, err := GetAbsolutePath(dn)
	if err != nil {
		return []string{""}, err
	}
	files, err := ioutil.ReadDir(dn)
	if err != nil {
		return []string{""}, err
	}
	for _, f := range files {
		flist = append(flist, f.Name())
	}
	return flist, nil
}

// Each calls fnc for each directory entry, optionally descending into real
// subdirectories. It does not follow directory symlinks. Directory handles are
// closed before callbacks or recursion; callback paths are absolute.
func Each(dirname string, recursive bool, fnc func(string, string, string, bool, os.FileInfo)) error {
	directory, err := os.Open(dirname)
	if err != nil {
		return err
	}
	list, readErr := directory.Readdir(-1)
	closeErr := directory.Close()
	if readErr != nil {
		return readErr
	}
	if closeErr != nil {
		return closeErr
	}
	// Resolve symlinks before Abs cleans the path: link/.. can name a different
	// directory from its lexical parent.
	absolute, err := filepath.EvalSymlinks(dirname)
	if err != nil {
		return err
	}
	absolute, err = filepath.Abs(absolute)
	if err != nil {
		return err
	}
	for _, fi := range list {
		path := filepath.Join(absolute, fi.Name())
		isDir := fi.IsDir()
		suffix := fi.Name()
		if dot := strings.LastIndexByte(suffix, '.'); dot >= 0 {
			suffix = suffix[dot+1:]
		}
		fnc(fi.Name(), suffix, path, isDir, fi)
		if recursive && isDir {
			if err := Each(path, recursive, fnc); err != nil {
				return err
			}
		}
	}
	return nil
}

var HomeDir string

// SetHomeDir sets the home directory
func SetHomeDir() string {
	if HomeDir == "#" || runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		HomeDir = home
		return home
	}
	HomeDir = os.Getenv("HOME")
	return os.Getenv("HOME")
}

// FakeHomeDir sets the home directory to a fake value
func FakeHomeDir(dir string) string {
	HomeDir = dir
	return dir
}

// GetHomeDir returns the home directory
func GetHomeDir() string {
	if HomeDir == "" {
		return SetHomeDir()
	}
	return HomeDir
}

func switchSymlink(path []byte, start int, link, after string) []byte {
	if link[0] == os.PathSeparator {
		return []byte(filepath.Join(link, after))
	}
	return []byte(filepath.Join(string(path[0:start]), link, after))
}

func nextComponent(path []byte, start int) []byte {
	v := bytes.IndexByte(path[start:], os.PathSeparator)
	if v < 0 {
		return path
	}
	return path[0 : start+v]
}

// GetAbsolutePath returns the absolute path to a file or dir
// if it is a relative path it is relative to the current working directory
func GetAbsolutePath(fn string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return getAbsolutePathHelper(fn, pwd)
}

// GetAbsolutePathByApp returns the absolute path to a file or dir
// if it is a relative path it is relative to the Go Application source or binary
func GetAbsolutePathByApp(fn string) (string, error) {
	pwd := gopath.Dir()
	return getAbsolutePathHelper(fn, pwd)
}

func getAbsolutePathHelper(fn string, pwd string) (string, error) {
	if len(fn) == 0 {
		return "", os.ErrInvalid
	}

	if fn[0] != os.PathSeparator {
		if fn[0] == '.' {
			fn = filepath.Join(pwd, fn)
		} else {
			fn = strings.Replace(fn, "~", GetHomeDir(), 1)
		}
	}

	path := []byte(fn)
	nlinks := 0
	start := 1
	prev := 1

	for start < len(path) {
		c := nextComponent(path, start)
		cur := c[start:]

		if len(cur) == 0 {
			copy(path[start:], path[start+1:])
			path = path[0 : len(path)-1]
		} else if len(cur) == 1 && cur[0] == '.' {
			if start+2 < len(path) {
				copy(path[start:], path[start+2:])
			}
			path = path[0 : len(path)-2]
		} else if len(cur) == 2 && cur[0] == '.' && cur[1] == '.' {
			copy(path[prev:], path[start+2:])
			path = path[0 : len(path)+prev-(start+2)]
			prev = 1
			start = 1
		} else {
			fi, err := os.Lstat(string(c))
			if err == nil {
				if isSymlink(fi) {
					nlinks++
					if nlinks > 8 {
						return "", os.ErrInvalid
					}

					var link string
					link, err := os.Readlink(string(c))
					if err != nil {
						return "", err
					}
					after := string(path[len(c):])
					path = switchSymlink(path, start, link, after)
					prev = 1
					start = 1
				} else {
					prev = start
					start = len(c) + 1
				}
			} else {
				return string(path), nil
			}
		}
	}

	for len(path) > 1 && path[len(path)-1] == os.PathSeparator {
		path = path[0 : len(path)-1]
	}
	return string(path), nil
}
