package simplecache

type SimpleCache struct {
	lookup   map[string]*dataNode //
	list     *dataList            //
	capacity int                  //
}

func Constructor(capacity int) SimpleCache {
	return SimpleCache{
		lookup:   map[string]*dataNode{},
		list:     CreateDoublyLinkedList(),
		capacity: capacity,
	}
}

func (this *SimpleCache) Get(key string) int {
	if _, exists := this.lookup[key]; !exists {
		return -1
	}

	node := this.lookup[key]
	this.list.MoveToTail(node)

	return node.GetValue()
}

func (this *SimpleCache) Set(key string, value, weight int) bool {
	if _, exists := this.lookup[key]; exists {
		return false
	}

	data := CreateNode(key, value, weight)
	if this.list.Size()+1 > this.capacity {
		new_data_score := float64(weight / -100)
		if deleted := this.list.CompareAndDeleteNode(new_data_score); !deleted {
			return false
		}
	}

	this.lookup[key] = data
	this.list.Append(data)

	return true
}

func (this *SimpleCache) GetAll() []int {
	return this.list.GetAll()
}
