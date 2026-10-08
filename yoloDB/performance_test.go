package yolodb

import "testing"

var benchmarkRows []map[string]interface{}

// Select must not spend quadratic work building an index that is never used.
func BenchmarkSelect(b *testing.B) {
	db := NewDatabase()
	if err := db.CreateTable("users", []*Column{{Name: "id", Typ: "int"}, {Name: "name", Typ: "string"}}); err != nil {
		b.Fatal(err)
	}
	for i := 0; i < 1000; i++ {
		if err := db.Insert("users", map[string]interface{}{"id": i, "name": "user"}); err != nil {
			b.Fatal(err)
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := db.Select("users", []string{"id"}, []string{"id > 990"})
		if err != nil || len(rows) != 9 {
			b.Fatalf("Select: %v, %d rows", err, len(rows))
		}
		benchmarkRows = rows
	}
}
