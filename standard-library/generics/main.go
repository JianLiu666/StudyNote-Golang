package main

import "fmt"

func main() {
	data1 := []int32{1, 2, 3, 4, 5}
	data2 := []float32{1.1, 2.2, 3.3, 4.4, 5.5}
	data3 := []string{"a", "b", "c", "d", "e"}

	sum1 := Sum(data1)
	sum2 := Sum(data2)
	sum3 := Sum(data3)

	fmt.Println(sum1)
	fmt.Println(sum2)
	fmt.Println(sum3)

	var m = map[int]string{1: "2", 2: "4", 4: "8"}
	fmt.Println("keys:", MapKeys[int, string](m))

	l := List[int]{}
	l.Push(10)
	l.Push(20)
	l.Push(30)
	fmt.Println("list:", l.GetAll())
}
