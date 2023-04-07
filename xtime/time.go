// xt is a collection of time-related functions and types
package xtime

import (
	"time"
)

// conv is a map of strftime format specifiers
var conv = map[rune]string{
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
	'%': "",
}

// StrfTime implements a subset of strftime
// http://man7.org/linux/man-pages/man3/strftime.3.html
func StrfTime(format string, t time.Time) string {
	ret := make([]byte, 0, len(format))
	for i := 0; i < len(format); i++ {
		if format[i] == '%' {
			if layout, ok := conv[rune(format[i+1])]; ok {
				ret = append(ret, []byte(t.Format(layout))...)
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

// TimeRange is a simple struct to represent a time range
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
	return tr.Start.Before(other.End) && tr.End.After(other.Start)
}

// Contains returns true if the time range contains another time range
func (tr TimeRange) Contains(other TimeRange) bool {
	return tr.Start.Before(other.Start) && tr.End.After(other.End)
}

// Within returns true if the time range is within another time range
func (tr TimeRange) Within(other TimeRange) bool {
	return other.Contains(tr)
}

// And returns the intersection of two time ranges
func And(tr1, tr2 TimeRange) TimeRange {
	if tr1.Overlaps(tr2) {
		return TimeRange{
			Start: tr1.Start,
			End:   tr2.End,
		}
	}
	return TimeRange{}
}

// Or returns the union of two time ranges
func Or(tr1, tr2 TimeRange) TimeRange {
	if tr1.Overlaps(tr2) {
		return TimeRange{
			Start: tr1.Start,
			End:   tr1.End,
		}
	}
	return TimeRange{
		Start: tr1.Start,
		End:   tr2.End,
	}
}

// Not returns the inverse of a time range
func Not(tr TimeRange, fullRange TimeRange) TimeRange {
	return TimeRange{
		Start: fullRange.Start,
		End:   fullRange.End,
	}.Subtract(tr)
}

// Subtract returns the difference of two time ranges
func (tr TimeRange) Subtract(other TimeRange) TimeRange {
	if !tr.Overlaps(other) {
		return tr
	}
	if tr.Contains(other) {
		return TimeRange{
			Start: tr.Start,
			End:   other.Start,
		}
	}
	if tr.Within(other) {
		return TimeRange{}
	}
	if tr.Start.Before(other.Start) {
		return TimeRange{
			Start: tr.Start,
			End:   other.Start,
		}
	}
	return TimeRange{
		Start: other.End,
		End:   tr.End,
	}
}
