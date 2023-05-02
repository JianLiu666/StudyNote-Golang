package main

import (
	"fmt"
	"testing"
)

func TestMutexProxy(t *testing.T) {
	db := NewDB()
	proxy := NewMutexProxy()

	for i := 0; i < 10000; i++ {
		go proxy.Execute(db, "key")
	}

	fmt.Println(db.Times())
}

func TestSyncmapProxy(t *testing.T) {
	db := NewDB()
	proxy := NewSyncMapProxy()

	for i := 0; i < 10000; i++ {
		go proxy.Execute(db, "key")
	}

	fmt.Println(db.Times())
}
