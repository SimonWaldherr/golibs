package yolodb_test

import (
	"fmt"

	yoloDB "simonwaldherr.de/go/golibs/yoloDB"
)

func Example() {
	// Erstellen einer neuen Datenbank
	db := yoloDB.NewDatabase()

	// Erstellen einer neuen Tabelle
	columns := []*yoloDB.Column{
		{Name: "id", Typ: "int"},
		{Name: "name", Typ: "string"},
		{Name: "age", Typ: "int"},
	}
	err := db.CreateTable("users", columns)
	if err != nil {
		panic(err)
	}

	// Einfügen von Datensätzen in die Tabelle
	data := map[string]interface{}{"id": 1, "name": "John Doe", "age": 30}
	err = db.Insert("users", data)
	if err != nil {
		panic(err)
	}
	data = map[string]interface{}{"id": 2, "name": "Jane Doe", "age": 32}
	err = db.Insert("users", data)
	if err != nil {
		panic(err)
	}

	// Auswählen von Datensätzen aus der Tabelle
	results, err := db.Select("users", []string{"name"}, []string{"age > 30"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Results:")
	for _, record := range results {
		fmt.Println(record)
	}

	// Output:
	// Results:
	// map[name:Jane Doe]
}
