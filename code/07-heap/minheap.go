package heap

// MinHeap 最小堆
type MinHeap struct {
	items []int
}

// NewMinHeap 创建最小堆
func NewMinHeap() *MinHeap {
	return &MinHeap{items: make([]int, 0)}
}

// Len 堆大小
func (h *MinHeap) Len() int {
	return len(h.items)
}

// Peek 查看最小值 O(1)
func (h *MinHeap) Peek() (int, bool) {
	if len(h.items) == 0 {
		return 0, false
	}
	return h.items[0], true
}

// Push 插入 O(log n)
func (h *MinHeap) Push(val int) {
	h.items = append(h.items, val)
	h.siftUp(len(h.items) - 1)
}

// Pop 删除最小值 O(log n)
func (h *MinHeap) Pop() (int, bool) {
	if len(h.items) == 0 {
		return 0, false
	}
	minVal := h.items[0]
	lastIdx := len(h.items) - 1
	h.items[0] = h.items[lastIdx]
	h.items = h.items[:lastIdx]
	if len(h.items) > 0 {
		h.siftDown(0)
	}
	return minVal, true
}

// siftUp 上浮
func (h *MinHeap) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.items[i] >= h.items[parent] {
			break
		}
		h.items[i], h.items[parent] = h.items[parent], h.items[i]
		i = parent
	}
}

// siftDown 下沉
func (h *MinHeap) siftDown(i int) {
	n := len(h.items)
	for {
		smallest := i
		left := 2*i + 1
		right := 2*i + 2

		if left < n && h.items[left] < h.items[smallest] {
			smallest = left
		}
		if right < n && h.items[right] < h.items[smallest] {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.items[i], h.items[smallest] = h.items[smallest], h.items[i]
		i = smallest
	}
}

// Items 返回底层数据
func (h *MinHeap) Items() []int {
	result := make([]int, len(h.items))
	copy(result, h.items)
	return result
}
