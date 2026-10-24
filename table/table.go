// Package table renders struct slices and parses JSON or CSV table data.
package table

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
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

// RenderUnicode renders the struct slice using Unicode borders.
func RenderUnicode(data interface{}, opt TableOption) (string, error) {
	return Render(data, TableOption{Rotate: opt.Rotate, Align: opt.Align, Format: "unicode"})
}

// RenderMarkdown renders the struct slice as a Markdown table
func RenderMarkdown(data interface{}, opt TableOption) (string, error) {
	return Render(data, TableOption{Rotate: opt.Rotate, Align: opt.Align, Format: "markdown"})
}

// RenderHTML renders the struct slice as an HTML table, escaping headers and
// cell values as text.
func RenderHTML(data interface{}, opt TableOption) (string, error) {
	return Render(data, TableOption{Rotate: opt.Rotate, Align: opt.Align, Format: "html"})
}

// Render formats a nonempty slice of structs. Struct fields must be exported;
// unsupported inputs return errors rather than exposing private fields.
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
	var buffer strings.Builder
	colWidths := calculateColumnWidths(headers, rows)

	buffer.WriteString("|")
	for i, header := range headers {
		buffer.WriteByte(' ')
		buffer.WriteString(alignText(header, colWidths[i], "left"))
		buffer.WriteString(" |")
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
			buffer.WriteByte(' ')
			buffer.WriteString(alignText(col, colWidths[i], alignment))
			buffer.WriteString(" |")
		}
		buffer.WriteString("\n")
	}

	return buffer.String(), nil
}

func renderHTML(headers []string, rows [][]string) (string, error) {
	var buffer strings.Builder

	buffer.WriteString("<table>\n<thead>\n<tr>")
	for _, header := range headers {
		buffer.WriteString("<th>")
		buffer.WriteString(html.EscapeString(header))
		buffer.WriteString("</th>")
	}
	buffer.WriteString("</tr>\n</thead>\n<tbody>\n")

	for _, row := range rows {
		buffer.WriteString("<tr>")
		for _, col := range row {
			buffer.WriteString("<td>")
			buffer.WriteString(html.EscapeString(col))
			buffer.WriteString("</td>")
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

// ParseCSVInput appends CSV records to the slice of structs pointed to by v.
// Headers name fields exactly; unknown columns are ignored. Destinations must be
// nonnil pointers to slices of struct values. Known columns must be exported,
// unique, and settable. Numeric values must fit the destination field's type.
// Validation and parsing errors leave the destination unchanged. The reader is
// left open. Records are streamed rather than retaining a second copy of all CSV rows.
func ParseCSVInput(r io.Reader, v interface{}) error {
	if r == nil {
		return fmt.Errorf("table: CSV reader is nil")
	}
	destination := reflect.ValueOf(v)
	if destination.Kind() != reflect.Ptr || destination.IsNil() || destination.Elem().Kind() != reflect.Slice || destination.Elem().Type().Elem().Kind() != reflect.Struct {
		return fmt.Errorf("table: destination must be a nonnil pointer to a slice of structs")
	}
	structSlice := destination.Elem()
	elemType := structSlice.Type().Elem()
	csvReader := csv.NewReader(r)
	csvReader.ReuseRecord = true
	headers, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("table: read CSV header: %w", err)
	}
	type column struct {
		name  string
		index []int
	}
	columns := make([]column, len(headers))
	seen := make(map[string]bool, len(headers))
	prototype := reflect.New(elemType).Elem()
	for i, name := range headers {
		field, found := elemType.FieldByName(name)
		if !found {
			continue
		}
		if field.PkgPath != "" {
			return fmt.Errorf("table: CSV column %q is not exported", name)
		}
		if seen[name] {
			return fmt.Errorf("table: duplicate CSV column %q", name)
		}
		seen[name] = true
		target, err := prototype.FieldByIndexErr(field.Index)
		if err != nil {
			return fmt.Errorf("table: CSV column %q: %w", name, err)
		}
		if !target.CanSet() {
			return fmt.Errorf("table: CSV column %q is not settable", name)
		}
		if !supportedCSVType(field.Type) {
			return fmt.Errorf("table: CSV column %q has unsupported field type %s", name, field.Type)
		}
		columns[i] = column{name: name, index: field.Index}
	}
	var parsed reflect.Value
	for row := 2; ; row++ {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("table: CSV row %d: %w", row, err)
		}
		elem := reflect.New(elemType).Elem()
		for i, column := range columns {
			if column.index == nil {
				continue
			}
			field := elem.FieldByIndex(column.index)
			if err := setFieldValue(field, record[i]); err != nil {
				return fmt.Errorf("table: CSV row %d column %q: %w", row, column.name, err)
			}
		}
		if !parsed.IsValid() {
			parsed = reflect.MakeSlice(structSlice.Type(), 0, 16)
		}
		parsed = reflect.Append(parsed, elem)
	}
	if parsed.IsValid() {
		structSlice.Set(reflect.AppendSlice(structSlice, parsed))
	}
	return nil
}

func supportedCSVType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return true
	case reflect.Struct:
		return typ == reflect.TypeOf(time.Time{})
	}
	return false
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, field.Type().Bits())
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
		} else {
			return fmt.Errorf("unsupported field type: %s", field.Type())
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

	if err := csvWriter.Write(headers); err != nil {
		return "", err
	}

	for _, row := range rows {
		if err := csvWriter.Write(row); err != nil {
			return "", err
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// parseStruct extracts headers and rows from a slice of exported struct fields.
func parseStruct(data interface{}) ([]string, [][]string, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil, nil, fmt.Errorf("data is not a slice")
	}
	if v.Len() == 0 {
		return nil, nil, fmt.Errorf("slice is empty")
	}
	elemType := v.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("element is not a struct")
	}
	headers := make([]string, elemType.NumField())
	for i := range headers {
		field := elemType.Field(i)
		if field.PkgPath != "" {
			return nil, nil, fmt.Errorf("field %q is not exported", field.Name)
		}
		headers[i] = field.Name
	}
	rows := make([][]string, v.Len())
	for i := range rows {
		elem := v.Index(i)
		row := make([]string, len(headers))
		for j := range row {
			row[j] = formatField(elem.Field(j).Interface())
		}
		rows[i] = row
	}
	return headers, rows, nil
}

// formatField formats field values, especially date fields according to ISO 8601
func formatField(field interface{}) string {
	switch v := field.(type) {
	case time.Time:
		return v.Format(time.RFC3339)
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
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

	// Complete numeric values avoid the scanner's per-cell allocations. The
	// scanner remains the fallback for partial numbers and its other legacy syntax.
	if len(value) > 0 && (value[0] >= '0' && value[0] <= '9' || value[0] == '+' || value[0] == '-' || value[0] == '.') {
		if _, err := strconv.ParseFloat(value, 64); err == nil {
			return "right"
		}
	}
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

// httpHeader is an internal struct used by RenderHTTPResponse.
type httpHeader struct {
	Key   string
	Value string
}

// RenderHTTPResponse renders the headers of an HTTP response as a table using the given options.
func RenderHTTPResponse(resp *http.Response, opt TableOption) (string, error) {
	var headers []httpHeader
	for key, values := range resp.Header {
		for _, value := range values {
			headers = append(headers, httpHeader{Key: key, Value: value})
		}
	}
	return Render(headers, opt)
}
