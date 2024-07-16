package table

import (
    "fmt"
)

// Beispiel-Datenstrukturen
type Person struct {
    Name  string
    Age   int
    Email string
}

type Product struct {
    ID    int
    Name  string
    Price float64
}

// Beispiel: ASCII-Tabelle ohne Rotation
func ExampleRenderASCII_noRotation() {
    people := []Person{
        {Name: "John", Age: 42, Email: "john@example.com"},
        {Name: "Jane", Age: 32, Email: "jane@example.com"},
    }

    asciiTable, _ := RenderASCII(people, TableOption{Rotate: false})
    fmt.Println(asciiTable)
    // Output:
    // +------+-----+------------------+
    // | Name | Age | Email            |
    // +------+-----+------------------+
    // | John | 42  | john@example.com |
    // | Jane | 32  | jane@example.com |
    // +------+-----+------------------+
}

// Beispiel: ASCII-Tabelle mit Rotation
func ExampleRenderASCII_withRotation() {
    people := []Person{
        {Name: "John", Age: 42, Email: "john@example.com"},
        {Name: "Jane", Age: 32, Email: "jane@example.com"},
    }

    rotatedASCIITable, _ := RenderASCII(people, TableOption{Rotate: true})
    fmt.Println(rotatedASCIITable)
    // Output:
    // +-------+------------------+------------------+
    // |       | Row 1            | Row 2            |
    // +-------+------------------+------------------+
    // | Name  | John             | Jane             |
    // | Age   | 42               | 32               |
    // | Email | john@example.com | jane@example.com |
    // +-------+------------------+------------------+
}

// Beispiel: Markdown-Tabelle ohne Rotation
func ExampleRenderMarkdown_noRotation() {
    people := []Person{
        {Name: "John", Age: 42, Email: "john@example.com"},
        {Name: "Jane", Age: 32, Email: "jane@example.com"},
    }

    markdownTable, _ := RenderMarkdown(people, TableOption{Rotate: false})
    fmt.Println(markdownTable)
    // Output:
    // | Name | Age | Email            |
    // |------|-----|------------------|
    // | John | 42  | john@example.com |
    // | Jane | 32  | jane@example.com |
}
