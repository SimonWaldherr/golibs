package xio

import (
	"bytes"
	"errors"
	"io"
	"math"
	"strings"
	"testing"
)

func TestReadLimit(t *testing.T) {
	for _, tc := range []struct {
		name, input, want string
		limit             int64
		err               error
	}{
		{"empty", "", "", 4, nil},
		{"below", "abc", "abc", 4, nil},
		{"exact", "abcd", "abcd", 4, nil},
		{"exceeded", "abcde", "abcd", 4, ErrLimitExceeded},
		{"zero empty", "", "", 0, nil},
		{"zero nonempty", "abc", "", 0, ErrLimitExceeded},
		{"negative", "abc", "", -1, ErrInvalidLimit},
		{"huge", "abc", "abc", math.MaxInt64, nil},
		{"huge empty", "", "", math.MaxInt64, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ReadLimit(strings.NewReader(tc.input), tc.limit)
			if string(got) != tc.want || !errors.Is(err, tc.err) {
				t.Fatalf("ReadLimit = %q, %v; want %q, %v", got, err, tc.want, tc.err)
			}
		})
	}
}

type readerFunc func([]byte) (int, error)

func (read readerFunc) Read(p []byte) (int, error) { return read(p) }

func TestReadLimitPreservesErrorsAndPartialData(t *testing.T) {
	wantErr := errors.New("reader failed")
	for _, limit := range []int64{10, math.MaxInt64} {
		reader := readerFunc(func(p []byte) (int, error) { return copy(p, "abc"), wantErr })
		data, err := ReadLimit(reader, limit)
		if string(data) != "abc" || !errors.Is(err, wantErr) {
			t.Fatalf("ReadLimit(%d) = %q, %v", limit, data, err)
		}
	}
	reader := readerFunc(func(p []byte) (int, error) { return copy(p, "abcde"), wantErr })
	data, err := ReadLimit(reader, 4)
	if string(data) != "abcd" || !errors.Is(err, ErrLimitExceeded) {
		t.Fatalf("overflow probe = %q, %v", data, err)
	}
	reader = readerFunc(func(p []byte) (int, error) { return copy(p, "abc"), io.EOF })
	data, err = ReadLimit(reader, 3)
	if string(data) != "abc" || err != nil {
		t.Fatalf("data with EOF = %q, %v", data, err)
	}
}

func TestReadLimitConsumesOnlyOneProbeByte(t *testing.T) {
	reader := bytes.NewBufferString("0123456789")
	data, err := ReadLimit(reader, 4)
	if string(data) != "0123" || !errors.Is(err, ErrLimitExceeded) || reader.String() != "56789" {
		t.Fatalf("bounded read = %q, %v; remaining = %q", data, err, reader.String())
	}
	if cap(data) != len(data) {
		t.Fatal("overflow result permits reslicing into the probe byte")
	}
}

func TestReadLimitInvalidDoesNotRead(t *testing.T) {
	reader := readerFunc(func([]byte) (int, error) {
		t.Fatal("read called for a negative limit")
		return 0, nil
	})
	if data, err := ReadLimit(reader, -1); data != nil || !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("invalid limit = %v, %v", data, err)
	}
}

type writerFunc func([]byte) (int, error)

func (write writerFunc) Write(p []byte) (int, error) { return write(p) }

func TestWriteFullRetriesShortWrites(t *testing.T) {
	var output bytes.Buffer
	var calls int
	writer := writerFunc(func(p []byte) (int, error) {
		calls++
		if len(p) > 2 {
			p = p[:2]
		}
		return output.Write(p)
	})
	input := []byte("abcdefg")
	n, err := WriteFull(writer, input)
	if n != len(input) || err != nil || output.String() != string(input) || calls != 4 {
		t.Fatalf("WriteFull = %d, %v; output %q, calls %d", n, err, output.String(), calls)
	}
	if string(input) != "abcdefg" {
		t.Fatal("WriteFull modified its input")
	}
}

func TestWriteFullErrorsAndImpossibleCounts(t *testing.T) {
	wantErr := errors.New("writer failed")
	for _, tc := range []struct {
		name     string
		count    int
		writeErr error
		wantN    int
		wantErr  error
	}{
		{"error without progress", 0, wantErr, 0, wantErr},
		{"partial error", 2, wantErr, 2, wantErr},
		{"complete with error", 4, wantErr, 4, wantErr},
		{"no progress", 0, nil, 0, io.ErrShortWrite},
		{"negative count", -1, nil, 0, io.ErrShortWrite},
		{"excess count", 5, nil, 0, io.ErrShortWrite},
		{"negative with error", -1, wantErr, 0, wantErr},
		{"excess with error", 5, wantErr, 0, wantErr},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			writer := writerFunc(func([]byte) (int, error) {
				calls++
				return tc.count, tc.writeErr
			})
			n, err := WriteFull(writer, []byte("abcd"))
			if n != tc.wantN || !errors.Is(err, tc.wantErr) || calls != 1 {
				t.Fatalf("WriteFull = %d, %v; calls %d", n, err, calls)
			}
		})
	}
}

func TestWriteFullCountsEarlierWritesOnFailure(t *testing.T) {
	for _, impossible := range []bool{false, true} {
		calls := 0
		writer := writerFunc(func(p []byte) (int, error) {
			calls++
			if calls == 1 {
				return 2, nil
			}
			if string(p) != "cd" {
				t.Fatalf("remaining data = %q", p)
			}
			if impossible {
				return -1, nil
			}
			return 0, nil
		})
		n, err := WriteFull(writer, []byte("abcd"))
		if n != 2 || !errors.Is(err, io.ErrShortWrite) || calls != 2 {
			t.Fatalf("WriteFull = %d, %v; calls %d", n, err, calls)
		}
	}
}

func TestWriteFullEmptyDoesNotWrite(t *testing.T) {
	writer := writerFunc(func([]byte) (int, error) {
		t.Fatal("write called for empty input")
		return 0, nil
	})
	if n, err := WriteFull(writer, nil); n != 0 || err != nil {
		t.Fatalf("empty write = %d, %v", n, err)
	}
}
