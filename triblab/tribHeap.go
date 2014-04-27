package triblab

import (
    "trib"
    "container/heap"
)

type tribHeap struct {
    sorter *ByTime
    capacity int
}

func (self tribHeap) Len() int { return self.sorter.len() }

func (self tribHeap) Less(i, j int) bool {
    // We want Pop to give us the highest, not lowest, priority so we use greater than here.
    return sorter.Less(j, i)
}

func (self tribHeap) Swap(i, j int) {
    self.sorter.Swap(i, j)
}

func (self *tribHeap) Push(x interface{}) {
    item := x.(*trib.Trib)
    //items := *(self.sorter)
    if self.sorter.Len() < self.capacity{
        heap.Push(self.sorter, item)
    } else {
        top := self.sorter.Peek()
        if top <  {
            self.sorter.Pop()
            heap.Push(self.sorter, item)
        }
    }
} 






// An Item is something we manage in a priority queue.
type Item struct {
    value    string // The value of the item; arbitrary.
    priority int    // The priority of the item in the queue.
    //index int
    
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

type MaxHeap struct{
    pq *PriorityQueue
    capacity int
}

func (self *MaxHeap) Push( x interface{}){
    //var item *Item
    item := x.(*Item)
    items := *(self.pq)
    if len(items) < self.capacity{
        heap.Push(self.pq, item)
    } else {
        top := self.pq.Peek()
        fmt.Printf("item is %d, top is %d\n",item.priority, top.(*Item).priority)
        if item.priority < top.(*Item).priority{
            top = self.pq.Pop()
            heap.Push(self.pq, item)
        }
    }
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    // We want Pop to give us the highest, not lowest, priority so we use greater than here.
    return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    //pq[i].index = i
    //pq[j].index = j
    
}

func (pq *PriorityQueue) Push(x interface{}) {
    //n := len(*pq)
    item := x.(*Item)
    //item.index = n
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[0]
    //item.index = -1 // for safety
    *pq = old[1 : n]
    return item
}

func (pq *PriorityQueue) Peek() interface{} {
    old := *pq
    //n := len(old)
    item := old[0]
    //item.index = -1 // for safety
    //*pq = old[0 : n]
    return item
    
}

func main(){
    item1 := &Item{"third",3}
    item2 := &Item{"forty",40}
    item3 := &Item{"first",1}
    item4 := &Item{"ten",10}
    item5 := &Item{"twenty",20}
    item6 := &Item{"twohundred",200}
    item7 := &Item{"second",2}
    item8 := &Item{"seven",7}
    
    pq := &PriorityQueue{}
    heap.Init(pq)
    mh := &MaxHeap{pq,5}
    mh.Push(item1)
    mh.Push(item2)
    mh.Push(item3)
    mh.Push(item4)
    mh.Push(item5)
    mh.Push(item6)
    mh.Push(item7)
    mh.Push(item8)
    
    
    for pq.Len() > 0 {
        item := heap.Pop(mh.pq).(*Item)
        fmt.Printf("%.2d:%s ", item.priority, item.value)
    }
    
}
