package simplecache

import (
	"sort"
)

type NormalCache struct {
	lookup   map[string]*dataNode // using hash table to store raw data
	size     int                  // num of data node in the linked list
	capacity int                  // maximum of data can hold
}

// Create normal cache strcuture
// @param capacity maximum of normal cache can hold
//
// @return normal Cache instance
func CreateNormalCache(capacity int) NormalCache {
	return NormalCache{
		lookup:   map[string]*dataNode{},
		size:     0,
		capacity: capacity,
	}
}

// Get data from normal cache
// @param key custom data uuid
//
// @return custom data
func (this *NormalCache) Get(key string) int {
	if _, exists := this.lookup[key]; !exists {
		return -1
	}

	node := this.lookup[key]
	return node.GetValue()
}

// Put data into normal cache
// @param key custom data uuid
// @param value custom data
// @param weight of custom data
//
// @return bool successful or not
func (this *NormalCache) Set(key string, value, weight int) bool {
	if _, exists := this.lookup[key]; exists {
		return false
	}

	data := CreateNode(key, value, weight)
	if this.size+1 > this.capacity {
		new_data_score := CalcScoreByWeight(weight)

		// find the least score of data
		least_cored_key := ""
		least_score := float64(0)
		for _, node := range this.lookup {
			current_score := CalcScoreByAccessTime(node.weight, node.lastAccessTime)

			if least_cored_key == "" || least_score > current_score {
				least_cored_key = node.key
				least_score = current_score
			}
		}

		if least_score < new_data_score {
			delete(this.lookup, least_cored_key)
			this.lookup[key] = data
			return true
		}
		return false

	} else {
		this.lookup[key] = data
		this.size++
	}

	return true
}

// Get all data with the accessing ordering
// Note: it won't change the last access time in each data
//
// @return []int data list
func (this *NormalCache) GetAll() []int {
	tmp := []*dataNode{}

	for _, node := range this.lookup {
		tmp = append(tmp, node)
	}

	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].lastAccessTime.UnixNano() < tmp[j].lastAccessTime.UnixNano()
	})

	res := []int{}
	for _, node := range tmp {
		res = append(res, node.value)
	}

	return res
}
