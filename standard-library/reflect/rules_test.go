package main

import (
	"fmt"
	"reflect"
	"testing"
)

/**
 * Test Example from: https://draveness.me/golang/docs/part2-foundation/ch04-basic/golang-reflect/
 */

// Reflection goes from interface value to reflection object.
// 反射可以從接口值得到反射對象
func Test_Rule1(t *testing.T) {
	author := "jian"
	fmt.Println("TypeOf author:", reflect.TypeOf(author))
	fmt.Println("ValueOf author:", reflect.ValueOf(author))
}

// Reflection goes from reflection object to interface value.
// 可以從反射對象得到接口值
func Test_Rule2(t *testing.T) {
	v := reflect.ValueOf(1)
	// 只有在需要將反射對象轉換回基本類型時才需要 explicit cast
	obj := v.Interface().(int)
	fmt.Println(obj)
}

// To modify a reflection object, the value must be settable.
// 要修改反射對象, 該值必須可以修改
func Test_Rule3(t *testing.T) {
	// panic: 因為我們得到的反射對象跟原本的 i 沒有關係 (deep copy)
	// i := 1
	// v := reflect.ValueOf(i)
	// v.SetInt(10)
	// fmt.Println(i)

	i := 1
	// 取得指向 i 的 pointer
	v := reflect.ValueOf(&i)
	// 取得 pointer i 實際指向的 variable 後並更新內容
	v.Elem().SetInt(10)
	fmt.Println(i)

	// 上面代碼等價於
	// i := 1
	// v := &i
	// *v = 10
}
