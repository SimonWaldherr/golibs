package structs

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// CustomValidatorFunc validates a field. Callbacks run without the registry lock
// and may register other validators. They must synchronize their own shared data.
type CustomValidatorFunc func(field reflect.Value) error

var (
	customValidators    = make(map[string]CustomValidatorFunc)
	customValidatorLock sync.RWMutex
)

// RegisterCustomValidator installs or replaces a validator safely across goroutines.
// A nil function unregisters the tag. Cached struct plans observe later registrations.
func RegisterCustomValidator(tag string, fn CustomValidatorFunc) {
	customValidatorLock.Lock()
	defer customValidatorLock.Unlock()
	if fn == nil {
		delete(customValidators, tag)
	} else {
		customValidators[tag] = fn
	}
}

func applyValidation(field reflect.Value, tags []string) error {
	for _, tag := range tags {
		switch {
		case tag == "required":
			if isEmptyValue(field) {
				return errors.New("field is required")
			}
		case strings.HasPrefix(tag, "gte=") || strings.HasPrefix(tag, "lte="):
			if field.Kind() < reflect.Int || field.Kind() > reflect.Int64 {
				return fmt.Errorf("rule %q requires a signed integer field", tag)
			}
			boundary, err := strconv.ParseInt(tag[4:], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid numeric rule %q: %w", tag, err)
			}
			if tag[:3] == "gte" && field.Int() < boundary {
				return errors.New("value is less than the minimum allowed")
			}
			if tag[:3] == "lte" && field.Int() > boundary {
				return errors.New("value is greater than the maximum allowed")
			}
		case tag == "email" || tag == "iscolor" || strings.HasPrefix(tag, "oneof=") ||
			strings.HasPrefix(tag, "minlen=") || strings.HasPrefix(tag, "maxlen=") || strings.HasPrefix(tag, "pattern="):
			if field.Kind() != reflect.String {
				return fmt.Errorf("rule %q requires a string field", tag)
			}
			text := field.String()
			switch {
			case tag == "email":
				if !isValidEmail(text) {
					return errors.New("invalid email address")
				}
			case tag == "iscolor":
				if !isValidColor(text) {
					return errors.New("invalid color value")
				}
			case strings.HasPrefix(tag, "oneof="):
				options := strings.Fields(tag[6:])
				if len(options) == 0 {
					return errors.New("oneof requires at least one option")
				}
				if !contains(options, text) {
					return errors.New("value is not one of the allowed options")
				}
			case strings.HasPrefix(tag, "minlen=") || strings.HasPrefix(tag, "maxlen="):
				length, err := parseLength(tag[7:])
				if err != nil {
					return err
				}
				if tag[:6] == "minlen" && len(text) < length {
					return errors.New("string length is less than the minimum allowed")
				}
				if tag[:6] == "maxlen" && len(text) > length {
					return errors.New("string length is greater than the maximum allowed")
				}
			case strings.HasPrefix(tag, "pattern="):
				expression, err := compilePattern(tag[8:])
				if err != nil {
					return fmt.Errorf("invalid validation pattern: %w", err)
				}
				if !expression.MatchString(text) {
					return errors.New("string does not match the required pattern")
				}
			}
		default:
			customValidatorLock.RLock()
			validator, exists := customValidators[tag]
			customValidatorLock.RUnlock()
			if !exists {
				return fmt.Errorf("unknown validation rule %q", tag)
			}
			if err := validator(field); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseLength(text string) (int, error) {
	length, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("invalid length %q: %w", text, err)
	}
	if length < 0 {
		return 0, errors.New("length must not be negative")
	}
	return length, nil
}
