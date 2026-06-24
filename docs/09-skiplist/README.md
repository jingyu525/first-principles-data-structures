# 09 SkipList（跳表）

## 第一性原理：为什么需要跳表？

### 问题

```
链表查找 O(n) → 太慢
BST 查找 O(log n) → 但需要旋转保持平衡

有没有一种结构：
- 查找 O(log n)
- 不需要旋转
- 实现比红黑树简单
- 支持范围查询

→ SkipList（跳表）
```

### 核心思想

```
在链表上建索引 → "跳着走" → 加速查找

普通链表的查找：一步一步走
原始: 1 → 3 → 6 → 9 → 12 → 17 → 19 → 21 → 25 → null

加了索引后：
L2:    1 → → → → → 12 → → → → → → → 25
L1:    1 → → → 6 → → → 12 → → → 19 → → → 25
L0:    1 → 3 → 6 → 9 → 12 → 17 → 19 → 21 → 25
```

这就是「空间换时间」——用额外的指针换取 O(log n) 的查找速度。

---

## 核心特性

### 多层索引结构

```
L3: head ──────────────→ 25
L2: head ──────→ 12 ──→ 25
L1: head ──→ 6 ──→ 12 ──→ 19 ──→ 25
L0: head → 1 → 3 → 6 → 9 → 12 → 17 → 19 → 21 → 25

找 17：
L3: head → 25 (25 > 17，不前进)
L2: head → 12 → 25 (25 > 17，停在 12)
L1: 12 → 19 (19 > 17，停在 12)
L0: 12 → 17 → 找到！
```

### 随机层高

```
每个节点随机决定"长多高"：

插入节点时：
1. 初始层高 = 1
2. 抛硬币：正面 → 层高 +1，继续抛；反面 → 停止
3. 概率：层高为 k 的概率 = (1/2)^k

期望：
- 层高为 1 的节点：50%
- 层高为 2 的节点：25%
- 层高为 3 的节点：12.5%
- ...

为什么随机？
→ 不需要旋转 → 概率上保证 O(log n)
→ 实现极简
```

### 时间复杂度

| 操作 | 平均 | 最坏 |
|------|------|------|
| 查找 | O(log n) | O(n) |
| 插入 | O(log n) | O(n) |
| 删除 | O(log n) | O(n) |

最坏情况是因为随机层高理论上可能造成不平衡，但概率极低。

### 与红黑树对比

|  | SkipList | 红黑树 |
|------|---------|--------|
| 实现难度 | 简单 | 复杂 |
| 范围查询 | 天然支持 | 需要额外处理 |
| 平衡维护 | 随机，不旋转 | 旋转 + 变色 |
| 并发友好 | 可做无锁 | 难做无锁 |
| 内存开销 | 指针多 | 每个节点一个颜色位 |

---

## Redis ZSet 源码解析

```
Redis ZSet（有序集合）= Dict + SkipList

为什么用两个结构？
- Dict（HashMap）→ 按 member 查找 score → O(1)
- SkipList → 按 score 排序、范围查询 → O(log n)
```

```c
// Redis SkipList 节点（src/server.h 简化）
typedef struct zskiplistNode {
    sds ele;        // 成员（member）
    double score;   // 分数
    struct zskiplistNode *backward; // 后退指针（L0 层）
    struct zskiplistLevel {
        struct zskiplistNode *forward; // 前进指针
        unsigned long span;            // 跨度（两个节点间的距离）
    } level[];  // 柔性数组，每个节点的层数不同
} zskiplistNode;

typedef struct zskiplist {
    struct zskiplistNode *header, *tail;
    unsigned long length;   // 节点总数
    int level;              // 最大层数
} zskiplist;
```

### Redis JumpList 的设计亮点

```
1. span（跨度）
   每个 forward 指针记录跳过的节点数
   → ZRANK 命令可以直接计算排名 O(log n)
   → 不需要遍历计数

2. backward 指针（后退）
   → 支持双向遍历
   → ZREVRANGE 反向范围查询

3. 层高上限 32
   → 2^32 个节点足够
   → 内存可控

4. 允许相同 score
   → 相同 score 按 member 字典序排序
```

```
Redis ZSet 命令与数据结构：

ZADD key score member   → 插入 SkipList + Dict
ZRANGE key 0 10         → SkipList 范围查询
ZRANK key member        → Dict O(1) 找 score + SkipList 算 span
ZSCORE key member       → Dict O(1)
```

---

## Go 代码实现

```go
// code/09-skiplist/skiplist.go

package skiplist

import (
    "math/rand"
)

const maxLevel = 32

// Node 跳表节点
type Node struct {
    key   int
    value interface{}
    next  []*Node // 每层的下一个节点
}

// SkipList 跳表
type SkipList struct {
    head  *Node
    level int // 当前最高层
}

// NewSkipList 创建跳表
func NewSkipList() *SkipList {
    return &SkipList{
        head:  &Node{next: make([]*Node, maxLevel)},
        level: 1,
    }
}

// randomLevel 随机层高
func randomLevel() int {
    level := 1
    for rand.Float64() < 0.5 && level < maxLevel {
        level++
    }
    return level
}

// Search 查找 O(log n)
func (sl *SkipList) Search(key int) (interface{}, bool) {
    cur := sl.head
    // 从最高层开始向下搜索
    for i := sl.level - 1; i >= 0; i-- {
        for cur.next[i] != nil && cur.next[i].key < key {
            cur = cur.next[i]
        }
    }
    // cur 是 key 的前驱节点
    if cur.next[0] != nil && cur.next[0].key == key {
        return cur.next[0].value, true
    }
    return nil, false
}

// Insert 插入 O(log n)
func (sl *SkipList) Insert(key int, value interface{}) {
    // 记录每层的前驱节点
    update := make([]*Node, maxLevel)
    cur := sl.head

    for i := sl.level - 1; i >= 0; i-- {
        for cur.next[i] != nil && cur.next[i].key < key {
            cur = cur.next[i]
        }
        update[i] = cur
    }

    // 如果 key 已存在，更新 value
    if cur.next[0] != nil && cur.next[0].key == key {
        cur.next[0].value = value
        return
    }

    // 生成新节点
    level := randomLevel()
    if level > sl.level {
        for i := sl.level; i < level; i++ {
            update[i] = sl.head
        }
        sl.level = level
    }

    newNode := &Node{
        key:   key,
        value: value,
        next:  make([]*Node, level),
    }

    // 插入到每一层
    for i := 0; i < level; i++ {
        newNode.next[i] = update[i].next[i]
        update[i].next[i] = newNode
    }
}

// Delete 删除 O(log n)
func (sl *SkipList) Delete(key int) bool {
    update := make([]*Node, maxLevel)
    cur := sl.head

    for i := sl.level - 1; i >= 0; i-- {
        for cur.next[i] != nil && cur.next[i].key < key {
            cur = cur.next[i]
        }
        update[i] = cur
    }

    target := cur.next[0]
    if target == nil || target.key != key {
        return false
    }

    for i := 0; i < sl.level; i++ {
        if update[i].next[i] != target {
            break
        }
        update[i].next[i] = target.next[i]
    }

    // 降低层数
    for sl.level > 1 && sl.head.next[sl.level-1] == nil {
        sl.level--
    }
    return true
}
```

---

## 小结

```
SkipList = 多层索引链表 + 随机层高

核心思想：
- 空间换时间 → 多层索引加速查找
- 随机平衡 → 不用旋转
- L0 层是完整链表 → 天然支持范围查询

为什么 Redis 选 SkipList 而不是红黑树？
1. 实现简单
2. 天然支持范围查询（ZRANGE）
3. span 字段支持排名计算（ZRANK）
4. 红黑树做不到这么简洁

工程价值：
- Redis ZSet → SkipList + Dict
- LevelDB → 内存中的 memtable 用 SkipList
```

---

**上一篇：[08 树](../08-tree/README.md)**
**下一篇：[10 B+Tree](../10-bplus-tree/README.md)**
