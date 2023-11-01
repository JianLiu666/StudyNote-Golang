package main

type List[T any] struct {
	head, tail *element[T]
}

type element[T any] struct {
	next *element[T]
	val  T
}

func (list *List[T]) Push(v T) {
	if list.tail == nil {
		list.head = &element[T]{val: v}
		list.tail = list.head
	} else {
		list.tail.next = &element[T]{val: v}
		list.tail = list.tail.next
	}
}

func (list *List[T]) GetAll() []T {
	var elems []T
	for e := list.head; e != nil; e = e.next {
		elems = append(elems, e.val)
	}
	return elems
}
