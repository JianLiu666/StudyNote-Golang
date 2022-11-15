package simplecache

import "time"

type dataNode struct {
	key            string    // uuid
	value          int       // custom data
	weight         int       // custom data
	lastAccessTime time.Time // user last accesst time
	next           *dataNode // point to next data node
	prev           *dataNode // point to previous data node
}

// create new data node
// @param key uuid
// @param value custom data
// @param weight custom data
//
// @return dataNode
func CreateNode(key string, value, weight int) *dataNode {
	return &dataNode{
		key:            key,
		value:          value,
		weight:         weight,
		lastAccessTime: time.Now(),
		next:           nil,
		prev:           nil,
	}
}

// Get value and update last access time
//
// @return int custom data
func (this *dataNode) GetValue() int {
	this.lastAccessTime = time.Now()

	return this.value
}
