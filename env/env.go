// Package env provides helpers for reading environment variables with defaults
// and type conversions.
package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// String returns the value of the environment variable named by key,
// or defaultVal if the variable is not set or is empty.
func String(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// Int returns the value of the environment variable named by key as an int.
// If the variable is not set, is empty, or cannot be parsed, defaultVal is returned.
func Int(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return i
}

// Int64 returns the value of the environment variable named by key as an int64.
// If the variable is not set, is empty, or cannot be parsed, defaultVal is returned.
func Int64(key string, defaultVal int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultVal
	}
	return i
}

// Float64 returns the value of the environment variable named by key as a float64.
// If the variable is not set, is empty, or cannot be parsed, defaultVal is returned.
func Float64(key string, defaultVal float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultVal
	}
	return f
}

// Bool returns the value of the environment variable named by key as a bool.
// Accepted true values (case-insensitive): "1", "t", "true", "yes", "on".
// Accepted false values (case-insensitive): "0", "f", "false", "no", "off".
// If the variable is not set, is empty, or cannot be parsed, defaultVal is returned.
func Bool(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	switch strings.ToLower(v) {
	case "1", "t", "true", "yes", "on":
		return true
	case "0", "f", "false", "no", "off":
		return false
	}
	return defaultVal
}

// Duration returns the value of the environment variable named by key as a time.Duration.
// The value must be a valid Go duration string (e.g. "5s", "1m30s", "2h").
// If the variable is not set, is empty, or cannot be parsed, defaultVal is returned.
func Duration(key string, defaultVal time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultVal
	}
	return d
}

// Slice returns the value of the environment variable named by key split on sep.
// If the variable is not set or is empty, defaultVal is returned.
func Slice(key, sep string, defaultVal []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	parts := strings.Split(v, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return defaultVal
	}
	return result
}

// MustString returns the value of the environment variable named by key.
// It panics if the variable is not set or is empty.
func MustString(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("required environment variable not set: " + key)
	}
	return v
}

// IsSet returns true if the environment variable named by key is set (even if empty).
func IsSet(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// Map returns a map of all environment variables that have the given prefix.
// The prefix is stripped from the keys in the returned map.
func Map(prefix string) map[string]string {
	result := make(map[string]string)
	for _, e := range os.Environ() {
		idx := strings.Index(e, "=")
		if idx < 0 {
			continue
		}
		k := e[:idx]
		v := e[idx+1:]
		if strings.HasPrefix(k, prefix) {
			result[strings.TrimPrefix(k, prefix)] = v
		}
	}
	return result
}
