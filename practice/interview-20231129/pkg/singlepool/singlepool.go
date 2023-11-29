package singlepool

import (
	"fmt"
	"interview20231129/model"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
)

type singlePool struct {
	mu sync.Mutex

	lookup map[string]string
	boys   *treemap.Map
	girls  *treemap.Map
}

func NewSinglePool() *singlePool {
	return &singlePool{

		lookup: map[string]string{},
		boys:   treemap.NewWithStringComparator(),
		girls:  treemap.NewWithStringComparator(),
	}
}

func (s *singlePool) AddSinglePersonAndMatch(user *model.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.lookup[user.Name]; exists {
		return
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
}

func (s *singlePool) RemoveSinglePerson(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.lookup[name]; !exists {
		return
	}

	s.boys.Remove(s.lookup[name])
	s.girls.Remove(s.lookup[name])
	delete(s.lookup, name)
}

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

func (s *singlePool) genHashKey(user *model.User) string {
	return fmt.Sprintf("%d-%s", user.Height, user.Name)
}
