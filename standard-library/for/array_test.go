package main

import (
	"fmt"
	"testing"
)

// 如果在遍歷 slice 的同時繼續加入元素到 slice, 是否會得到一個永遠不停止的循環?
func TestForloop(t *testing.T) {
	arr := []int{1, 2, 3}
	for _, v := range arr {
		arr = append(arr, v)
	}
	fmt.Println(arr)
}

// 指針問題
func TestForloopReference(t *testing.T) {
	arr := []int{1, 2, 3}
	newArr := []*int{}
	for _, v := range arr {
		// 正確的做法應該是把 &arr[i] 賦值進去
		newArr = append(newArr, &v)
	}
	for _, v := range newArr {
		fmt.Println(*v)
	}
}
