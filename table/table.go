package table

import (
    "bytes"
    "encoding/json"
    "encoding/xml"
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

    buffer.WriteString("+")
    for _, header := range headers {
        buffer.WriteString(strings.Repeat("-", len(header)+2) + "+")
    }
    buffer.WriteString("\n|")
    for _, header := range headers {
        buffer.WriteString(fmt.Sprintf(" %s |", header))
    }
    buffer.WriteString("\n+")
    for range headers {
        buffer.WriteString(strings.Repeat("-", len(headers[0])+2) + "+")
    }
    buffer.WriteString("\n")

    for _, row := range rows {
        buffer.WriteString("|")
        for _, col := range row {
            buffer.WriteString(fmt.Sprintf(" %s |", col))
        }
        buffer.WriteString("\n")
    }
    buffer.WriteString("+")
    for range headers {
        buffer.WriteString(strings.Repeat("-", len(headers[0])+2) + "+")
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

    buffer.WriteString("|")
    for _, header := range headers {
        buffer.WriteString(fmt.Sprintf(" %s |", header))
    }
    buffer.WriteString("\n|")
    for range headers {
        buffer.WriteString(" --- |")
    }
    buffer.WriteString("\n")

    for _, row := range rows {
        buffer.WriteString("|")
        for _, col := range row {
            buffer.WriteString(fmt.Sprintf(" %s |", col))
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
    var newHeaders []string
    var newRows [][]string

    for _, header := range headers {
        newHeaders = append(newHeaders, header)
    }

    for i := 0; i < len(rows[0]); i++ {
        var newRow []string
        for j := 0; j < len(rows); j++ {
            newRow = append(newRow, rows[j][i])
        }
        newRows = append(newRows, newRow)
    }

    return newHeaders, newRows
}

func ConvertToJSON(data interface{}) (string, error) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    return string(jsonData), nil
}

func ConvertToXML(data interface{}) (string, error) {
    xmlData, err := xml.Marshal(data)
    if err != nil {
        return "", err
    }
    return string(xmlData), nil
}
