package structs

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type fieldPlan struct {
	index  int
	tags   []string
	nested bool
}

type structPlan struct {
	fields []fieldPlan
	err    error
}

var (
	structPlans  sync.Map // reflect.Type -> *structPlan; types and tags are immutable.
	patternPlans sync.Map // string -> patternPlan
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	colorPattern = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)
)

// ValidateAndModify applies validate and modify tags in their declared order.
// s must be a nonnil pointer to a struct. Invalid rules return errors rather
// than panicking. Fields without these tags are ignored, including private fields.
// Earlier modifications remain applied if a subsequent rule fails. Callers own
// the value and must synchronize concurrent access to that same value.
func ValidateAndModify(s interface{}) error {
	return validateAndModify(s, nil)
}

func validateAndModify(s interface{}, visited map[interface{}]bool) error {
	value := reflect.ValueOf(s)
	if value.Kind() != reflect.Ptr || value.IsNil() || value.Elem().Kind() != reflect.Struct {
		return errors.New("invalid type, expected a pointer to a struct")
	}
	value = value.Elem()
	typ := value.Type()
	// Keep visited objects alive while traversing the graph; pointer interfaces
	// retain type and identity without relying on uintptr addresses.
	identity := value.Addr().Interface()
	if visited != nil {
		if visited[identity] {
			return nil
		}
		visited[identity] = true
	}
	cached, ok := structPlans.Load(typ)
	if !ok {
		plan := &structPlan{}
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			tags, err := getCombinedTags(field.Tag)
			if err != nil {
				plan.err = fmt.Errorf("structs: field %s: %w", field.Name, err)
				break
			}
			if len(tags) == 0 {
				continue
			}
			if field.PkgPath != "" {
				plan.err = fmt.Errorf("structs: tagged field %s is not exported", field.Name)
				break
			}
			kind := field.Type.Kind()
			nested := (kind == reflect.Struct || kind == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct) && contains(tags, "validate:required")
			plan.fields = append(plan.fields, fieldPlan{index: i, tags: tags, nested: nested})
		}
		cached, _ = structPlans.LoadOrStore(typ, plan)
	}
	plan := cached.(*structPlan)
	if plan.err != nil {
		return plan.err
	}
	for _, field := range plan.fields {
		current := value.Field(field.index)
		if err := processTags(current, field.tags); err != nil {
			return fmt.Errorf("structs: field %s: %w", typ.Field(field.index).Name, err)
		}
		if field.nested {
			if visited == nil {
				visited = map[interface{}]bool{identity: true}
			}
			if current.Kind() == reflect.Struct {
				current = current.Addr()
			}
			if err := validateAndModify(current.Interface(), visited); err != nil {
				return fmt.Errorf("structs: field %s: %w", typ.Field(field.index).Name, err)
			}
		}
	}
	return nil
}

// getCombinedTags reads standard quoted Go tags, preserving validate/modify
// order and ignoring unrelated metadata. Malformed metadata must not spin or panic.
func getCombinedTags(tag reflect.StructTag) ([]string, error) {
	var result []string
	remaining := string(tag)
	for remaining != "" {
		remaining = strings.TrimLeft(remaining, " ")
		if remaining == "" {
			break
		}
		i := 0
		for i < len(remaining) && remaining[i] > ' ' && remaining[i] != ':' && remaining[i] != '"' && remaining[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(remaining) || remaining[i] != ':' || remaining[i+1] != '"' {
			return nil, errors.New("malformed struct tag")
		}
		key := remaining[:i]
		remaining = remaining[i+1:]
		end := 1
		for end < len(remaining) && remaining[end] != '"' {
			if remaining[end] == '\\' {
				end++
			}
			end++
		}
		if end >= len(remaining) {
			return nil, errors.New("unterminated struct tag")
		}
		text, err := strconv.Unquote(remaining[:end+1])
		if err != nil {
			return nil, fmt.Errorf("invalid quoted struct tag: %w", err)
		}
		remaining = remaining[end+1:]
		if key == "validate" || key == "modify" {
			rules, err := parseTag(text)
			if err != nil {
				return nil, err
			}
			for _, rule := range rules {
				if rule != "" {
					result = append(result, key+":"+rule)
				}
			}
		}
	}
	return result, nil
}

// parseTag permits commas inside regex groups, character classes and quantifiers.
func parseTag(tag string) ([]string, error) {
	var result []string
	start, depth, classStart := 0, 0, 0
	inClass, quoted := false, false
	pattern := isPatternRule(tag)
	for i := 0; i < len(tag); i++ {
		if pattern {
			if quoted {
				if tag[i] == '\\' && i+1 < len(tag) && tag[i+1] == 'E' {
					quoted = false
					i++
				}
				continue
			}
			if tag[i] == '\\' {
				if !inClass && i+1 < len(tag) && tag[i+1] == 'Q' {
					quoted = true
				}
				i++
				continue
			}
			if inClass {
				if tag[i] == '[' && i+1 < len(tag) && tag[i+1] == ':' {
					if end := strings.Index(tag[i+2:], ":]"); end >= 0 {
						i += end + 3
						continue
					}
				}
				firstClosing := i == classStart+1 || i == classStart+2 && tag[classStart+1] == '^'
				if tag[i] == ']' && !firstClosing {
					inClass = false
				}
				continue
			}
			switch tag[i] {
			case '[':
				inClass, classStart = true, i
			case '(':
				depth++
			case ')':
				if depth > 0 {
					depth--
				}
			case '{':
				if end, ok := repeatEnd(tag, i); ok {
					i = end
				}
			}
		}
		if tag[i] == ',' && depth == 0 && !inClass {
			result = append(result, normalizeRule(tag[start:i]))
			start = i + 1
			pattern = isPatternRule(tag[start:])
		}
	}
	if quoted {
		return nil, errors.New("unterminated pattern literal quote; close it with \\E")
	}
	return append(result, normalizeRule(tag[start:])), nil
}

func isPatternRule(rule string) bool {
	separator := strings.IndexAny(rule, "=,")
	return separator >= 0 && rule[separator] == '=' && strings.TrimSpace(rule[:separator]) == "pattern"
}

// Whitespace around rule names is cosmetic; parameter bytes are significant
// for regular expressions, default values and formatting strings.
func normalizeRule(rule string) string {
	if separator := strings.IndexByte(rule, '='); separator >= 0 {
		return strings.TrimSpace(rule[:separator]) + rule[separator:]
	}
	return strings.TrimSpace(rule)
}

// A literal opening brace must not hide subsequent rules. Only complete
// numeric repetition syntax may contain a comma outside a class or group.
func repeatEnd(tag string, start int) (int, bool) {
	i := start + 1
	first := i
	for i < len(tag) && tag[i] >= '0' && tag[i] <= '9' {
		i++
	}
	if i == first || i-first > 1 && tag[first] == '0' {
		return 0, false
	}
	if i < len(tag) && tag[i] == ',' {
		i++
		first = i
		for i < len(tag) && tag[i] >= '0' && tag[i] <= '9' {
			i++
		}
		if i-first > 1 && tag[first] == '0' {
			return 0, false
		}
	}
	return i, i < len(tag) && tag[i] == '}'
}

func processTags(field reflect.Value, tags []string) error {
	for _, tag := range tags {
		if strings.HasPrefix(tag, "validate:") {
			if err := applyValidation(field, []string{tag[len("validate:"):]}); err != nil {
				return err
			}
		} else if strings.HasPrefix(tag, "modify:") {
			if err := applyModification(field, []string{tag[len("modify:"):]}); err != nil {
				return err
			}
		}
	}
	return nil
}

type patternPlan struct {
	expression *regexp.Regexp
	err        error
}

func compilePattern(pattern string) (*regexp.Regexp, error) {
	if cached, ok := patternPlans.Load(pattern); ok {
		plan := cached.(patternPlan)
		return plan.expression, plan.err
	}
	expression, err := regexp.Compile(pattern)
	cached, _ := patternPlans.LoadOrStore(pattern, patternPlan{expression, err})
	plan := cached.(patternPlan)
	return plan.expression, plan.err
}

// Helper functions

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Struct:
		return v.IsZero()
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
	return emailPattern.MatchString(email)
}

func isValidColor(color string) bool {
	return colorPattern.MatchString(color)
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
