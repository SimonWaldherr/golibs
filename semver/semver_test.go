package semver

import (
	"testing"
)

func TestParse(t *testing.T) {
	version, err := Parse("1.2.3-alpha+001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if version.Major != 1 || version.Minor != 2 || version.Patch != 3 || version.PreRelease != "alpha" || version.BuildMetadata != "001" {
		t.Fatalf("parsed version incorrect: %v", version)
	}
}

func TestString(t *testing.T) {
	version, _ := Parse("1.2.3-alpha+001")
	if version.String() != "v1.2.3-alpha+001" {
		t.Fatalf("expected version string 'v1.2.3-alpha+001', got %s", version.String())
	}
}

func TestCompare(t *testing.T) {
	v1, _ := Parse("1.2.3")
	v2, _ := Parse("2.0.0")
	if v1.Compare(v2) >= 0 {
		t.Fatalf("expected v1 < v2")
	}
}

func TestConstraint(t *testing.T) {
	constraint, err := ParseConstraint(">= 1.2.3")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	version, _ := Parse("1.2.3")
	if !constraint.Matches(version) {
		t.Fatalf("expected version to match constraint")
	}
}
