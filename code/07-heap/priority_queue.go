package heap

import (
	"container/heap"
)

// Item 优先队列元素
type Item struct {
	Value    string
	Priority int
	index    int // heap 内部使用的索引
}

// PriorityQueue 实现 heap.Interface
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// 小顶堆：priority 小的优先
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	item.index = -1 // 标记为已移除
	*pq = old[0 : n-1]
	return item
}

// PushItem 插入元素
func (pq *PriorityQueue) PushItem(item *Item) {
	heap.Push(pq, item)
}

// PopItem 弹出最小优先级元素
func (pq *PriorityQueue) PopItem() *Item {
	return heap.Pop(pq).(*Item)
}
