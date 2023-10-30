package main

import (
	"fmt"
	"reflect"
	"testing"
)

func Add(a, b int) int {
	return a + b
}

func TestDynamicInvoke(t *testing.T) {
	v := reflect.ValueOf(Add)
	if v.Kind() != reflect.Func {
		return
	}

	typ := v.Type()
	argv := make([]reflect.Value, typ.NumIn())
	for i := range argv {
		if typ.In(i).Kind() != reflect.Int {
			return
		}
		argv[i] = reflect.ValueOf(i)
	}

	result := v.Call(argv)
	if len(result) != 1 || result[0].Kind() != reflect.Int {
		return
	}

	fmt.Println(result[0].Int())
}
