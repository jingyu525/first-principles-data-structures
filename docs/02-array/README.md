# 02 数组

## 第一性原理：为什么需要数组？

### 问题

```
内存是线性的（地址连续增长）

如果我要存 5 个整数，最简单的方式是什么？

→ 把它们排在一起
→ 每个元素占用固定大小
→ 通过起始地址 + 偏移量直接定位
```

这就是数组。

---

## 核心特性

### 连续内存分配

```
内存地址:  1000  1004  1008  1012  1016
            ↓     ↓     ↓     ↓     ↓
数组元素:  [1]   [2]   [3]   [4]   [5]
索引:       0     1     2     3     4
```

```
元素地址 = 基地址 + 索引 × 元素大小
arr[3] 地址 = 1000 + 3 × 4 = 1012
```

### 随机访问 O(1)

```go
arr := [5]int{10, 20, 30, 40, 50}
// arr[2] 直接计算地址，一次内存访问
fmt.Println(arr[2]) // 30
```

### 插入/删除 O(n)

```go
// 在索引 2 插入 25
// 需要将索引 2~4 的元素全部后移
// [10, 20, _, 30, 40, 50]
//          ↑ 插入 25
// [10, 20, 25, 30, 40, 50]
```

---

## 时间复杂度

| 操作 | 复杂度 | 原因 |
|------|--------|------|
| 随机访问 | O(1) | 地址公式直接计算 |
| 尾部插入 | O(1) | 已知末尾位置 |
| 中间插入 | O(n) | 需要移动后续元素 |
| 删除 | O(n) | 需要移动后续元素 |
| 搜索（无序） | O(n) | 逐个比较 |
| 搜索（有序） | O(log n) | 二分查找 |

---

## Go 代码实现

```go
// code/02-array/array.go

package array

// StaticArray 静态数组
type StaticArray [5]int

// DynamicArray 动态数组（类似 Go Slice 的简化实现）
type DynamicArray struct {
    data     []int
    size     int // 当前元素数量
    capacity int // 容量
}

// NewDynamicArray 创建动态数组
func NewDynamicArray(capacity int) *DynamicArray {
    return &DynamicArray{
        data:     make([]int, capacity),
        size:     0,
        capacity: capacity,
    }
}

// Get 随机访问 O(1)
func (a *DynamicArray) Get(index int) (int, bool) {
    if index < 0 || index >= a.size {
        return 0, false
    }
    return a.data[index], true
}

// Append 尾部插入 O(1) 均摊
func (a *DynamicArray) Append(val int) {
    if a.size == a.capacity {
        a.resize(a.capacity * 2)
    }
    a.data[a.size] = val
    a.size++
}

// Insert 中间插入 O(n)
func (a *DynamicArray) Insert(index, val int) bool {
    if index < 0 || index > a.size {
        return false
    }
    if a.size == a.capacity {
        a.resize(a.capacity * 2)
    }
    // 后移元素
    for i := a.size; i > index; i-- {
        a.data[i] = a.data[i-1]
    }
    a.data[index] = val
    a.size++
    return true
}

// Delete 删除 O(n)
func (a *DynamicArray) Delete(index int) bool {
    if index < 0 || index >= a.size {
        return false
    }
    // 前移元素
    for i := index; i < a.size-1; i++ {
        a.data[i] = a.data[i+1]
    }
    a.size--
    return true
}

// Search 线性查找 O(n)
func (a *DynamicArray) Search(val int) int {
    for i := 0; i < a.size; i++ {
        if a.data[i] == val {
            return i
        }
    }
    return -1
}

// BinarySearch 二分查找 O(log n)，要求数组有序
func (a *DynamicArray) BinarySearch(val int) int {
    lo, hi := 0, a.size-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if a.data[mid] == val {
            return mid
        } else if a.data[mid] < val {
            lo = mid + 1
        } else {
            hi = mid - 1
        }
    }
    return -1
}

func (a *DynamicArray) resize(newCap int) {
    newData := make([]int, newCap)
    copy(newData, a.data[:a.size])
    a.data = newData
    a.capacity = newCap
}
```

---

## 工程案例

### 1. Go Slice

Go 的切片底层就是数组：

```go
s := make([]int, 0, 10) // len=0, cap=10
// 底层结构：
// type slice struct {
//     array unsafe.Pointer  // 指向底层数组
//     len   int
//     cap   int
// }

s = append(s, 1, 2, 3) // 尾部追加，O(1) 均摊

// 扩容规则（Go 1.18+）：
// cap < 256 → 翻倍
// cap >= 256 → 约 1.25 倍增长
```

**为什么 Go Slice 这样设计？**

- 连续内存 → CPU 缓存友好
- 动态扩容 → 自动管理内存
- 共享底层数组 → 零拷贝切片操作

### 2. Kafka 缓冲区

```
Kafka Producer 收集消息 → 批量发送

为什么用数组（队列）作为缓冲区？

1. 消息顺序写入 → 连续内存高效
2. 批量读取 → 顺序 IO
3. 数组的预分配 → 减少 GC 压力
```

```
Producer → [msg1][msg2][msg3]... → Batch Send → Broker
                ↑
            数组缓冲区（RecordAccumulator）
```

### 3. 分页查询

```sql
-- MySQL LIMIT 本质是数组的 offset + limit
SELECT * FROM users LIMIT 10 OFFSET 20;
-- 底层：跳过前 20 条，取 10 条
```

```go
// 分页查询的数组视角
func paginate(arr []int, page, pageSize int) []int {
    start := (page - 1) * pageSize
    end := start + pageSize
    if start >= len(arr) {
        return nil
    }
    if end > len(arr) {
        end = len(arr)
    }
    return arr[start:end]
}
```

---

## 练习题

### 1. 反转数组

```go
// 要求：原地反转，O(1) 额外空间
func reverse(arr []int) {
    for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
        arr[i], arr[j] = arr[j], arr[i]
    }
}
```

### 2. 二分查找

实现二分查找，返回目标值的索引，不存在则返回 -1。

详见 [`exercises/02-array/`](../../exercises/02-array/)

---

## 小结

```
数组 = 连续内存 + 随机访问 O(1)

工程价值：
- Go Slice → 动态数组
- Kafka → 批量缓冲区
- MySQL → 分页查询

核心权衡：
- 优点：访问快、缓存友好
- 缺点：插入删除慢、需要连续内存
```

---

**上一篇：[01 第一性原理](../01-first-principle/README.md)**
**下一篇：[03 链表](../03-linked-list/README.md)**
