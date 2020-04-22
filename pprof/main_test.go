package main

import "testing"

func BenchmarkSliceInsert(b *testing.B) {
	var s1 []uint32
	b.ResetTimer()
	var i uint32
	for i = 0; i < uint32(b.N); i++ {
		s1 = append(s1, i)
	}
}
