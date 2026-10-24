package table

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRenderCSVFlushesOutput(t *testing.T) {
	data := []struct{ Value string }{{Value: "a,b\n\"c\""}}
	output, err := RenderCSV(data)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := csv.NewReader(strings.NewReader(output)).ReadAll()
	if err != nil || !reflect.DeepEqual(rows, [][]string{{"Value"}, {"a,b\n\"c\""}}) {
		t.Fatalf("CSV output = %q, rows = %v, error = %v", output, rows, err)
	}
}

func TestRenderHTMLRetainsEscaping(t *testing.T) {
	output, err := RenderHTML([]struct{ Value string }{{Value: `<img src=x onerror="alert(1)">&'`}}, TableOption{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output, "<img") || !strings.Contains(output, "&lt;img src=x onerror=&#34;alert(1)&#34;&gt;&amp;&#39;") {
		t.Fatalf("HTML cell was not escaped: %s", output)
	}
}

func TestRenderInvalidInputs(t *testing.T) {
	for _, input := range []interface{}{nil, 1, []int{1}, []struct{}{}, []struct{ hidden string }{{hidden: "secret"}}} {
		for _, format := range []string{"ascii", "unicode", "markdown", "html"} {
			if output, err := Render(input, TableOption{Format: format}); err == nil || output != "" {
				t.Errorf("Render(%T, %s) = %q, %v", input, format, output, err)
			}
		}
	}
	if _, err := Render([]struct{ Value string }{{"x"}}, TableOption{Format: "unsupported"}); err == nil {
		t.Fatal("unknown format accepted")
	}
}

func TestParseCSVInvalidDestinations(t *testing.T) {
	var nilSlice *[]benchmarkRow
	for _, destination := range []interface{}{nil, nilSlice, []benchmarkRow{}, new(int), new([]int), new([]*benchmarkRow)} {
		if err := ParseCSVInput(strings.NewReader("Name\nAlice\n"), destination); err == nil {
			t.Errorf("accepted destination %T", destination)
		}
	}
	var rows []benchmarkRow
	if err := ParseCSVInput(nil, &rows); err == nil {
		t.Fatal("nil reader accepted")
	}
}

type failingCSVReader struct{ err error }

func (reader failingCSVReader) Read([]byte) (int, error) { return 0, reader.err }

func TestParseCSVReaderErrorAndHeaderOnly(t *testing.T) {
	want := errors.New("read failed")
	rows := []benchmarkRow{{Name: "existing", Count: 7}}
	reader := io.MultiReader(strings.NewReader("Name,Count,Price\nvalid,1,2\n"), failingCSVReader{err: want})
	if err := ParseCSVInput(reader, &rows); !errors.Is(err, want) || len(rows) != 1 || rows[0].Name != "existing" {
		t.Fatalf("read error changed destination or was lost: %v, %v", rows, err)
	}
	if err := ParseCSVInput(strings.NewReader("Name,Count,Price\n"), &rows); err != nil || len(rows) != 1 {
		t.Fatalf("header-only input: %v, %v", rows, err)
	}
	var empty []benchmarkRow
	if err := ParseCSVInput(strings.NewReader("Name,Count,Price\n"), &empty); err != nil || empty != nil {
		t.Fatalf("header-only input allocated rows: %v, %v", empty, err)
	}
}

func TestParseCSVValidationIsTransactional(t *testing.T) {
	for _, input := range []string{
		"", "Name,Count,Price\nvalid,1,2\ninvalid,no,3\n",
		"Name,Count,Price\nvalid,1,2\nmissing,2\n",
		"Name,Count,Price\nvalid,1,2\ninvalid,2,not-a-number\n",
		"Name,Count,Price\nvalid,1,2\n\"unfinished",
	} {
		original := []benchmarkRow{{Name: "existing", Count: 7, Price: 1.25}}
		rows := append([]benchmarkRow(nil), original...)
		if err := ParseCSVInput(strings.NewReader(input), &rows); err == nil || !reflect.DeepEqual(rows, original) {
			t.Errorf("ParseCSVInput(%q) = %v, rows = %v", input, err, rows)
		}
	}
}

func TestParseCSVNumericBoundsAndHiddenFields(t *testing.T) {
	type row struct {
		Count  int8
		Price  float32
		hidden string
	}
	for _, input := range []string{
		"Count,Price\n128,1\n", "Count,Price\n-129,1\n", "Count,Price\n1,1e100\n",
		"hidden\nsecret\n", "Count,Count\n1,2\n",
	} {
		var rows []row
		if err := ParseCSVInput(strings.NewReader(input), &rows); err == nil || rows != nil {
			t.Errorf("accepted invalid fields %q: %v, %v", input, rows, err)
		}
	}
	var rows []row
	if err := ParseCSVInput(strings.NewReader("Count,Price,Ignored\n127,3.5,anything\n-128,0,ignored\n"), &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].Count != 127 || rows[1].Count != -128 || rows[0].Price != 3.5 {
		t.Fatalf("boundary values changed: %+v", rows)
	}
}

func TestParseCSVSupportedAndUnsupportedFields(t *testing.T) {
	type supported struct {
		Name string
		When time.Time
		Live bool
	}
	rows := []supported{{Name: "existing"}}
	if err := ParseCSVInput(strings.NewReader("Name,When,Live\nAlice,2026-10-24T12:00:00Z,true\n"), &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[1].Name != "Alice" || !rows[1].Live || rows[1].When.Hour() != 12 {
		t.Fatalf("supported fields changed: %v", rows)
	}
	for _, input := range []string{"Name,When,Live\nAlice,invalid,true\n", "Name,When,Live\nAlice,2026-10-24T12:00:00Z,invalid\n"} {
		copy := append([]supported(nil), rows...)
		if err := ParseCSVInput(strings.NewReader(input), &copy); err == nil || !reflect.DeepEqual(copy, rows) {
			t.Fatalf("invalid row changed existing rows: %v, %v", copy, err)
		}
	}
	for _, input := range []string{"Value\n", "Value\nanything\n"} {
		unsupported := []struct{ Value struct{ X int } }{{Value: struct{ X int }{7}}}
		if err := ParseCSVInput(strings.NewReader(input), &unsupported); err == nil || len(unsupported) != 1 || unsupported[0].Value.X != 7 {
			t.Fatalf("unsupported struct field: %v, %v", unsupported, err)
		}
		sliceField := []struct{ Value []byte }{{Value: []byte("existing")}}
		if err := ParseCSVInput(strings.NewReader(input), &sliceField); err == nil || len(sliceField) != 1 || string(sliceField[0].Value) != "existing" {
			t.Fatalf("unsupported slice field: %v, %v", sliceField, err)
		}
	}
}

func TestParseCSVPromotedNilPointer(t *testing.T) {
	type Embedded struct{ Name string }
	for _, input := range []string{"Name\n", "Name\nAlice\n"} {
		var rows []struct{ *Embedded }
		if err := ParseCSVInput(strings.NewReader(input), &rows); err == nil || rows != nil {
			t.Fatalf("nil embedded pointer: %v, %v", rows, err)
		}
	}
}

func TestFormatFieldRetainsPrimitiveFormatting(t *testing.T) {
	for _, value := range []interface{}{
		"text", true, false, int(-1), int8(-2), int16(-3), int32(-4), int64(-5),
		uint(1), uint8(2), uint16(3), uint32(4), uint64(math.MaxUint64),
		float32(1.0 / 3), float32(1e-7), float64(1.0 / 3), 1e25,
		math.Copysign(0, -1), math.NaN(), math.Inf(1), []int{1, 2},
	} {
		if got, want := formatField(value), fmt.Sprintf("%v", value); got != want {
			t.Errorf("%T formatting = %q, want %q", value, got, want)
		}
	}
}

func TestColumnAlignmentMatchesLegacyScanner(t *testing.T) {
	for _, value := range []string{
		"", " ", "text", "row-001", "Name", "Count", "Price",
		"0", "-0", "+0", "123", "-123", "+123", ".5", "-.5", "+.5", "1.", "1.25",
		" 12", "12 ", "\t12", "\n12", "12\n", "\u200312", "12\u2003",
		"12kg", "12abc", "12 abc", "1.2.3", "1e2", "1e+2", "1e-2", "1e", "1e+", "1e9999", "1e-9999",
		"1p2", "1.2p4", "1P2", "1_000.5", "1__000", "0b101", "0x1p2", "+0x1.fp2", "0X.8P2", "0x_1p2", "0xdeadbeef",
		"Inf", "+Inf", "-Inf", "inf", "Infinity", "+Infinity", "-Infinity", "iNfInItY", "Inf trailing",
		"NaN", "nan", "NAN", "+NaN", "-NaN", "NaN trailing", "+", "-", ".", "+.", "-."} {
		want := "left"
		if _, err := fmt.Sscanf(value, "%f", new(float64)); err == nil {
			want = "right"
		}
		if got := getColumnAlignment(0, nil, value); got != want {
			t.Errorf("alignment(%q) = %s, want legacy %s", value, got, want)
		}
		if got := getColumnAlignment(0, []string{"center"}, value); got != "center" {
			t.Errorf("explicit alignment(%q) = %s", value, got)
		}
	}
}

func TestPlainNumericAlignmentDoesNotAllocate(t *testing.T) {
	if allocations := testing.AllocsPerRun(100, func() {
		if got := getColumnAlignment(0, nil, "12345.625"); got != "right" {
			t.Fatal("numeric alignment changed")
		}
	}); allocations != 0 {
		t.Fatalf("plain numeric alignment allocated %g times", allocations)
	}
}

func FuzzParseCSVInput(f *testing.F) {
	for _, seed := range []string{"", "Name,Count,Price\nAlice,1,2.5", "Name,Count\nBob,no", "\"unfinished"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		original := []benchmarkRow{{Name: "existing", Count: 7}}
		rows := append([]benchmarkRow(nil), original...)
		if err := ParseCSVInput(strings.NewReader(input), &rows); err != nil && !reflect.DeepEqual(rows, original) {
			t.Fatalf("failed parse changed destination: %v", rows)
		}
	})
}
