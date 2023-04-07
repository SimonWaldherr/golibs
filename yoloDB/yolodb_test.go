package yolodb

import (
	"testing"
)

func Test_NewDatabase(t *testing.T) {
	db := NewDatabase()
	if db == nil {
		t.Fatalf("NewDatabase() returned nil")
	}
}

func Test_CreateTable(t *testing.T) {
	db := NewDatabase()
	columns := []*Column{
		{Name: "id", Typ: "int"},
		{Name: "name", Typ: "string"},
		{Name: "age", Typ: "int"},
	}
	err := db.CreateTable("users", columns)
	if err != nil {
		t.Fatalf("CreateTable() failed: %v", err)
	}
}
