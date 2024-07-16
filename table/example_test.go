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
    // | Name | Age |      Email       |
    // +------+-----+------------------+
    // | John |  42 | john@example.com |
    // | Jane |  32 | jane@example.com |
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
    // +------+------------------+
    // | Name | John             |
    // +------+------------------+
    // | Age  | 42               |
    // +------+------------------+
    // | Email| john@example.com |
    // +------+------------------+
    // | Name | Jane             |
    // +------+------------------+
    // | Age  | 32               |
    // +------+------------------+
    // | Email| jane@example.com |
    // +------+------------------+
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
    // | Name | Age | Email |
    // | --- | --- | --- |
    // | John | 42 | john@example.com |
    // | Jane | 32 | jane@example.com |
}

// Beispiel: Markdown-Tabelle mit Rotation
func ExampleRenderMarkdown_withRotation() {
    people := []Person{
        {Name: "John", Age: 42, Email: "john@example.com"},
        {Name: "Jane", Age: 32, Email: "jane@example.com"},
    }

    rotatedMarkdownTable, _ := RenderMarkdown(people, TableOption{Rotate: true})
    fmt.Println(rotatedMarkdownTable)
    // Output:
    // | Name | John | Jane |
    // | --- | --- | --- |
    // | Age | 42 | 32 |
    // | Email | john@example.com | jane@example.com |
}

// Beispiel: ASCII-Tabelle für eine andere Struct
func ExampleRenderASCII_products() {
    products := []Product{
        {ID: 1, Name: "Laptop", Price: 999.99},
        {ID: 2, Name: "Smartphone", Price: 499.49},
    }

    productASCIITable, _ := RenderASCII(products, TableOption{Rotate: false})
    fmt.Println(productASCIITable)
    // Output:
    // +----+------------+--------+
    // | ID |    Name    | Price  |
    // +----+------------+--------+
    // |  1 |   Laptop   | 999.99 |
    // |  2 | Smartphone | 499.49 |
    // +----+------------+--------+
}

// Beispiel: JSON-Ausgabe
func ExampleConvertToJSON() {
    people := []Person{
        {Name: "John", Age: 42, Email: "john@example.com"},
        {Name: "Jane", Age: 32, Email: "jane@example.com"},
    }

    peopleJSON, _ := ConvertToJSON(people)
    fmt.Println(peopleJSON)
    // Output:
    // [{"Name":"John","Age":42,"Email":"john@example.com"},{"Name":"Jane","Age":32,"Email":"jane@example.com"}]
}

// Beispiel: XML-Ausgabe
func ExampleConvertToXML() {
    people := []Person{
        {Name: "John", Age: 42, Email: "john@example.com"},
        {Name: "Jane", Age: 32, Email: "jane@example.com"},
    }

    peopleXML, _ := ConvertToXML(people)
    fmt.Println(peopleXML)
    // Output:
    // <[]Person><Person><Name>John</Name><Age>42</Age><Email>john@example.com</Email></Person><Person><Name>Jane</Name><Age>32</Age><Email>jane@example.com</Email></Person></[]Person>
}
