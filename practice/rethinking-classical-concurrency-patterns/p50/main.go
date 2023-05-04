package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

var addr net.Addr
var addrOnce sync.Once

type Pool struct {
	mu       sync.Mutex
	cond     sync.Cond
	numConns int
	limit    int
	idle     []net.Conn
}

func NewPool(limit int) *Pool {
	p := &Pool{
		limit: limit,
	}
	p.cond.L = &p.mu

	return p
}

func (p *Pool) Acquire() (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 透過 condition variable 確認 somthing changed 時是否滿足可繼續執行的條件
	// 只有等到有閒置的 conn 或是 conns 總數未達上限時才可以繼續
	for len(p.idle) == 0 && p.numConns >= p.limit {
		p.cond.Wait()
	}

	if len(p.idle) > 0 {
		c := p.idle[len(p.idle)-1]
		p.idle = p.idle[:len(p.idle)-1]
		return c, nil
	}

	c, err := dial()
	if err == nil {
		p.numConns++
	}

	return c, err
}

func (p *Pool) Release(c net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.idle = append(p.idle, c)
	p.cond.Signal()
}

func (p *Pool) Hijack() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.numConns--
	p.cond.Signal()
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

			conn, err := p.Acquire()
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

			conn, err := p.Acquire()
			if err != nil {
				panic(err)
			}
			defer p.Hijack()

			fmt.Fprintf(conn, "Goobye from hijacked connection %p!\n", conn)
		}()
	}

	wg.Wait()
}
