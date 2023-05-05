package main

import (
	"fmt"
	"sync"
)

type Notifier struct {
	mu      sync.Mutex
	changed sync.Cond
	seq     int64
}

func NewNotifier() *Notifier {
	r := new(Notifier)
	r.changed.L = &r.mu
	return r
}

func (n *Notifier) AwaitChange(seq int64) (newSeq int64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for n.seq == seq {
		n.changed.Wait()
	}
	return n.seq
}

func (n *Notifier) NotifyChange() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.seq++
	n.changed.Broadcast()
}

func main() {
	r := NewNotifier()

	var wg sync.WaitGroup
	for n := 10; n > 0; n-- {
		wg.Add(1)
		go func(n int) {
			r.AwaitChange(0)
			wg.Done()
		}(n)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	n := 0
	for {
		select {
		case <-done:
			fmt.Printf("Done after %d reload(s).\n", n)
			return
		default:
			r.NotifyChange()
			n++
		}
	}
}
