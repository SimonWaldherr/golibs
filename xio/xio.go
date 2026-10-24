// Package xio provides bounded reads and complete writes for generic streams.
// The caller owns the reader or writer; these helpers never close it.
package xio

import (
	"errors"
	"io"
	"math"
)

var (
	// ErrLimitExceeded means that a stream contained more bytes than permitted.
	ErrLimitExceeded = errors.New("xio: read limit exceeded")
	// ErrInvalidLimit means that a negative read limit was supplied.
	ErrInvalidLimit = errors.New("xio: read limit must be nonnegative")
)

// ReadLimit reads at most limit bytes, returning partial data on a read error.
// If the stream is larger, it consumes one additional probe byte and returns
// the allowed prefix with ErrLimitExceeded. This error takes precedence over
// a read error accompanying that probe byte. A zero limit accepts only an empty
// stream. A negative limit returns ErrInvalidLimit without reading.
// The limit does not cause a proportional allocation before data is read.
func ReadLimit(reader io.Reader, limit int64) ([]byte, error) {
	if limit < 0 {
		return nil, ErrInvalidLimit
	}
	if limit == math.MaxInt64 {
		// A byte slice cannot exceed MaxInt64; avoid overflowing limit + 1.
		return io.ReadAll(reader)
	}
	data, err := io.ReadAll(io.LimitReader(reader, limit+1))
	if int64(len(data)) > limit {
		return data[:limit:limit], ErrLimitExceeded
	}
	return data, err
}

// WriteFull writes all of data, retrying positive short writes without errors.
// It returns the total accepted bytes and the first writer error. A zero-byte
// write without an error returns io.ErrShortWrite. Impossible byte counts also
// return io.ErrShortWrite unless the writer supplied its own error; such counts
// are not included in the total. An empty input succeeds without calling writer.
func WriteFull(writer io.Writer, data []byte) (written int, err error) {
	for len(data) > 0 {
		n, err := writer.Write(data)
		if n < 0 || n > len(data) {
			if err != nil {
				return written, err
			}
			return written, io.ErrShortWrite
		}
		written += n
		data = data[n:]
		if err != nil {
			return written, err
		}
		if n == 0 {
			return written, io.ErrShortWrite
		}
	}
	return written, nil
}
