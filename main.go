package main

import "log"

func main() {
	case2()
}

func case2() {
	a := [][]int{{1, 2, 3}}
	b := make([][]int, 1)
	// b = append(b, make([]int, 0))
	copy(b, a)

	b[0][0] = 4

	log.Print(a)
	log.Print(b)
}

func case1() {
	a := [][]int{{1, 2, 3}}
	b := make([][]int, 0)
	copy(b, a)

	b[0][0] = 4

	log.Print(a)
	log.Print(b)
}
