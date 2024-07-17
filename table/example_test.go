package table

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
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
	// +-------+------------------+------------------+
	// |       | Row 1            | Row 2            |
	// +-------+------------------+------------------+
	// | Name  | John             | Jane             |
	// | Age   |               42 |               32 |
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
	// | John |  42 | john@example.com |
	// | Jane |  32 | jane@example.com |
}

type TestStruct struct {
	Name      string
	Age       int
	Salary    float64
	StartDate time.Time
}

func ExampleParseCSVInput() {
	csvData := `Name,Age,Salary,StartDate
Alice,30,70000.50,2021-05-01T00:00:00Z
Bob,22,48000.00,2022-07-15T00:00:00Z`
	var data []TestStruct
	err := ParseCSVInput(strings.NewReader(csvData), &data)
	if err != nil {
		fmt.Println(err)
	}
	markdownTable, _ := RenderMarkdown(data, TableOption{Rotate: false})
	fmt.Println(markdownTable)

	// Output:
	// | Name  | Age | Salary  | StartDate            |
	// |-------|-----|---------|----------------------|
	// | Alice |  30 | 70000.5 | 2021-05-01T00:00:00Z |
	// | Bob   |  22 |   48000 | 2022-07-15T00:00:00Z |
}

// Header represents a single HTTP header key-value pair
type Header struct {
	Key   string
	Value string
}

func firstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Worldn"))
}

// ExampleRenderHTTPResponse demonstrates how to render HTTP response headers
func ExampleRenderHTTPResponse() {
	req := httptest.NewRequest("GET", "https://github.com", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()

	defer resp.Body.Close()

	// Convert headers to a slice of Header structs
	var headers []Header
	for key, values := range resp.Header {
		for _, value := range values {
			headers = append(headers, Header{Key: key, Value: firstN(value, 30)})
		}
	}

	// Render the headers as a Markdown table
	markdownTable, err := RenderMarkdown(headers, TableOption{Rotate: false})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(markdownTable)
	// Output:
	// | Key          | Value                     |
	// |--------------|---------------------------|
	// | Content-Type | text/plain; charset=utf-8 |
}
