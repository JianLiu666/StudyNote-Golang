package trading

import "interview20231208/model"

type CustomHeap struct {
	quantites int
	orders    []*model.Order
	compFunc  func(i, j *model.Order) bool
}

func NewCustomHeap(f func(i, j *model.Order) bool) *CustomHeap {
	return &CustomHeap{
		orders:   make([]*model.Order, 0),
		compFunc: f,
	}
}

func (h *CustomHeap) TotalQuantity() int {
	return h.quantites
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
	e := v.(*model.Order)
	h.orders = append(h.orders, e)
	h.quantites += e.Quantity
}

func (h *CustomHeap) Pop() (v any) {
	h.quantites -= h.orders[h.Len()-1].Quantity
	v, h.orders = h.orders[h.Len()-1], h.orders[:h.Len()-1]
	return
}
