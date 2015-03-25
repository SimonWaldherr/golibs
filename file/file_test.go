package file

import (
	"log"
	"os"
	"strings"
	"testing"
)

func Test_ReadUntil(t *testing.T) {
	content, lastChar, pos, err := ReadUntil("file.go", []string{"\n", "\r"})
	if err != nil || content != "package file" || lastChar != "\n" || pos != 12 {
		t.Fatalf("file.Each Test failed")
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
		t.Fatalf("file.Each Test failed")
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
		t.Fatalf("file.X Test failed")
	}
}
