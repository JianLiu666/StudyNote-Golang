package treemap

import (
	"fmt"
	"interview20231129/model"
	"interview20231129/pkg/e"
	"interview20231129/pkg/singlepool"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
)

// Single Pool
//
// 1. 使用 goDS 的 treemap (based on red-black tree) 作為配對池的主要資料結構:
//   - add single person 時的時間複雜度約在 O(logn)
//   - match 基於 iterator 操作, 時間複雜度為 O(n), 需要思考如何解決 (TODO)
//
// 2. 使用 hash table (lookup) 檢查重複創建的用戶:
//   - 每次檢查的時間複雜度為 O(1)
type singlePool struct {
	mu sync.Mutex

	lookup map[string]string
	boys   *treemap.Map
	girls  *treemap.Map
}

func newTreemapSinglePool() *singlePool {
	return NewTreemapSinglePool().(*singlePool)
}

func NewTreemapSinglePool() singlepool.SinglePool {
	return &singlePool{

		lookup: map[string]string{},
		boys:   treemap.NewWithStringComparator(),
		girls:  treemap.NewWithStringComparator(),
	}
}

// AddSinglePersonAndMatch 加入新用戶且根據配對規則進行配對與更新用戶狀態
//
// @param user 用戶資訊
func (s *singlePool) AddSinglePersonAndMatch(user *model.User) e.CODE {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.lookup[user.Name]; exists {
		return e.ERROR_ADD_DUPLICATED_USER
	}

	user.UUID = s.genHashKey(user)

	if user.Gender == 1 {
		s.match(user, s.boys, s.girls, func(a, b int) bool {
			return a > b
		})
	} else {
		s.match(user, s.girls, s.boys, func(a, b int) bool {
			return a < b
		})
	}

	return e.SUCCESS
}

// RemoveSinglePerson 移除指定用戶
//
// @param name 用戶姓名
func (s *singlePool) RemoveSinglePerson(name string) e.CODE {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.lookup[name]; !exists {
		return e.ERROR_USER_NOT_FOUND
	}

	s.boys.Remove(s.lookup[name])
	s.girls.Remove(s.lookup[name])
	delete(s.lookup, name)

	return e.SUCCESS
}

// match 根據配對規則進行配對
//
//  1. 將 user 與 candidatePool 進行匹配, 只要符合規則就將雙方的約會次數(NumDates) 同時減 1
//  2. 直到 candidatePool 的遍歷流程結束後, user 還有剩餘的約會次數時就將 user 加入到 userPool 內
//
// @param user 用戶資訊
//
// @param userPool 與用戶相同性別的配對池
//
// @param candidatePool 與用戶不同性別的配對池
//
// @param comp 配對規則 (true: 配對成功, false: 配對失敗)
func (s *singlePool) match(user *model.User, userPool, candidatePool *treemap.Map, comp func(a, b int) bool) {
	it := candidatePool.Iterator()
	for it.Next() {
		candidateName, candidateInfo := it.Key().(string), it.Value().(*model.User)
		if !comp(user.Height, candidateInfo.Height) {
			continue
		}

		user.NumDates--
		candidateInfo.NumDates--

		if candidateInfo.NumDates == 0 {
			candidatePool.Remove(candidateName)
			delete(s.lookup, candidateInfo.Name)
		}
	}

	if user.NumDates > 0 {
		userPool.Put(user.UUID, user)
		s.lookup[user.Name] = user.UUID
	}
}

// genHashKey 根據用戶資訊產生 unique key
//
// @param user 用戶資訊
func (s *singlePool) genHashKey(user *model.User) string {
	return fmt.Sprintf("%d-%s", user.Height, user.Name)
}
