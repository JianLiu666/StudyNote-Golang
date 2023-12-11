package trading

import "interview20231208/model"

type CustomHeap struct {
	orders   []*model.Order
	compFunc func(i, j *model.Order) bool
}

func NewCustomHeap(f func(i, j *model.Order) bool) *CustomHeap {
	return &CustomHeap{
		orders:   make([]*model.Order, 0),
		compFunc: f,
	}
}

func (h *CustomHeap) Peek() *model.Order {
	return h.orders[0]
}

func (h *CustomHeap) Len() int {
	return len(h.orders)
}

func (h *CustomHeap) Swap(i, j int) {
	h.orders[i], h.orders[j] = h.orders[j], h.orders[i]
}

func (h *CustomHeap) Less(i, j int) bool {
	return h.compFunc(h.orders[i], h.orders[j])
}

func (h *CustomHeap) Push(v any) {
	h.orders = append(h.orders, v.(*model.Order))
}

func (h *CustomHeap) Pop() (v any) {
	v, h.orders = h.orders[h.Len()-1], h.orders[:h.Len()-1]
	return
}
