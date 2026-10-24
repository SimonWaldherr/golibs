package structs

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestPatternSyntaxKeepsFollowingRules(t *testing.T) {
	for _, tc := range []struct {
		pattern string
		value   string
	}{
		{`foo{bar`, "foo{bar,minlen=100"},
		{`\Q(\E`, "(,minlen=100"},
		{`^\Qa,b\E$`, "a,b"},
		{`(a}b,c)`, "a}b,c"},
		{`^[]a,]$`, "]"},
	} {
		t.Run(tc.pattern, func(t *testing.T) {
			expression, err := regexp.Compile(tc.pattern)
			if err != nil || !expression.MatchString(tc.value) {
				t.Fatalf("invalid test pattern/value: %v", err)
			}
			metadata := reflect.StructTag("validate:" + strconv.Quote("pattern="+tc.pattern+",minlen=100"))
			tags, err := getCombinedTags(metadata)
			want := []string{"validate:pattern=" + tc.pattern, "validate:minlen=100"}
			if err != nil || !reflect.DeepEqual(tags, want) {
				t.Fatalf("pattern consumed or split the following rule: %v, %v; want %v", tags, err, want)
			}
		})
	}
}

func TestRuleParameterWhitespaceIsPreserved(t *testing.T) {
	for _, tc := range []struct{ key, rule string }{
		{"validate", "pattern=foo "},
		{"modify", "default= value "},
		{"modify", "sprintf=%s "},
	} {
		metadata := reflect.StructTag(tc.key + ":" + strconv.Quote(tc.rule))
		tags, err := getCombinedTags(metadata)
		want := []string{tc.key + ":" + tc.rule}
		if err != nil || !reflect.DeepEqual(tags, want) {
			t.Errorf("significant parameter whitespace changed: %v, %v; want %v", tags, err, want)
		}
	}
}

func FuzzCombinedTagMetadata(f *testing.F) {
	for _, seed := range []string{
		"", "validate:", `validate:"`, `validate:"required" modify:"trim"`,
		`json:"name,omitempty" modify:"trim,upper" validate:"pattern=^[A-Z]{2,4}$"`,
		`validate:"pattern=\\Q(\\E,minlen=100"`,
		`validate:"pattern=foo{bar,minlen=100"`,
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		// Exercise arbitrary malformed metadata directly, without synthesizing
		// reflect.StructOf types that Go retains permanently.
		tags, err := getCombinedTags(reflect.StructTag(raw))
		if err == nil {
			for _, tag := range tags {
				if !(strings.HasPrefix(tag, "validate:") || strings.HasPrefix(tag, "modify:")) {
					t.Fatalf("unexpected rule prefix: %q", tag)
				}
			}
		}
		// Quoting arbitrary unrelated metadata must not inject, swallow, or
		// reorder the following known tags.
		metadata := "unrelated:" + strconv.Quote(raw) + ` modify:"trim" validate:"required"`
		tags, err = getCombinedTags(reflect.StructTag(metadata))
		want := []string{"modify:trim", "validate:required"}
		if err != nil || !reflect.DeepEqual(tags, want) {
			t.Fatalf("quoted metadata changed ordered rules: %v, %v", tags, err)
		}
	})
}
