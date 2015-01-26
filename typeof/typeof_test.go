package typeof

import (
	"testing"
)

func println(t *testing.T, str, type1, type2 string) {
	t1 := Type(str)
	t2 := DBType(str)

	if t1 != type1 {
		t.Fatalf("Type Test for %v failed", str)
	}
	if t2 != type2 {
		t.Fatalf("DBType Test for %v failed", str)
	}
}

func Test_Type(t *testing.T) {
	println(t, "true", "bool", "bool")
	println(t, "1000", "int", "int")
	println(t, "2^16", "int", "int")
	println(t, "10E24", "int", "int")
	println(t, "10 EUR", "price", "string")
	println(t, "10.23 £", "price", "string")
	println(t, "mail@example.tld", "email", "string")
	println(t, "978-3499626005", "isbn", "string")
	println(t, "127.0.0.1", "ipv4", "string")
	println(t, "2a01:4f8:192:34aa::10", "ipv6", "string")
	println(t, "b8:c7:5d:c6:6c:c6", "mac", "string")
	println(t, "#fefefe", "color", "string")
}
