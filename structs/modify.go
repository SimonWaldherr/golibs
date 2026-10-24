package structs

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"unicode/utf8"
)

// applyModification applies rules in order. String operations require string
// fields; abs requires a signed integer. Invalid rules and bounds return errors.
func applyModification(field reflect.Value, tags []string) error {
	for _, tag := range tags {
		if tag == "abs" {
			if field.Kind() < reflect.Int || field.Kind() > reflect.Int64 {
				return fmt.Errorf("rule %q requires a signed integer field", tag)
			}
			if field.Int() < 0 {
				if field.Int() == -int64(1)<<(field.Type().Bits()-1) {
					return fmt.Errorf("absolute value overflows %s", field.Type())
				}
				field.SetInt(-field.Int())
			}
			continue
		}
		if field.Kind() != reflect.String {
			return fmt.Errorf("rule %q requires a string field", tag)
		}
		text := field.String()
		switch tag {
		case "upper":
			field.SetString(strings.ToUpper(text))
		case "lower":
			field.SetString(strings.ToLower(text))
		case "trim":
			field.SetString(strings.TrimSpace(text))
		case "asciiify":
			field.SetString(toASCII(text))
		case "first_upper":
			if text != "" {
				_, size := utf8.DecodeRuneInString(text)
				field.SetString(strings.ToUpper(text[:size]) + text[size:])
			}
		case "remove_leading_zeros":
			field.SetString(strings.TrimLeft(text, "0"))
		case "base64_encode":
			field.SetString(base64.StdEncoding.EncodeToString([]byte(text)))
		case "base64_decode":
			decoded, err := base64.StdEncoding.DecodeString(text)
			if err != nil {
				return err
			}
			field.SetString(string(decoded))
		case "url_encode":
			field.SetString(url.QueryEscape(text))
		case "url_decode":
			decoded, err := url.QueryUnescape(text)
			if err != nil {
				return err
			}
			field.SetString(decoded)
		default:
			switch {
			case strings.HasPrefix(tag, "default="):
				if text == "" {
					field.SetString(tag[8:])
				}
			case strings.HasPrefix(tag, "truncate="):
				length, err := parseLength(tag[9:])
				if err != nil {
					return err
				}
				if len(text) > length {
					field.SetString(text[:length])
				}
			case strings.HasPrefix(tag, "mask="):
				startText, endText, ok := strings.Cut(tag[5:], "-")
				if !ok {
					return fmt.Errorf("invalid mask range %q", tag[5:])
				}
				start, err := parseLength(startText)
				if err != nil {
					return err
				}
				end, err := parseLength(endText)
				if err != nil {
					return err
				}
				if start > end || end > len(text) {
					return fmt.Errorf("mask range %d-%d exceeds string length %d", start, end, len(text))
				}
				field.SetString(text[:start] + strings.Repeat("*", end-start) + text[end:])
			case strings.HasPrefix(tag, "pattern="):
				expression, err := compilePattern(tag[8:])
				if err != nil {
					return err
				}
				if match := expression.FindString(text); match != "" {
					field.SetString(match)
				}
			case strings.HasPrefix(tag, "hash="):
				algorithm := tag[5:]
				hash := Hash(text, algorithm)
				if hash == "" {
					return fmt.Errorf("unsupported hash algorithm %q", algorithm)
				}
				field.SetString(hash)
			case strings.HasPrefix(tag, "sprintf="):
				field.SetString(fmt.Sprintf(tag[8:], text))
			default:
				return fmt.Errorf("unknown modification rule %q", tag)
			}
		}
	}
	return nil
}
