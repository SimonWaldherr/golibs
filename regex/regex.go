// regex is a wrapper for the standard regexp package.
// It automates the regexp.Compile process for you.
package regex

import (
	"regexp"
)

var regexArray = make(map[string]*regexp.Regexp)

// ReplaceAllString returns a copy of src, replacing matches of the regular expression
// with the replacement string replace.  Inside replace, $ signs are interpreted as
// in Expand, so for instance $1 represents the text of the first submatch.
func ReplaceAllString(src, regex, replace string) string {
	if re, ok := regexArray[regex]; ok {
		return re.ReplaceAllString(src, replace)
	}
	regexArray[regex], _ = regexp.Compile(regex)
	return ReplaceAllString(src, regex, replace)
}

// ReplaceAllStringFunc returns a copy of src in which all matches of the
// regular expression have been replaced by the return value of function replace applied
// to the matched substring.  The replacement returned by replace is substituted
// directly, without using Expand.
func ReplaceAllStringFunc(src, regex string, replace func(s string) string) string {
	if re, ok := regexArray[regex]; ok {
		return re.ReplaceAllStringFunc(src, replace)
	}
	regexArray[regex], _ = regexp.Compile(regex)
	return ReplaceAllStringFunc(src, regex, replace)
}

// FindAllString returns a slice of all strings holding the text of the leftmost
// match in src of the regular expression.  If there is no match, the return value is nil.
// It will be empty if the regular expression successfully matches an empty string.
func FindAllString(src, regex string) []string {
	if re, ok := regexArray[regex]; ok {
		return re.FindAllString(src, -1)
	}
	regexArray[regex], _ = regexp.Compile(regex)
	return FindAllString(src, regex)
}

// FindAllStringSubmatch returns a slice of a slice of strings holding the text of the
// leftmost match of the regular expression in src and the matches.
// A return value of nil indicates no match.
func FindAllStringSubmatch(src, regex string) [][]string {
	if re, ok := regexArray[regex]; ok {
		return re.FindAllStringSubmatch(src, -1)
	}
	regexArray[regex], _ = regexp.Compile(regex)
	return FindAllStringSubmatch(src, regex)
}
