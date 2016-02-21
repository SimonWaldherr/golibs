package file

import (
	"log"
	"os"
	"runtime"
	"simonwaldherr.de/go/golibs/as"
	"strings"
	"testing"
)

type pathTest struct {
	in  string
	out string
}

var pathTests []pathTest
var realHome string

func init() {
	os.Symlink("../examples/all.go", "../examples/symlink")
	os.Symlink("../examples/1", "../examples/2")
	os.Symlink("../examples/2", "../examples/1")
	pathTests = []pathTest{
		{
			in:  "/Users/simonwaldherr/git/golibs/file/test.txt",
			out: "/Users/simonwaldherr/git/golibs/file/test.txt",
		},
		{
			in:  "~/git/golibs/file/test.txt",
			out: "/Users/simonwaldherr/git/golibs/file/test.txt",
		},
		{
			in:  "./test.txt",
			out: "/Users/simonwaldherr/git/golibs/file/test.txt",
		},
		{
			in:  "../file/test.txt",
			out: "/Users/simonwaldherr/git/golibs/file/test.txt",
		},
		{
			in:  "./../file/test.txt",
			out: "/Users/simonwaldherr/git/golibs/file/test.txt",
		},
		{
			in:  "../file/../file/test.txt",
			out: "/Users/simonwaldherr/git/golibs/file/test.txt",
		},
		{
			in:  "../file/../file/../examples/symlink",
			out: "/Users/simonwaldherr/git/golibs/examples/all.go",
		},
		{
			in:  "../file/../file/../examples/2",
			out: "",
		},
	}
	realHome = GetHomeDir()
	FakeHomeDir("/Users/simonwaldherr")
}

func Test_Read(t *testing.T) {
	str, err := Read("test.txt")
	if err != nil {
		t.Fatalf("file.Read Test failed")
	}
	if str != "Lorem ipsum dolor sit amet, consectetur adipisici elit, sed eiusmod tempor incidunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquid ex ea commodi consequat. Quis aute iure reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint obcaecat cupiditat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." {
		t.Fatalf("file.Read Test failed")
	}
	_, err = Read("😃")
	if err == nil {
		t.Fatalf("file.Read Test failed")
	}
}

func Test_ReadUntil(t *testing.T) {
	content, lastChar, pos, err := ReadUntil("file_test.go", []string{"\n", "\r"})
	if err != nil || content != "package file" || lastChar != "\n" || pos != 12 {
		t.Fatalf("file.ReadUntil Test #1 failed")
	}
	content, lastChar, pos, err = ReadUntil("file1.go", []string{"\n", "\r"})
	if err == nil {
		t.Fatalf("file.ReadUntil Test #2 failed")
	}
	content, lastChar, pos, err = ReadUntil("file.go", []string{"😃"})
	if err != nil || lastChar != "" {
		t.Fatalf("file.ReadUntil Test #3 failed: %v %v", err, lastChar)
	}
}

func Test_ReadBlocks(t *testing.T) {
	var (
		lines1 int
		lines2 int
	)
	content, err := ReadBlocks("file.go", []string{"\n", "\r"}, func(x string) (string, error) {
		lines1++
		return x, nil
	})
	lines2 = strings.Count(content, "\n")
	log.Printf("x:%v x:%v\n", lines1, lines2)
	if err != nil || (lines1 != lines2+1) {
		t.Fatalf("file.ReadBlocks1 Test failed")
	}

	content, err = ReadBlocks("", []string{"\n", "\r"}, func(x string) (string, error) {
		lines1++
		return x, nil
	})
	if err == nil {
		t.Fatalf("file.ReadBlocks2 Test failed")
	}
}

func Test_ReadDir(t *testing.T) {
	dir, err := ReadDir(".")
	if err != nil || len(dir) < 2 {
		t.Fatalf("file.ReadDir Test failed")
	}
}

func Test_Each(t *testing.T) {
	var count int
	var fsize int64
	err := Each("..", true, func(filename, extension, filepath string, dir bool, fileinfo os.FileInfo) {
		if extension == "go" && !dir {
			size, sizeErr := Size(filepath)
			fsize = fsize + size
			//t.Logf("%v, %v, %v, %v\n", filename, filepath, dir, fileinfo)
			count++
			if sizeErr != nil {
				t.Fatalf("file.Each Test failed")
			}
		}
	})
	if fsize < 100000 {
		t.Fatalf("file.Each Size Test failed, measured size is: %v", fsize)
	}
	if err != nil {
		t.Fatalf("file.Each Test failed")
	}
	if count < 10 {
		t.Fatalf("file.Each Count Test failed")
	}
}

func Test_Write(t *testing.T) {
	str := "Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Maecenas sed diam eget risus varius blandit sit amet non magna. Morbi leo risus, porta ac consectetur ac, vestibulum at eros."
	err := Write("writetest.log", str, false)
	if err != nil {
		t.Fatalf("file.Write Test failed")
	}
}

func Test_Size(t *testing.T) {
	size, err := Size("writetest.log")
	if err != nil || size != 193 {
		t.Fatalf("file.Size Test failed")
	}
}

func Test_Time(t *testing.T) {
	atime, mtime, ctime, err := Time("writetest.log")
	now := as.Time("10.10.2015")
	a := as.Int(atime)
	m := as.Int(mtime)
	c := as.Int(ctime)
	n := as.Int(now)
	if err != nil || a < n || m < n || c < n {
		t.Fatalf("file.Time Test failed: \natime:\t%v\nctime:\t%v\natime:\t%v\nnow:\t%v\n", a, m, c, n)
	}
	_, _, _, err = Time("foobar")
	if err == nil {
		t.Fatalf("file.Time Test 2 failed")
	}
}

func Test_Write2(t *testing.T) {
	str := "Aenean lacinia bibendum nulla sed consectetur. Maecenas sed diam eget risus varius blandit sit amet non magna. Maecenas faucibus mollis interdum."
	err := Write("writetest.log", str, true)
	if err != nil {
		t.Fatalf("file.Write Test failed")
	}
}

func Test_Size2(t *testing.T) {
	size, err := Size("writetest.log")
	if err != nil || size != 338 {
		t.Fatalf("file.Size Test failed")
	}
}

func Test_Clean(t *testing.T) {
	err := Clean("writetest.log")
	if err != nil {
		t.Fatalf("file.Clean Test failed")
	}
}

func Test_Size3(t *testing.T) {
	size, err := Size("writetest.log")
	if err != nil || size != 0 {
		t.Fatalf("file.Size Test failed")
	}
}

func Test_Rename(t *testing.T) {
	err := Rename("writetest.log", "writetest2.log")
	if err != nil {
		t.Fatalf("file.Rename Test failed")
	}
}

func Test_X(t *testing.T) {
	_, err := Size("writetest.log")
	if err == nil {
		t.Fatalf("file.X Test failed")
	}
}

func Test_Copy(t *testing.T) {
	err := Copy("writetest2.log", "writetest.log")
	if err != nil {
		t.Fatalf("file.X Test failed")
	}
}

func Test_Delete(t *testing.T) {
	err := Delete("writetest.log")
	if err != nil {
		t.Fatalf("file.Delete Test failed")
	}
	err = Delete("writetest2.log")
	if err != nil {
		t.Fatalf("file.Delete Test failed")
	}
}

func Test_X2(t *testing.T) {
	_, err := Size("writetest2.log")
	if err == nil {
		t.Fatalf("file.X2 Test failed")
	}
}

func Test_X3(t *testing.T) {
	x := Exists("test.txt")
	if x != true {
		t.Fatalf("file.Exists Test failed")
	}
	x = IsDir("test.txt")
	if x != false {
		t.Fatalf("file.IsDir Test failed")
	}
	x = IsDir("../file/")
	if x != true {
		t.Fatalf("file.IsDir Test failed")
	}
	x = IsFile("../file/")
	if x == true {
		t.Fatalf("file.IsFile Test failed")
	}
	x = IsFile("test.txt")
	if x != true {
		t.Fatalf("file.IsFile2 Test failed")
	}
	x = IsSymlink("test.txt")
	if x != false {
		t.Fatalf("file.IsSymlink1 Test failed")
	}
}

func Test_Ψ(t *testing.T) {
	x := Exists("😃")
	if x == true {
		t.Fatalf("file.Exists2 Test failed")
	}
	x = IsFile("😃")
	if x == true {
		t.Fatalf("file.IsFile3 Test failed")
	}
	x = IsSymlink("😃")
	if x != false {
		t.Fatalf("file.IsSymlink2 Test failed")
	}
	x = IsDir("😃")
	if x != false {
		t.Fatalf("file.IsDir2 Test failed")
	}
	_, err := Size("😃")
	if err == nil {
		t.Fatalf("file.Size Test failed")
	}
	err = Rename("😃", "😃")
	if err == nil {
		t.Fatalf("file.Rename Test failed")
	}
	err = Delete("😃")
	if err == nil {
		t.Fatalf("file.Delete Test failed")
	}
	err = Write("", "", false)
	if err == nil {
		t.Fatalf("file.Write Test failed")
	}
	err = Each("", true, func(filename, extension, filepath string, dir bool, fileinfo os.FileInfo) {
		if extension == "go" && !dir {
			t.Logf("%v, %v, %v, %v\n", filename, filepath, dir, fileinfo)
		}
	})
	if err == nil {
		t.Fatalf("file.Each Test failed")
	}
	_, err = GetAbsolutePath("//\\//\\......//////..\\\\//😃\\")
	_, err = GetAbsolutePath("")
	if err == nil {
		t.Errorf("file.GetAbsolutePath Test failed")
	}
}

func TestHomeDir(t *testing.T) {
	if HomeDir == "" {
		t.Errorf("file.[SG]etHomeDir test #1 failed: %v", HomeDir)
	}
	if realHome == "" {
		t.Errorf("file.[SG]etHomeDir test #2 failed: %v", realHome)
	}
}

func changeWD(str string) string {
	var orig, wd string
	wd, _ = os.Getwd()
	wd = strings.Replace(wd+"--", "file--", "", 1)
	orig = "/Users/simonwaldherr/git/golibs/"

	switch {
	case strings.Contains(wd, "/home/travis/gopath/"):
		// https://travis-ci.org/
		return strings.Replace(str, orig, "/home/travis/gopath/src/github.com/SimonWaldherr/golibs/", 1)
	case strings.Contains(wd, "/home/magnum/golibs/"):
		// https://magnum-ci.com/
		return strings.Replace(str, orig, "/home/magnum/golibs/", 1)
	case strings.Contains(wd, "c:\\gopath\\src\\github.com\\simonwaldherr\\golibs\\"):
		// https://ci.appveyor.com/
		return strings.Replace(str, orig, "c:\\gopath\\src\\github.com\\simonwaldherr\\golibs\\", 1)
	case strings.Contains(wd, "/home/runner/workspace/"):
		// https://semaphoreci.com/
		return strings.Replace(str, orig, "/home/runner/golibs/", 1)
	}

	return strings.Replace(str, orig, wd, 1)
}

func TestGetAbsolutePath(t *testing.T) {
	var err error
	var converted string
	var failed bool = false

	for i, te := range pathTests {
		input := te.in
		expected := te.out
		if i > 1 {

			expected = changeWD(expected)

			if runtime.GOOS == "windows" {
				expected = strings.Replace(expected, "/", "\\", -1)
			}
		} else if runtime.GOOS == "windows" {
			input = strings.Replace(input, "/", "\\", -1)
		}
		converted, err = GetAbsolutePath(te.in)

		if !((runtime.GOOS == "windows" && i == 6) || (strings.Contains(converted, "/var/lib/wercker/"))) {
			if converted != expected {
				failed = true
				t.Errorf("GetAbsolutePath test #%d failed \nIn: \"%s\"\nExpected: \"%s\"\nActually: \"%s\"\nError: \"%v\"", i, input, expected, converted, err)
			}
		}
	}

	if failed {
		wd, _ := os.Getwd()
		t.Errorf("os.Getwd() => %s\n", wd)
	}

	if runtime.GOOS != "windows" {
		converted, err = GetAbsolutePath("")
		if err == nil {
			t.Errorf("GetAbsolutePath test #8 failed: %v", converted)
		}
		converted, err = GetAbsolutePath("///")
		if converted != "/" {
			t.Errorf("GetAbsolutePath test #9 failed: %v", converted)
		}
		converted, err = GetAbsolutePath("/././")
		if converted != "/" {
			t.Errorf("GetAbsolutePath test #10 failed: %v", converted)
		}
		converted, err = GetAbsolutePath("/../")
		if converted != "/" {
			t.Errorf("GetAbsolutePath test #11 failed: %v", converted)
		}
	}
}

func TestSetHomeDir(t *testing.T) {
	FakeHomeDir("#")
	SetHomeDir()
	str := GetHomeDir()
	if str == "" {
		t.Errorf("file.[SG]etHomeDir test #3 failed: %v", str)
	}
}
