package main

import (
	"sync"
	"time"
)

type Idler struct {
	mu    sync.Mutex
	idle  sync.Cond
	busy  bool
	idles int64
}

func NewIdler() *Idler {
	i := new(Idler)
	i.idle.L = &i.mu
	return i
}

func (i *Idler) AwaitIdle() {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 用一個變數紀錄 wait 瞬間的 i.idles(i.e. sequence number)
	// 如果沒有這樣做的話, 當 Idler 狀態從 busy->idle 的瞬間又發生 idle->busy 後, 很可能這條 goroutine 才剛被喚醒
	// 於是就錯過這個短暫的 idle status 繼續休眠了
	// 因此只有在 status 仍在 busy 且 idles 沒有變化時才會進入休眠
	idles := i.idles
	for i.busy && idles == i.idles {
		i.idle.Wait()
	}
}

func (i *Idler) SetBusy(b bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	wasBusy := i.busy
	i.busy = b

	// 表示狀態從 busy -> idle 的時候
	if wasBusy && !i.busy {
		// 累加 idles 表示狀態從 busy -> idle 的次數
		i.idles++
		i.idle.Broadcast()
	}
}

func main() {
	i := NewIdler()
	i.SetBusy(true)
	go func() {
		time.Sleep(1 * time.Second)
		i.SetBusy(false)
	}()
	i.AwaitIdle()
	i.AwaitIdle()
}
