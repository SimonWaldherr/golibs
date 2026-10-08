// Package yolodb provides a small in-memory database for int and string columns.
// Conditions use "column operator value", with = for both types and > or < for
// integers. Multiple conditions are ANDed. Database methods are safe to call
// concurrently. Data is held in memory and is not persisted.
package yolodb

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Database holds named tables. Its zero value is ready to use.
// A Database must not be copied after first use.
type Database struct {
	mu     sync.RWMutex
	tables map[string]*Table
}

// Table holds a schema and its records, managed through Database methods.
type Table struct {
	columns     []*Column
	columnIndex map[string]int
	records     []*Record
}

// Column describes a named column. Typ must be "int" or "string".
type Column struct {
	Name string
	Typ  string
}

// Record holds values in schema order, managed through Database methods.
type Record struct {
	data []interface{}
}

// NewDatabase returns an empty database.
func NewDatabase() *Database {
	return &Database{tables: make(map[string]*Table)}
}

// CreateTable creates a table with a copied schema. Column names must be unique
// and nonempty; nil columns and unsupported types return errors.
func (db *Database) CreateTable(name string, columns []*Column) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.tables[name]; ok {
		return fmt.Errorf("table %s already exists", name)
	}
	table := &Table{columns: make([]*Column, len(columns)), columnIndex: make(map[string]int, len(columns))}
	for i, col := range columns {
		if col == nil || col.Name == "" {
			return fmt.Errorf("column %d must have a name", i)
		}
		if col.Typ != "int" && col.Typ != "string" {
			return fmt.Errorf("invalid column type %s", col.Typ)
		}
		if _, ok := table.columnIndex[col.Name]; ok {
			return fmt.Errorf("duplicate column %s", col.Name)
		}
		copy := *col
		table.columns[i] = &copy
		table.columnIndex[col.Name] = i
	}
	if db.tables == nil {
		db.tables = make(map[string]*Table)
	}
	db.tables[name] = table
	return nil
}

// Insert adds a row. Data must contain every schema column with its exact Go
// type (int or string), and no additional columns. Validation happens before writing.
func (db *Database) Insert(name string, data map[string]interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	table, ok := db.tables[name]
	if !ok {
		return fmt.Errorf("table %s not found", name)
	}
	values, err := table.row(data)
	if err != nil {
		return err
	}
	table.records = append(table.records, &Record{data: values})
	return nil
}

// Select returns independent maps in insertion order. Each row includes only the
// requested columns, which must exist. Empty conditions match all rows. Empty
// columns produce empty maps for matching rows; no matches return a nil slice.
func (db *Database) Select(name string, columns []string, whereClauses []string) ([]map[string]interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	table, ok := db.tables[name]
	if !ok {
		return nil, fmt.Errorf("table %s not found", name)
	}
	indexes := make([]int, len(columns))
	for i, name := range columns {
		index, ok := table.columnIndex[name]
		if !ok {
			return nil, fmt.Errorf("column %s not found", name)
		}
		indexes[i] = index
	}
	conditions, err := table.conditions(whereClauses)
	if err != nil {
		return nil, err
	}
	var results []map[string]interface{}
	for _, record := range table.records {
		if !matches(record, conditions) {
			continue
		}
		row := make(map[string]interface{}, len(columns))
		for i, index := range indexes {
			row[columns[i]] = record.data[index]
		}
		results = append(results, row)
	}
	return results, nil
}

// Update replaces each matching row with a complete row, using Insert's schema
// rules. Empty conditions update every row. Invalid data or conditions leave all
// rows unchanged, even when the table is empty or no rows match.
func (db *Database) Update(name string, data map[string]interface{}, whereClauses []string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	table, ok := db.tables[name]
	if !ok {
		return fmt.Errorf("table %s not found", name)
	}
	values, err := table.row(data)
	if err != nil {
		return err
	}
	conditions, err := table.conditions(whereClauses)
	if err != nil {
		return err
	}
	for _, record := range table.records {
		if matches(record, conditions) {
			copy(record.data, values)
		}
	}
	return nil
}

// Delete removes matching rows, preserving the order of remaining rows.
// Empty conditions remove all rows. Invalid conditions leave the table unchanged.
func (db *Database) Delete(name string, whereClauses []string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	table, ok := db.tables[name]
	if !ok {
		return fmt.Errorf("table %s not found", name)
	}
	conditions, err := table.conditions(whereClauses)
	if err != nil {
		return err
	}
	remaining := table.records[:0]
	for _, record := range table.records {
		if !matches(record, conditions) {
			remaining = append(remaining, record)
		}
	}
	// Release removed records instead of keeping them in the backing array.
	for i := len(remaining); i < len(table.records); i++ {
		table.records[i] = nil
	}
	table.records = remaining
	return nil
}

func (table *Table) row(data map[string]interface{}) ([]interface{}, error) {
	if len(data) != len(table.columns) {
		return nil, fmt.Errorf("data does not match table schema")
	}
	values := make([]interface{}, len(table.columns))
	for i, col := range table.columns {
		value, ok := data[col.Name]
		if !ok {
			return nil, fmt.Errorf("missing value for column %s", col.Name)
		}
		valid := false
		switch col.Typ {
		case "int":
			_, valid = value.(int)
		case "string":
			_, valid = value.(string)
		}
		if !valid {
			return nil, fmt.Errorf("invalid type for column %s", col.Name)
		}
		values[i] = value
	}
	return values, nil
}

type condition struct {
	column   int
	operator byte
	value    interface{}
}

func (table *Table) conditions(clauses []string) ([]condition, error) {
	conditions := make([]condition, 0, len(clauses))
	for _, clause := range clauses {
		parts := strings.Fields(clause)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid where clause: %s", clause)
		}
		index, ok := table.columnIndex[parts[0]]
		if !ok {
			return nil, fmt.Errorf("column %s not found", parts[0])
		}
		operator := parts[1]
		if operator != "=" && operator != ">" && operator != "<" {
			return nil, fmt.Errorf("invalid operator %s", operator)
		}
		var value interface{} = parts[2]
		if table.columns[index].Typ == "int" {
			parsed, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, fmt.Errorf("invalid integer for column %s: %w", parts[0], err)
			}
			value = parsed
		} else if operator != "=" {
			return nil, fmt.Errorf("operator %s requires an int column", operator)
		}
		conditions = append(conditions, condition{column: index, operator: operator[0], value: value})
	}
	return conditions, nil
}

func matches(record *Record, conditions []condition) bool {
	for _, condition := range conditions {
		value := record.data[condition.column]
		switch condition.operator {
		case '=':
			if value != condition.value {
				return false
			}
		case '>':
			if value.(int) <= condition.value.(int) {
				return false
			}
		case '<':
			if value.(int) >= condition.value.(int) {
				return false
			}
		}
	}
	return true
}
