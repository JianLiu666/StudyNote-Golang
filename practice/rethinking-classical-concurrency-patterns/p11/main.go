package main

import "fmt"

func main() {
	a := fetch("a")
	b := fetch("b")

	// 如果這樣寫的話, 就變成順序執行, 不是 concurrently
	// a := fetch("a")
	// b := fetch("b")

	consume(<-a, <-b)
}

// FUTURE: API
// 返回一個 proxy object (channel) 使調用者晚一點拿到結果
func fetch(s string) <-chan string {
	c := make(chan string, 1)
	go func() {
		// do something...
		c <- s
	}()
	return c
}

// FUTURE: CALL SITE
func consume(a, b string) {
	fmt.Println(a, b)
}
