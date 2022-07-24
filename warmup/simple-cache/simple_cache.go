package simplecache

import (
	"math"
	"time"
)

type dataNode struct {
	key            string    //
	val            int       //
	weight         int       //
	lastAccessTime time.Time //
	next           *dataNode //
	prev           *dataNode //
}

type SimpleCache struct {
	lookup   map[string]*dataNode //
	head     *dataNode            //
	tail     *dataNode            //
	size     int                  //
	capacity int                  //
}

func Constructor(capacity int) SimpleCache {
	return SimpleCache{
		lookup:   map[string]*dataNode{},
		head:     nil,
		tail:     nil,
		size:     0,
		capacity: capacity,
	}
}

func (this *SimpleCache) Get(key string) int {
	if _, exists := this.lookup[key]; !exists {
		return -1
	}

	this.lookup[key].lastAccessTime = time.Now()
	node := this.lookup[key]

	// 維護 linked list, 把最新的資料移動到最末端
	if this.size > 1 && node != this.tail {
		if node == this.head {
			this.head = this.head.next
			this.head.prev = nil
		} else {
			node.prev.next = node.next
			node.next.prev = node.prev
		}
		this.tail.next = node
		node.prev = this.tail
		node.next = nil
		this.tail = node
	}

	return this.lookup[key].val
}

func (this *SimpleCache) Set(key string, value, weight int) bool {
	if _, exists := this.lookup[key]; exists {
		return false
	}

	// step.1 prepare data
	data := &dataNode{
		key:            key,
		val:            value,
		weight:         weight,
		lastAccessTime: time.Now(),
		next:           nil,
		prev:           nil,
	}

	// step.2 檢查是否超過上限
	if this.size+1 > this.capacity {
		// step.3 找到要刪除的資料
		new_data_score := weight / -100
		current := this.head

		deleted := false
		for current.next != nil && !deleted {
			dividend := float64(time.Now().UnixMilli() - current.lastAccessTime.UnixMilli() + 1)
			old_data_score := int(float64(weight) / math.Log(dividend))

			if old_data_score <= new_data_score {
				// do something
				if current == this.head {
					this.head = this.head.next
					this.head.prev = nil
				} else {
					current.prev.next = current.next
					current.next.prev = current.prev
				}
				deleted = true
			} else {
				// 這個不能刪, 繼續往下找
				current = current.next
			}
		}

		if !deleted {
			return false
		}
	}

	this.lookup[key] = data

	// 把資料加到 linked list
	if this.size == 0 {
		this.head = data
		this.tail = data
	} else {
		this.tail.next = data
		data.prev = this.tail
		this.tail = data
	}
	this.size++

	return true
}

func (this *SimpleCache) GetAll() []int {
	res := []int{}

	current := this.head
	for current != nil {
		res = append(res, current.val)
		current = current.next
	}

	return res
}
