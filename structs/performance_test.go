package structs

import "testing"

var benchmarkValidationError error

type benchmarkAccount struct {
	Name  string `modify:"trim,upper" validate:"required,minlen=3"`
	Email string `validate:"required,email"`
	Color string `validate:"iscolor"`
	Role  string `validate:"oneof=admin user"`
	Code  string `validate:"pattern=^[A-Z]{2}[0-9]{3}$"`
	Age   int    `validate:"gte=0,lte=130"`
}

func BenchmarkValidateAndModify(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		account := benchmarkAccount{" alice ", "alice@example.com", "#abcdef", "admin", "AB123", 30}
		benchmarkValidationError = ValidateAndModify(&account)
		if benchmarkValidationError != nil {
			b.Fatal(benchmarkValidationError)
		}
	}
}
