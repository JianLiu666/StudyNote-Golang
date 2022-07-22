package benchmark

import (
	"testing"
)

func BenchmarkSliceInsert(b *testing.B) {
	var s1 []int

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s1 = append(s1, i)
	}
}
