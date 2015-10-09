// file wraps around the standard functions to simplify reading and writing on disk
package file

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"simonwaldherr.de/go/golibs/gopath"
	"strings"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Exists(fn string) bool {
	if _, err := os.Stat(fn); err == nil {
		return true
	}
	return false
}

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

func IsSymlink(fn string) bool {
	file, err := os.Lstat(fn)
	if err != nil {
		return false
	}
	return isSymlink(file)
}

func Read(fn string) (string, error) {
	var file *os.File
	var err error

	buf := bytes.NewBuffer(nil)

	file, err = os.Open(fn)

	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err = io.Copy(buf, file); err != nil {
		return "", err
	}

	s := string(buf.Bytes())
	return s, nil
}

func ReadUntil(fn string, delim []string) (string, string, int, error) {
	file, err := os.Open(fn)

	defer file.Close()

	if err != nil {
		return "", "", 0, err
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanRunes)
	pos := 0
	buf := ""

	for scanner.Scan() {
		char := scanner.Text()
		if contains(delim, char) {
			return buf, char, pos, nil
		}
		buf += char
		pos++
	}
	return buf, "", pos, nil
}

func ReadBlocks(fn string, delim []string, fnc func(string) (string, error)) (string, error) {
	file, err := os.Open(fn)

	defer file.Close()

	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanRunes)
	pos := 0
	buf := ""
	ret := ""
	str := ""

	for scanner.Scan() {
		char := scanner.Text()
		if contains(delim, char) {
			str, err = fnc(buf)
			ret += str
			buf = ""
		}
		buf += char
		pos++
	}
	return ret, nil
}

func Write(fn, str string, append bool) error {
	var file *os.File
	var err error

	if append {
		file, err = os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0666))
	} else {
		file, err = os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0666))
	}

	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(str); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}
	return nil
}

func Size(fn string) (int64, error) {
	file, err := os.Open(fn)

	defer file.Close()

	if err != nil {
		return -1, err
	}

	fi, err := file.Stat()

	if err != nil {
		return -1, err
	}

	return fi.Size(), nil
}

func Clean(fn string) error {
	return Write(fn, "", false)
}

func Rename(from, to string) error {
	err := os.Rename(from, to)

	if err != nil {
		return err
	}
	return nil
}

func Copy(from, to string) error {
	r, err := os.Open(from)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(to)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

func Delete(fn string) error {
	err := os.Remove(fn)

	if err != nil {
		return err
	}
	return nil
}

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

func Each(dirname string, recursive bool, fnc func(string, string, string, bool, os.FileInfo)) error {
	file, err := os.Open(dirname)

	defer file.Close()

	if err != nil {
		return err
	}
	list, err := file.Readdir(-1)

	if err != nil {
		return err
	}
	for _, fi := range list {
		path, _ := filepath.Abs(dirname + string(os.PathSeparator) + fi.Name())
		isDir := fi.IsDir()
		split := strings.Split(fi.Name(), ".")
		suffix := split[len(split)-1:][0] // if you want the real filetype, use http.DetectContentType
		fnc(fi.Name(), suffix, path, isDir, fi)
		if recursive && isDir {
			err = Each(path, recursive, fnc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var HomeDir string

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

func FakeHomeDir(dir string) string {
	HomeDir = dir
	return dir
}

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
					link, err = os.Readlink(string(c))
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
