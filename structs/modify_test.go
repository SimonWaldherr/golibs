package structs

import (
	"testing"
)

type Person struct {
	Name  string `modify:"upper"`
	Age   int    `modify:"abs"`
	Email string `modify:"lower"`
}

func TestModify(t *testing.T) {
	person := &Person{
		Name:  "john",
		Age:   -25,
		Email: "JOHN.DOE@EXAMPLE.COM",
	}

	if err := ValidateAndModify(person); err != nil {
		t.Errorf("Modification failed: %v", err)
	}

	if person.Name != "JOHN" {
		t.Errorf("Expected Name to be 'JOHN', got '%s'", person.Name)
	}

	if person.Age != 25 {
		t.Errorf("Expected Age to be 25, got %d", person.Age)
	}

	if person.Email != "john.doe@example.com" {
		t.Errorf("Expected Email to be 'john.doe@example.com', got '%s'", person.Email)
	}
}
