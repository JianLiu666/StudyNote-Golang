package main

import "sync"

type syncMapProxy struct {
	cache sync.Map
}

func NewSyncMapProxy() IProxy {
	return &syncMapProxy{
		cache: sync.Map{},
	}
}

func (p *syncMapProxy) Execute(db *MockDB, key string) (value any) {
	var wg sync.WaitGroup
	wg.Add(1)

	callback := func() any {
		wg.Wait()
		// value 因為 closure 關係, 所以可以拿到從 db 取回的資料
		return value
	}

	// 瞬時當下所有的併發請求都會被 sync.map 擋下來
	fn, isSecond := p.cache.LoadOrStore(key, callback)
	if isSecond {
		return fn.(func() any)()
	}

	defer func() {
		// 當所有瞬時併發的 function call 都結束後就刪掉一切
		wg.Done()
		p.cache.Delete(key)
	}()

	return db.Execute(key)
}
