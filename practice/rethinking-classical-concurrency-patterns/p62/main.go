package main

import (
	"fmt"
	"sync"
)

type Item = int

type Queue struct {
	items chan []Item // non-empty slices only
	empty chan bool   // holds true if the queue is empty
}

func NewQueue() *Queue {
	items := make(chan []Item, 1)
	empty := make(chan bool, 1)
	empty <- true
	return &Queue{
		items: items,
		empty: empty,
	}
}

func (q *Queue) Get() Item {
	// 一旦有其中一個 goroutine 拿到 items 後, 代表其他 goroutines 需要繼續等待 items 被再次放回 channel
	items := <-q.items
	item := items[0]
	items = items[1:]
	if len(items) == 0 {
		// semaphore 表示 items 空了
		q.empty <- true
	} else {
		// 還有資料的話將 items 放回 channel 讓其他 goroutines 使用
		q.items <- items
	}
	return item
}

func (q *Queue) Put(item Item) {
	var items []Item

	// 只有其中一種可能
	// 1. 等到某個 goroutine 用完 items 放回 channel
	// 2. 收到 semaphore 信號得知 items 內空了
	select {
	case items = <-q.items:
		// 繼續使用原本的 slice reference
	case <-q.empty:
	}
	items = append(items, item)
	q.items <- items
}

func main() {
	q := NewQueue()

	var wg sync.WaitGroup
	for n := 20; n > 0; n-- {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			item := q.Get()
			fmt.Printf("%2d: %2d\n", n, item)
		}(n)
	}

	for i := 0; i < 100; i++ {
		q.Put(i)
	}

	wg.Wait()
}
