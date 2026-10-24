package structs

import (
	"errors"
	"math"
	"reflect"
	"sync"
	"testing"
	"time"
)

func taggedValue(tag string, initial interface{}) interface{} {
	typ := reflect.StructOf([]reflect.StructField{{Name: "Value", Type: reflect.TypeOf(initial), Tag: reflect.StructTag(tag)}})
	value := reflect.New(typ)
	value.Elem().Field(0).Set(reflect.ValueOf(initial))
	return value.Interface()
}

func TestInvalidRulesReturnErrors(t *testing.T) {
	for _, tc := range []struct {
		tag   string
		value interface{}
	}{
		{`validate:"pattern=["`, "abc"},
		{`validate:"unknown_rule"`, "abc"},
		{`validate:"gte=oops"`, 5},
		{`validate:"lte=999999999999999999999999999999"`, 5},
		{`validate:"gte=1"`, "abc"},
		{`validate:"email"`, 5},
		{`validate:"minlen=-1"`, "abc"},
		{`validate:"maxlen=oops"`, "abc"},
		{`validate:"oneof="`, ""},
		{`modify:"truncate=-1"`, "abc"},
		{`modify:"truncate=oops"`, "abc"},
		{`modify:"mask=1"`, "abc"},
		{`modify:"mask=-1-2"`, "abc"},
		{`modify:"mask=2-1"`, "abc"},
		{`modify:"mask=0-99"`, "abc"},
		{`modify:"pattern=["`, "abc"},
		{`modify:"upper"`, 5},
		{`modify:"abs"`, "abc"},
		{`modify:"abs"`, int8(math.MinInt8)},
		{`modify:"abs"`, int64(math.MinInt64)},
		{`modify:"hash=unsupported"`, "abc"},
		{`modify:"unknown_rule"`, "abc"},
	} {
		t.Run(tc.tag+reflect.TypeOf(tc.value).String(), func(t *testing.T) {
			if err := ValidateAndModify(taggedValue(tc.tag, tc.value)); err == nil {
				t.Fatal("invalid rule was silently accepted")
			}
		})
	}
}

func TestMalformedMetadataTerminates(t *testing.T) {
	for _, tag := range []string{`validate`, `validate:`, `validate:"`, `validate:"bad`, `validate:"\q"`} {
		value := taggedValue(tag, "abc")
		result := make(chan error, 1)
		go func() { result <- ValidateAndModify(value) }()
		select {
		case err := <-result:
			if err == nil {
				t.Fatalf("accepted malformed tag %q", tag)
			}
		case <-time.After(3 * time.Second):
			t.Fatalf("malformed tag %q did not terminate", tag)
		}
	}
}

func TestRuleOrderEscapingAndPatternCommas(t *testing.T) {
	value := struct {
		Name string `json:"name,omitempty" modify:"trim, upper" validate:"required,pattern=^[A-Z]{2,4}$"`
		Code string `validate:"pattern=^a\"b$"`
		Mark string `validate:"pattern=^[({]$,maxlen=1"`
	}{" ab ", `a"b`, "{"}
	if err := ValidateAndModify(&value); err != nil || value.Name != "AB" {
		t.Fatalf("ordered rules: %+v, %v", value, err)
	}
	value.Name = "cd"
	if err := ValidateAndModify(&value); err != nil || value.Name != "CD" {
		t.Fatalf("cached plan retained old field values: %+v, %v", value, err)
	}
	first := struct {
		Value string `validate:"minlen=4" modify:"trim"`
	}{" ab "}
	if err := ValidateAndModify(&first); err != nil || first.Value != "ab" {
		t.Fatalf("declaration order changed: %+v, %v", first, err)
	}
	second := struct {
		Value string `modify:"trim" validate:"minlen=4"`
	}{" ab "}
	if err := ValidateAndModify(&second); err == nil {
		t.Fatal("validation did not observe the preceding modification")
	}
}

func TestPrivateFieldsAndInvalidInputs(t *testing.T) {
	value := struct {
		private string
		Public  string `modify:"first_upper"`
	}{"unchanged", "äpfel"}
	if err := ValidateAndModify(&value); err != nil || value.Public != "Äpfel" || value.private != "unchanged" {
		t.Fatalf("private or Unicode field: %+v, %v", value, err)
	}
	private := struct {
		private string `modify:"upper"`
	}{"untouched"}
	if err := ValidateAndModify(&private); err == nil || private.private != "untouched" {
		t.Fatal("private tagged field did not return an error")
	}
	var nilStruct *User
	for _, input := range []interface{}{nil, nilStruct, User{}, new(int)} {
		if err := ValidateAndModify(input); err == nil {
			t.Errorf("accepted invalid input %T", input)
		}
	}
}

func TestConcurrentValidatorRegistration(t *testing.T) {
	const tag = "concurrent_hardening"
	type value struct {
		Text string `validate:"concurrent_hardening"`
	}
	RegisterCustomValidator(tag, func(reflect.Value) error { return nil })
	defer RegisterCustomValidator(tag, nil)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for i := 0; i < 100; i++ {
				RegisterCustomValidator(tag, func(reflect.Value) error { return nil })
				if err := ValidateAndModify(&value{"text"}); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	group.Wait()
}

func TestCachedPlansObserveValidatorChanges(t *testing.T) {
	const tag = "late_hardening"
	type value struct {
		Text string `validate:"late_hardening"`
	}
	input := &value{"text"}
	if err := ValidateAndModify(input); err == nil {
		t.Fatal("unregistered validator was silently ignored")
	}
	want := errors.New("validator failure")
	RegisterCustomValidator(tag, func(reflect.Value) error {
		// Reentrant registration must not deadlock under the read lock.
		RegisterCustomValidator("nested_hardening", func(reflect.Value) error { return nil })
		return want
	})
	defer RegisterCustomValidator(tag, nil)
	defer RegisterCustomValidator("nested_hardening", nil)
	if err := ValidateAndModify(input); !errors.Is(err, want) {
		t.Fatalf("custom error was lost: %v", err)
	}
	RegisterCustomValidator(tag, nil)
	if err := ValidateAndModify(input); err == nil {
		t.Fatal("cached plan did not observe removal")
	}
}

func TestRequiredNestedFieldsAreValidated(t *testing.T) {
	type address struct {
		Street string `validate:"required"`
		City   string `validate:"required"`
	}
	for _, nested := range []interface{}{address{City: "Berlin"}, &address{City: "Berlin"}} {
		value := taggedValue(`validate:"required"`, nested)
		if err := ValidateAndModify(value); err == nil {
			t.Fatalf("nested required rule was bypassed for %T", nested)
		}
	}
	type linked struct {
		Name string  `validate:"required"`
		Next *linked `validate:"required"`
	}
	value := &linked{Name: "root"}
	value.Next = value
	if err := ValidateAndModify(value); err != nil {
		t.Fatalf("cyclic object did not terminate: %v", err)
	}
}

func TestSharedNestedFieldsAreModifiedOnce(t *testing.T) {
	type child struct {
		Value string `modify:"sprintf=[%s]"`
	}
	type childPointer *child
	shared := &child{Value: "value"}
	value := struct {
		First  *child       `validate:"required"`
		Second childPointer `validate:"required"`
	}{shared, shared}
	if err := ValidateAndModify(&value); err != nil || shared.Value != "[value]" {
		t.Fatalf("shared child was revisited: %q, %v", shared.Value, err)
	}
}
