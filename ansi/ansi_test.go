package ansi

import (
	"strconv"
	"testing"
)

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
	if Color(3*3, Magenta) != "\033[35m9\033[0m" {
		t.Fatalf("Color Test failed")
	}
}

func Test_BgColor(t *testing.T) {
	if BgColor("X", Red) != "\033[41mX\033[0m" {
		t.Fatalf("BgColor Test failed")
	}
	if BgColor("X", Yellow) != "\033[43mX\033[0m" {
		t.Fatalf("BgColor Test failed")
	}
	if BgColor("X", Green) != "\033[42mX\033[0m" {
		t.Fatalf("BgColor Test failed")
	}
	if BgColor("X", Blue) != "\033[44mX\033[0m" {
		t.Fatalf("BgColor Test failed")
	}
	if BgColor("X", Magenta) != "\033[45mX\033[0m" {
		t.Fatalf("BgColor Test failed")
	}
	if BgColor(3*3, Magenta) != "\033[45m9\033[0m" {
		t.Fatalf("BgColor Test failed")
	}
}

func Test_Bold(t *testing.T) {
	if Bold("X") != "\033[1mX\033[0m" {
		t.Fatalf("Bold Test failed")
	}
	if Bold(33) != "\033[1m33\033[0m" {
		t.Fatalf("Bold Test failed")
	}
}

func Test_Underline(t *testing.T) {
	if Underline("X") != "\033[4mX\033[0m" {
		t.Fatalf("Underline Test failed")
	}
	if Underline(42.23) != "\033[4m42.23\033[0m" {
		t.Fatalf("Underline Test failed")
	}
}

func Test_Log(t *testing.T) {
	if Log(strconv.ParseInt("42", 10, 0)).(int64) != int64(42) {
		t.Fatalf("Log Test failed %v", Log(strconv.ParseInt("42", 10, 0)))
	}
	if Log(strconv.ParseInt("Foobar", 10, 0)).(int64) != int64(0) {
		t.Fatalf("Log Test failed %v", Log(strconv.ParseInt("Foobar", 10, 0)))
	}
}

func Test_Set(t *testing.T) {
	if Set("01234567890", TCrossedOut, FgGreen, BgCyan, TBold) != "\033[1m\033[46m\033[32m\033[9m01234567890\033[0m\033[0m\033[0m\033[0m" {
		t.Fatalf(Set("01234567890", TCrossedOut, FgGreen, BgCyan, TBold))
	}
}
