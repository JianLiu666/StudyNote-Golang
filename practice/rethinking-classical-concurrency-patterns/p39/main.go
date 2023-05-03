package main

import (
	"fmt"
	"sync"
)

type Item = int

type Queue struct {
	mu        sync.Mutex
	items     []Item
	itemAdded sync.Cond
}

func NewQueue() *Queue {
	q := new(Queue)
	q.itemAdded.L = &q.mu
	return q
}

func (q *Queue) Get() Item {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) == 0 {
		q.itemAdded.Wait()
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue) Put(item Item) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
	q.itemAdded.Signal()
}

func (q *Queue) GetMany(n int) []Item {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) < n {
		q.itemAdded.Wait()
	}
	items := q.items[:n:n]
	q.items = q.items[n:]
	return items
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
		q.Put(i)
	}

	wg.Wait()
}
