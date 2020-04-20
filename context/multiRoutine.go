package main

import (
	"context"
	"log"
	"strconv"
	"time"
)

func multiRoutine() {
	// context.Background 創建一個空的 context(root)
	// context.WithCancel 返回一個可取消的 sub context, 當作參數傳給 goroutine 使用
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 2; i++ {
		go multiRoutineEpoch(ctx, strconv.Itoa(i))
	}

	time.Sleep(5 * time.Second)
	// 對所有被 ctx 監聽的 goroutine 發出中斷信號
	cancel()
	log.Print("All of goroutines stopping...")

	// 緩衝一段時間確認所有 Goroutine 是否正確退出
	time.Sleep(5 * time.Second)
}

/** Goroutine 迭代事件
 *
 * @param ctx
 * @param instance */
func multiRoutineEpoch(ctx context.Context, instance string) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Goroutine:%s stop.", instance)
			return
		default:
			log.Printf("Goroutine:%s running...", instance)
			time.Sleep(1 * time.Second)
		}
	}
}
