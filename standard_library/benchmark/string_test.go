package pprof

import (
	"strings"
	"testing"
)

func BenchmarkStringBuilder(b *testing.B) {
	var str strings.Builder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str.WriteString("s")
	}
}

func BenchmarkStringConcat(b *testing.B) {
	str := ""
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str += "s"
	}
}
