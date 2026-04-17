package structs

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"hash/crc32"
	"net/url"
	"strconv"
)

// JSON validates and modifies the given struct pointer, then marshals it to a JSON string.
func JSON(s interface{}) (string, error) {
	if err := ValidateAndModify(s); err != nil {
		return "", err
	}
	b, err := json.Marshal(s)
	return string(b), err
}

// NestedStruct validates and modifies the given struct pointer.
// Nested struct fields are validated according to their own tags
// when accessed via the validate:"required" constraint.
func NestedStruct(s interface{}) error {
	return ValidateAndModify(s)
}

// Base64 encodes the given string using standard base64 encoding.
func Base64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Hash computes a hash of the given string using the specified algorithm.
// Supported algorithms: "md5", "sha256", "sha512", "crc32".
// Note: MD5 and CRC32 are not cryptographically secure; use sha256 or sha512
// for security-sensitive operations.
// Returns an empty string for unsupported algorithms.
func Hash(s, algorithm string) string {
	switch algorithm {
	case "md5":
		h := md5.Sum([]byte(s))
		return hex.EncodeToString(h[:])
	case "sha256":
		h := sha256.Sum256([]byte(s))
		return hex.EncodeToString(h[:])
	case "sha512":
		h := sha512.Sum512([]byte(s))
		return hex.EncodeToString(h[:])
	case "crc32":
		h := crc32.ChecksumIEEE([]byte(s))
		return strconv.FormatUint(uint64(h), 10)
	}
	return ""
}

// URLEncode encodes the given string for safe use in URL query parameters.
func URLEncode(s string) string {
	return url.QueryEscape(s)
}

// CustomValidation is a type alias for CustomValidatorFunc,
// used when registering custom field validators via RegisterCustomValidator.
type CustomValidation = CustomValidatorFunc
