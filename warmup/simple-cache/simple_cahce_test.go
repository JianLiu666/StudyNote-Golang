package simplecache

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestSimepleCache(t *testing.T) {
	cache := Constructor(100)

	for i := 0; i < 100; i++ {
		cache.Set(strconv.Itoa(i), i, 100)
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println(cache.GetAll())
	fmt.Println(cache.Get("50"))

	time.Sleep(10 * time.Millisecond)
	fmt.Println(cache.Set("100", 100, -10000))

	fmt.Println(cache.GetAll())
}
