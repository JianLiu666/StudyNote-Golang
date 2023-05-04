package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

var addr net.Addr
var addrOnce sync.Once

type token struct {
}

type Pool struct {
	sem  chan token    // semaphore pattern
	idle chan net.Conn //
}

func NewPool(limit int) *Pool {
	return &Pool{
		sem:  make(chan token, limit),
		idle: make(chan net.Conn, limit),
	}
}

func (p *Pool) Acquire(ctx context.Context) (net.Conn, error) {
	select {
	case conn := <-p.idle:
		// 用來等待 pool 中有閒置的 conn 可以用
		return conn, nil
	case p.sem <- token{}:
		// semaphore pattern 控制 conn pool 數量上限
		conn, err := dial()
		if err != nil {
			<-p.sem
		}
		return conn, err
	case <-ctx.Done():
		// 請求超時處理
		return nil, ctx.Err()
	}
}

func (p *Pool) Release(c net.Conn) {
	p.idle <- c
}

func (p *Pool) Hijack(c net.Conn) {
	<-p.sem
}

func dial() (net.Conn, error) {
	addrOnce.Do(func() {
		ln, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		addr = ln.Addr()
		go func() {
			for {
				in, err := ln.Accept()
				if err != nil {
					return
				}
				go io.Copy(os.Stdout, in)
			}
		}()
	})
	return net.Dial(addr.Network(), addr.String())
}

func main() {
	p := NewPool(3)

	var wg sync.WaitGroup
	for n := 10; n > 0; n-- {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			conn, err := p.Acquire(context.Background())
			if err != nil {
				panic(err)
			}
			defer p.Release(conn)

			fmt.Fprintf(conn, "Hello from goroutine %d on connection %p!\n", n, conn)
		}(n)
	}

	for n := 4; n > 0; n-- {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := p.Acquire(context.Background())
			if err != nil {
				panic(err)
			}
			defer p.Hijack(conn)

			fmt.Fprintf(conn, "Goodbye from hijacked connection %p!\n", conn)
		}()
	}

	wg.Wait()
}
