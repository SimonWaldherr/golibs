// Package xtime provides strftime formatting and operations on time ranges.
package xtime

import (
	"time"
)

// conv is a map of strftime format specifiers
var conv = map[byte]string{
	'a': "Mon",
	'A': "Monday",
	'b': "Jan",
	'B': "January",
	'd': "02",
	'D': "01-02-2006",
	'F': "2006-01-02",
	'H': "15",
	'I': "03",
	'L': ".000",
	'm': "01",
	'M': "04",
	'p': "PM",
	'S': "05",
	'x': "2006-01-02",
	'X': "15:04:05",
	'y': "06",
	'Y': "2006",
	'z': "-0700",
	'Z': "MST",
	'0': "",
}

// StrfTime implements a subset of strftime. Supported directives are
// %a %A %b %B %d %D %F %H %I %L %m %M %p %S %x %X %y %Y %z %Z and %%.
// Unknown directives and trailing percent signs are preserved literally.
// The legacy %0 directive produces no output. Formatting uses t's location.
func StrfTime(format string, t time.Time) string {
	ret := make([]byte, 0, len(format))
	for i := 0; i < len(format); i++ {
		if format[i] == '%' && i+1 < len(format) {
			if format[i+1] == '%' {
				ret = append(ret, '%')
				i++
			} else if layout, ok := conv[format[i+1]]; ok {
				ret = t.AppendFormat(ret, layout)
				i++
			} else {
				ret = append(ret, format[i])
			}
		} else {
			ret = append(ret, format[i])
		}
	}
	return string(ret)
}

// Fmt (xtime.Fmt) is an alias for StrfTime
func Fmt(format string, t time.Time) string {
	return StrfTime(format, t)
}

// FmtNow is like StrfTime, but automatically with the current local time
func FmtNow(format string) string {
	t := time.Now()
	return StrfTime(format, t)
}

// TimeRange represents [Start, End). Start is included and End is excluded.
// A range is empty when End is not after Start.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Duration returns the duration of the time range
func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// Overlaps returns true if the time range overlaps with another time range
func (tr TimeRange) Overlaps(other TimeRange) bool {
	return !tr.empty() && !other.empty() && tr.Start.Before(other.End) && tr.End.After(other.Start)
}

// Contains reports whether a nonempty range contains another nonempty range.
// Equal boundaries are permitted, so a nonempty range contains itself.
func (tr TimeRange) Contains(other TimeRange) bool {
	return !tr.empty() && !other.empty() && !tr.Start.After(other.Start) && !tr.End.Before(other.End)
}

// Within returns true if the time range is within another time range
func (tr TimeRange) Within(other TimeRange) bool {
	return other.Contains(tr)
}

// And returns the intersection of two time ranges
func And(tr1, tr2 TimeRange) TimeRange {
	if tr1.Overlaps(tr2) {
		start, end := tr1.Start, tr1.End
		if tr2.Start.After(start) {
			start = tr2.Start
		}
		if tr2.End.Before(end) {
			end = tr2.End
		}
		return TimeRange{Start: start, End: end}
	}
	return TimeRange{}
}

// Or returns the smallest range covering both nonempty ranges. For disjoint
// ranges this includes the gap, because a single TimeRange cannot express a union.
// Empty ranges are ignored.
func Or(tr1, tr2 TimeRange) TimeRange {
	if tr1.empty() {
		if tr2.empty() {
			return TimeRange{}
		}
		return tr2
	}
	if tr2.empty() {
		return tr1
	}
	start, end := tr1.Start, tr1.End
	if tr2.Start.Before(start) {
		start = tr2.Start
	}
	if tr2.End.After(end) {
		end = tr2.End
	}
	return TimeRange{Start: start, End: end}
}

// Not returns fullRange.Subtract(tr), keeping the earliest remaining segment.
func Not(tr TimeRange, fullRange TimeRange) TimeRange {
	return fullRange.Subtract(tr)
}

// Subtract removes other from tr. When removal splits tr, only the earliest
// remaining segment is returned. A completely removed or empty range returns zero.
func (tr TimeRange) Subtract(other TimeRange) TimeRange {
	if tr.empty() {
		return TimeRange{}
	}
	if !tr.Overlaps(other) {
		return tr
	}
	if tr.Start.Before(other.Start) {
		return TimeRange{
			Start: tr.Start,
			End:   other.Start,
		}
	}
	if !other.End.Before(tr.End) {
		return TimeRange{}
	}
	return TimeRange{
		Start: other.End,
		End:   tr.End,
	}
}

func (tr TimeRange) empty() bool { return !tr.Start.Before(tr.End) }
