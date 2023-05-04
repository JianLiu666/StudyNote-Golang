package main

import (
	"math/rand"
	"sync"
	"testing"
)

var hugeSlice = make([]task, 10000)
var hang = make(chan bool)
var forever = make(chan struct{})
var mu sync.RWMutex

const limit = 100

type task struct {
}

func perform(t task) {
	if rand.Int31()%100 == 0 {
		mu.Lock()
		// 即使我們在任務完成的時候會結束掉 worker
		// 但也不能保證 worker 在執行的過程中是否會發生問題
		// 例如現在這個位置因為沒有人 consume chanel 導致無法在繼續寫入而卡住
		<-forever
		mu.Unlock()
	}
}

func prepare(t task) task {
	mu.Lock()
	mu.Unlock()
	return t
}

func TestDeadlock(t *testing.T) {
	// start the workers
	work := make(chan task)
	var wg sync.WaitGroup
	for n := limit; n > 0; n-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range work {
				perform(task)
			}
		}()
	}

	// send the work
	for _, task := range hugeSlice {
		// 上面的 channel 卡住事件連帶影響沒辦法再繼續派發任務給 worker
		// 造成 deadlock 僵局
		work <- prepare(task)
	}

	// signal end of work
	close(work)

	// wait for completion
	wg.Wait()
}
