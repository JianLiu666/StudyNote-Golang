package main

import (
	"fmt"
	"sync"
)

type Item = int

type waiter struct {
	n int
	c chan []Item
}

type state struct {
	items []Item   //
	wait  []waiter // 用來管理 goroutines 取得資料的 metadata
}

type Queue struct {
	s chan state
}

func NewQueue() *Queue {
	s := make(chan state, 1)
	s <- state{}
	return &Queue{
		s: s,
	}
}

func (q *Queue) Put(item Item) {
	// 等同於 lock
	s := <-q.s
	s.items = append(s.items, item)

	// 檢查是否有 goroutines 正在等待資料進入
	for len(s.wait) > 0 {
		w := s.wait[0]
		if len(s.items) < w.n {
			break
		}
		// 打包資料從 state.items 到 waiter.c
		w.c <- s.items[:w.n:w.n]
		s.items = s.items[w.n:]
		s.wait = s.wait[1:]
	}

	// 等同於 unlock
	q.s <- s
}

func (q *Queue) GetMany(n int) []Item {
	s := <-q.s

	// 數量足夠且沒有人在等待時, 就直接取走需要的數量
	if len(s.wait) == 0 && len(s.items) >= n {
		items := s.items[:n:n]
		s.items = s.items[n:]
		q.s <- s
		return items
	}

	// 數量不夠或有人在等待, 就往後排隊
	// 很巧妙的是這裡使用的是 unbuffered channel, 確保一定會等到資料進來
	c := make(chan []Item)
	s.wait = append(s.wait, waiter{n, c})
	q.s <- s

	return <-c
}

func main() {
	q := NewQueue()

	var wg sync.WaitGroup
	for n := 10; n > 0; n-- {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			items := q.GetMany(n)
			fmt.Printf("%2d: %2d\n", n, items)
		}(n)
	}

	for i := 0; i < 100; i++ {
		q.Put(1)
	}

	wg.Wait()
}
