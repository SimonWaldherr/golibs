package yolodb

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func userDatabase(t *testing.T) *Database {
	t.Helper()
	db := NewDatabase()
	if err := db.CreateTable("users", []*Column{{Name: "id", Typ: "int"}, {Name: "name", Typ: "string"}}); err != nil {
		t.Fatal(err)
	}
	return db
}

func user(id int, name string) map[string]interface{} {
	return map[string]interface{}{"id": id, "name": name}
}

func TestCRUD(t *testing.T) {
	db := userDatabase(t)
	for i, name := range []string{"Alice", "Bob", "Bob", "Diana"} {
		if err := db.Insert("users", user(i+1, name)); err != nil {
			t.Fatal(err)
		}
	}
	cases := []struct {
		where []string
		ids   []int
	}{
		{nil, []int{1, 2, 3, 4}},
		{[]string{"name = Bob"}, []int{2, 3}},
		{[]string{"id > 1", "id < 4"}, []int{2, 3}},
		{[]string{"id = 2", "id = 3"}, nil},
		{[]string{"  id   =\t4  "}, []int{4}},
		{[]string{"name = Nobody"}, nil},
	}
	for _, tc := range cases {
		rows, err := db.Select("users", []string{"id", "id"}, tc.where)
		if err != nil {
			t.Fatal(err)
		}
		var got []int
		for _, row := range rows {
			got = append(got, row["id"].(int))
			if len(row) != 1 {
				t.Fatalf("projection = %v", row)
			}
		}
		if !reflect.DeepEqual(got, tc.ids) {
			t.Errorf("Select(%v) = %v, want %v", tc.where, got, tc.ids)
		}
	}
	rows, err := db.Select("users", nil, nil)
	if err != nil || len(rows) != 4 || len(rows[0]) != 0 {
		t.Fatalf("empty projection: %v, %v", rows, err)
	}
	rows, err = db.Select("users", []string{"id", "name"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	rows[0]["name"] = "changed externally"
	if err := db.Update("users", user(8, "Updated"), []string{"id = 2"}); err != nil {
		t.Fatal(err)
	}
	if err := db.Delete("users", []string{"id = 3"}); err != nil {
		t.Fatal(err)
	}
	rows, err = db.Select("users", []string{"id", "name"}, nil)
	want := []map[string]interface{}{user(1, "Alice"), user(8, "Updated"), user(4, "Diana")}
	if err != nil || !reflect.DeepEqual(rows, want) {
		t.Fatalf("after update/delete: %v, %v; want %v", rows, err, want)
	}
	if err := db.Update("users", user(9, "Everyone"), nil); err != nil {
		t.Fatal(err)
	}
	if err := db.Delete("users", []string{"id = 100"}); err != nil {
		t.Fatal(err)
	}
	rows, err = db.Select("users", []string{"id"}, []string{"name = Everyone"})
	if err != nil || len(rows) != 3 {
		t.Fatalf("update all: %v, %v", rows, err)
	}
	if err := db.Delete("users", nil); err != nil {
		t.Fatal(err)
	}
	rows, err = db.Select("users", []string{"id"}, nil)
	if err != nil || rows != nil {
		t.Fatalf("delete all: %v, %v", rows, err)
	}
}

func TestDuplicateAndUpdatedRows(t *testing.T) {
	db := userDatabase(t)
	for _, id := range []int{1, 2} {
		if err := db.Insert("users", user(id, "shared")); err != nil {
			t.Fatal(err)
		}
	}
	if err := db.Delete("users", []string{"id = 1"}); err != nil {
		t.Fatal(err)
	}
	rows, err := db.Select("users", []string{"id"}, []string{"name = shared"})
	if err != nil || len(rows) != 1 || rows[0]["id"] != 2 {
		t.Fatalf("deleting a duplicate lost the survivor: %v, %v", rows, err)
	}
	if err := db.Update("users", user(7, "updated"), []string{"id = 2"}); err != nil {
		t.Fatal(err)
	}
	if err := db.Delete("users", []string{"id = 7"}); err != nil {
		t.Fatal(err)
	}
	// Querying an old value after deleting an updated row must not follow stale
	// index entries into a missing record.
	rows, err = db.Select("users", []string{"id"}, []string{"id = 2"})
	if err != nil || rows != nil {
		t.Fatalf("query after update/delete: %v, %v", rows, err)
	}
}

func TestSchemaAndRowValidation(t *testing.T) {
	cases := [][]*Column{
		{nil}, {{Name: "", Typ: "int"}}, {{Name: "id", Typ: "float64"}},
		{{Name: "id", Typ: "int"}, {Name: "id", Typ: "string"}},
	}
	for _, schema := range cases {
		db := NewDatabase()
		if err := db.CreateTable("invalid", schema); err == nil {
			t.Errorf("accepted schema %v", schema)
		}
	}
	var db Database // zero value
	schema := []*Column{{Name: "id", Typ: "int"}, {Name: "name", Typ: "string"}}
	if err := db.CreateTable("users", schema); err != nil {
		t.Fatal(err)
	}
	schema[0].Name, schema[0].Typ, schema[1] = "changed", "string", nil
	if err := db.CreateTable("users", nil); err == nil {
		t.Fatal("accepted duplicate table")
	}
	if err := db.Insert("users", user(1, "Alice")); err != nil {
		t.Fatalf("schema was aliased: %v", err)
	}
	type namedInt int
	for _, row := range []map[string]interface{}{
		nil, {"id": 1}, {"id": 1, "extra": "x"},
		{"id": nil, "name": "Alice"}, {"id": namedInt(1), "name": "Alice"},
		{"id": 1, "name": 2}, {"id": 1, "name": "Alice", "extra": 2},
	} {
		if err := db.Insert("users", row); err == nil {
			t.Errorf("Insert accepted invalid row: %v", row)
		}
		if err := db.Update("users", row, nil); err == nil {
			t.Errorf("Update accepted invalid row: %v", row)
		}
	}
	rows, err := db.Select("users", []string{"id", "name"}, nil)
	if err != nil || !reflect.DeepEqual(rows, []map[string]interface{}{user(1, "Alice")}) {
		t.Fatalf("failed writes changed records: %v, %v", rows, err)
	}
	if err := db.CreateTable("empty_schema", nil); err != nil {
		t.Fatal(err)
	}
	if err := db.Insert("empty_schema", nil); err != nil {
		t.Fatal(err)
	}
}

func TestQueryValidation(t *testing.T) {
	db := userDatabase(t)
	// Validate predicates even before records have been inserted.
	for _, where := range [][]string{
		{""}, {"id"}, {"id ="}, {"id = 1 extra"}, {"missing = 1"},
		{"id != 1"}, {"id == 1"}, {"id = invalid"}, {"id = 99999999999999999999999999"},
		{"name > Alice"}, {"name < Alice"}, {"id = 1", "invalid"},
	} {
		if rows, err := db.Select("users", []string{"id"}, where); err == nil || rows != nil {
			t.Errorf("Select accepted %v: %v, %v", where, rows, err)
		}
		if err := db.Update("users", user(1, "Alice"), where); err == nil {
			t.Errorf("Update accepted %v", where)
		}
		if err := db.Delete("users", where); err == nil {
			t.Errorf("Delete accepted %v", where)
		}
	}
	if _, err := db.Select("users", []string{"missing"}, nil); err == nil {
		t.Fatal("accepted invalid projection")
	}
	if err := db.Insert("missing", user(1, "Alice")); err == nil {
		t.Fatal("Insert accepted missing table")
	}
	if _, err := db.Select("missing", nil, nil); err == nil {
		t.Fatal("Select accepted missing table")
	}
	if err := db.Update("missing", nil, nil); err == nil {
		t.Fatal("Update accepted missing table")
	}
	if err := db.Delete("missing", nil); err == nil {
		t.Fatal("Delete accepted missing table")
	}
	if err := db.Insert("users", user(1, "Alice")); err != nil {
		t.Fatal(err)
	}
	if err := db.Update("users", user(2, "changed"), []string{"id = 1", "bad"}); err == nil {
		t.Fatal("accepted invalid update condition")
	}
	if err := db.Delete("users", []string{"bad"}); err == nil {
		t.Fatal("accepted invalid delete condition")
	}
	rows, err := db.Select("users", []string{"id", "name"}, nil)
	if err != nil || !reflect.DeepEqual(rows, []map[string]interface{}{user(1, "Alice")}) {
		t.Fatalf("failed queries changed records: %v, %v", rows, err)
	}
}

func TestConcurrentCRUD(t *testing.T) {
	db := userDatabase(t)
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			if err := db.CreateTable(fmt.Sprintf("table%d", worker), nil); err != nil {
				t.Error(err)
			}
			for i := 0; i < 20; i++ {
				id := worker*20 + i
				where := []string{fmt.Sprintf("id = %d", id)}
				if err := db.Insert("users", user(id, "inserted")); err != nil {
					t.Error(err)
				}
				if err := db.Update("users", user(id, "updated"), where); err != nil {
					t.Error(err)
				}
				if _, err := db.Select("users", []string{"id"}, nil); err != nil {
					t.Error(err)
				}
				if i%2 == 0 {
					if err := db.Delete("users", where); err != nil {
						t.Error(err)
					}
				}
			}
		}(worker)
	}
	wg.Wait()
	rows, err := db.Select("users", []string{"id", "name"}, nil)
	if err != nil || len(rows) != 80 {
		t.Fatalf("concurrent CRUD: %d rows, %v", len(rows), err)
	}
	for _, row := range rows {
		if row["id"].(int)%2 != 1 || row["name"] != "updated" {
			t.Errorf("unexpected row: %v", row)
		}
	}
}
