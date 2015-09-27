package gopath

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf8"
)

type FileType int

const (
	DontExist FileType = iota
	Directory
	ASCII
	Binary
)

func GetFileType(name string) FileType {
	if file, err := os.Open(name); err == nil {
		if fileStat, err := file.Stat(); err == nil && fileStat.IsDir() {
			return Directory
		} else if err == nil {
			data := make([]byte, 100)
			_, err := file.Read(data)
			if err != nil {
				return DontExist
			}
			if utf8.Valid(data) {
				return ASCII
			}
			return Binary
		}
	}
	return DontExist
}

func Compiler() string {
	return runtime.Compiler
}

func GOARCH() string {
	return runtime.GOARCH
}

func GOOS() string {
	return runtime.GOOS
}

func GOROOT() string {
	return runtime.GOROOT()
}

func GOPATH() string {
	return os.Getenv("GOPATH")
}

func WD() string {
	dir, _ := os.Getwd()
	return dir
}

func Compiled() bool {
	if strings.HasPrefix(os.Args[0], "/var/folders/") ||
		strings.HasPrefix(os.Args[0], "/tmp/go-build") ||
		strings.Contains(os.Args[0], "\\AppData\\Local\\Temp\\") {
		return false
	}
	return true
}

func Name() string {
  var filename string
  if Compiled() == false {
    _, filename, _, _ = runtime.Caller(1)
  } else {
    filename, _ = filepath.Abs(filepath.Join(WD(), os.Args[0]))
  }
  return filename
}

func Path() string {
  return filepath.Dir(Name())
}
