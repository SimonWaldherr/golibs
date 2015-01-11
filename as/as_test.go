package as

import "testing"

func Test_Bool(t *testing.T) {
	if Bool("t") != true {
		t.Fatalf("Bool Test failed")
	}
	if Bool("true") != true {
		t.Fatalf("Bool Test failed")
	}
	if Bool(1) != true {
		t.Fatalf("Bool Test failed")
	}
	if Bool("1") != true {
		t.Fatalf("Bool Test failed")
	}
	if Bool("FALSE") != false {
		t.Fatalf("Bool Test failed")
	}
}

func Test_Bytes(t *testing.T) {
	if String([]byte{'l', 'o', 'r', 'e', 'm'}) != String(Bytes("lorem")) {
		t.Fatalf("Bytes Test failed")
	}
	if String([]byte("true")) != String(Bytes(true)) {
		t.Fatalf("Bytes Test failed")
	}
}

func Test_Duration(t *testing.T) {
	if String(Duration(42)) != "42ns" {
		t.Fatalf("Duration Test failed")
	}
}

func Test_Float(t *testing.T) {
	if Float("13.37") != float64(13.37) {
		t.Fatalf("Float Test failed")
	}
	if Float(false) != float64(0) {
		t.Fatalf("Float Test failed")
	}
	if Float(int(42)) != float64(42) {
		t.Fatalf("Float Test failed")
	}
	if Float(int8(43)) != float64(43) {
		t.Fatalf("Float Test failed")
	}
	if Float(int16(44)) != float64(44) {
		t.Fatalf("Float Test failed")
	}
	if Float(int32(45)) != float64(45) {
		t.Fatalf("Float Test failed")
	}
	if Float(int64(4500000)) != float64(4500000) {
		t.Fatalf("Float Test failed")
	}
	if Float(Time("01.01.1970")) != 0 {
		t.Fatalf("Float Test failed")
	}
}

func Test_FloatFromXString(t *testing.T) {
	if FloatFromXString("13.000.000,00") != float64(13000000) {
		t.Fatalf("FloatFromXString Test failed")
	}
	if FloatFromXString("13.000,00") != float64(13000) {
		t.Fatalf("FloatFromXString Test failed")
	}
	if FloatFromXString("13,000,000.00") != float64(13000000) {
		t.Fatalf("FloatFromXString Test failed")
	}
	if FloatFromXString("13,00") != float64(13) {
		t.Fatalf("FloatFromXString Test failed")
	}
}

func Test_Int(t *testing.T) {
	if Int(32.23) != int64(32) {
		t.Fatalf("Int Test failed")
	}
	if Int("32.73") != int64(33) {
		t.Fatalf("Int Test failed")
	}
	if Int(true) != int64(1) {
		t.Fatalf("Int Test failed")
	}
	if Int(Time("01.01.1970")) != int64(0) {
		t.Fatalf("Int Test failed")
	}
}

func Test_String(t *testing.T) {
	if String(1) != "1" {
		t.Fatalf("String Test failed")
	}
	if String(true) != "true" {
		t.Fatalf("String Test failed")
	}
	if String([]byte{'i', 'p', 's', 'u', 'm'}) != "ipsum" {
		t.Fatalf("String Test failed")
	}
	if String(42.0001) != "42.0001" {
		t.Fatalf("String Test failed")
	}
	if String(false) != "false" {
		t.Fatalf("String Test failed")
	}
}

func Test_Time(t *testing.T) {
	if String(Time("11.01.2015")) != "2015-01-11T00:00:00Z" {
		t.Fatalf("Time Test failed")
	}
}

func Test_Uint(t *testing.T) {
	if Uint(32.23) != uint64(32) {
		t.Fatalf("Int Test failed")
	}
	if Uint("32.73") != uint64(33) {
		t.Fatalf("Int Test failed")
	}
	if Uint(true) != uint64(1) {
		t.Fatalf("Int Test failed")
	}
	if Uint(Time("01.01.1970")) != uint64(0) {
		t.Fatalf("Int Test failed")
	}
}


