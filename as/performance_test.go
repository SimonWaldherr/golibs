package as

import "testing"

var (
	benchmarkValid bool
	benchmarkType  string
	benchmarkError error
)

func BenchmarkValidation(b *testing.B) {
	for _, tc := range []struct {
		name  string
		input string
		check func(string) bool
	}{
		{"email_valid", "person+tag@example.com", IsEmailValid},
		{"email_invalid", "not an email", IsEmailValid},
		{"url_valid", "https://example.com/docs/index.html", IsURLValid},
		{"url_invalid", "invalid URL", IsURLValid},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchmarkValid = tc.check(tc.input)
			}
		})
	}
}

func BenchmarkClassify(b *testing.B) {
	for _, tc := range []struct{ name, input string }{
		{"boolean", "true"},
		{"integer", "12345"},
		{"email", "person@example.com"},
		{"unicode", "😃"},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchmarkType, benchmarkError = Type(tc.input)
			}
		})
	}
}
