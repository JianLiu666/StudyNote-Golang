package simplecache

type SimpleCache struct {
	lookup   map[string]*dataNode // using hash table to store raw data
	list     *dataList            // using doubly linked list to keep data ordering
	capacity int                  // maximum of data can hold
}

// Create simple cache strcuture
// @param capacity maximum of simple cache can hold
//
// @return Simple Cache instance
func CreateSimpleCache(capacity int) SimpleCache {
	return SimpleCache{
		lookup:   map[string]*dataNode{},
		list:     CreateDoublyLinkedList(),
		capacity: capacity,
	}
}

// Get data from simple cache
// @param key custom data uuid
//
// @return custom data
func (this *SimpleCache) Get(key string) int {
	if _, exists := this.lookup[key]; !exists {
		return -1
	}

	node := this.lookup[key]
	this.list.MoveToTail(node)

	return node.GetValue()
}

// Put data into simple cache
// @param key custom data uuid
// @param value custom data
// @param weight of custom data
//
// @return bool successful or not
func (this *SimpleCache) Set(key string, value, weight int) bool {
	if _, exists := this.lookup[key]; exists {
		return false
	}

	data := CreateNode(key, value, weight)
	if this.list.Size()+1 > this.capacity {
		new_data_score := CalcScoreByWeight(weight)
		if deleted := this.list.CompareAndDeleteNode(new_data_score); !deleted {
			return false
		}
	}

	this.lookup[key] = data
	this.list.Append(data)

	return true
}

// Get all data with the accessing ordering
// Note: it won't change the last access time in each data
//
// @return []int data list
func (this *SimpleCache) GetAll() []int {
	return this.list.GetAll()
}
