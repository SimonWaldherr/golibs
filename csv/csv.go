// Package csv reads semicolon-separated CSV data with a header row.
package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// LoadCSVfromFile parses a CSV from a file and returns two maps
// one map with the data and another one with the csv head. It returns the rows
// read before an error; use ReadFile to receive the error as well.
func LoadCSVfromFile(filename string) (map[int][]string, map[string]int) {
	data, head, _ := ReadFile(filename)
	return data, head
}

// LoadCSVfromString parses a CSV from a string and returns two maps
// one map with the data and another one with the csv head. It returns the rows
// read before an error; use ReadString to receive the error as well.
func LoadCSVfromString(csv string) (map[int][]string, map[string]int) {
	data, head, _ := ReadString(csv)
	return data, head
}

// Read parses semicolon-separated CSV from reader. The first record supplies
// the header names, and data rows are indexed starting at 1. All records must
// have the same number of fields as the header. On error, Read returns the
// successfully read rows and header together with the error.
func Read(reader io.Reader) (map[int][]string, map[string]int, error) {
	data := map[int][]string{}
	head := map[string]int{}

	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	for row := 0; ; row++ {
		record, err := csvReader.Read()
		if err == io.EOF {
			return data, head, nil
		}
		if err != nil {
			return data, head, fmt.Errorf("csv: read record %d: %w", row+1, err)
		}

		if row == 0 {
			head = GetHead(map[int][]string{0: record})
		} else {
			data[row] = record
		}
	}
}

// ReadFile parses semicolon-separated CSV from a file. It follows the row and
// partial-result conventions of Read, and reports file-open errors.
func ReadFile(filename string) (map[int][]string, map[string]int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return map[int][]string{}, map[string]int{}, err
	}
	defer func() { _ = f.Close() }()
	return Read(f)
}

// ReadString parses semicolon-separated CSV from a string.
func ReadString(input string) (map[int][]string, map[string]int, error) {
	return Read(strings.NewReader(input))
}

// GetHead returns a map which represents the head of a csv map
func GetHead(data map[int][]string) map[string]int {
	head := make(map[string]int, len(data[0]))
	for pos, name := range data[0] {
		head[name] = pos
	}
	return head
}
