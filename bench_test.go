package router

import "testing"

func BenchmarkParse(b *testing.B) {
	urlPath := "/1/classes/go/123456789"
	b.ReportAllocs()
	b.ResetTimer()
    var part []byte
    var rest string 
	// var part string
	for i := 0; i < b.N; i++ {
		part, rest = parse(urlPath)
		for len(part) != 0 {
			part, rest = parse(rest)
		}
	}
}
