# 03 链表

## 第一性原理：为什么需要链表？

### 问题

```
数组的问题：
- 插入/删除中间元素需要 O(n) 移动
- 需要连续的内存空间
- 扩容时需要拷贝整个数组

如果数据的修改频率远高于读取频率呢？

→ 放弃连续内存
→ 用指针串联数据
→ 插入删除只需改指针
```

这就是链表。

---

## 核心特性

### 非连续内存 + 指针串联

```
数组：
[1] → [2] → [3] → [4] → [5]   （内存连续）

链表：
[1|●]──→  [3|●]──→  [5|●]──→  nil   （内存分散）
 ↑          ↑          ↑
 head       │          │
            ← 指针      ← 指针
```

### 插入 O(1)

```
在节点 A 和 B 之间插入新节点 C：

插入前：A → B
插入后：A → C → B

只需改两个指针，不需要移动任何数据
```

### 删除 O(1)

```
删除节点 B：

删除前：A → B → C
删除后：A → C

只需改一个指针
```

### 查找 O(n)

```
链表不支持随机访问
找到第 k 个节点需要从头遍历 k 次

arr[100] → O(1)
list.Get(100) → O(100)
```

---

## 时间复杂度

| 操作 | 数组 | 链表 | 原因 |
|------|------|------|------|
| 随机访问 | O(1) | O(n) | 链表无地址公式 |
| 头部插入 | O(n) | O(1) | 链表改指针即可 |
| 尾部插入 | O(1) | O(1)/O(n) | 无尾指针时需要遍历 |
| 中间插入 | O(n) | O(1) | 链表只改指针 |
| 删除 | O(n) | O(1) | 链表只改指针 |
| 查找 | O(n) | O(n) | 都需要遍历 |

---

## Go 代码实现

```go
// code/03-linked-list/list.go

package list

// ListNode 单链表节点
type ListNode struct {
    Val  int
    Next *ListNode
}

// LinkedList 单链表
type LinkedList struct {
    head *ListNode
    size int
}

// NewLinkedList 创建链表
func NewLinkedList() *LinkedList {
    return &LinkedList{head: nil, size: 0}
}

// InsertHead 头部插入 O(1)
func (l *LinkedList) InsertHead(val int) {
    node := &ListNode{Val: val, Next: l.head}
    l.head = node
    l.size++
}

// InsertTail 尾部插入 O(n)
func (l *LinkedList) InsertTail(val int) {
    node := &ListNode{Val: val}
    if l.head == nil {
        l.head = node
        l.size++
        return
    }
    cur := l.head
    for cur.Next != nil {
        cur = cur.Next
    }
    cur.Next = node
    l.size++
}

// Delete 删除指定值的第一个节点 O(n)
func (l *LinkedList) Delete(val int) bool {
    if l.head == nil {
        return false
    }
    if l.head.Val == val {
        l.head = l.head.Next
        l.size--
        return true
    }
    cur := l.head
    for cur.Next != nil {
        if cur.Next.Val == val {
            cur.Next = cur.Next.Next
            l.size--
            return true
        }
        cur = cur.Next
    }
    return false
}

// Get 获取第 k 个元素 O(k)
func (l *LinkedList) Get(k int) (int, bool) {
    if k < 0 || k >= l.size {
        return 0, false
    }
    cur := l.head
    for i := 0; i < k; i++ {
        cur = cur.Next
    }
    return cur.Val, true
}

// Search 查找 O(n)
func (l *LinkedList) Search(val int) int {
    cur := l.head
    idx := 0
    for cur != nil {
        if cur.Val == val {
            return idx
        }
        cur = cur.Next
        idx++
    }
    return -1
}
```

### 双向链表

```go
// code/03-linked-list/double_list.go

// DoubleNode 双向链表节点
type DoubleNode struct {
    Val  int
    Prev *DoubleNode
    Next *DoubleNode
}

// DoubleLinkedList 双向链表
type DoubleLinkedList struct {
    head *DoubleNode
    tail *DoubleNode
    size int
}

// InsertTail O(1) — 有尾指针
func (l *DoubleLinkedList) InsertTail(val int) {
    node := &DoubleNode{Val: val}
    if l.tail == nil {
        l.head = node
        l.tail = node
    } else {
        l.tail.Next = node
        node.Prev = l.tail
        l.tail = node
    }
    l.size++
}

// Delete O(1) — 有前驱指针
func (l *DoubleLinkedList) Delete(node *DoubleNode) {
    if node.Prev != nil {
        node.Prev.Next = node.Next
    } else {
        l.head = node.Next
    }
    if node.Next != nil {
        node.Next.Prev = node.Prev
    } else {
        l.tail = node.Prev
    }
    l.size--
}
```

---

## 工程案例

### 1. LRU 缓存（链表核心应用）

```
LRU = HashMap + 双向链表

HashMap → O(1) 查找
双向链表 → O(1) 移动节点到头部（最近使用）
```

详细见 [11 LRU](../11-lru/README.md)

### 2. 浏览器历史记录

```
后退/前进 = 双向链表

[页面A] ⇄ [页面B] ⇄ [页面C]
                     ↑
                  当前页面

后退 → Prev
前进 → Next
```

### 3. Redis QuickList

```
Redis 3.2+ 的 List 类型使用 QuickList

QuickList = 双向链表 + 压缩列表(zipList)

结构：
head ⇄ [zipList1] ⇄ [zipList2] ⇄ [zipList3] ⇄ tail

为什么这样设计？
- 纯链表 → 内存碎片多、指针开销大
- 纯 zipList → 大列表插入慢
- QuickList → 折中方案：链表的灵活性 + 压缩列表的内存效率
```

```
Redis List 的工程权衡：

| 版本 | 结构 | 问题 |
|------|------|------|
| Redis 2.x | linkedlist + ziplist | 切换不灵活 |
| Redis 3.2+ | quicklist | 统一结构，自适应 |
```

---

## 练习题

### 1. 反转链表

```go
// 迭代法
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode
    cur := head
    for cur != nil {
        next := cur.Next // 暂存下一个
        cur.Next = prev  // 反转指针
        prev = cur       // prev 前进
        cur = next       // cur 前进
    }
    return prev
}
```

### 2. 链表环检测（Floyd 判圈算法）

```go
func hasCycle(head *ListNode) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            return true
        }
    }
    return false
}
```

---

## 小结

```
链表 = 非连续内存 + 指针串联

链表出现的原因 → 解决数组插入/删除 O(n) 的问题

工程价值：
- LRU → 双向链表 + HashMap
- Redis QuickList → 双向链表 + 压缩列表
- 浏览器历史 → 双向链表

有没有注意到一个模式？
工程中很少单独用链表，都是「链表 + 其他结构」的组合。
```

---

**上一篇：[02 数组](../02-array/README.md)**
**下一篇：[04 栈](../04-stack/README.md)**
