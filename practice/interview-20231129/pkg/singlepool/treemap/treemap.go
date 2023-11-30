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
//   - match 基於 iterator 操作, 時間複雜度為 O(n), 需要思考如何優化 (TODO)
//   - query 也基於 iterator 操作, 時間複雜度為 O(n), 需要思考如何優化 (TODO)
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
//
// @return e.Code 執行結果狀態碼
func (s *singlePool) AddSinglePersonAndMatch(user *model.User) e.CODE {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.lookup[user.Name]; exists {
		return e.ERROR_ADD_DUPLICATED_USER
	}

	user.UUID = s.genHashKey(user)

	if user.Gender == e.BOY {
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
//
// @return e.Code 執行結果狀態碼
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

// QuerySinglePeople 根據查詢條件返回符合條件的用戶
//
// @param limit 取回的用戶數量
//
// @param opts 查詢條件
//
// @return []*model.User 符合條件的用戶
//
// @return e.Code 執行結果狀態碼
func (s *singlePool) QuerySinglePeople(limit int, opts *singlepool.QueryOpts) ([]*model.User, e.CODE) {
	result := make([]*model.User, 0, limit)

	// 初始化 iterator
	boysIt := s.boys.Iterator()
	boyExists := false
	if opts.Gender != e.GIRL {
		boyExists = boysIt.Next()
	}

	girlsIt := s.girls.Iterator()
	girlExists := false
	if opts.Gender != e.BOY {
		girlExists = girlsIt.Next()
	}

	// 取得符合查詢條件的用戶, 直到查詢數量已滿或遍歷完所有的用戶
	for limit > 0 && (boyExists || girlExists) {
		if opts.Gender != e.GIRL {
			// 只要性別條件不是限定女生, 就需要查找男生用戶
			// 移動 position 到符合條件的 node, 或直到 end 為止
			for boyExists {
				boy := boysIt.Value().(*model.User)
				if s.validQueryOpts(boy, opts) {
					break
				}
				boyExists = boysIt.Next()
			}
		}
		if opts.Gender != e.BOY {
			// 只要性別條件不是限定男生, 就需要查找女生用戶
			// 移動 position 到符合條件的 node, 或直到 end 為止
			for girlExists {
				girl := girlsIt.Value().(*model.User)
				if s.validQueryOpts(girl, opts) {
					break
				}
				girlExists = girlsIt.Next()
			}
		}

		if boyExists && girlExists {
			// 如果同時有兩個用戶都滿足條件時, 優先取回身高較矮的用戶
			boy := boysIt.Value().(*model.User)
			girl := girlsIt.Value().(*model.User)
			if boy.Height < girl.Height {
				result = append(result, boy)
				boyExists = boysIt.Next()
			} else {
				result = append(result, girl)
				girlExists = girlsIt.Next()
			}

		} else if boyExists {
			boy := boysIt.Value().(*model.User)
			result = append(result, boy)
			boyExists = boysIt.Next()

		} else if girlExists {
			girl := girlsIt.Value().(*model.User)
			result = append(result, girl)
			girlExists = girlsIt.Next()
		}

		limit--
	}

	return result, e.SUCCESS
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

// validQueryOpts 檢查指定用戶是否完全符合查詢條件
//
// @param user 指定用戶
//
// @param opts 查詢條件
//
// @return bool 是否符合條件
func (s *singlePool) validQueryOpts(user *model.User, opts *singlepool.QueryOpts) bool {
	if opts.Name != "" && user.Name != opts.Name {
		return false
	}
	if (opts.Gender == e.BOY || opts.Gender == e.GIRL) && user.Gender != opts.Gender {
		return false
	}
	if opts.MinHeight != -1 && user.Height < opts.MinHeight {
		return false
	}
	if opts.MaxHeight != -1 && user.Height > opts.MaxHeight {
		return false
	}
	if opts.MinNumDates != -1 && user.NumDates < opts.MinNumDates {
		return false
	}
	if opts.MaxNumDates != -1 && user.NumDates > opts.MaxNumDates {
		return false
	}
	return true
}
