package main

import (
	"context"
	"log"
	"time"
)

var key string = "name"

func withValue() {
	ctx, cancel := context.WithCancel(context.Background())

	// 附加一個 Key-Value mapping, 這裡的 key 必須是可比較性的，Value 必須符合多執行序安全
	valueCtx := context.WithValue(ctx, key, "watch1")
	go withValueEpoch(valueCtx)

	time.Sleep(5 * time.Second)
	// 對所有被 ctx 監聽的 goroutine 發出中斷信號
	cancel()
	log.Print("All of goroutines stopping...")

	// 緩衝一段時間確認所有 goroutine 是否正確退出
	time.Sleep(5 * time.Second)

}

// Goroutine 迭代事件
//
// @param ctx
func withValueEpoch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// 取出值
			log.Printf("%v: goroutine stopped.", ctx.Value(key))
			return
		default:
			// 取出值
			log.Printf("%v: goroutine running...", ctx.Value(key))
			time.Sleep(1 * time.Second)
		}
	}
}
