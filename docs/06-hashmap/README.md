# 06 HashMap

> **这是整个项目中最重要的一章**
>
> HashMap 是所有「快速查找」的基础，理解它就理解了 Go Map、Redis Hash、数据库索引的底层逻辑。

## 第一性原理：为什么需要 HashMap？

### 问题

```
数组 → 按索引访问 O(1)，但索引必须是数字
链表 → 按顺序查找 O(n)

如何实现「按任意 Key 访问，O(1)」？

→ 把 Key 转换成数字
→ 用这个数字作为数组索引
→ 直接定位
```

```
Key → Hash函数 → 数字 → 数组索引 → O(1)
```

这就是 HashMap。

---

## 核心三要素

### 1. Hash 函数

```
把任意输入映射到固定范围的数字

要求：
- 确定性：同样的输入 → 同样的输出
- 均匀性：不同的输入 → 均匀分布
- 高效性：计算要快
```

```go
// Go 对 string 的 hash 函数（runtime/hash64.go）
// 使用 AES 指令加速

// 简化的 hash 函数示例
func simpleHash(key string, bucketCount int) int {
    h := 0
    for _, ch := range key {
        h = h*31 + int(ch)
    }
    return h % bucketCount
}
```

### 2. 冲突解决

```
不同的 Key 可能产生相同的 Hash → 冲突

两种解决方式：

拉链法 (Chaining)：
  Bucket 0 → [k1,v1] → [k5,v5]
  Bucket 1 → [k2,v2]
  Bucket 2 → [k3,v3] → [k6,v6]
  ...

开放寻址法 (Open Addressing)：
  冲突了 → 找下一个空位置
  Bucket[i] → Bucket[i+1] → Bucket[i+2] → ...
```

| 方法 | 优点 | 缺点 | 使用场景 |
|------|------|------|---------|
| 拉链法 | 简单、装填因子可 >1 | 额外指针开销 | Go Map、Java HashMap |
| 开放寻址 | 内存紧凑、缓存友好 | 装填因子 <1、删除复杂 | Python dict、Redis Hash(小) |

### 3. 扩容 (Rehashing)

```
装填因子 (Load Factor) = 元素数量 / 桶数量

当装填因子超过阈值 → 扩容

Go Map 扩容阈值：6.5（即每个桶平均 6.5 个元素）
```

```
扩容步骤：
1. 分配新的、更大的桶数组（通常是 2 倍）
2. 重新 Hash 所有元素
3. 迁移到新桶

代价：O(n)，但均摊到每次插入 → O(1)
```

---

## 时间复杂度

| 操作 | 平均 | 最坏 |
|------|------|------|
| 查找 | O(1) | O(n) |
| 插入 | O(1) | O(n) |
| 删除 | O(1) | O(n) |

最坏情况：所有 key 都冲突到同一个桶 → 退化成链表。

**工程中如何避免最坏？**
- 好的 hash 函数
- 链表太长时转为红黑树（Java HashMap >8）
- 适当的扩容策略

---

## Go Map 源码分析

```go
// Go Map 底层结构（runtime/map.go 简化版）

// 桶的结构
type bmap struct {
    // tophash: 存储 key hash 的高 8 位，用于快速比较
    tophash [8]uint8
    // keys: 存储 8 个 key
    // values: 存储 8 个 value
    // overflow: 溢出桶指针
}

// Map 底层结构
type hmap struct {
    count     int    // 元素数量
    flags     uint8
    B         uint8  // 桶数量的对数（桶数 = 2^B）
    noverflow uint16 // 溢出桶数量
    hash0     uint32 // hash seed（随机，防止 hash 碰撞攻击）
    buckets    unsafe.Pointer // 桶数组
    oldbuckets unsafe.Pointer // 扩容时的旧桶
    nevacuate  uintptr        // 扩容进度
}
```

### Go Map 的设计亮点

```
1. 每个桶存 8 个 key-value（不是 1 个）
   → 减少指针开销
   → 提高缓存命中率

2. tophash 快速过滤
   → 先比较 hash 高 8 位
   → 不匹配直接跳过
   → 不用比较完整的 key

3. 渐进式扩容
   → 不一次性迁移所有数据
   → 每次读写时迁移一部分
   → 避免 STW（Stop The World）

4. 随机 hash seed
   → 每次启动生成随机 seed
   → 防止 hash flooding 攻击
```

```
Go Map 查找流程：

1. 计算 hash(key)
2. 定位桶：buckets[hash & (2^B - 1)]
3. 比较 tophash（hash 高 8 位）
4. tophash 匹配 → 比较完整 key
5. 桶内 8 个都查完 → 查 overflow 桶
```

---

## Redis Hash 分析

```
Redis Hash 有两种内部编码：

1. ziplist（压缩列表）
   条件：元素少且小（<512 个，每个 < 64 字节）
   本质：连续内存上的紧凑数组
   操作：O(n)，但 n 小时比 hash table 更快

2. hashtable（哈希表）
   条件：元素多或大
   本质：拉链法 hash table
   操作：O(1)
```

```c
// Redis dict 结构（src/dict.h 简化）
typedef struct dict {
    dictType *type;
    dictht ht[2];       // 两个 hash table（用于渐进式 rehash）
    long rehashidx;     // rehash 进度，-1 表示不在 rehash
} dict;

typedef struct dictht {
    dictEntry **table;  // 桶数组
    unsigned long size; // 桶数量
    unsigned long sizemask; // size - 1
    unsigned long used; // 元素数量
} dictht;

typedef struct dictEntry {
    void *key;
    union { void *val; uint64_t u64; int64_t s64; } v;
    struct dictEntry *next; // 拉链法
} dictEntry;
```

```
Redis 渐进式 rehash：

ht[0] → 旧表
ht[1] → 新表（2 倍大小）
rehashidx → 当前迁移到哪个桶

每次操作时迁移一个桶
→ 查找：先查 ht[0]，再查 ht[1]
→ 插入：只插入 ht[1]
→ 不阻塞服务
```

---

## 练习题

### Two Sum

```go
func twoSum(nums []int, target int) []int {
    m := make(map[int]int)
    for i, v := range nums {
        if j, ok := m[target-v]; ok {
            return []int{j, i}
        }
        m[v] = i
    }
    return nil
}
```

时间复杂度：O(n)，只用一次遍历。
不用 HashMap 则需要 O(n²)。

---

## 小结

```
HashMap = Hash函数 + 冲突解决 + 扩容

三个关键问题：
1. 怎么把 Key 变成数字？→ Hash 函数
2. 冲突了怎么办？→ 拉链法/开放寻址
3. 满了怎么办？→ 扩容 + rehash

工程价值：
- Go Map → 渐进式扩容、tophash 优化
- Redis Dict → 双表渐进式 rehash
- 数据库索引 → Hash 索引（等值查询快，不支持范围）

关键洞察：
O(1) 不是魔法，是空间换时间。
预分配的空间 + 好的 hash 函数 → 接近 O(1)。
```

---

**上一篇：[05 队列](../05-queue/README.md)**
**下一篇：[07 Heap](../07-heap/README.md)**
