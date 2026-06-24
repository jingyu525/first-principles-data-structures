# 12 一致性哈希

## 第一性原理：为什么需要一致性哈希？

### 问题

```
分布式缓存：多台服务器

数据如何分布？
→ Hash(key) % N（N = 服务器数量）

问题来了：加一台服务器 → N 变了 → 几乎所有 key 都要重新分布！

Hash("user:100") % 3 = 1  →  Server 1
Hash("user:100") % 4 = 2  →  Server 2  （加了服务器，位置变了！）

缓存命中率 → 接近 0% → 缓存雪崩
```

**需要一种哈希方案：增加/减少节点时，只影响少量 key。**

这就是一致性哈希。

---

## 核心设计

### Hash Ring（哈希环）

```
把所有哈希值映射到一个环上：

                    0
              ● ← ← ← ← ●
            ↗             ↘
           ●               ●
   2^32-1  ●     哈希空间    ●
           ●    [0, 2^32)   ●
           ●               ●
            ↘             ↗
              ● → → → → ●

hash 范围：[0, 2^32 - 1]
```

### 节点与 Key 的映射

```
1. 服务器节点映射到环上
   Hash(ServerA) → 位置 A
   Hash(ServerB) → 位置 B
   Hash(ServerC) → 位置 C

2. Key 映射到环上
   Hash("key1") → 位置 K1

3. 从 K1 顺时针找到第一个节点 → 这个节点负责 key1

        0
   K1 ●
       ↘
        ● ServerA
         \
          ● K2
         /
   K3 ●
       ↗
        ● ServerB
```

```
Hash Ring：

        0
       ●
    /     \
   ●  N1   ●  N2
  /  K1 K2  \
 ●           ●
  \  K3     /
   ●  K4   ●  N3
    \     /
       ●

Key1 → 顺时针到 N1 → N1 负责
Key2 → 顺时针到 N1 → N1 负责
Key3 → 顺时针到 N3 → N3 负责
Key4 → 顺时针到 N3 → N3 负责
```

### 增加/删除节点的影响

```
增加 ServerD：
1. ServerD 落在环上某位置
2. 只有 ServerD 和下一个节点之间的 Key 需要迁移

受影响的 Key 比例：1/N（N = 节点数）
→ 而不是 Hash % N 方案中的 (N-1)/N
```

### 虚拟节点 (Virtual Nodes)

```
问题：节点太少 → 分布不均

解决：每个物理节点对应多个虚拟节点

ServerA → 150 个虚拟节点，均匀分布在环上
ServerB → 150 个虚拟节点，均匀分布在环上
ServerC → 150 个虚拟节点，均匀分布在环上

虚拟节点越多 → 分布越均匀
```

```
环上有虚拟节点后：

  ●A1  ●B1  ●C1  ●A2  ●B2  ●C2  ●A3  ●B3  ●C3 ...

Key 落到 A1、A2、A3 → 都由 ServerA 处理
→ 虚拟节点均匀分布 → 负载均衡
```

---

## Go 代码实现

```go
// code/12-consistent-hash/chash.go

package chash

import (
    "hash/crc32"
    "sort"
    "strconv"
)

// HashFunc hash 函数类型
type HashFunc func(data []byte) uint32

// ConsistentHash 一致性哈希
type ConsistentHash struct {
    hashFunc   HashFunc
    replicas   int            // 虚拟节点倍数
    hashRing   []int          // 排序的 hash 环
    hashMap    map[int]string // hash → 节点名
}

// New 创建一致性哈希
func New(replicas int, fn HashFunc) *ConsistentHash {
    if fn == nil {
        fn = crc32.ChecksumIEEE
    }
    return &ConsistentHash{
        hashFunc: fn,
        replicas: replicas,
        hashMap:  make(map[int]string),
    }
}

// AddNode 添加节点
func (ch *ConsistentHash) AddNode(nodes ...string) {
    for _, node := range nodes {
        for i := 0; i < ch.replicas; i++ {
            // 生成虚拟节点的 hash
            virtualKey := node + "#" + strconv.Itoa(i)
            hash := int(ch.hashFunc([]byte(virtualKey)))
            ch.hashRing = append(ch.hashRing, hash)
            ch.hashMap[hash] = node
        }
    }
    sort.Ints(ch.hashRing)
}

// RemoveNode 删除节点
func (ch *ConsistentHash) RemoveNode(node string) {
    for i := 0; i < ch.replicas; i++ {
        virtualKey := node + "#" + strconv.Itoa(i)
        hash := int(ch.hashFunc([]byte(virtualKey)))
        delete(ch.hashMap, hash)
    }
    // 重建 hashRing
    ch.hashRing = make([]int, 0, len(ch.hashMap))
    for h := range ch.hashMap {
        ch.hashRing = append(ch.hashRing, h)
    }
    sort.Ints(ch.hashRing)
}

// GetNode 根据 key 获取节点
func (ch *ConsistentHash) GetNode(key string) string {
    if len(ch.hashRing) == 0 {
        return ""
    }
    hash := int(ch.hashFunc([]byte(key)))
    // 二分查找：找到第一个 ≥ hash 的节点
    idx := sort.Search(len(ch.hashRing), func(i int) bool {
        return ch.hashRing[i] >= hash
    })
    // 环状：如果超出范围，回到第一个
    if idx == len(ch.hashRing) {
        idx = 0
    }
    return ch.hashMap[ch.hashRing[idx]]
}
```

---

## 工程案例

### 1. Redis Cluster

```
Redis Cluster 使用哈希槽（Hash Slot）：

总共 16384 个槽（不是一致性哈希环）
Hash(key) → CRC16(key) % 16384 → slot
每个节点负责一部分 slot

Server A: slots 0-5460
Server B: slots 5461-10922
Server C: slots 10923-16383

为什么用哈希槽而不是一致性哈希？
→ 槽的总数固定，方便管理和迁移
→ 每个节点负责连续的槽范围
→ MIGRATE 命令可以精确迁移槽
```

### 2. Kafka Partition 分配

```
Kafka 消费者组分配分区：

同一个 Consumer Group 中：
Partition 0 → Consumer 1
Partition 1 → Consumer 2
Partition 2 → Consumer 3

Consumer 增减：
→ 触发 rebalance
→ 只重新分配部分 partition
→ 类似一致性哈希的思想
```

### 3. Cassandra / DynamoDB

```
一致性哈希的经典应用：

节点均匀分布在环上
数据按一致性哈希分配到节点
每个数据复制到环上后续 N 个节点（N=3 即三个副本）
```

---

## 小结

```
一致性哈希 = Hash Ring + 虚拟节点 + 顺时针查找

解决的问题：
Hash % N → 增减节点影响几乎所有 key
一致性哈希 → 只影响 1/N 的 key

核心思想：
把节点和 key 都映射到同一个环上
key 顺时针找最近的节点

虚拟节点的作用：
解决节点数量少时的不均匀分布

工程价值：
- Redis Cluster → 哈希槽（变形）
- Kafka → Consumer 分区分配
- 分布式缓存 → 水平扩展
```

---

**上一篇：[11 LRU](../11-lru/README.md)**
**下一篇：[13 图](../13-graph/README.md)**
