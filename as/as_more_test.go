package as

import (
	"reflect"
	"testing"
)

func TestHexAndRGB(t *testing.T) {
	tests := []struct {
		name                string
		hex                 string
		wantR, wantG, wantB int
		wantErr             bool
	}{
		{"valid1", "#ff5733", 255, 87, 51, false},
		{"valid2", "00ff00", 0, 255, 0, false},
		{"invalid_len", "#abc", 0, 0, 0, true},
	}

	for _, tc := range tests {
		r, g, b, err := HexToRGB(tc.hex)
		if (err != nil) != tc.wantErr {
			t.Fatalf("%s: unexpected error state: %v", tc.name, err)
		}
		if err == nil {
			if r != tc.wantR || g != tc.wantG || b != tc.wantB {
				t.Fatalf("%s: got R:%d G:%d B:%d, want R:%d G:%d B:%d", tc.name, r, g, b, tc.wantR, tc.wantG, tc.wantB)
			}
		}
	}
}

func TestRGBToHex(t *testing.T) {
	if RGBToHex(255, 87, 51) != "#ff5733" {
		t.Fatalf("RGBToHex failed")
	}
}

func TestMapToStructAndStructToMap(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	m := map[string]interface{}{"name": "Alice", "age": 30}
	var p Person
	if err := MapToStruct(m, &p); err != nil {
		t.Fatalf("MapToStruct error: %v", err)
	}
	if p.Name != "Alice" || p.Age != 30 {
		t.Fatalf("MapToStruct produced wrong struct: %+v", p)
	}

	// type mismatch yields an error
	m2 := map[string]interface{}{"name": "Bob", "age": "25"}
	var p2 Person
	if err := MapToStruct(m2, &p2); err == nil {
		t.Fatalf("MapToStruct should have returned error on type mismatch")
	}

	// StructToMap
	p3 := Person{Name: "Carol", Age: 22}
	mm := StructToMap(p3)
	want := map[string]interface{}{"Name": "Carol", "Age": 22}
	if !reflect.DeepEqual(mm, want) {
		t.Fatalf("StructToMap failed: got %v want %v", mm, want)
	}
}

func TestValidators(t *testing.T) {
	emails := map[string]bool{
		"test@example.com": true,
		"invalid-email":    false,
		"name@domain":      false,
	}
	for e, want := range emails {
		if IsEmailValid(e) != want {
			t.Fatalf("IsEmailValid(%s) = %v, want %v", e, IsEmailValid(e), want)
		}
	}

	urls := map[string]bool{
		"http://example.com":     true,
		"https://www.google.com": true,
		"ftp://example.com":      false,
		"invalid":                false,
	}
	for u, want := range urls {
		if IsURLValid(u) != want {
			t.Fatalf("IsURLValid(%s) = %v, want %v", u, IsURLValid(u), want)
		}
	}

	cards := map[string]bool{
		"4111111111111111": true,
		"5500000000000004": true,
		"1234567890123456": false,
		"378282246310005":  true,
	}
	for c, want := range cards {
		if IsCreditCardValid(c) != want {
			t.Fatalf("IsCreditCardValid(%s) = %v, want %v", c, IsCreditCardValid(c), want)
		}
	}
}

func TestJSONHelpers(t *testing.T) {
	v := map[string]interface{}{"key": "value"}
	s, err := ToJSON(v)
	if err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := FromJSON(s, &out); err != nil {
		t.Fatalf("FromJSON error: %v", err)
	}
	if !reflect.DeepEqual(v, out) {
		t.Fatalf("JSON roundtrip mismatch: %v vs %v", v, out)
	}

	// invalid JSON
	if err := FromJSON("{invalid}", &out); err == nil {
		t.Fatalf("FromJSON should error on invalid input")
	}
}
