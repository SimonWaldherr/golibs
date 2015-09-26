package gopath

import (
	"testing"
)

func Test_Compiler(t *testing.T) {
	t.Log(Compiler())
}

func Test_GOARCH(t *testing.T) {
	t.Log(GOARCH())
}

func Test_GOOS(t *testing.T) {
	t.Log(GOOS())
}

func Test_GOROOT(t *testing.T) {
	t.Log(GOROOT())
}

func Test_GOPATH(t *testing.T) {
	t.Log(GOPATH())
}

func Test_Name(t *testing.T) {
	t.Log(Name())
}

func Test_X(t *testing.T) {
	t.Log(Ga())
}
