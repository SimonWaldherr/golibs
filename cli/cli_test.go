package cli

import "testing"
import "os/exec"

func Test_Color(t *testing.T) {
	if Color("X", Red) != "\033[31mX\033[0m" {
		t.Fatalf("Color Test failed")
	}
	if Color("X", Yellow) != "\033[33mX\033[0m" {
		t.Fatalf("Color Test failed")
	}
	if Color("X", Green) != "\033[32mX\033[0m" {
		t.Fatalf("Color Test failed")
	}
	if Color("X", Blue) != "\033[34mX\033[0m" {
		t.Fatalf("Color Test failed")
	}
	if Color("X", Magenta) != "\033[35mX\033[0m" {
		t.Fatalf("Color Test failed")
	}
}

func Test_Bold(t *testing.T) {
	if Bold("X") != "\033[1mX\033[0m" {
		t.Fatalf("Bold Test failed")
	}
}

func Test_Underline(t *testing.T) {
	if Underline("X") != "\033[4mX\033[0m" {
		t.Fatalf("Underline Test failed")
	}
}
