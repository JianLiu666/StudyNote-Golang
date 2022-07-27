package simplecache

import (
	"strconv"
	"testing"
)

func BenchmarkNormalSet_Size_1000(b *testing.B) {
	num := 1000
	cache := CreateNormalCache(num)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i), i, 100000)
	}
}

func BenchmarkNormalGet_Size_1000(b *testing.B) {
	num := 1000
	cache := CreateNormalCache(num)
	for i := 0; i < num; i++ {
		cache.Set(strconv.Itoa(i), i, 100000)
	}

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(strconv.Itoa(i % num))
	}
}

func BenchmarkNormalSet_Size_10000(b *testing.B) {
	num := 10000
	cache := CreateNormalCache(num)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i), i, 100000)
	}
}

func BenchmarkNormalGet_Size_10000(b *testing.B) {
	num := 10000
	cache := CreateNormalCache(num)
	for i := 0; i < num; i++ {
		cache.Set(strconv.Itoa(i), i, 100000)
	}

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(strconv.Itoa(i % num))
	}
}
