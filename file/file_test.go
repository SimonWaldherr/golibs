package file

import (
	"log"
	"os"
	"strings"
	"testing"
)

func Test_Read(t *testing.T) {
	str, err := Read("test.txt")
	if err != nil {
		t.Fatalf("file.Read Test failed")
	}
	if str != "Lorem ipsum dolor sit amet, consectetur adipisici elit, sed eiusmod tempor incidunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquid ex ea commodi consequat. Quis aute iure reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint obcaecat cupiditat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." {
		t.Fatalf("file.Read Test failed")
	}
}

func Test_ReadUntil(t *testing.T) {
	content, lastChar, pos, err := ReadUntil("file_test.go", []string{"\n", "\r"})
	if err != nil || content != "package file" || lastChar != "\n" || pos != 12 {
		t.Fatalf("file.ReadUntil Test failed")
	}
	content, lastChar, pos, err = ReadUntil("file1.go", []string{"\n", "\r"})
	if err == nil {
		t.Fatalf("file.ReadUntil Test failed")
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
		t.Fatalf("file.ReadBlocks Test failed")
	}
}

func Test_Each(t *testing.T) {
	err := Each("..", true, func(filename, extension, filepath string, dir bool, fileinfo os.FileInfo) {
		if extension == "go" && !dir {
			log.Printf("%v, %v, %v, %v\n", filename, filepath, dir, fileinfo)
		}
	})
	if err != nil {
		t.Fatalf("file.Each Test failed")
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

func Test_Delete(t *testing.T) {
	err := Delete("writetest2.log")
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
	x = IsFile("test.txt")
	if x != true {
		t.Fatalf("file.IsFile Test failed")
	}
	x = IsFile("😃")
	if x == true {
		t.Fatalf("file.IsFile Test failed")
	}
	x = IsSymlink("test.txt")
	if x != false {
		t.Fatalf("file.IsSymlink1 Test failed")
	}
	x = IsSymlink("😃")
	if x != false {
		t.Fatalf("file.IsSymlink2 Test failed")
	}
}
