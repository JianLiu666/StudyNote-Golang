package simplecache

import (
	"math"
	"time"
)

type dataList struct {
	head *dataNode // point to list head
	tail *dataNode // point to list tail
	size int       // num of data node in the linked list
}

// create doubly linked list
//
// @return dataList
func CreateDoublyLinkedList() *dataList {
	return &dataList{
		head: nil,
		tail: nil,
		size: 0,
	}
}

// get num of data in the doubly linked list
//
// @return int num of data
func (this *dataList) Size() int {
	return this.size
}

// append data node to list tail
// @param node data node
func (this *dataList) Append(node *dataNode) {
	if this.size == 0 {
		this.head = node
		this.tail = node
	} else {
		this.tail.next = node
		node.prev = this.tail
		this.tail = node
	}

	this.size++
}

// move specific node to list tail
// @param node specific node
func (this *dataList) MoveToTail(node *dataNode) {
	if this.size == 1 || node == this.tail {
		return
	}

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

// compare data node's score with given threshold and find specific node
// @param score the score of threshold
//
// @return bool deleted or not
func (this *dataList) CompareAndDeleteNode(score float64) bool {
	deleted := false

	current := this.head
	for current != nil && !deleted {
		dividend := float64(time.Now().UnixMilli() - current.lastAccessTime.UnixMilli() + 1)
		current_score := float64(current.weight) / math.Log(dividend)

		if current_score >= score {
			current = current.next
			continue
		}

		// delete specific node
		if current == this.head {
			this.head = this.head.next
			this.head.prev = nil
		} else if current == this.tail {
			this.tail = this.tail.prev
			this.tail.next = nil
		} else {
			current.prev.next = current.next
			current.next.prev = current.prev
		}
		current.next = nil
		current.prev = nil
		deleted = true
	}

	return deleted
}

// Get all data with the accessing ordering
//
// @return []int data list
func (this *dataList) GetAll() []int {
	res := []int{}

	current := this.head
	for current != nil {
		res = append(res, current.value)
		current = current.next
	}

	return res
}
