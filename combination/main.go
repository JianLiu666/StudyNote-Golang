package main

import (
	"fmt"
)

type ITalk interface {
	Hello()
}

type A struct {
}

func (a *A) Hello() {
	fmt.Println("this is a")
}

type B struct {
	ITalk
}

func (b *B) Hello() {
	fmt.Println("this is b")
}

type C struct {
	ITalk
}

type D struct {
}

func (d *D) Hello() {
	fmt.Println("this is d")
}

func main() {
	a := A{}
	a.Hello()

	fmt.Println("-")

	// B 中實現了與 ITalk 同名的方法 Hello, 此時直接訪問 b.Hello 是訪問到 b 自己的方法,
	// 想要訪問到 A 的方法需要透過 b.ITalk.Hello() 來完成.
	b := B{&A{}}
	b.Hello()
	b.ITalk.Hello()

	fmt.Println("-")

	c := C{&A{}}
	c.Hello()

	c.ITalk = &D{}
	c.Hello()
}
