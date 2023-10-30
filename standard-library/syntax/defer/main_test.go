package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_PreCompute(t *testing.T) {
	startedAt := time.Now()
	// time.Since(startedAt) 不是在函數退出之前計算的, 而是在 defer 關鍵字調用時計算
	// 所以 defer 關鍵字會立刻 copy 函數中引用的外部參數
	// defer fmt.Println(time.Since(startedAt))

	// 改成匿名函數的方式, defer 調用時會 copy 匿名函數的 pointer, 等到函數返回前才實際執行匿名函數結果
	defer func() { fmt.Println(time.Since(startedAt)) }()

	time.Sleep(time.Second)
}
