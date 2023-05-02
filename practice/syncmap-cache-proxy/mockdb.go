package main

type MockDB struct {
	times int
}

func NewDB() *MockDB {
	return &MockDB{
		times: 0,
	}
}

func (db *MockDB) Execute(key string) any {
	db.times++

	return 1
}

func (db *MockDB) Times() int {
	return db.times
}
