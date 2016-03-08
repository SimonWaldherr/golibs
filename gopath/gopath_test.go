package gopath

import (
	"os"
	"strings"
	"testing"
)

var archlist = map[string]bool{
	"386":   true,
	"amd64": true,
	"arm":   true,
}

var oslist = map[string]bool{
	"darwin":  true,
	"linux":   true,
	"windows": true,
}

func isDir(fn string) bool {
	file, err := os.Stat(fn)
	if err != nil {
		return false
	}
	if file.IsDir() {
		return true
	}
	return false
}

func Test_Compiler(t *testing.T) {
	if v := Compiler(); v != "gc" {
		t.Fatalf("Compiler should be \"gc\" but is: %v", v)
	}
}

func Test_GOARCH(t *testing.T) {
	v := GOARCH()
	if _, ok := archlist[v]; !ok {
		t.Fatalf("GOARCH should be \"386\", \"amd64\" or \"arm\" but is: %v", v)
	}
}

func Test_GOOS(t *testing.T) {
	v := GOOS()
	if _, ok := oslist[v]; !ok {
		t.Fatalf("GOOS should be \"darwin\", \"linux\" or \"windows\" but is: %v", v)
	}
}

func Test_GOROOT(t *testing.T) {
	if v := GOROOT(); !isDir(v) {
		t.Fatalf("GOROOT should be a directory but is: %v", v)
	}
}

func Test_GOPATH(t *testing.T) {
	if v := GOPATH(); !isDir(v) {
		t.Fatalf("GOPATH should be a directory but is: %v", v)
	}
}

func Test_Compiled(t *testing.T) {
	if v := Compiled(); v == true {
		t.Fatalf("Compiled should be false but is: %v", v)
	}
}

func Test_Path(t *testing.T) {
	if v := Path(); !strings.HasSuffix(v, "gopath_test.go") {
		t.Fatalf("Path should end with \"gopath_test.go\" but is: %v", v)
	}
}

func Test_Name(t *testing.T) {
	if v := Name(); !strings.HasSuffix(v, "gopath_test.go") {
		t.Fatalf("Name should end with \"gopath_test.go\" but is: %v", v)
	}
}

func Test_Dir(t *testing.T) {
	if v := Dir(); !isDir(v) {
		t.Fatalf("Dir should be a directory but is: %v", v)
	}
}

func Test_GetFileType(t *testing.T) {
	if v := GetFileType(Dir() + "404"); v != 0 {
		t.Fatalf("GetFileType should be \"0\" but is: %v", v)
	}
	if v := GetFileType(Dir()); v != 1 {
		t.Fatalf("GetFileType should be \"1\" but is: %v", v)
	}
	if v := GetFileType(Name()); v != 2 {
		t.Fatalf("GetFileType should be \"2\" but is: %v", v)
	}
	if v := GetFileType(os.Args[0]); v != 3 {
		t.Fatalf("GetFileType should be \"3\" but is: %v", v)
	}
}

func Test_WD(t *testing.T) {
	if v := WD(); !isDir(v) {
		t.Fatalf("WD should be a directory but is: %v", v)
	}
}
