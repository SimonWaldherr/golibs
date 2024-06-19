package structs_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"testing"
	"unicode"
	
	"simonwaldherr.de/go/golibs/structs"
)

type LIPS struct {
	VBELN string
	POSNR string
	MATNR string
	MATKL string
	ARKTX string
	EANNR string
	LGORT string
	LFIMG string
	VRKME string
	VKBUR string
}

func ExampleReflect() {
	var lips LIPS
	x := structs.Reflect(lips)

	keys := []string{}
	for k := range x {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("key: %v, value: %v\n", k, x[k])
	}

	// Output:
	// key: ARKTX, value: string
	// key: EANNR, value: string
	// key: LFIMG, value: string
	// key: LGORT, value: string
	// key: MATKL, value: string
	// key: MATNR, value: string
	// key: POSNR, value: string
	// key: VBELN, value: string
	// key: VKBUR, value: string
	// key: VRKME, value: string
}

func ExampleReflectHelper() {
	var lips LIPS
	v := reflect.ValueOf(lips)
	t := reflect.TypeOf(lips)

	structs.ReflectHelper(v, t, 0, func(name string, vtype string, value interface{}, depth int) {
		fmt.Printf("%v - %v - %v - %v\n", name, vtype, value, depth)
	})

	// Output:
	// VBELN - string -  - 0
	// POSNR - string -  - 0
	// MATNR - string -  - 0
	// MATKL - string -  - 0
	// ARKTX - string -  - 0
	// EANNR - string -  - 0
	// LGORT - string -  - 0
	// LFIMG - string -  - 0
	// VRKME - string -  - 0
	// VKBUR - string -  - 0
}

type ExampleUser struct {
	Name     string `json:"name" modify:"upper" validate:"required"`
	Age      int    `json:"age" modify:"abs" validate:"gte=0,lte=130"`
	Email    string `json:"email" modify:"lower" modify:"trim" validate:"required,email"`
	Password string `json:"password" modify:"trim" validate:"minlen=8"`
	Status   string `json:"status" validate:"oneof=active inactive"`
	Zip      string `modify:"pattern=\d+"`
}

func ExampleValidateAndModify() {
	user := &ExampleUser{
		Name:     "john",
		Age:      -25,
		Email:    "JOHN.DOE@EXAMPLE.COM ",
		Password: "  password123  ",
		Status:   "active",
		Zip:      "94522 Wallersdorf",
	}

	if err := structs.ValidateAndModify(user); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Validation and modification passed!")
		fmt.Printf("%+v\n", user)
	}

	// Output:
	// Validation and modification passed!
	// &{Name:JOHN Age:25 Email:john.doe@example.com Password:password123 Status:active Zip:94522}
}

func TestExampleUserValidation(t *testing.T) {
	validUser := &ExampleUser{
		Name:     "Alice",
		Age:      28,
		Email:    "alice@example.com",
		Password: "securepassword",
		Status:   "active",
	}

	if err := structs.ValidateAndModify(validUser); err != nil {
		t.Errorf("Expected valid user to pass validation, but got error: %v", err)
	}

	invalidUser := &ExampleUser{
		Name:     "",
		Age:      150,
		Email:    "invalid-email",
		Password: "short",
		Status:   "unknown",
	}

	if err := structs.ValidateAndModify(invalidUser); err == nil {
		t.Errorf("Expected invalid user to fail validation, but got no error")
	}
}

func ExampleJSON() {
	jsonString := `{
		"name": "john",
		"age": -25,
		"email": "JOHN.DOE@EXAMPLE.COM ",
		"password": "  password123  ",
		"status": "active"
	}`

	var user ExampleUser

	// Decode JSON string
	if err := json.Unmarshal([]byte(jsonString), &user); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Validate and modify struct fields
	if err := structs.ValidateAndModify(&user); err != nil {
		log.Fatalf("Validation/Modification error: %v", err)
	}

	// Encode struct back to JSON
	modifiedJSON, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}

	fmt.Println("Modified JSON:")
	fmt.Println(string(modifiedJSON))

	// Output:
	// Modified JSON:
	//{
	//     "name": "JOHN",
	//     "age": 25,
	//     "email": "john.doe@example.com",
	//     "password": "password123",
	//     "status": "active",
	//     "Zip": ""
	//}
}

// Example of using the validator with a struct that has nested fields
func ExampleNestedStruct() {
	type Address struct {
		Street  string `validate:"required"`
		City    string `validate:"required"`
		ZipCode string `validate:"required"`
	}

	type User struct {
		Name    string  `validate:"required"`
		Age     int     `validate:"gte=0,lte=130"`
		Email   string  `validate:"required,email"`
		Address Address `validate:"required"`
	}

	user := &User{
		Name:  "J. Doe",
		Age:   28,
		Email: "j.doe@example.tld",
		Address: Address{
			Street:  "123 Main St",
			City:    "Springfield",
			ZipCode: "12345",
		},
	}

	if err := structs.ValidateAndModify(user); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Validation and modification passed!")
		fmt.Printf("%+v\n", user)
	}

	// Output:
	// Validation and modification passed!
	// &{Name:J. Doe Age:28 Email:j.doe@example.tld Address:{Street:123 Main St City:Springfield ZipCode:12345}}
}

// Example of using the validator with a struct that has nested fields
func ExampleBase64() {
	type User struct {
		Name string `validate:"required"`
		Data string `modify:"base64_decode" validate:"required"`
	}

	user := &User{
		Name: "Simon",
		Data: "SGVsbG8gV29ybGQ=",
	}

	if err := structs.ValidateAndModify(user); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Validation and modification passed!")
		fmt.Printf("%+v\n", user)
	}

	// Output:
	// Validation and modification passed!
	// &{Name:Simon Data:Hello World}
}

// Example of using the hashing feature of the validator
func ExampleHash() {
	type User struct {
		Name string `validate:"required"`
		Data string `modify:"hash=md5"`
	}

	user := &User{
		Name: "Simon",
		Data: "Hello World",
	}

	if err := structs.ValidateAndModify(user); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Validation and modification passed!")
		fmt.Printf("%+v\n", user)
	}

	// Output:
	// Validation and modification passed!
	// &{Name:Simon Data:b10a8db164e0754105b7a99be72e3fe5}
}

// Example of using the URL encoding feature of the validator
func ExampleURLEncode() {
	type Req struct {
		Data string `modify:"url_encode"`
	}

	req := &Req{
		Data: "Hello World",
	}

	if err := structs.ValidateAndModify(req); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Validation and modification passed!")
		fmt.Printf("%+v\n", req)
	}

	// Output:
	// Validation and modification passed!
	// &{Data:Hello+World}
}

// CustomValidatorFunc ist eine benutzerdefinierte Validierungsfunktion, die überprüft, ob ein String nur Buchstaben enthält.
func CustomValidatorFuncOnlyLetters(field reflect.Value) error {
	if field.Kind() == reflect.String {
		str := field.String()
		for _, char := range str {
			if !unicode.IsLetter(char) {
				return errors.New("field contains non-letter characters")
			}
		}
	}
	return nil
}

// Example with a custom validation function
func ExampleCustomValidation() {
	type User struct {
		Name string `validate:"required,onlyletters"`
	}

	// Register custom validation function
	structs.RegisterCustomValidator("onlyletters", CustomValidatorFuncOnlyLetters)

	user := &User{
		Name: "Abc",
	}

	if err := structs.ValidateAndModify(user); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Validation and modification passed!")
		fmt.Printf("%+v\n", user)
	}

	// Output:
	// Validation and modification passed!
	// &{Name:Abc}
}
