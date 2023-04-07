// yolodb is a simple in-memory database that supports basic CRUD operations.
package yolodb

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Database is the main struct of the database.
type Database struct {
	tables map[string]*Table
}

// Table is a struct that holds the table information.
type Table struct {
	name     string
	columns  []*Column
	records  []*Record
	nextID   int
	indexMap map[string]map[interface{}][]int
}

// Column is a struct that holds the column information.
type Column struct {
	Name string
	Typ  string
}

// Record is a struct that holds the record information.
type Record struct {
	id   int
	data []interface{}
}

// NewDatabase returns a pointer to a new Database.
func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*Table),
	}
}

// CreateTable creates a new table in the database.
func (db *Database) CreateTable(name string, columns []*Column) error {
	if _, ok := db.tables[name]; ok {
		return fmt.Errorf("Table %s already exists", name)
	}

	for _, col := range columns {
		if col.Typ != "int" && col.Typ != "string" {
			return fmt.Errorf("Invalid column type %s", col.Typ)
		}
	}

	db.tables[name] = &Table{
		name:     name,
		columns:  columns,
		records:  []*Record{},
		nextID:   1,
		indexMap: make(map[string]map[interface{}][]int),
	}

	return nil
}

// Insert inserts a new record into the table.
func (db *Database) Insert(name string, data map[string]interface{}) error {
	table, ok := db.tables[name]
	if !ok {
		return fmt.Errorf("Table %s not found", name)
	}

	if len(data) != len(table.columns) {
		return errors.New("Data does not match table schema")
	}

	recordData := make([]interface{}, len(table.columns))
	for i, col := range table.columns {
		if val, ok := data[col.Name]; ok {
			if reflect.TypeOf(val).String() != col.Typ {
				return fmt.Errorf("Invalid type for column %s", col.Name)
			}
			recordData[i] = val
		} else {
			return fmt.Errorf("Missing value for column %s", col.Name)
		}
	}

	record := &Record{
		id:   table.nextID,
		data: recordData,
	}

	for i, col := range table.columns {
		if _, ok := table.indexMap[col.Name]; !ok {
			table.indexMap[col.Name] = make(map[interface{}][]int)
		}
		table.indexMap[col.Name][recordData[i]] = append(table.indexMap[col.Name][recordData[i]], table.nextID)
	}

	table.records = append(table.records, record)
	table.nextID++

	return nil
}

// Select selects records from the table.
func (db *Database) Select(name string, columns []string, whereClauses []string) ([]map[string]interface{}, error) {
	table, ok := db.tables[name]
	if !ok {
		return nil, fmt.Errorf("Table %s not found", name)
	}

	var results []map[string]interface{}

	columnSet := make(map[string]bool)
	for _, col := range columns {
		columnSet[col] = true
	}

	var indexMap map[string]map[interface{}][]int
	if len(whereClauses) > 0 {
		indexMap = make(map[string]map[interface{}][]int)
		for _, clause := range whereClauses {
			parts := strings.Split(clause, " ")
			if len(parts) != 3 {
				return nil, fmt.Errorf("Invalid where clause: %s", clause)
			}

			colName := parts[0]
			op := parts[1]
			valStr := parts[2]
			val, err := db.parseValue(valStr, table.getColumnType(colName))
			if err != nil {
				return nil, err
			}

			if _, ok := indexMap[colName]; !ok {
				if _, ok := table.indexMap[colName]; !ok {
					return nil, fmt.Errorf("Column %s is not indexed", colName)
				}
				indexMap[colName] = table.indexMap[colName]
			}

			var ids []int
			switch op {
			case "=":
				if vals, ok := indexMap[colName][val]; ok {
					ids = vals
				}
			case ">":
				for key, vals := range indexMap[colName] {
					if key.(int) > val.(int) {
						ids = append(ids, vals...)
					}
				}
			case "<":
				for key, vals := range indexMap[colName] {
					if key.(int) < val.(int) {
						ids = append(ids, vals...)
					}
				}
			}

			newIndexMap := make(map[interface{}][]int)
			for _, id := range ids {
				record := table.getRecordByID(id)
				if _, ok := indexMap[colName][record.data[table.getColumnIndex(colName)]]; ok {
					newIndexMap[record.data[table.getColumnIndex(colName)]] = append(newIndexMap[record.data[table.getColumnIndex(colName)]], id)
				}
			}

			indexMap[colName] = newIndexMap
		}
	} else {
		indexMap = table.indexMap
	}

	for _, record := range table.records {
		if len(whereClauses) > 0 {
			match := true
			for _, clause := range whereClauses {
				parts := strings.Split(clause, " ")
				colName := parts[0]
				op := parts[1]
				valStr := parts[2]
				val, err := db.parseValue(valStr, table.getColumnType(colName))
				if err != nil {
					return nil, err
				}
				colIndex := table.getColumnIndex(colName)
				if op == "=" && record.data[colIndex] != val {
					match = false
					break
				}
				if op == ">" && record.data[colIndex].(int) <= val.(int) {
					match = false
					break
				}
				if op == "<" && record.data[colIndex].(int) >= val.(int) {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}

		row := make(map[string]interface{})
		for i, col := range table.columns {
			if columnSet[col.Name] {
				row[col.Name] = record.data[i]
			}
		}
		results = append(results, row)
	}

	return results, nil
}

// Update updates records in the table.
func (db *Database) Update(name string, data map[string]interface{}, whereClauses []string) error {
	table, ok := db.tables[name]
	if !ok {
		return fmt.Errorf("Table %s not found", name)
	}

	if len(data) != len(table.columns) {
		return errors.New("Data does not match table schema")
	}

	indexMap := make(map[string]map[interface{}][]int)
	for _, clause := range whereClauses {
		parts := strings.Split(clause, " ")
		if len(parts) != 3 {
			return fmt.Errorf("Invalid where clause: %s", clause)
		}

		colName := parts[0]
		op := parts[1]
		valStr := parts[2]
		val, err := db.parseValue(valStr, table.getColumnType(colName))
		if err != nil {
			return err
		}

		if _, ok := indexMap[colName]; !ok {
			if _, ok := table.indexMap[colName]; !ok {
				return fmt.Errorf("Column %s is not indexed", colName)
			}
			indexMap[colName] = table.indexMap[colName]
		}

		var ids []int
		switch op {
		case "=":
			if vals, ok := indexMap[colName][val]; ok {
				ids = vals
			}
		case ">":
			for key, vals := range indexMap[colName] {
				if key.(int) > val.(int) {
					ids = append(ids, vals...)
				}
			}
		case "<":
			for key, vals := range indexMap[colName] {
				if key.(int) < val.(int) {
					ids = append(ids, vals...)
				}
			}
		}

		newIndexMap := make(map[interface{}][]int)
		for _, id := range ids {
			record := table.getRecordByID(id)
			if _, ok := indexMap[colName][record.data[table.getColumnIndex(colName)]]; ok {
				newIndexMap[record.data[table.getColumnIndex(colName)]] = append(newIndexMap[record.data[table.getColumnIndex(colName)]], id)
			}
		}

		indexMap[colName] = newIndexMap
	}

	for _, record := range table.records {
		if len(whereClauses) > 0 {
			match := true
			for _, clause := range whereClauses {
				parts := strings.Split(clause, " ")
				colName := parts[0]
				op := parts[1]
				valStr := parts[2]
				val, err := db.parseValue(valStr, table.getColumnType(colName))
				if err != nil {
					return err
				}
				colIndex := table.getColumnIndex(colName)
				if op == "=" && record.data[colIndex] != val {
					match = false
					break
				}
				if op == ">" && record.data[colIndex].(int) <= val.(int) {
					match = false
					break
				}
				if op == "<" && record.data[colIndex].(int) >= val.(int) {
					match = false
					break
				}

			}
			if !match {
				continue
			}
		}

		for colName, val := range data {
			//colIndex := table.getColumnIndex(colName.(string))
			colIndex := table.getColumnIndex(colName)
			if reflect.TypeOf(val).String() != table.columns[colIndex].Typ {
				return fmt.Errorf("Invalid type for column %s", colName)
			}
			record.data[colIndex] = val
			if _, ok := table.indexMap[colName]; ok {
				if _, ok := table.indexMap[colName][record.data[colIndex]]; !ok {
					table.indexMap[colName][record.data[colIndex]] = []int{}
				}
				table.indexMap[colName][record.data[colIndex]] = append(table.indexMap[colName][record.data[colIndex]], record.id)
			}
		}
	}

	return nil
}

// Delete deletes records from the table.
func (db *Database) Delete(name string, whereClauses []string) error {
	table, ok := db.tables[name]
	if !ok {
		return fmt.Errorf("Table %s not found", name)
	}

	indexMap := make(map[string]map[interface{}][]int)
	for _, clause := range whereClauses {
		parts := strings.Split(clause, " ")
		if len(parts) != 3 {
			return fmt.Errorf("Invalid where clause: %s", clause)
		}

		colName := parts[0]
		op := parts[1]
		valStr := parts[2]
		val, err := db.parseValue(valStr, table.getColumnType(colName))
		if err != nil {
			return err
		}

		if _, ok := indexMap[colName]; !ok {
			if _, ok := table.indexMap[colName]; !ok {
				return fmt.Errorf("Column %s is not indexed", colName)
			}
			indexMap[colName] = table.indexMap[colName]
		}

		var ids []int
		switch op {
		case "=":
			if vals, ok := indexMap[colName][val]; ok {
				ids = vals
			}
		case ">":
			for key, vals := range indexMap[colName] {
				if key.(int) > val.(int) {
					ids = append(ids, vals...)
				}
			}
		case "<":
			for key, vals := range indexMap[colName] {
				if key.(int) < val.(int) {
					ids = append(ids, vals...)
				}
			}
		}

		newIndexMap := make(map[interface{}][]int)
		for _, id := range ids {
			record := table.getRecordByID(id)
			if _, ok := indexMap[colName][record.data[table.getColumnIndex(colName)]]; ok {
				newIndexMap[record.data[table.getColumnIndex(colName)]] = append(newIndexMap[record.data[table.getColumnIndex(colName)]], id)
			}
		}

		indexMap[colName] = newIndexMap
	}

	newRecords := []*Record{}
	for _, record := range table.records {
		if len(whereClauses) > 0 {
			match := true
			for _, clause := range whereClauses {
				parts := strings.Split(clause, " ")
				colName := parts[0]
				op := parts[1]
				valStr := parts[2]
				val, err := db.parseValue(valStr, table.getColumnType(colName))
				if err != nil {
					return err
				}
				colIndex := table.getColumnIndex(colName)
				if op == "=" && record.data[colIndex] != val {
					match = false
					break
				}
				if op == ">" && record.data[colIndex].(int) <= val.(int) {
					match = false
					break
				}
				if op == "<" && record.data[colIndex].(int) >= val.(int) {
					match = false
					break
				}
			}
			if !match {
				newRecords = append(newRecords, record)
				continue
			}
		}

		for i, col := range table.columns {
			if _, ok := table.indexMap[col.Name]; ok {
				delete(table.indexMap[col.Name], record.data[i])
			}
		}
	}

	table.records = newRecords

	return nil
}

// parseValue parses a value string into the appropriate type.
func (db *Database) parseValue(valStr string, typ string) (interface{}, error) {
	switch typ {
	case "int":
		return strconv.Atoi(valStr)
	case "string":
		return valStr, nil
	default:
		return nil, fmt.Errorf("Invalid column type %s", typ)
	}
}

// getRecordByID returns the record with the given id.
func (t *Table) getRecordByID(id int) *Record {
	for _, record := range t.records {
		if record.id == id {
			return record
		}
	}
	return nil
}

// getColumnType returns the type of the column with the given name.
func (t *Table) getColumnType(name string) string {
	for _, col := range t.columns {
		if col.Name == name {
			return col.Typ
		}
	}
	return ""
}

// getColumnIndex returns the index of the column with the given name.
func (t *Table) getColumnIndex(name string) int {
	for i, col := range t.columns {
		if col.Name == name {
			return i
		}
	}
	return -1
}
