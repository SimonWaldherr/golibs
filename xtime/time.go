package xtime

import (
	"time"
)

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
