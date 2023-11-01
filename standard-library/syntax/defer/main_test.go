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

func Test_TrickyCases(t *testing.T) {
	// 匿名函式會對返回值 `r` 的值 +1, 因此在函式返回 `0` 時,
	// 在返回前被 defer 語句 +1
	f1 := func() (r int) {
		defer func() {
			r++
		}()
		return 0
	}

	// 函式返回 `5`, defer 函式增加的對象是 `t` 不是 `r`
	// 所以直接返回 `5`
	f2 := func() (r int) {
		t := 5
		defer func() {
			t = t + 5
		}()
		return t
	}

	// 函式返回 `1`, 且在返回前被 defer +5
	// 但要注意的是 defer 函式裡面的 `r` 是函式參數不是返回值的 `r`
	// 所以還是返回 1
	f3 := func() (r int) {
		defer func(r int) {
			r = r + 5
		}(r)
		return 1
	}

	fmt.Println(f1())
	fmt.Println(f2())
	fmt.Println(f3())
}
