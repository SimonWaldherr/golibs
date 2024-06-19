package structs

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type CustomValidatorFunc func(field reflect.Value) error

var customValidators = make(map[string]CustomValidatorFunc)

func RegisterCustomValidator(tag string, fn CustomValidatorFunc) {
	customValidators[tag] = fn
}

// applyValidation wendet die einzelnen Validierungsregeln auf das Feld an
func applyValidation(field reflect.Value, tags []string) error {
	for _, tag := range tags {
		switch {
		case tag == "required":
			if isEmptyValue(field) {
				return errors.New("field is required")
			}
		case tag == "email":
			if !isValidEmail(field.String()) {
				return errors.New("invalid email address")
			}
		case tag == "iscolor":
			if !isValidColor(field.String()) {
				return errors.New("invalid color value")
			}
		case strings.HasPrefix(tag, "gte="):
			min, _ := strconv.Atoi(tag[4:])
			if field.Int() < int64(min) {
				return errors.New("value is less than the minimum allowed")
			}
		case strings.HasPrefix(tag, "lte="):
			max, _ := strconv.Atoi(tag[4:])
			if field.Int() > int64(max) {
				return errors.New("value is greater than the maximum allowed")
			}
		case strings.HasPrefix(tag, "oneof="):
			options := regexp.MustCompile(`\s+`).Split(tag[6:], -1)
			if !contains(options, field.String()) {
				return errors.New("value is not one of the allowed options")
			}
		case strings.HasPrefix(tag, "minlen="):
			minLen, _ := strconv.Atoi(tag[7:])
			if len(field.String()) < minLen {
				return errors.New("string length is less than the minimum allowed")
			}
		case strings.HasPrefix(tag, "maxlen="):
			maxLen, _ := strconv.Atoi(tag[7:])
			if len(field.String()) > maxLen {
				return errors.New("string length is greater than the maximum allowed")
			}
		case strings.HasPrefix(tag, "pattern="):
			pattern := tag[8:]
			if !regexp.MustCompile(pattern).MatchString(field.String()) {
				return errors.New("string does not match the required pattern")
			}
		default:
			if customValidator, exists := customValidators[tag]; exists {
				if err := customValidator(field); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
