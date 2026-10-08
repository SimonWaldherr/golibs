// Package semver parses and compares semantic versions.
package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	semverRegex = regexp.MustCompile(`^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)
)

// Version represents a semantic version. Build metadata does not affect precedence.
type Version struct {
	Major         int
	Minor         int
	Patch         int
	PreRelease    string
	BuildMetadata string
}

// Parse parses a semantic version with an optional leading "v".
// Core version components must fit in an int.
func Parse(version string) (*Version, error) {
	matches := semverRegex.FindStringSubmatch(version)
	if matches == nil {
		return nil, errors.New("invalid semantic version")
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %w", err)
	}
	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %w", err)
	}
	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %w", err)
	}
	for _, identifier := range strings.Split(matches[4], ".") {
		if isNumeric(identifier) && len(identifier) > 1 && identifier[0] == '0' {
			return nil, errors.New("numeric pre-release identifier has a leading zero")
		}
	}

	return &Version{
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		PreRelease:    matches[4],
		BuildMetadata: matches[5],
	}, nil
}

// String returns the version with a leading "v".
func (v *Version) String() string {
	ver := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		ver += "-" + v.PreRelease
	}
	if v.BuildMetadata != "" {
		ver += "+" + v.BuildMetadata
	}
	return ver
}

// Compare returns -1, 0, or 1 if v precedes, equals, or follows other.
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		return compareInt(v.Major, other.Major)
	}
	if v.Minor != other.Minor {
		return compareInt(v.Minor, other.Minor)
	}
	if v.Patch != other.Patch {
		return compareInt(v.Patch, other.Patch)
	}
	return comparePreRelease(v.PreRelease, other.PreRelease)
}

func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func comparePreRelease(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return 1
	}
	if b == "" {
		return -1
	}
	aParts, bParts := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		aPart, bPart := aParts[i], bParts[i]
		if aPart == bPart {
			continue
		}
		aNumeric, bNumeric := isNumeric(aPart), isNumeric(bPart)
		if aNumeric && bNumeric {
			// Numeric identifiers have no leading zeroes, so compare their
			// lengths before their digits without limiting them to int size.
			if len(aPart) != len(bPart) {
				return compareInt(len(aPart), len(bPart))
			}
		} else if aNumeric {
			return -1
		} else if bNumeric {
			return 1
		}
		return strings.Compare(aPart, bPart)
	}
	return compareInt(len(aParts), len(bParts))
}

func isNumeric(identifier string) bool {
	if identifier == "" {
		return false
	}
	for i := 0; i < len(identifier); i++ {
		if identifier[i] < '0' || identifier[i] > '9' {
			return false
		}
	}
	return true
}

// Constraint compares a version against a single version requirement.
type Constraint struct {
	Operator string
	Version  *Version
}

// ParseConstraint parses =, !=, <, <=, >, or >= followed by a version.
// Whitespace between the operator and version is optional.
func ParseConstraint(constraint string) (*Constraint, error) {
	constraint = strings.TrimSpace(constraint)
	operator := ""
	for _, candidate := range []string{">=", "<=", "!=", "=", ">", "<"} {
		if strings.HasPrefix(constraint, candidate) {
			operator = candidate
			break
		}
	}
	if operator == "" {
		return nil, errors.New("invalid constraint format")
	}

	version, err := Parse(strings.TrimSpace(strings.TrimPrefix(constraint, operator)))
	if err != nil {
		return nil, err
	}

	return &Constraint{
		Operator: operator,
		Version:  version,
	}, nil
}

// Matches reports whether v satisfies the constraint.
func (c *Constraint) Matches(v *Version) bool {
	switch c.Operator {
	case "=":
		return v.Compare(c.Version) == 0
	case "!=":
		return v.Compare(c.Version) != 0
	case "<":
		return v.Compare(c.Version) < 0
	case "<=":
		return v.Compare(c.Version) <= 0
	case ">":
		return v.Compare(c.Version) > 0
	case ">=":
		return v.Compare(c.Version) >= 0
	default:
		return false
	}
}
