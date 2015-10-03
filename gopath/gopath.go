// gopath provides an easy way to get system information
package gopath

import (
	"os"
	"path"
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

// GetFileType returns the type of a given file.
// The return values can be Directory, ASCII, Binary
// or DontExist.
func GetFileType(name string) FileType {
	if file, err := os.Open(name); err == nil {
		if fileStat, err := file.Stat(); err == nil && fileStat.IsDir() {
			return Directory
		} else if err == nil {
			data := make([]byte, 100)
			file.Read(data)
			if utf8.Valid(data) {
				return ASCII
			}
			return Binary
		}
	}
	return DontExist
}

// Compiler returns the name of the used compiler
func Compiler() string {
	return runtime.Compiler
}

// GOARCH returns the system architecture
func GOARCH() string {
	return runtime.GOARCH
}

// GOOS returns the name of the Operating System
func GOOS() string {
	return runtime.GOOS
}

// GOROOT returns the $GOROOT environment variable
func GOROOT() string {
	return runtime.GOROOT()
}

// GOPATH returns the $GOPATH environment variable
func GOPATH() string {
	return os.Getenv("GOPATH")
}

// WD returns the working directory, this is
// the folder where you were last.
func WD() string {
	dir, _ := os.Getwd()
	return dir
}

// Compiled returns a boolean which tells if you run
// a already compiled binary or "go run" the source.
func Compiled() bool {
	if strings.HasPrefix(os.Args[0], "/var/folders/") ||
		strings.HasPrefix(os.Args[0], "/tmp/go-build") ||
		strings.Contains(os.Args[0], "\\AppData\\Local\\Temp\\") {
		return false
	}
	return true
}

// gpath is the helper function for the public functions
// Path(), Name() and Dir()
func gpath(i int) string {
	var filename string
	if Compiled() == false {
		_, filename, _, _ = runtime.Caller(i)
	} else {
		filename, _ = filepath.Abs(filepath.Join(WD(), os.Args[0]))
	}
	return filename
}

// Path returns the full path of the binary or source file
// which is currently running.
func Path() string {
	return gpath(2)
}

// Name returns the filename of the binary or source file
// which is currently running.
func Name() string {
	_, f := path.Split(gpath(2))
	return f
}

// Dir returns the directory in which the binary or source file
// is stored. This is useful if you have config files in the
// same folder.
func Dir() string {
	return filepath.Dir(gpath(2))
}
