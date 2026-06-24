# 07 Heap（堆）

## 第一性原理：为什么需要堆？

### 问题

```
HashMap → O(1) 查找特定元素

但如果是这些需求呢？

- 快速找到最大值/最小值 → O(1)
- 快速删除最大值/最小值 → O(log n)
- 动态插入 → O(log n)

数组：找最值 O(n)
排序数组：找最值 O(1)，但插入 O(n)
HashMap：找最值 O(n)
```

**需要一种新结构：专门为「最值」服务**

这就是堆。

---

## 核心特性

### 堆的定义

```
最小堆：
- 完全二叉树
- 父节点 ≤ 子节点

       1        ← 最小值（根节点）
     /   \
    3     6
   / \   /
  5   9 8
```

```
最大堆：
- 完全二叉树
- 父节点 ≥ 子节点

       9        ← 最大值（根节点）
     /   \
    6     8
   / \   /
  3   5 1
```

### 为什么用数组存储？

```
完全二叉树 → 可以紧凑地存在数组中

       1         索引 0
     /   \
    3     6      索引 1, 2
   / \   /
  5   9 8        索引 3, 4, 5

数组：[1, 3, 6, 5, 9, 8]

父子关系公式：
parent(i) = (i - 1) / 2
left(i)   = 2*i + 1
right(i)  = 2*i + 2
```

---

## 时间复杂度

| 操作 | 复杂度 | 说明 |
|------|--------|------|
| Peek (取最值) | O(1) | 根节点就是 |
| Push (插入) | O(log n) | 从底部上浮 |
| Pop (删除最值) | O(log n) | 用最后一个替换根，然后下沉 |
| Build (建堆) | O(n) | 从后往前 heapify |

---

## Go 代码实现

### 自定义最小堆

```go
// code/07-heap/minheap.go

package heap

// MinHeap 最小堆
type MinHeap struct {
    items []int
}

// NewMinHeap 创建最小堆
func NewMinHeap() *MinHeap {
    return &MinHeap{items: make([]int, 0)}
}

// Peek 查看最小值 O(1)
func (h *MinHeap) Peek() (int, bool) {
    if len(h.items) == 0 {
        return 0, false
    }
    return h.items[0], true
}

// Push 插入元素 O(log n)
func (h *MinHeap) Push(val int) {
    h.items = append(h.items, val)
    h.siftUp(len(h.items) - 1)
}

// Pop 删除并返回最小值 O(log n)
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
```

### Go 标准库 container/heap

```go
// Go 标准库的堆接口
// 只要实现了这 5 个方法，就可以用 heap 包

import "container/heap"

type Item struct {
    value    string
    priority int
    index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x any)        { *pq = append(*pq, x.(*Item)) }
func (pq *PriorityQueue) Pop() any {
    old := *pq
    n := len(old)
    item := old[n-1]
    *pq = old[0 : n-1]
    return item
}

// 使用
pq := &PriorityQueue{}
heap.Init(pq)
heap.Push(pq, &Item{value: "task1", priority: 3})
heap.Push(pq, &Item{value: "task2", priority: 1})
item := heap.Pop(pq).(*Item) // priority=1 的先出
```

---

## 工程案例

### 1. Go Timer（时间轮+堆）

```go
// Go 的 time.Timer 底层使用四叉堆

// runtime/time.go
// timer 存储在四叉堆中，按触发时间排序
// 四叉堆比二叉堆更宽，减少堆高度

// 这棵"时间树"的最小值就是最近的定时器
// runtime 循环检查堆顶 → 触发了就执行

// 本质：用堆管理定时器，O(1) 获取最近到期时间
```

```
Go Timer 的堆：

        10ms       ← 最近到期
     /  |  |  \
  50ms 30ms 80ms 100ms

定时器轮询：
1. 取堆顶 → 10ms
2. 10ms 到期 → Pop，执行
3. 新的堆顶 → 30ms
4. 等待 20ms...
```

### 2. K8S Scheduler

```
kube-scheduler 的核心：优先级队列

待调度的 Pod 放入队列，按优先级出队

Pod1 (priority: 100)  ← 优先调度
Pod2 (priority: 50)
Pod3 (priority: 10)

实现：最大堆（priority 高的在堆顶）
```

### 3. TopK 问题

```go
// 找最大的 K 个元素 → 用最小堆
func topK(nums []int, k int) []int {
    h := NewMinHeap()
    for _, v := range nums {
        h.Push(v)
        if h.Len() > k {
            h.Pop() // 弹出最小的，保持堆中只有 K 个最大元素
        }
    }
    return h.Items()
}
```

**为什么找最大 K 个用最小堆？**
- 最小堆的堆顶是最小值
- 堆大小保持 K
- 新元素比堆顶大 → 弹出堆顶，插入新元素
- 最终堆中就是最大的 K 个

---

## 练习题

### TopK 和优先队列

详见 [`exercises/07-heap/`](../../exercises/07-heap/)

---

## 小结

```
堆 = 完全二叉树 + 父≤子(或≥) + 数组存储

为什么不能用二叉搜索树代替堆？
- BST 找最值 → O(log n) 要走到最左/最右
- 堆找最值 → O(1) 就是根节点
- BST 平衡成本高 → AVL/红黑树需要旋转
- 堆维护简单 → 只上浮/下沉

堆专为「最值」场景设计：简单、高效。

工程价值：
- Go Timer → 四叉堆管理定时器
- K8S Scheduler → 优先队列调度 Pod
- TopK → 海量数据中找热门
```

---

**上一篇：[06 HashMap](../06-hashmap/README.md)**
**下一篇：[08 树](../08-tree/README.md)**
