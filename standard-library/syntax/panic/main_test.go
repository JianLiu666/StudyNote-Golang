package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_CrossGoroutines(t *testing.T) {
	// 發生 panic 時只會調用同一個 goroutine 的 defer function
	defer println("in main")
	go func() {
		defer println("in goroutine")
		panic("")
	}()

	time.Sleep(time.Second)
}

func Test_NestedPanic(t *testing.T) {
	defer fmt.Println("in main")
	defer func() {
		defer func() {
			panic("panic again and again")
		}()
		panic("panic again")
	}()
	panic("panic once")
}

func Test_Recover(t *testing.T) {
	a := 0
	f(&a)
	fmt.Println(a)
}

func f(a *int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("catched: ", r)
		}
	}()

	*a = 1

	panic("panic!")
}
