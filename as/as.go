// Package as helps by converting various data types to various other types
package as

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ree struct {
	typ string
	re  string
}

type Dynamic struct {
	Type string
	Data interface{}
}

// regex (regular expression) in a slice of ree struct
var regex = [...]ree{
	{
		typ: "bool",
		re:  "(0|1|true|false)",
	}, {
		typ: "int",
		re:  "\\d+",
	}, {
		typ: "int",
		re:  "\\d+[\\^eE]\\d+",
	}, {
		typ: "float",
		re:  "[-+]?[0-9]*[\\.,]?[0-9]+([eE][-+]?[0-9]+)?",
	}, {
		typ: "price",
		re:  "((€|\\$|¢|£|EURO?|DM|USD) ?)?[-+]?[0-9]*[\\.,]?[0-9]+([eE][-+]?[0-9]+)?( ?(€|\\$|¢|£|EURO?|DM|USD))?",
	}, {
		typ: "url",
		re:  "(https?:\\/\\/)?([\\da-z\\.-]+)\\.([a-z\\.]{2,6})([\\/\\w \\.-]*)*\\/?",
	}, {
		typ: "ipv4",
		re:  "(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])",
	}, {
		typ: "ipv6",
		re:  "([0-9A-Fa-f]{1,4}:([0-9A-Fa-f]{1,4}:([0-9A-Fa-f]{1,4}:([0-9A-Fa-f]{1,4}:([0-9A-Fa-f]{1,4}:[0-9A-Fa-f]{0,4}|:[0-9A-Fa-f]{1,4})?|(:[0-9A-Fa-f]{1,4}){0,2})|(:[0-9A-Fa-f]{1,4}){0,3})|(:[0-9A-Fa-f]{1,4}){0,4})|:(:[0-9A-Fa-f]{1,4}){0,5})((:[0-9A-Fa-f]{1,4}){2}|:(25[0-5]|(2[0-4]|1[0-9]|[1-9])?[0-9])(\\.(25[0-5]|(2[0-4]|1[0-9]|[1-9])?[0-9])){3})|(([0-9A-Fa-f]{1,4}:){1,6}|:):[0-9A-Fa-f]{0,4}|([0-9A-Fa-f]{1,4}:){7}:",
	}, {
		typ: "mac",
		re:  "([0-9a-f]{1,2}:){5}([0-9a-f]{1,2})",
	}, {
		typ: "email",
		re:  "[a-z0-9\\._%\\+\\-]+\\@[a-z0-9\\.\\-]+\\.[a-z]{2,9}",
	}, {
		typ: "creditcard",
		re:  "(?:4\\d{12}(?:\\d{3})?|5[1-5]\\d{14}|3[47]\\d{13}|3(?:0[0-5]|[68]\\d)\\d{11}|6(?:011|5\\d{2})\\d{12}|(?:2131|1800|35\\d{3})\\d{11})",
	}, {
		typ: "color",
		re:  "#[a-f0-9]{2,6}",
	}, {
		typ: "color",
		re:  "(rgb|hsl|yuv)\\( *[\\d\\.%]+ *, *[\\d\\.%]+ *, *[\\d\\.%]+ *\\)",
	}, {
		typ: "color",
		re:  "(rgba|cmyk)\\( *\\d+ *, *\\d+ *, *\\d+ *, *\\d+ *\\)",
	}, {
		typ: "isbn",
		re:  "(1(?:(0)|3))?:?[- ]?(\\s)*[0-9]+[- ]?[0-9]+[- ]?[0-9]+[- ]?[0-9]*[- ]*[xX0-9]",
	}, {
		typ: "date",
		re:  "(?i:([MDCLXVI])((M{0,3})((C[DM])|(D?C{0,3}))?((X[LC])|(L?XX{0,2})|L)?((I[VX])|(V?(II{0,2}))|V)?))",
	}, {
		typ: "alpha",
		re:  "[a-zA-Z]+",
	}, {
		typ: "alphanumeric",
		re:  "[a-zA-Z0-9]+",
	}, {
		typ: "base64",
		re:  "(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})",
	}, {
		typ: "string",
		re:  "[[:print:]]+",
	}, {
		typ: "ascii",
		re:  "[[:ascii:]]+",
	},
}

// timeformats contains the supported time formats
// for the convertion to time.Time.
var timeformats = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"02.01.06",
	"01/02/06",
	"2006-01-02",
	"2006/01/02",
	"01/02/2006",
	"02.01.2006",
	"01/02/06 15:04",
	"2006-01-02 15:04",
	"2006-01-02T15:04",
	"01/02/2006 15:04",
	"02.01.06 15:04:05",
	"01/02/06 15:04:05",
	"01/02/2006 15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"02.01.2006 15:04:05",
	"2006-01-02 15:04:05 -0700 MST",
}

// Bool returns a boolean value.
// It mainly depends on the output of strconv.ParseBool,
// but also checks for integer values.
func Bool(valuea ...interface{}) bool {
	value := valuea[0]
	if Int(value) > 0 {
		return true
	}
	b, _ := strconv.ParseBool(String(value))
	return b
}

// Bytes returns a slice of bytes.
func Bytes(valuea ...interface{}) []byte {
	value := valuea[0]
	if value == nil {
		return []byte{}
	}

	switch val := value.(type) {
	case bool:
		if val == true {
			return []byte("true")
		}
		return []byte("false")
	case string:
		return []byte(val)
	case []byte:
		return val
	default:
		return []byte(fmt.Sprintf("%v", value))
	}
}

// Duration converts input values to time.Duration.
// It mainly depends on time.ParseDuration.
func Duration(valuea ...interface{}) time.Duration {
	value := valuea[0]
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return time.Duration(Int(value))
	case uint, uint8, uint16, uint32, uint64:
		return time.Duration(Int(value))
	case float32, float64:
		return time.Duration(Int(value))
	default:
		dur, _ := time.ParseDuration(String(value))
		return dur
	}
}

// FixedLengthAfter appends spacer chars after a string
func FixedLengthAfter(str string, spacer string, length int) string {
	spacer = spacer[:1]
	l := length - len(str)
	if l > 0 {
		return str + strings.Repeat(spacer, l)
	}
	if l == 0 {
		return str
	}
	return str[:length]
}

// FixedLengthBefore prepends spacer chars before a string
func FixedLengthBefore(str string, spacer string, length int) string {
	spacer = spacer[:1]
	l := length - len(str)
	if l > 0 {
		return strings.Repeat(spacer, l) + str
	}
	if l == 0 {
		return str
	}
	return str[:length]
}

// FixedLengthCenter adds spacer chars after and before a string
func FixedLengthCenter(str string, spacer string, length int) string {
	spacer = spacer[:1]
	l := length - len(str)
	if l > 0 {
		if l%2 == 0 {
			l = l / 2
			return strings.Repeat(spacer, l) + str + strings.Repeat(spacer, l)
		}
		l = (l + 1) / 2
		return strings.Repeat(spacer, l) + str + strings.Repeat(spacer, l-1)
	}
	if l == 0 {
		return str
	}
	return str[:length]
}

// Float converts it's input to type float64.
// int, uint and float gets converted as expected,
// time is transformed to a float of the corresponding timestamp.
// strings and byte slices gets converted via strconv.ParseFloat.
func Float(valuea ...interface{}) float64 {
	value := valuea[0]
	switch val := value.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return float64(val)
	case time.Time:
		return float64(val.Unix())
	case bool:
		if val == true {
			return float64(1)
		}
		return float64(0)
	default:
		f, _ := strconv.ParseFloat(String(value), 64)
		return float64(f)
	}
}

// FloatFromXString converts strings to float64.
// Most values can be converted to float via Float(),
// but floats as strings in e.g. german spelling
// should be converted with this function.
func FloatFromXString(valuea ...string) float64 {
	value := valuea[0]
	value = strings.Trim(value, "\t\n\r¢§$€ ")
	var float float64
	c := strings.Count(value, ",")
	p := strings.Count(value, ".")
	fc := strings.Index(value, ",")
	fp := strings.Index(value, ".")
	if c == 0 && p == 1 {
		float, _ = strconv.ParseFloat(value, 64)
	} else if c == 1 && p == 0 {
		value = strings.Replace(value, ",", ".", 1)
		float, _ = strconv.ParseFloat(value, 64)
	} else if c == 0 && p == 0 {
		intx, _ := strconv.ParseInt(value, 0, 64)
		float = float64(intx)
	} else if c > 1 && p < 2 {
		value = strings.Replace(value, ",", "", -1)
		float, _ = strconv.ParseFloat(value, 64)
	} else if c < 2 && p > 1 {
		value = strings.Replace(value, ".", "", -1)
		value = strings.Replace(value, ",", ".", 1)
		float, _ = strconv.ParseFloat(value, 64)
	} else if c == 1 && p == 1 {
		if fp < fc {
			value = strings.Replace(value, ".", "", -1)
			value = strings.Replace(value, ",", ".", 1)
		} else {
			value = strings.Replace(value, ",", "", -1)
		}
		float, _ = strconv.ParseFloat(value, 64)
	} else {
		value = "0"
		float, _ = strconv.ParseFloat(value, 64)
	}
	return float64(float)
}

// Int returns an int64 of the input value.
// Float values and float values in strings will be rounded via
// "round half towards positive infinity".
// strings get converted via strconv.ParseFloat.
func Int(valuea ...interface{}) int64 {
	value := valuea[0]
	switch val := value.(type) {
	case int:
		return int64(val)
	case int8:
		return int64(val)
	case int16:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return int64(val)
	case uint:
		return int64(val)
	case uint8:
		return int64(val)
	case uint16:
		return int64(val)
	case uint32:
		return int64(val)
	case uint64:
		return int64(val)
	case float32:
		return int64(val + 0.5)
	case float64:
		return int64(val + 0.5)
	case time.Time:
		return int64(val.Unix())
	case bool:
		if val == true {
			return int64(1)
		}
		return int64(0)
	default:
		i, _ := strconv.ParseFloat(String(value), 64)
		return int64(i + 0.5)
	}
}

// String converts input values to string.
// Time and Duration gets converted via standard functions.
// Most types gets "converted" via fmt.Sprintf.
func String(valuea ...interface{}) string {
	value := valuea[0]
	if value == nil {
		return ""
	}

	switch val := value.(type) {
	case bool:
		if value.(bool) == true {
			return "true"
		}
		return "false"
	case time.Duration:
		return string(val.String())
	case time.Time:
		return string(val.Format(time.RFC3339))
	case string:
		return val
	case []byte:
		return string(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Time converts inputs values to time.Time.
// Time formats in the variable timeformats can be used.
func Time(valuea ...interface{}) time.Time {
	value := valuea[0]
	s := String(value)
	for _, format := range timeformats {
		r, err := time.Parse(format, s)
		if err == nil {
			return r
		}
	}
	return time.Time{}
}

// Trimmed takes the first given value, converts it to
// a string, trims the whitespace an returns it.
func Trimmed(valuea ...interface{}) string {
	value := valuea[0]
	return strings.TrimSpace(String(value))
}

// Uint returns an uint64 of the input value.
// Float values and float values in strings will be rounded via
// "round half towards positive infinity".
// strings get converted via strconv.ParseFloat.
func Uint(valuea ...interface{}) uint64 {
	value := valuea[0]
	switch val := value.(type) {
	case int:
		return uint64(val)
	case int8:
		return uint64(val)
	case int16:
		return uint64(val)
	case int32:
		return uint64(val)
	case int64:
		return uint64(val)
	case uint:
		return uint64(val)
	case uint8:
		return uint64(val)
	case uint16:
		return uint64(val)
	case uint32:
		return uint64(val)
	case uint64:
		return uint64(val)
	case float32:
		return uint64(val + 0.5)
	case float64:
		return uint64(val + 0.5)
	case time.Time:
		return uint64(val.Unix())
	case bool:
		if val == true {
			return uint64(1)
		}
		return uint64(0)
	default:
		i, _ := strconv.ParseFloat(String(value), 64)
		return uint64(i + 0.5)
	}
}

// Type returns a type (string) of a string.
func Type(valuea ...interface{}) (string, error) {
	var err error
	str := strings.Trim(String(valuea[0]), " \t\n\r")
	if !Time(str).IsZero() {
		return "date", nil
	}
	for _, b := range regex {
		var match bool
		re := "(?i)^" + b.re + "$"
		if match, err = regexp.MatchString(re, str); match == true {
			//fmt.Printf("%v tested for %v with %v; result: %v\n", str, b.typ, b.re, match)
			return b.typ, nil
		}
		//fmt.Printf("%v tested for %v with %v; result: %v\n", str, b.typ, b.re, match)
	}
	return "", err
}

// DBType returns a Database Type of a string.
func DBType(str string) string {
	t, err := Type(str)
	if err != nil {
		return "string"
	}
	switch t {
	case "bool", "int", "string", "float":
		return t
	default:
		return "string"
	}
}

// DBTypeMultiple returns the lowest common denominator of a type for all inserted DBTypes
func DBTypeMultiple(val []string) string {
	var typeint int
	for _, typ := range val {
		for i, b := range regex {
			if b.typ == typ {
				if typeint < i {
					typeint = i
				}
			}
		}
	}
	return regex[typeint].typ
}
