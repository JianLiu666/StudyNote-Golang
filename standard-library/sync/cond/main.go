package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

var status int64

func main() {
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 10; i++ {
		go listen(i, c)
	}

	time.Sleep(time.Second)
	atomic.StoreInt64(&status, 1)
	// 只會隨機喚醒一個正在等待的 goroutine
	c.Signal()

	time.Sleep(2 * time.Second)
	// 喚醒所有正在等待的 goroutines
	c.Broadcast()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func listen(id int, c *sync.Cond) {
	c.L.Lock()
	defer c.L.Unlock()

	for atomic.LoadInt64(&status) != 1 {
		c.Wait()
	}

	log.Println("listen: ", id)
}
