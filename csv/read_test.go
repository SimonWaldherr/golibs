package csv

import (
	stdcsv "encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestReadQuotedFields(t *testing.T) {
	data, head, err := ReadString("id;name;notes\n1;\"Alice;Bob\";\"first\nsecond\"\n2;Zoë;plain\n")
	wantData := map[int][]string{
		1: {"1", "Alice;Bob", "first\nsecond"},
		2: {"2", "Zoë", "plain"},
	}
	wantHead := map[string]int{"id": 0, "name": 1, "notes": 2}
	if err != nil || !reflect.DeepEqual(data, wantData) || !reflect.DeepEqual(head, wantHead) {
		t.Fatalf("ReadString() = %v, %v, %v", data, head, err)
	}
}

func TestReadEmptyAndHeaderOnly(t *testing.T) {
	for _, input := range []string{"", "id;name\n"} {
		t.Run(input, func(t *testing.T) {
			data, head, err := ReadString(input)
			if err != nil || data == nil || head == nil || len(data) != 0 {
				t.Fatalf("ReadString() = %v, %v, %v", data, head, err)
			}
			if input != "" && !reflect.DeepEqual(head, map[string]int{"id": 0, "name": 1}) {
				t.Errorf("header = %v", head)
			}
		})
	}
}

func TestReadParseErrorPreservesSuccessfulRows(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{"field count", "id;name\n1;Alice\n2\n", stdcsv.ErrFieldCount},
		{"unterminated quote", "id;name\n1;Alice\n2;\"broken", stdcsv.ErrQuote},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, head, err := ReadString(tt.input)
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
			var parseError *stdcsv.ParseError
			if !errors.As(err, &parseError) {
				t.Fatalf("error does not retain csv.ParseError: %v", err)
			}
			if !reflect.DeepEqual(data, map[int][]string{1: {"1", "Alice"}}) ||
				!reflect.DeepEqual(head, map[string]int{"id": 0, "name": 1}) {
				t.Fatalf("partial result = %v, %v", data, head)
			}
			legacyData, legacyHead := LoadCSVfromString(tt.input)
			if !reflect.DeepEqual(legacyData, data) || !reflect.DeepEqual(legacyHead, head) {
				t.Fatalf("legacy wrapper did not return successful rows: %v, %v", legacyData, legacyHead)
			}
		})
	}
}

type failingReader struct {
	err   error
	calls int
}

func (r *failingReader) Read(p []byte) (int, error) {
	r.calls++
	return 0, r.err
}

func TestReadReaderError(t *testing.T) {
	wantErr := errors.New("reader failed")
	failing := &failingReader{err: wantErr}
	reader := io.MultiReader(strings.NewReader("id;name\n1;Alice\n"), failing)
	data, head, err := Read(reader)
	if !errors.Is(err, wantErr) || failing.calls != 1 {
		t.Fatalf("error = %v, reads after error = %d", err, failing.calls)
	}
	if !reflect.DeepEqual(data, map[int][]string{1: {"1", "Alice"}}) || head["name"] != 1 {
		t.Fatalf("partial result = %v, %v", data, head)
	}
}

func TestReadFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "data.csv")
	if err := os.WriteFile(filename, []byte("id;name\n1;Alice\n"), 0600); err != nil {
		t.Fatal(err)
	}
	data, head, err := ReadFile(filename)
	if err != nil || !reflect.DeepEqual(data, map[int][]string{1: {"1", "Alice"}}) || head["name"] != 1 {
		t.Fatalf("ReadFile() = %v, %v, %v", data, head, err)
	}

	missing := filepath.Join(t.TempDir(), "missing.csv")
	data, head, err = ReadFile(missing)
	if !errors.Is(err, os.ErrNotExist) || data == nil || head == nil || len(data) != 0 || len(head) != 0 {
		t.Fatalf("ReadFile(missing) = %v, %v, %v", data, head, err)
	}
	data, head = LoadCSVfromFile(missing)
	if data == nil || head == nil || len(data) != 0 || len(head) != 0 {
		t.Fatalf("LoadCSVfromFile(missing) = %v, %v", data, head)
	}
}
