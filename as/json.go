package as

import (
	"strconv"
)

type Int64 int64

// MarshalJSON implements json.Marshaler interface for the Int64 Type.
func (i Int64) MarshalJSON() ([]byte, error) {
	return Bytes(i), nil
}

// UnmarshalJSON implements json.Unmarshaler inferface for the Int64 Type.
func (i *Int64) UnmarshalJSON(buf []byte) error {
	var i64 Int64
	str, err := strconv.Unquote(string(buf))
	if err == nil {
		i64 = Int64(Int(str))
	} else {
		i64 = Int64(Int(buf))
	}

	*i = i64
	return nil
}

// UnmarshalJSON implements json.Unmarshaler inferface for the Dynamic Type.
func (d *Dynamic) UnmarshalJSON(buf []byte) error {
	var dyn Dynamic
	str, err := strconv.Unquote(string(buf))
	if err != nil {
		str = string(buf)
	}

	dyn.Type, _ = Type(str)

	switch dyn.Type {
	case "date":
		dyn.Data = Time(str)
	default:
		dyn.Data = str
	}

	*d = dyn
	return nil
}
