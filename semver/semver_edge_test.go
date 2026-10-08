package semver

import (
	"strconv"
	"strings"
	"testing"
)

func TestParseValidIdentifiers(t *testing.T) {
	for _, input := range []string{
		"0.0.0", "v1.2.3", "1.2.3-0", "1.2.3-alpha.1", "1.2.3-x-y-z.--",
		"1.2.3+001", "1.2.3-alpha+build.001", "1.2.3-999999999999999999999999999999",
		strconv.Itoa(int(^uint(0)>>1)) + ".0.0",
	} {
		t.Run(input, func(t *testing.T) {
			v, err := Parse(input)
			if err != nil {
				t.Fatalf("Parse(%q): %v", input, err)
			}
			if got := v.String(); got != "v"+strings.TrimPrefix(input, "v") {
				t.Errorf("String() = %q, want %q", got, "v"+strings.TrimPrefix(input, "v"))
			}
		})
	}
}

func TestParseRejectsInvalidVersions(t *testing.T) {
	for _, input := range []string{
		"", "1", "1.2", "1.2.3.4", "1.2.-3", "01.2.3", "1.02.3", "1.2.03",
		"1.2.3-01", "1.2.3-alpha.01", "1.2.3-alpha_1", "1.2.3-alpha..1",
		"1.2.3-.alpha", "1.2.3-alpha.", "1.2.3-", "1.2.3+", "1.2.3+build_1",
		"1.2.3+build..1", "1.2.3+build.", "1.2.3+é", " 1.2.3", "1.2.3 ",
		"999999999999999999999999999999.0.0", "1.999999999999999999999999999999.0",
		"1.2.999999999999999999999999999999",
	} {
		t.Run(input, func(t *testing.T) {
			if v, err := Parse(input); err == nil || v != nil {
				t.Fatalf("Parse(%q) = %v, %v; want nil and an error", input, v, err)
			}
		})
	}
}

func TestComparePrecedence(t *testing.T) {
	// SemVer's canonical precedence chain includes numeric and text identifiers.
	inputs := []string{
		"1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-alpha.beta", "1.0.0-beta",
		"1.0.0-beta.2", "1.0.0-beta.11", "1.0.0-rc.1", "1.0.0", "1.0.1", "1.1.0", "2.0.0",
	}
	versions := make([]*Version, len(inputs))
	for i, input := range inputs {
		var err error
		versions[i], err = Parse(input)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i, a := range versions {
		for j, b := range versions {
			want := compareInt(i, j)
			if got := a.Compare(b); got != want {
				t.Errorf("%s.Compare(%s) = %d, want %d", inputs[i], inputs[j], got, want)
			}
		}
	}
}

func TestCompareLargeNumericIdentifiers(t *testing.T) {
	pairs := [][2]string{
		{"1.0.0-9", "1.0.0-10"},
		{"1.0.0-999999999999999999999999999999", "1.0.0-1000000000000000000000000000000"},
		{"1.0.0-1000000000000000000000000000000", "1.0.0-1000000000000000000000000000001"},
		{"1.0.0-999999999999999999999999999999", "1.0.0--"},
	}
	for _, pair := range pairs {
		a, err := Parse(pair[0])
		if err != nil {
			t.Fatal(err)
		}
		b, err := Parse(pair[1])
		if err != nil {
			t.Fatal(err)
		}
		if a.Compare(b) != -1 || b.Compare(a) != 1 {
			t.Errorf("expected %s < %s", a, b)
		}
	}
}

func TestCompareIgnoresBuildMetadata(t *testing.T) {
	a, err := Parse("1.2.3-alpha+build.1")
	if err != nil {
		t.Fatal(err)
	}
	b, err := Parse("1.2.3-alpha+build.2")
	if err != nil {
		t.Fatal(err)
	}
	if a.Compare(b) != 0 {
		t.Errorf("build metadata changed precedence: %s, %s", a, b)
	}
}

func TestConstraintOperators(t *testing.T) {
	versions := []string{"1.2.2", "1.2.3", "1.2.4"}
	tests := []struct {
		operator string
		want     [3]bool
	}{
		{"=", [3]bool{false, true, false}},
		{"!=", [3]bool{true, false, true}},
		{"<", [3]bool{true, false, false}},
		{"<=", [3]bool{true, true, false}},
		{">", [3]bool{false, false, true}},
		{">=", [3]bool{false, true, true}},
	}
	for _, tt := range tests {
		for _, separator := range []string{"", " ", "\t"} {
			input := " \t" + tt.operator + separator + "v1.2.3 \n"
			c, err := ParseConstraint(input)
			if err != nil {
				t.Fatalf("ParseConstraint(%q): %v", input, err)
			}
			for i, version := range versions {
				v, err := Parse(version)
				if err != nil {
					t.Fatal(err)
				}
				if got := c.Matches(v); got != tt.want[i] {
					t.Errorf("constraint %q matches %s = %v, want %v", input, version, got, tt.want[i])
				}
			}
		}
	}
}

func TestParseConstraintRejectsInvalidInput(t *testing.T) {
	for _, input := range []string{"", "1.2.3", "== 1.2.3", "~ 1.2.3", "nonsense 1.2.3", ">=", "> =1.2.3", ">= 1.2.3 extra", ">= 1.2.3-01"} {
		if c, err := ParseConstraint(input); err == nil || c != nil {
			t.Errorf("ParseConstraint(%q) = %v, %v; want nil and an error", input, c, err)
		}
	}
}
