package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	semverRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([\w\.\-]+))?(?:\+([\w\.\-]+))?$`)
)

type Version struct {
	Major         int
	Minor         int
	Patch         int
	PreRelease    string
	BuildMetadata string
}

func Parse(version string) (*Version, error) {
	matches := semverRegex.FindStringSubmatch(version)
	if matches == nil {
		return nil, errors.New("invalid semantic version")
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	return &Version{
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		PreRelease:    matches[4],
		BuildMetadata: matches[5],
	}, nil
}

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
	return strings.Compare(a, b)
}

type Constraint struct {
	Operator string
	Version  *Version
}

func ParseConstraint(constraint string) (*Constraint, error) {
	parts := strings.Fields(constraint)
	if len(parts) != 2 {
		return nil, errors.New("invalid constraint format")
	}

	version, err := Parse(parts[1])
	if err != nil {
		return nil, err
	}

	return &Constraint{
		Operator: parts[0],
		Version:  version,
	}, nil
}

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
