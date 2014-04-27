package triblab

import (
	"container/heap"
	"trib"
	//"fmt"
)

//heap with capacity to improve performance
type tribHeap struct {
	sorter   *ByTime
	capacity int
}

func (self tribHeap) Len() int { return self.sorter.Len() }

func (self tribHeap) Less(i, j int) bool {
	return self.sorter.Less(i, j)
}

func (self tribHeap) Swap(i, j int) {
	self.sorter.Swap(i, j)
}

func (self *tribHeap) Push(x interface{}) {
	//fmt.Println("push x")
	item := x.(*trib.Trib)
	if self.sorter.Len() < self.capacity {
		heap.Push(self.sorter, item)
	} else {
		top := heap.Pop(self.sorter).(*trib.Trib)
		// We want Pop to give us the highest, not lowest, priority so we use greater than here.
		if (top == nil) || (!compare(item, top)) {
			heap.Pop(self.sorter)
			heap.Push(self.sorter, item)
		} else {
			heap.Push(self.sorter, top)
		}
	}
}

func (self *tribHeap) Pop() interface{} {
	//fmt.Println("heap Pop once")
	return heap.Pop(self.sorter)
}

func (self *ByTime) Push(x interface{}) {
	item := x.(*trib.Trib)
	*self = append(*self, item)
}

func (self *ByTime) Pop() interface{} {
	old := *self
	n := len(old)
	x := old[n-1]
	*self = old[0 : n-1]
	//fmt.Println("sorter pop", x, "from", n)
	return x
}
