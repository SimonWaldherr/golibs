package regex

import (
	"regexp"
)

var regexArray = make(map[string]*regexp.Regexp)

func ReplaceAllString(str, regex, replace string) string {
	if re, ok := regexArray[regex]; ok {
		return re.ReplaceAllString(str, replace)
	}
	regexArray[regex], _ = regexp.Compile(regex)
	return ReplaceAllString(str, regex, replace)
}

func ReplaceAllStringFunc(str, regex string, replace func(s string) string) string {
	if re, ok := regexArray[regex]; ok {
		return re.ReplaceAllStringFunc(str, replace)
	}
	regexArray[regex], _ = regexp.Compile(regex)
	return ReplaceAllStringFunc(str, regex, replace)
}

func FindAllString(str, regex string) []string {
	if re, ok := regexArray[regex]; ok {
		return re.FindAllString(str, -1)
	}
	regexArray[regex], _ = regexp.Compile(regex)
	return FindAllString(str, regex)
}

func FindAllStringSubmatch(str, regex string) [][]string {
	if re, ok := regexArray[regex]; ok {
		return re.FindAllStringSubmatch(str, -1)
	}
	regexArray[regex], _ = regexp.Compile(regex)
	return FindAllStringSubmatch(str, regex)
}
