package benchmark

import (
	"strings"
	"testing"
)

func BenchmarkStringBuilder(b *testing.B) {
	var str strings.Builder

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str.WriteString("s")
	}
}

func BenchmarkStringConcat(b *testing.B) {
	str := ""

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str += "s"
	}
}
