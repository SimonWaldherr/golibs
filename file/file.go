// "file" wraps around the standard functions to simplify reading and writing on disk
package file

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
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

func IsSymlink(fn string) bool {
	file, err := os.Lstat(fn)
	if err != nil {
		return false
	}
	if file.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	}
	return false
}

func Read(fn string) (string, error) {
	buf := bytes.NewBuffer(nil)
	file, err := os.Open(fn)

	defer file.Close()

	if err != nil {
		return "", err
	}
	io.Copy(buf, file)

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
	var (
		file   *os.File
		err error
	)
	if append {
		file, err = os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0666))
	} else {
		file, err = os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0666))
	}

	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.WriteString(str)
	if err != nil {
		return err
	}

	file.Sync()
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

func Rename(fn, fn2 string) error {
	err := os.Rename(fn, fn2)

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
