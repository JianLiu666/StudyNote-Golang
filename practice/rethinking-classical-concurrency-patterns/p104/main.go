package main

import (
	"context"
	"fmt"
	"sync"
)

type state struct {
	seq     int64
	changed chan struct{}
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
	st := <-n.st
	if st.seq != seq {
		n.st <- st
		return st.seq
	}
	if st.changed == nil {
		st.changed = make(chan struct{})
	}
	n.st <- st

	select {
	case <-ctx.Done():
		return seq
	case <-st.changed:
		return seq + 1
	}
}

func (n *Notifier) NotifyChange() {
	st := <-n.st
	// 表示 busy -> idle 的狀態變換, 透過關閉 channel 的方式通知所有 goroutines
	// 跟 p102 的差別在於所有 goroutines 監聽同一個 channel
	// 換句話說只需維護一個 channel, 但也沒辦法針對不同的 goroutine 發送不同事件
	if st.changed != nil {
		close(st.changed)
	}
	n.st <- state{st.seq + 1, nil}
}

func main() {
	notifier := NewNotifier()

	var wg sync.WaitGroup
	for n := 10; n > 0; n-- {
		wg.Add(1)
		go func(n int) {
			notifier.AwaitChange(context.Background(), 0)
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
			notifier.NotifyChange()
			n++
		}
	}
}
