package gopath

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

func Name() string {
	_, dir, _, _ := runtime.Caller(1)
	return dir
}

func Ga() string {
	var x string
	x += fmt.Sprintln("")
	x += fmt.Sprintln("os.Args[0]:", os.Args[0])
	x += fmt.Sprintln("filepath.Dir(os.Args[0]):", filepath.Dir(os.Args[0]))
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	x += fmt.Sprintln("filepath.Abs:", dir)
	dir, _ = os.Getwd()
	x += fmt.Sprintln("Getwd:", dir)
	_, dir, _, _ = runtime.Caller(0)
	x += fmt.Sprintln("runtime.Caller:", dir)
	x += fmt.Sprintln("os.Args FileType:", GetFileType(os.Args[0]))
	x += fmt.Sprintln("Name FileType:", GetFileType(Name()))
	return x
}

func Type() string {
	return os.Getenv("GOPATH")
}
