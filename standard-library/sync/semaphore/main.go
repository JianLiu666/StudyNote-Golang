package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/semaphore"
)

var (
	maxResource = int64(32)
	sem         = semaphore.NewWeighted(maxResource) // 給定 semaphore 最大的資源數量
	out         = make([]int, 32)
)

func main() {
	ctx := context.Background()

	for i := range out {
		// 如果目前 semaphore 的剩餘資源不夠的話就會進入到 waiting 直到有足夠的資源後才會被喚醒
		if err := sem.Acquire(ctx, int64(i)); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		// 模擬業務調用
		go func(i int) {
			defer sem.Release(int64(i))
			time.Sleep(time.Second)
			log.Println("done")
			out[i] = i
		}(i)
	}

	// 確認所有資源都已經正常歸位
	if err := sem.Acquire(ctx, maxResource); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
	}

	fmt.Println(out)
}
