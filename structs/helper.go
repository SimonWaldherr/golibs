package structs

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

// parseCombinedTags teilt die kombinierten Tags in eine Liste von einzelnen Tags auf
func parseCombinedTags(tag string) []string {
	return regexp.MustCompile(`;\s*`).Split(tag, -1)
}

// parseTag teilt einen Untertag in einzelne Regeln auf
func parseTag(tag string) []string {
	return regexp.MustCompile(`,\s*`).Split(tag, -1)
}

// ValidateAndModify prüft und modifiziert die struct-Felder basierend auf den Tags in der angegebenen Reihenfolge
func ValidateAndModify(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("invalid type, expected a pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)

		combinedTags := getCombinedTags(structField.Tag)
		tags := parseCombinedTags(combinedTags)

		if err := processTags(field, tags); err != nil {
			return err
		}
	}
	return nil
}

// getCombinedTags kombiniert die Tags "validate" und "modify" unter Beibehaltung der Reihenfolge
func getCombinedTags(tag reflect.StructTag) string {
	var combinedTags []string
	tagString := string(tag)

	for tagString != "" {
		var key, value string
		idx := strings.Index(tagString, ":")
		if idx != -1 {
			key = tagString[:idx]
			tagString = tagString[idx+1:]
			if tagString[0] == '"' {
				endIdx := strings.Index(tagString[1:], "\"") + 1
				value = tagString[1:endIdx]
				tagString = tagString[endIdx+1:]
			}
		}
		if key != "" && value != "" {
			parts := strings.Split(value, ",")
			for _, part := range parts {
				combinedTags = append(combinedTags, key+":"+part)
			}
		}
		tagString = strings.TrimLeft(tagString, " ")
	}

	return strings.Join(combinedTags, ";")
}

// processTags wendet die einzelnen Validierungs- und Modifikationsregeln in der angegebenen Reihenfolge an
func processTags(field reflect.Value, tags []string) error {
	for _, tag := range tags {
		if strings.HasPrefix(tag, "validate:") {
			validationTags := parseTag(tag[len("validate:"):])
			if err := applyValidation(field, validationTags); err != nil {
				return err
			}
		} else if strings.HasPrefix(tag, "modify:") {
			modificationTags := parseTag(tag[len("modify:"):])
			if err := applyModification(field, modificationTags); err != nil {
				return err
			}
		}
	}
	return nil
}

// Hilfsfunktionen

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidColor(color string) bool {
	re := regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)
	return re.MatchString(color)
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func toASCII(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r > unicode.MaxASCII {
			result.WriteRune('?')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
