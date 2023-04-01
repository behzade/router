package router

import "testing"

func BenchmarkParse(b *testing.B) {
	urlPath := "/1/classes/go/123456789"
	b.ReportAllocs()
	b.ResetTimer()
	// var part string
	for i := 0; i < b.N; i++ {
		part, rest := parse(urlPath)
		for part != "" {
			part, rest = parse(rest)
		}
	}
}

func BenchmarkParseStruct(b *testing.B) {
	urlPath := "/1/classes/go/123456789"
	b.ReportAllocs()
	b.ResetTimer()
    p := path{path: urlPath}
	// var part string
	for i := 0; i < b.N; i++ {
        p.offset = 0
		for p.offset < len(p.path) {
			_ = p.parse()
		}
	}
}
