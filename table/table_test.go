package table

import (
    "fmt"
    "testing"
)

type person struct {
    Name string
    Age  int
}

func TestRenderASCII(t *testing.T) {
    people := []person{
        {Name: "John", Age: 42},
        {Name: "Jane", Age: 32},
    }

    table, err := RenderASCII(people, TableOption{Rotate: false})
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    fmt.Println(table)

    table, err = RenderASCII(people, TableOption{Rotate: true})
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    fmt.Println(table)
}

func TestRenderMarkdown(t *testing.T) {
    people := []person{
        {Name: "John", Age: 42},
        {Name: "Jane", Age: 32},
    }

    table, err := RenderMarkdown(people, TableOption{Rotate: false})
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    fmt.Println(table)

    table, err = RenderMarkdown(people, TableOption{Rotate: true})
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    fmt.Println(table)
}
