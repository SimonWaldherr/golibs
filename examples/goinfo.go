// +build local

package main

import (
	"fmt"
	"simonwaldherr.de/go/golibs/gopath"
)

func main() {
	fmt.Println("Compiler:    ", gopath.Compiler())
	fmt.Println("GOARCH:      ", gopath.GOARCH())
	fmt.Println("GOOS:        ", gopath.GOOS())
	fmt.Println("GOROOT:      ", gopath.GOROOT())
	fmt.Println("GOPATH:      ", gopath.GOPATH())
	fmt.Println("WD:          ", gopath.WD())
	fmt.Println("Compiled:    ", gopath.Compiled())
	fmt.Println("Path:        ", gopath.Path())
	fmt.Println("Name:        ", gopath.Name())
	fmt.Println("Dir:         ", gopath.Dir())
	fmt.Println("GetFileType: ", gopath.GetFileType(gopath.Path()))
}
