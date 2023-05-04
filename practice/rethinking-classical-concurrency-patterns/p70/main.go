package main

import (
	"context"
	"fmt"
	"time"
)

type Idler struct {
	// Holds a channel to be closed when idle,
	// or nil if idle already.
	next chan chan struct{}
}

func NewIdler() *Idler {
	next := make(chan chan struct{}, 1)
	// 給定最初的 status = idle
	next <- nil
	return &Idler{
		next: next,
	}
}

func (i *Idler) AwaitIdle(ctx context.Context) error {
	// 取得 Idler 現在的 status
	idle := <-i.next
	i.next <- idle

	// 如果 status = busy 時只有兩個選項
	if idle != nil {
		select {
		case <-ctx.Done():
			// 交由外部決定等待超時
			return ctx.Err()
		case <-idle:
			// Idler status 從 busy -> idle
		}
	}
	return nil
}

func (i *Idler) SetBusy(b bool) {
	// 取出當前的 status
	idle := <-i.next

	if idle == nil && b {
		// idle -> busy 情境
		// 建立一個供其他 goroutines 監聽的 status struct
		idle = make(chan struct{})

	} else if idle != nil && !b {
		// busy -> idle 情境
		// 將其他 goroutines 正在監聽的 channel 關閉, 表示狀態切換成 idle 了
		close(idle)
		idle = nil
	}

	i.next <- idle
}

func main() {
	i := NewIdler()
	i.SetBusy(true)
	go func() {
		time.Sleep(1 * time.Second)
		i.SetBusy(false)
	}()

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// 等待到超時, 因為 busy -> idle 的時間需要 1s, context 只等待 0.01s
	fmt.Println(i.AwaitIdle(ctx))
	fmt.Println(time.Since(start))

	// 等待到 Idler 從 busy -> idle
	fmt.Println(i.AwaitIdle(context.Background()))
	fmt.Println(time.Since(start))
}
