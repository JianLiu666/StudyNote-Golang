package main

import (
	"context"
	"fmt"
	"sync"
)

type state struct {
	seq  int64
	wait []chan<- int64
}

type Notifier struct {
	st chan state
}

func NewNotifier() *Notifier {
	st := make(chan state, 1)
	st <- state{seq: 0}
	return &Notifier{st}
}

func (n *Notifier) AwaitChange(ctx context.Context, seq int64) (newSeq int64) {
	c := make(chan int64, 1)
	st := <-n.st
	if st.seq == seq {
		st.wait = append(st.wait, c)
	} else {
		c <- st.seq
	}
	n.st <- st

	select {
	case <-ctx.Done():
		return seq
	case newSeq = <-c:
		return newSeq
	}
}

func (n *Notifier) NotifyChange() {
	st := <-n.st
	for _, c := range st.wait {
		c <- st.seq + 1
	}
	n.st <- state{st.seq + 1, nil}
}

func main() {
	notifier := NewNotifier()

	var wg sync.WaitGroup
	for n := 10; n > 0; n-- {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			notifier.AwaitChange(context.Background(), 0)
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
			notifier.NotifyChange()
			n++
		}
	}
}
