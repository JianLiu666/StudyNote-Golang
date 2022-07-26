package simplecache

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimepleCache(t *testing.T) {
	ast := assert.New(t)

	cache := CreateSimpleCache(10)
	for i := 0; i < 10; i++ {
		cache.Set(strconv.Itoa(i), i, 100)
		time.Sleep(50 * time.Millisecond)
	}

	ast.Equal([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, cache.GetAll())
	ast.Equal(5, cache.Get("5"))
	ast.Equal([]int{0, 1, 2, 3, 4, 6, 7, 8, 9, 5}, cache.GetAll())
	ast.Equal(-1, cache.Get("100"))

	time.Sleep(10 * time.Millisecond)

	ast.Equal(true, cache.Set("100", 100, -10000))
	ast.Equal([]int{1, 2, 3, 4, 6, 7, 8, 9, 5, 100}, cache.GetAll())
}

func BenchmarkSimpleSet(b *testing.B) {
	cache := CreateSimpleCache(1000)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i), i, 100)
	}
}

func BenchmarkSimpleGet(b *testing.B) {
	cache := CreateSimpleCache(1000)
	for i := 0; i < 1000; i++ {
		cache.Set(strconv.Itoa(i), i, 100)
	}

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i%1000), 0, 100)
	}
}
