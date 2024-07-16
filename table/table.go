package table

import (
    "bytes"
    "fmt"
    "reflect"
    "strings"
)

// TableOption defines the configuration for table output
type TableOption struct {
    Rotate bool
}

// RenderASCII renders the struct slice as an ASCII table
func RenderASCII(data interface{}, opt TableOption) (string, error) {
    var buffer bytes.Buffer
    headers, rows := parseStruct(data)

    if opt.Rotate {
        headers, rows = rotate(headers, rows)
    }

    colWidths := calculateColumnWidths(headers, rows)

    // Create header line
    buffer.WriteString("+")
    for _, width := range colWidths {
        buffer.WriteString(strings.Repeat("-", width+2) + "+")
    }
    buffer.WriteString("\n|")
    for i, header := range headers {
        buffer.WriteString(fmt.Sprintf(" %-*s |", colWidths[i], header))
    }
    buffer.WriteString("\n+")
    for _, width := range colWidths {
        buffer.WriteString(strings.Repeat("-", width+2) + "+")
    }
    buffer.WriteString("\n")

    // Create rows
    for _, row := range rows {
        buffer.WriteString("|")
        for i, col := range row {
            buffer.WriteString(fmt.Sprintf(" %-*s |", colWidths[i], col))
        }
        buffer.WriteString("\n")
    }
    buffer.WriteString("+")
    for _, width := range colWidths {
        buffer.WriteString(strings.Repeat("-", width+2) + "+")
    }
    buffer.WriteString("\n")

    return buffer.String(), nil
}

// RenderMarkdown renders the struct slice as a Markdown table
func RenderMarkdown(data interface{}, opt TableOption) (string, error) {
    var buffer bytes.Buffer
    headers, rows := parseStruct(data)

    if opt.Rotate {
        headers, rows = rotate(headers, rows)
    }

    colWidths := calculateColumnWidths(headers, rows)

    // Create header line
    buffer.WriteString("|")
    for i, header := range headers {
        buffer.WriteString(fmt.Sprintf(" %-*s |", colWidths[i], header))
    }
    buffer.WriteString("\n|")
    for i := range headers {
        buffer.WriteString(strings.Repeat("-", colWidths[i]+2) + "|")
    }
    buffer.WriteString("\n")

    // Create rows
    for _, row := range rows {
        buffer.WriteString("|")
        for i, col := range row {
            buffer.WriteString(fmt.Sprintf(" %-*s |", colWidths[i], col))
        }
        buffer.WriteString("\n")
    }

    return buffer.String(), nil
}

func parseStruct(data interface{}) ([]string, [][]string) {
    v := reflect.ValueOf(data)
    if v.Kind() != reflect.Slice {
        return nil, nil
    }

    var headers []string
    var rows [][]string

    for i := 0; i < v.Len(); i++ {
        elem := v.Index(i)
        if elem.Kind() != reflect.Struct {
            continue
        }

        var row []string
        elemType := elem.Type()
        for j := 0; j < elem.NumField(); j++ {
            if i == 0 {
                headers = append(headers, elemType.Field(j).Name)
            }
            row = append(row, fmt.Sprintf("%v", elem.Field(j).Interface()))
        }
        rows = append(rows, row)
    }

    return headers, rows
}

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
