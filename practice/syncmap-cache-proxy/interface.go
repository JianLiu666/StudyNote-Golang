package main

type IProxy interface {
	Execute(db *MockDB, key string) any
}
