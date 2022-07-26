package simplecache

import (
	"math"
	"sort"
	"time"
)

type SampleCache struct {
	lookup   map[string]*dataNode // using hash table to store raw data
	size     int                  // num of data node in the linked list
	capacity int                  // maximum of data can hold
}

// Create sample cache strcuture
// @param capacity maximum of sample cache can hold
//
// @return Sample Cache instance
func CreateSampleCache(capacity int) SampleCache {
	return SampleCache{
		lookup:   map[string]*dataNode{},
		size:     0,
		capacity: capacity,
	}
}

// Get data from sample cache
// @param key custom data uuid
//
// @return custom data
func (this *SampleCache) Get(key string) int {
	if _, exists := this.lookup[key]; !exists {
		return -1
	}

	node := this.lookup[key]
	return node.GetValue()
}

// Put data into sample cache
// @param key custom data uuid
// @param value custom data
// @param weight of custom data
//
// @return bool successful or not
func (this *SampleCache) Set(key string, value, weight int) bool {
	if _, exists := this.lookup[key]; exists {
		return false
	}

	data := CreateNode(key, value, weight)
	if this.size+1 > this.capacity {
		new_data_score := float64(weight / -100)

		target_key := ""
		target_score := float64(0)
		for _, node := range this.lookup {
			dividend := float64(time.Now().UnixMilli() - node.lastAccessTime.UnixMilli() + 1)
			current_score := float64(node.weight) / math.Log(dividend)

			if target_key == "" || target_score > current_score {
				target_key = node.key
				target_score = current_score
			}
		}

		if target_score < new_data_score {
			delete(this.lookup, target_key)
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
func (this *SampleCache) GetAll() []int {
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
