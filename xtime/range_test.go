package xtime

import (
	"strings"
	"testing"
	"time"
)

func TestStrfTimeDirectives(t *testing.T) {
	tm := time.Date(2026, 10, 8, 15, 4, 5, 123000000, time.FixedZone("CEST", 2*60*60))
	cases := map[string]string{
		"%a": "Thu", "%A": "Thursday", "%b": "Oct", "%B": "October",
		"%d": "08", "%D": "10-08-2026", "%F": "2026-10-08", "%H": "15",
		"%I": "03", "%L": ".123", "%m": "10", "%M": "04", "%p": "PM",
		"%S": "05", "%x": "2026-10-08", "%X": "15:04:05", "%y": "26",
		"%Y": "2026", "%z": "+0200", "%Z": "CEST", "%0": "", "%%": "%",
		"%": "%", "%q": "%q", "%%%Y%": "%2026%", "日%Y年": "日2026年", "": "",
	}
	for format, want := range cases {
		if got := StrfTime(format, tm); got != want {
			t.Errorf("StrfTime(%q) = %q, want %q", format, got, want)
		}
	}
}

func hours(start, end int) TimeRange {
	base := time.Date(2026, 10, 8, 0, 0, 0, 0, time.UTC)
	return TimeRange{Start: base.Add(time.Duration(start) * time.Hour), End: base.Add(time.Duration(end) * time.Hour)}
}

func TestRangeOperations(t *testing.T) {
	cases := []struct {
		name                               string
		a, b                               TimeRange
		overlaps, contains                 bool
		intersection, envelope, difference TimeRange
	}{
		{"overlap right", hours(0, 4), hours(2, 6), true, false, hours(2, 4), hours(0, 6), hours(0, 2)},
		{"overlap left", hours(2, 6), hours(0, 4), true, false, hours(2, 4), hours(0, 6), hours(4, 6)},
		{"inner", hours(0, 6), hours(2, 4), true, true, hours(2, 4), hours(0, 6), hours(0, 2)},
		{"outer", hours(2, 4), hours(0, 6), true, false, hours(2, 4), hours(0, 6), TimeRange{}},
		{"equal", hours(0, 4), hours(0, 4), true, true, hours(0, 4), hours(0, 4), TimeRange{}},
		{"same start", hours(0, 6), hours(0, 4), true, true, hours(0, 4), hours(0, 6), hours(4, 6)},
		{"same end", hours(0, 6), hours(2, 6), true, true, hours(2, 6), hours(0, 6), hours(0, 2)},
		{"touching", hours(0, 2), hours(2, 4), false, false, TimeRange{}, hours(0, 4), hours(0, 2)},
		{"disjoint reverse", hours(4, 6), hours(0, 2), false, false, TimeRange{}, hours(0, 6), hours(4, 6)},
		{"empty left", hours(2, 2), hours(0, 4), false, false, TimeRange{}, hours(0, 4), TimeRange{}},
		{"empty right", hours(0, 4), hours(2, 2), false, false, TimeRange{}, hours(0, 4), hours(0, 4)},
		{"inverted", hours(4, 2), hours(3, 6), false, false, TimeRange{}, hours(3, 6), TimeRange{}},
		{"both empty", TimeRange{}, hours(4, 2), false, false, TimeRange{}, TimeRange{}, TimeRange{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.a.Overlaps(tc.b); got != tc.overlaps || got != tc.b.Overlaps(tc.a) {
				t.Errorf("Overlaps = %v, want %v and symmetry", got, tc.overlaps)
			}
			if got := tc.a.Contains(tc.b); got != tc.contains || got != tc.b.Within(tc.a) {
				t.Errorf("Contains = %v, want %v and inverse Within", got, tc.contains)
			}
			if got := And(tc.a, tc.b); got != tc.intersection || got != And(tc.b, tc.a) {
				t.Errorf("And = %v, want %v and symmetry", got, tc.intersection)
			}
			if got := Or(tc.a, tc.b); got != tc.envelope || got != Or(tc.b, tc.a) {
				t.Errorf("Or = %v, want %v and symmetry", got, tc.envelope)
			}
			if got := tc.a.Subtract(tc.b); got != tc.difference || got != Not(tc.b, tc.a) {
				t.Errorf("Subtract = %v, want %v and equivalent Not", got, tc.difference)
			}
		})
	}
	if got := hours(0, 6).Duration(); got != 6*time.Hour {
		t.Fatalf("Duration = %v", got)
	}
}

func FuzzStrfTimeLiterals(f *testing.F) {
	for _, seed := range []string{"", "%", "%%", "日%Y年", "trailing%"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, literal string) {
		escaped := strings.ReplaceAll(literal, "%", "%%")
		if got := StrfTime(escaped, time.Unix(0, 0)); got != literal {
			t.Fatalf("escaped literal %q became %q", literal, got)
		}
	})
}
