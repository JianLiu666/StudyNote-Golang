package main

import "fmt"

func main() {
	// Producer-Consumer Queue: Call Site
	for r := range glob("[ab]*") {
		fmt.Println(r)
	}
}

// Producer-Consumer Queue: API
func glob(pattern string) <-chan rune {
	c := make(chan rune)
	go func() {
		defer close(c)
		for _, r := range pattern {
			c <- r
		}
	}()
	return c
}
