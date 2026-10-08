// Package regex is a wrapper for the standard regexp package.
// It automates the regexp.Compile process for you.
package regex

import (
	"regexp"
	"sync"
)

var regexArray = make(map[string]*regexp.Regexp)
var regexLock sync.RWMutex

// compile reuses cached expressions without holding a lock during matching or
// user callbacks. Invalid expressions are never added to the cache.
func compile(pattern string) (*regexp.Regexp, error) {
	regexLock.RLock()
	re, ok := regexArray[pattern]
	regexLock.RUnlock()
	if ok {
		return re, nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	regexLock.Lock()
	if cached, ok := regexArray[pattern]; ok {
		re = cached
	} else {
		regexArray[pattern] = re
	}
	regexLock.Unlock()
	return re, nil
}

// CheckRegex checks if the regex is valid
func CheckRegex(regex string) error {
	_, err := regexp.Compile(regex)
	return err
}

// CacheRegex caches the compiled regex
func CacheRegex(regex string) error {
	_, err := compile(regex)
	return err
}

// MatchString reports whether the string src contains any match of the regular expression.
func MatchString(src, regex string) (bool, error) {
	re, err := compile(regex)
	if err != nil {
		return false, err
	}
	return re.MatchString(src), nil
}

// ReplaceAllString returns a copy of src, replacing matches of the regular expression
// with the replacement string replace.  Inside replace, $ signs are interpreted as
// in Expand, so for instance $1 represents the text of the first submatch.
func ReplaceAllString(src, regex, replace string) (string, error) {
	re, err := compile(regex)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(src, replace), nil
}

// ReplaceAllStringFunc returns a copy of src in which all matches of the
// regular expression have been replaced by the return value of function replace applied
// to the matched substring.  The replacement returned by replace is substituted
// directly, without using Expand.
func ReplaceAllStringFunc(src, regex string, replace func(s string) string) (string, error) {
	re, err := compile(regex)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllStringFunc(src, replace), nil
}

// FindAllString returns a slice of all strings holding the text of the leftmost
// match in src of the regular expression.  If there is no match, the return value is nil.
// It will be empty if the regular expression successfully matches an empty string.
func FindAllString(src, regex string) ([]string, error) {
	re, err := compile(regex)
	if err != nil {
		return []string{}, err
	}
	return re.FindAllString(src, -1), nil
}

// FindAllStringSubmatch returns a slice of a slice of strings holding the text of the
// leftmost match of the regular expression in src and the matches.
// A return value of nil indicates no match.
func FindAllStringSubmatch(src, regex string) ([][]string, error) {
	re, err := compile(regex)
	if err != nil {
		return [][]string{}, err
	}
	return re.FindAllStringSubmatch(src, -1), nil
}

// Count returns the number of cached regex
func Count() int {
	regexLock.RLock()
	defer regexLock.RUnlock()
	return len(regexArray)
}

// Cleanup cleans the cache
func Cleanup() {
	regexLock.Lock()
	defer regexLock.Unlock()
	regexArray = make(map[string]*regexp.Regexp)
}
