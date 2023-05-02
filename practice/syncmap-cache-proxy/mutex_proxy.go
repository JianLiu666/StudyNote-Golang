package main

import "sync"

type mutexProxy struct {
	cache sync.Map
	mu    sync.Mutex
}

func NewMutexProxy() IProxy {
	return &mutexProxy{
		cache: sync.Map{},
		mu:    sync.Mutex{},
	}
}

func (p *mutexProxy) Execute(db *MockDB, key string) (value any) {
	p.mu.Lock()
	defer p.mu.Unlock()

	value, ok := p.cache.Load(key)
	if !ok {
		value = db.Execute(key)
		p.cache.Store(key, value)
	}

	return value
}
