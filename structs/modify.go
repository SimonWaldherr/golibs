package structs

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// applyModification wendet die einzelnen Modifikationsregeln auf das Feld an
func applyModification(field reflect.Value, tags []string) error {
	for _, tag := range tags {
		switch tag {
		case "upper":
			if field.Kind() == reflect.String {
				field.SetString(strings.ToUpper(field.String()))
			}
		case "lower":
			if field.Kind() == reflect.String {
				field.SetString(strings.ToLower(field.String()))
			}
		case "trim":
			if field.Kind() == reflect.String {
				field.SetString(strings.TrimSpace(field.String()))
			}
		case "abs":
			if field.Kind() == reflect.Int || field.Kind() == reflect.Int8 || field.Kind() == reflect.Int16 ||
				field.Kind() == reflect.Int32 || field.Kind() == reflect.Int64 {
				if field.Int() < 0 {
					field.SetInt(-field.Int())
				}
			}
		case "asciiify":
			if field.Kind() == reflect.String {
				field.SetString(toASCII(field.String()))
			}
		case "first_upper":
			if field.Kind() == reflect.String && len(field.String()) > 0 {
				field.SetString(strings.ToUpper(string(field.String()[0])) + field.String()[1:])
			}
		case "remove_leading_zeros":
			if field.Kind() == reflect.String {
				field.SetString(strings.TrimLeft(field.String(), "0"))
			}
		default:
			if strings.HasPrefix(tag, "default=") {
				defaultValue := tag[8:]
				if isEmptyValue(field) {
					if field.Kind() == reflect.String {
						field.SetString(defaultValue)
					}
				}
			} else if strings.HasPrefix(tag, "truncate=") {
				maxLen, _ := strconv.Atoi(tag[9:])
				if len(field.String()) > maxLen {
					field.SetString(field.String()[:maxLen])
				}
			} else if strings.HasPrefix(tag, "mask=") {
				indices := strings.Split(tag[5:], "-")
				start, _ := strconv.Atoi(indices[0])
				end, _ := strconv.Atoi(indices[1])
				if start < len(field.String()) && end <= len(field.String()) {
					masked := field.String()[:start] + strings.Repeat("*", end-start) + field.String()[end:]
					field.SetString(masked)
				}
			} else if strings.HasPrefix(tag, "pattern=") {
				pattern := tag[8:]
				re, err := regexp.Compile(pattern)
				if err != nil {
					return err
				}
				match := re.FindString(field.String())
				if match != "" {
					field.SetString(match)
				}
			} else if strings.HasPrefix(tag, "base64_encode") {
				if field.Kind() == reflect.String {
					field.SetString(base64.StdEncoding.EncodeToString([]byte(field.String())))
				}
			} else if strings.HasPrefix(tag, "base64_decode") {
				if field.Kind() == reflect.String {
					decoded, err := base64.StdEncoding.DecodeString(field.String())
					if err != nil {
						return err
					}
					field.SetString(string(decoded))
				}
			} else if strings.HasPrefix(tag, "hash=") {
				hashType := tag[5:]
				switch hashType {
				case "md5":
					hash := md5.Sum([]byte(field.String()))
					field.SetString(hex.EncodeToString(hash[:]))
				case "sha256":
					hash := sha256.Sum256([]byte(field.String()))
					field.SetString(hex.EncodeToString(hash[:]))
				case "sha512":
					hash := sha512.Sum512([]byte(field.String()))
					field.SetString(hex.EncodeToString(hash[:]))
				case "crc32":
					hash := crc32.ChecksumIEEE([]byte(field.String()))
					field.SetString(strconv.FormatUint(uint64(hash), 10))
				}
			} else if strings.HasPrefix(tag, "url_encode") {
				if field.Kind() == reflect.String {
					field.SetString(url.QueryEscape(field.String()))
				}
			} else if strings.HasPrefix(tag, "url_decode") {
				if field.Kind() == reflect.String {
					decoded, err := url.QueryUnescape(field.String())
					if err != nil {
						return err
					}
					field.SetString(decoded)
				}
			} else if strings.HasPrefix(tag, "sprintf=") {
				if field.Kind() == reflect.String {
					format := tag[8:]
					field.SetString(fmt.Sprintf(format, field.String()))
				}
			}
		}
	}
	return nil
}
