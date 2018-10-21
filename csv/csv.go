package csv

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

// LoadCSVfromFile parses a CSV from a file and returns two maps
// one map with the data and another one with the csv head
func LoadCSVfromFile(filename string) (map[int][]string, map[string]int) {
	fp, _ := os.Open(filename)
	return loadCSV(bufio.NewReader(fp))
}

// LoadCSVfromString parses a CSV from a string and returns two maps
// one map with the data and another one with the csv head
func LoadCSVfromString(csv string) (map[int][]string, map[string]int) {
	fp := strings.NewReader(csv)
	return loadCSV(fp)
}

func loadCSV(reader io.Reader) (map[int][]string, map[string]int) {
	var row int
	var head = map[int][]string{}
	var data = map[int][]string{}

	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if row == 0 {
			head[row] = record
		} else {
			data[row] = record
		}
		row++
	}
	return data, GetHead(head)
}

// GetHead returns a map which represents the head of a csv map
func GetHead(data map[int][]string) map[string]int {
	head := make(map[string]int, len(data[0]))
	for pos, name := range data[0] {
		head[name] = pos
	}
	return head
}
