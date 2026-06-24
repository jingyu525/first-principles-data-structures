# Heap 练习题

## 1. 数组中的第 K 个最大元素

### 思路

用大小为 K 的最小堆。遍历数组，维护堆中保持 K 个最大元素。堆顶即第 K 大。

### 实现

```go
func findKthLargest(nums []int, k int) int {
    h := &MinHeap{}

    for _, v := range nums {
        heap.Push(h, v)
        if h.Len() > k {
            heap.Pop(h)
        }
    }
    return heap.Pop(h).(int)
}

type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool   { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x any)          { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}
```

## 2. 前 K 个高频元素

### 思路

先统计频率（HashMap），然后用大小为 K 的最小堆（按频率排序）。

### 实现

```go
func topKFrequent(nums []int, k int) []int {
    freq := make(map[int]int)
    for _, v := range nums {
        freq[v]++
    }

    h := &FreqHeap{}
    heap.Init(h)

    for num, count := range freq {
        heap.Push(h, Item{num, count})
        if h.Len() > k {
            heap.Pop(h)
        }
    }

    result := make([]int, k)
    for i := k - 1; i >= 0; i-- {
        result[i] = heap.Pop(h).(Item).num
    }
    return result
}

type Item struct {
    num   int
    count int
}

type FreqHeap []Item

func (h FreqHeap) Len() int            { return len(h) }
func (h FreqHeap) Less(i, j int) bool   { return h[i].count < h[j].count }
func (h FreqHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *FreqHeap) Push(x any)          { *h = append(*h, x.(Item)) }
func (h *FreqHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}
```
