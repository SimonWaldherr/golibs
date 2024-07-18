package table

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// TableOption defines the configuration for table output
type TableOption struct {
	Rotate bool
	Align  []string // can be "left", "right", or "center"
	Format string   // can be "ascii", "unicode", "markdown", or "html"
}

// RenderASCII renders the struct slice as an ASCII table
func RenderASCII(data interface{}, opt TableOption) (string, error) {
	return Render(data, TableOption{Rotate: opt.Rotate, Align: opt.Align, Format: "ascii"})
}

// RenderASCII renders the struct slice as an ASCII table
func RenderUnicode(data interface{}, opt TableOption) (string, error) {
	return Render(data, TableOption{Rotate: opt.Rotate, Align: opt.Align, Format: "unicode"})
}

// RenderMarkdown renders the struct slice as a Markdown table
func RenderMarkdown(data interface{}, opt TableOption) (string, error) {
	return Render(data, TableOption{Rotate: opt.Rotate, Align: opt.Align, Format: "markdown"})
}

// RenderHTML renders the struct slice as an HTML table
func RenderHTML(data interface{}, opt TableOption) (string, error) {
	return Render(data, TableOption{Rotate: opt.Rotate, Align: opt.Align, Format: "html"})
}

func Render(data interface{}, opt TableOption) (string, error) {
	headers, rows, err := parseStruct(data)
	if err != nil {
		return "", err
	}

	if opt.Rotate {
		headers, rows = rotate(headers, rows)
	}

	switch opt.Format {
	case "ascii":
		return renderASCII(headers, rows, opt, basicBorders)
	case "unicode":
		return renderASCII(headers, rows, opt, niceBorders)
	case "markdown":
		return renderMarkdown(headers, rows, opt)
	case "html":
		return renderHTML(headers, rows)
	default:
		return "", fmt.Errorf("unsupported format: %s", opt.Format)
	}
}

type borders struct {
	topLeft, topMid, topRight, midLeft, midMid, midRight, botLeft, botMid, botRight, hor, ver string
}

var basicBorders = borders{"+", "+", "+", "+", "+", "+", "+", "+", "+", "-", "|"}
var niceBorders = borders{"┏", "┳", "┓", "┡", "╇", "┩", "└", "┴", "┘", "━", "┃"}

func renderASCII(headers []string, rows [][]string, opt TableOption, brd borders) (string, error) {
	var buffer bytes.Buffer
	colWidths := calculateColumnWidths(headers, rows)

	buffer.WriteString(brd.topLeft)
	for i, width := range colWidths {
		if i > 0 {
			buffer.WriteString(brd.topMid)
		}
		buffer.WriteString(strings.Repeat(brd.hor, width+2))
	}
	buffer.WriteString(brd.topRight + "\n" + brd.ver)
	for i, header := range headers {
		if i > 0 {
			buffer.WriteString(brd.ver)
		}
		alignment := getColumnAlignment(i, opt.Align, header)
		buffer.WriteString(" " + alignText(header, colWidths[i], alignment) + " ")
	}
	buffer.WriteString(brd.ver + "\n" + brd.midLeft)
	for i, width := range colWidths {
		if i > 0 {
			buffer.WriteString(brd.midMid)
		}
		buffer.WriteString(strings.Repeat(brd.hor, width+2))
	}
	buffer.WriteString(brd.midRight + "\n")

	for _, row := range rows {
		buffer.WriteString(brd.ver)
		for i, col := range row {
			if i > 0 {
				buffer.WriteString(brd.ver)
			}
			alignment := getColumnAlignment(i, opt.Align, col)
			buffer.WriteString(" " + alignText(col, colWidths[i], alignment) + " ")
		}
		buffer.WriteString(brd.ver + "\n")
	}
	buffer.WriteString(brd.botLeft)
	for i, width := range colWidths {
		if i > 0 {
			buffer.WriteString(brd.botMid)
		}
		buffer.WriteString(strings.Repeat(brd.hor, width+2))
	}
	buffer.WriteString(brd.botRight + "\n")

	return buffer.String(), nil
}

func renderMarkdown(headers []string, rows [][]string, opt TableOption) (string, error) {
	var buffer bytes.Buffer
	colWidths := calculateColumnWidths(headers, rows)

	buffer.WriteString("|")
	for i, header := range headers {
		buffer.WriteString(fmt.Sprintf(" %-*s |", colWidths[i], header))
	}
	buffer.WriteString("\n|")
	for i := range headers {
		buffer.WriteString(strings.Repeat("-", colWidths[i]+2) + "|")
	}
	buffer.WriteString("\n")

	for _, row := range rows {
		buffer.WriteString("|")
		for i, col := range row {
			alignment := getColumnAlignment(i, opt.Align, col)
			buffer.WriteString(fmt.Sprintf(" %s |", alignText(col, colWidths[i], alignment)))
		}
		buffer.WriteString("\n")
	}

	return buffer.String(), nil
}

func renderHTML(headers []string, rows [][]string) (string, error) {
	var buffer bytes.Buffer

	buffer.WriteString("<table>\n<thead>\n<tr>")
	for _, header := range headers {
		buffer.WriteString(fmt.Sprintf("<th>%s</th>", html.EscapeString(header)))
	}
	buffer.WriteString("</tr>\n</thead>\n<tbody>\n")

	for _, row := range rows {
		buffer.WriteString("<tr>")
		for _, col := range row {
			buffer.WriteString(fmt.Sprintf("<td>%s</td>", html.EscapeString(col)))
		}
		buffer.WriteString("</tr>\n")
	}
	buffer.WriteString("</tbody>\n</table>\n")

	return buffer.String(), nil
}

// ParseJSONInput parses JSON input into a slice of structs
func ParseJSONInput(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(v)
}

// ParseCSVInput parses CSV input into a slice of structs
func ParseCSVInput(r io.Reader, v interface{}) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	headers := records[0]
	structSlice := reflect.ValueOf(v).Elem()
	elemType := structSlice.Type().Elem()

	for _, record := range records[1:] {
		elem := reflect.New(elemType).Elem()
		for i, header := range headers {
			field := elem.FieldByName(header)
			if field.IsValid() {
				err := setFieldValue(field, record[i])
				if err != nil {
					return err
				}
			}
		}
		structSlice.Set(reflect.Append(structSlice, elem))
	}

	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			timeValue, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(timeValue))
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

// RenderCSV renders the struct slice as a CSV output
func RenderCSV(data interface{}) (string, error) {
	var buffer bytes.Buffer
	headers, rows, err := parseStruct(data)
	if err != nil {
		return "", err
	}

	csvWriter := csv.NewWriter(&buffer)
	defer csvWriter.Flush()

	if err := csvWriter.Write(headers); err != nil {
		return "", err
	}

	for _, row := range rows {
		if err := csvWriter.Write(row); err != nil {
			return "", err
		}
	}

	return buffer.String(), nil
}

// parseStruct extracts headers and rows from a slice of structs
func parseStruct(data interface{}) ([]string, [][]string, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil, nil, fmt.Errorf("data is not a slice")
	}

	if v.Len() == 0 {
		return nil, nil, fmt.Errorf("slice is empty")
	}

	var headers []string
	var rows [][]string

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() != reflect.Struct {
			return nil, nil, fmt.Errorf("element is not a struct")
		}

		var row []string
		elemType := elem.Type()
		for j := 0; j < elem.NumField(); j++ {
			if i == 0 {
				headers = append(headers, elemType.Field(j).Name)
			}
			field := elem.Field(j).Interface()
			row = append(row, formatField(field))
		}
		rows = append(rows, row)
	}

	return headers, rows, nil
}

// formatField formats field values, especially date fields according to ISO 8601
func formatField(field interface{}) string {
	switch v := field.(type) {
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// rotate swaps rows and columns
func rotate(headers []string, rows [][]string) ([]string, [][]string) {
	newHeaders := make([]string, len(rows))
	for i := range rows {
		newHeaders[i] = fmt.Sprintf("Row %d", i+1)
	}

	newRows := make([][]string, len(headers))
	for i := range headers {
		newRow := make([]string, len(rows))
		for j := range rows {
			newRow[j] = rows[j][i]
		}
		newRows[i] = newRow
	}

	// Add original headers to newRows
	for i, header := range headers {
		newRows[i] = append([]string{header}, newRows[i]...)
	}

	// Add an empty column header for the row identifiers
	newHeaders = append([]string{""}, newHeaders...)

	return newHeaders, newRows
}

// calculateColumnWidths calculates the width of each column
func calculateColumnWidths(headers []string, rows [][]string) []int {
	colWidths := make([]int, len(headers))

	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, col := range row {
			if len(col) > colWidths[i] {
				colWidths[i] = len(col)
			}
		}
	}

	return colWidths
}

// getColumnAlignment returns the alignment for a given column index
func getColumnAlignment(index int, align []string, value string) string {
	if index < len(align) {
		return align[index]
	}

	// Default alignment based on value type
	if _, err := fmt.Sscanf(value, "%f", new(float64)); err == nil {
		return "right"
	}

	return "left"
}

// Alignment function to align text based on the specified alignment type
func alignText(text string, width int, alignment string) string {
	switch alignment {
	case "center":
		spaces := width - len(text)
		left := spaces / 2
		right := spaces - left
		return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
	case "right":
		return fmt.Sprintf("%*s", width, text)
	default:
		return fmt.Sprintf("%-*s", width, text)
	}
}
