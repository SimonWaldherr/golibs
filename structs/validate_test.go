package structs

import (
	"testing"
)

type User struct {
	FirstName string `validate:"required"`
	Age       int    `validate:"gte=0,lte=130"`
	Email     string `validate:"required,email"`
}

func TestValidate(t *testing.T) {
	user := &User{
		FirstName: "John",
		Age:       30,
		Email:     "john.doe@example.com",
	}

	if err := ValidateAndModify(user); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	invalidUser := &User{
		FirstName: "",
		Age:       150,
		Email:     "invalid-email",
	}

	if err := ValidateAndModify(invalidUser); err == nil {
		t.Errorf("Validation should have failed for invalidUser")
	}
}
