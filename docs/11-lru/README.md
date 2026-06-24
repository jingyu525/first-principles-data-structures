# 11 LRU 缓存

## 第一性原理：为什么需要 LRU？

### 问题

```
缓存空间有限 → 满了要淘汰旧的

淘汰策略：
- FIFO：淘汰最早放入的 — 不管是否常用
- LFU：淘汰使用次数最少的 — 需要维护计数
- LRU：淘汰最久未使用的 — 兼顾时间和频率

为什么 LRU 是最常用的？
→ 符合时间局部性原理
→ 最近使用的数据，更可能再次被使用
```

---

## 核心设计

### LRU = HashMap + 双向链表

```
为什么需要两个结构？

HashMap → O(1) 查找缓存数据
双向链表 → O(1) 移动/删除节点（维护使用顺序）

单用 HashMap：无法维护顺序
单用链表：查找 O(n)
两者结合：各取所长
```

```
LRU 结构示意：

HashMap:                   双向链表（最近 → 最久）：
┌──────────┐              ┌───┐    ┌───┐    ┌───┐
│ key1 → ●─┼─────────────→│ k1│←──→│ k2│←──→│ k3│
├──────────┤              └───┘    └───┘    └───┘
│ key2 → ●─┼─────────────→  head              tail
├──────────┤                 ↑                 ↑
│ key3 → ●─┼─────────────→ 最近使用         最久未使用
└──────────┘
```

---

## 操作分析

### Get(key)

```
1. HashMap 查找 → O(1)
2. 找到 → 移动到链表头部（最近使用） → O(1)
3. 未找到 → 返回空

因为 HashMap 存了链表节点指针 → 移动只需改几个指针
```

### Put(key, value)

```
1. key 已存在 → 更新 value → 移动到头部 → O(1)
2. key 不存在：
   a. 缓存未满 → 创建新节点 → 放到头部 → HashMap 记录 → O(1)
   b. 缓存已满 → 删除尾部节点（最久未使用） → 创建新节点 → 放到头部 → O(1)
```

---

## 时间复杂度

| 操作 | 复杂度 |
|------|--------|
| Get | O(1) |
| Put | O(1) |

全部 O(1) — 这是靠 HashMap 和双向链表的精巧配合。

---

## Go 代码实现

```go
// code/11-lru/lru.go

package lru

// Node 双向链表节点
type Node struct {
    key, value int
    prev, next *Node
}

// LRUCache LRU 缓存
type LRUCache struct {
    capacity int
    cache    map[int]*Node
    head     *Node // 哨兵头（最近使用）
    tail     *Node // 哨兵尾（最久未使用）
}

// Constructor 创建 LRU
func Constructor(capacity int) LRUCache {
    l := LRUCache{
        capacity: capacity,
        cache:    make(map[int]*Node),
        head:     &Node{},
        tail:     &Node{},
    }
    l.head.next = l.tail
    l.tail.prev = l.head
    return l
}

// Get 获取值 O(1)
func (l *LRUCache) Get(key int) int {
    if node, ok := l.cache[key]; ok {
        l.moveToHead(node)
        return node.value
    }
    return -1
}

// Put 插入值 O(1)
func (l *LRUCache) Put(key, value int) {
    if node, ok := l.cache[key]; ok {
        node.value = value
        l.moveToHead(node)
        return
    }

    newNode := &Node{key: key, value: value}
    l.cache[key] = newNode
    l.addToHead(newNode)

    if len(l.cache) > l.capacity {
        removed := l.removeTail()
        delete(l.cache, removed.key)
    }
}

// addToHead 添加到头部（最近使用）
func (l *LRUCache) addToHead(node *Node) {
    node.prev = l.head
    node.next = l.head.next
    l.head.next.prev = node
    l.head.next = node
}

// removeNode 删除节点
func (l *LRUCache) removeNode(node *Node) {
    node.prev.next = node.next
    node.next.prev = node.prev
}

// moveToHead 移到头部
func (l *LRUCache) moveToHead(node *Node) {
    l.removeNode(node)
    l.addToHead(node)
}

// removeTail 删除尾部（最久未使用）
func (l *LRUCache) removeTail() *Node {
    node := l.tail.prev
    l.removeNode(node)
    return node
}
```

---

## 工程案例

### 1. Redis 内存淘汰

```
Redis maxmemory-policy 中的 LRU 相关策略：

- volatile-lru：在设置过期时间的 key 中，淘汰 LRU
- allkeys-lru：在所有 key 中，淘汰 LRU
- volatile-lfu：LFU 版本

Redis 的 LRU 是近似 LRU：
- 不是维护全局链表（内存开销大）
- 而是采样 N 个 key，淘汰其中时间戳最老的
- maxmemory-samples 控制采样数量（默认 5）
```

```c
// Redis 近似 LRU
// 每个对象有一个 lru 字段，记录上次访问时间

// 淘汰时：
// 1. 随机采样 maxmemory-samples 个 key
// 2. 淘汰其中 lru 值最小的（即最久未访问）
// 3. 不够继续采样淘汰

// 不是精确 LRU，但性能好、内存开销小
```

### 2. 浏览器缓存

```
Chrome 的 HTTP 缓存：

请求资源 → 先查缓存（内存/磁盘）
命中 → 直接使用
未命中 → 网络请求 → 存入缓存

淘汰策略：LRU（有限缓存空间）
```

### 3. CDN 缓存

```
CDN 边缘节点空间有限：

热点内容 → 保留在缓存
冷门内容 → LRU 淘汰

为什么不用 LFU？
→ LFU 维护计数开销大
→ 老的热点内容会一直占据空间
→ LRU 能自然淘汰过气内容
```

---

## 练习题

见 [`exercises/11-lru/`](../../exercises/11-lru/)

---

## 小结

```
LRU = HashMap（查找）+ 双向链表（顺序维护）

思路：
1. 用 HashMap 实现 O(1) 查找
2. 用双向链表维护使用顺序
3. 链表头部 = 最近使用，尾部 = 最久未使用
4. 满了就删尾部

这就是工程中的「组合设计」：
不是发明新结构，而是组合已知结构解决问题。

工程价值：
- Redis 近似 LRU → 内存淘汰
- 浏览器 → HTTP 缓存
- CDN → 边缘缓存
```

---

**上一篇：[10 B+Tree](../10-bplus-tree/README.md)**
**下一篇：[12 一致性哈希](../12-consistent-hash/README.md)**
