# 14 从数据结构看系统设计

> **这是整个项目的高潮。**
>
> 用前面学过的所有数据结构，理解真实系统的设计。

---

## 核心论点

```
优秀的系统 = 数据结构 × 系统约束

不是凭空设计，而是在合适的地方用合适的数据结构。
```

---

## Redis

### 整体架构

```
Redis = HashTable + SkipList + QuickList + 字符串 + 集合

数据类型 → 底层结构：

String   → SDS (Simple Dynamic String)
List     → QuickList (双向链表 + 压缩列表)
Set      → HashTable / IntSet
Hash     → HashTable / ZipList
ZSet     → Dict(HashTable) + SkipList
Stream   → Radix Tree
```


### 为什么 Redis 这么快？

```
1. 内存存储 → 主要数据结构都在内存中
2. 数据结构选择精准：
   - HashTable → O(1) 查找
   - SkipList → O(log n) 排序 + 范围查询
3. 单线程 + IO 多路复用 → 无锁竞争
4. 精心优化的编码：
   - 小数据用 ZipList（紧凑内存）
   - 大数据自动切换 HashTable/SkipList
```

### ZSet 为什么用 Dict + SkipList？

```
Dict（HashTable）：
- member → score 的映射 → O(1)
- ZSCORE key member → 直接查 Dict

SkipList：
- 按 score 排序 → O(log n)
- ZRANGE key 0 10 → 查 SkipList
- ZRANK key member → 利用 span 字段

两个结构互补：
- Dict 解决「按成员查找分数」
- SkipList 解决「按分数排序/排名」
```

---

## MySQL

### 整体架构

```
MySQL InnoDB = B+Tree + BufferPool + RedoLog + UndoLog

索引层：
  聚簇索引（主键） → B+Tree（叶子存完整行）
  二级索引 → B+Tree（叶子存主键值）

缓存层：
  BufferPool → LRU 变体（管理磁盘页缓存）

事务层：
  UndoLog → 多版本并发控制 (MVCC)
  RedoLog → 崩溃恢复
```

### 为什么 InnoDB 用 B+Tree 而不是红黑树/跳表/Hash？

| 结构 | 为什么不适合 MySQL | 适合什么 |
|------|-------------------|---------|
| 红黑树 | 每个节点一个值，树高 → 磁盘 IO 多 | 内存索引 |
| SkipList | 指针多 → 空间大；内存结构，不适合页式存储 | Redis ZSet |
| Hash | 不支持范围查询 → `WHERE id > 100` 无法高效 | 等值查询 |
| **B+Tree** | 节点多值 → 树矮 → IO 少；叶子链表 → 范围查询快 | **磁盘索引** |

### SELECT 全过程

```sql
SELECT * FROM users WHERE name = 'Alice';
```

```
1. 查询优化器 → 选择 idx_name 索引
2. idx_name 的 B+Tree → 二分查找 'Alice' → O(log n)
   找到叶子节点：Alice → pk=25
3. 主键 B+Tree → 用 pk=25 查聚簇索引 → O(log n)
   找到叶子节点：id=25 的完整行数据
4. 返回结果

两次 B+Tree 查找（回表）
总 IO：4-6 次（每棵树 2-3 层）
```

---

## Kafka

### 整体架构

```
Kafka = Queue + Partition + Log Segment

队列抽象：
  Topic → 消息的逻辑队列
  Partition → Topic 的物理分片

存储结构：
  Partition → 有序的 Log Segment 文件
  Segment → 顺序写文件（append-only）
```

### 为什么 Kafka 这么快？

```
1. 顺序 IO
   partition 内是 append-only 日志
   → 不修改中间数据
   → 利用磁盘顺序写（600MB/s+）

2. Page Cache（操作系统页缓存）
   → 依赖 OS 而不是自建缓存
   → 读写走 Page Cache

3. Zero Copy（零拷贝）
   → sendfile() 系统调用
   → 数据从磁盘 → Page Cache → 网卡，不经过用户态

4. 批量处理
   → Producer 批量发送
   → Consumer 批量拉取
   → 减少网络开销
```

```
Partition 的内部结构：

┌───────────────── Partition ───────────────────┐
│                                                │
│  Segment 0       Segment 1       Segment 2     │
│  ┌────────┐     ┌────────┐     ┌────────┐     │
│  │ offset │     │ offset │     │ offset │     │
│  │   0    │     │ 1000   │     │ 2000   │     │
│  │  ...   │     │  ...   │     │  ...   │     │
│  │  999   │     │ 1999   │     │ 2999   │     │
│  └────────┘     └────────┘     └────────┘     │
│                                                │
│  写入方向 →                                     │
│                                                │
└────────────────────────────────────────────────┘

Consumer Offset: pointer 指向已消费位置
→ 本质是数组索引的概念
```

---

## Go Runtime

### 整体架构

```
Go Runtime = Queue + Heap + Map + Stack

Channel   → 循环队列 + 等待队列
Timer     → 四叉堆（按时间排序）
Scheduler → 优先队列（runqueue）
Map       → HashTable（渐进式扩容）
Goroutine → 分段栈（stack）
GC        → 三色标记（图遍历） + 写屏障
```

### Channel

```
Channel = 循环队列（hchan.buf）+ 等待队列（sendq/recvq）

ch := make(chan int, 10)  → 容量 10 的循环队列

发送 (ch <- v)：
  buf 不满 → 写入循环队列
  buf 满 → goroutine 加入 sendq，阻塞

接收 (v := <-ch)：
  buf 不空 → 从循环队列读取
  buf 空 → goroutine 加入 recvq，阻塞
```

### Timer

```
Go Timer 底层是四叉堆：

     [100ms]       ← 最近的定时器
   /  |  |  \
 ... ... ... ...

runtime 的调度循环：
1. 取堆顶定时器时间 t
2. 计算等待时间 t - now
3. 如果 ≤ 0 → 触发，继续取下一个
4. 如果 > 0 → 调度其他 goroutine 或休眠
```

### Goroutine 调度

```
GMP 模型：

G (Goroutine) → 用户态线程
M (Machine)   → 系统线程
P (Processor) → 逻辑处理器（调度上下文）

每个 P 有一个本地 runqueue（队列）：
  P1: [G1] [G2] [G3]    ← 循环队列
  P2: [G4] [G5]
  P3: [G6] [G7] [G8] [G9]

工作窃取 (Work Stealing)：
  P 的队列为空 → 从其他 P 偷一半 goroutine

这也是队列的应用
```

### GC 三色标记

```
三色标记 = 图的 BFS 遍历

初始：所有对象白色
从根对象开始 BFS：
  1. 访问到的对象 → 灰色
  2. 对象的子对象访问完 → 黑色
  3. 最终白色对象 → 垃圾，回收

栈 → 根对象
队列 → BFS 的待处理节点（灰色集合）
```

---

## Kubernetes

### 整体架构

```
Kubernetes = Heap + Queue + Graph

Scheduler    → 优先队列（Pod 按优先级调度）
Controller   → 队列（事件队列 + 工作队列）
API Server   → 图（资源对象之间的关系）
etcd         → B+Tree（MVCC 存储）
```

### K8S Scheduler

```
Pod 调度 = 优先队列 + 过滤 + 打分

1. Pod 进入调度队列（优先队列）
   高优先级 Pod → 先出队

2. 过滤阶段
   排除不满足条件的 Node：
   - 资源不足
   - 亲和性不满足
   - 污点不匹配

3. 打分阶段
   对剩余 Node 打分 → 选最高分
   - 资源均衡
   - 亲和性加分
   - 镜像本地化

队列 + 堆 → 调度基础
```

### K8S 资源依赖图

```
Deployment
   │
   ├──→ ReplicaSet
   │       │
   │       └──→ Pod ──→ PVC ──→ PV
   │               │
   │               ├──→ ConfigMap
   │               ├──→ Secret
   │               └──→ Service
   │
   └──→ HorizontalPodAutoscaler

这是有向图
- 控制器通过 OwnerReference 建立关系
- 删除 Deployment → 级联删除所有下游
- 依赖分析 → 拓扑排序
```

### Controller 模式

```
Controller = 事件队列 + 工作队列

1. Watch API Server → 事件入队
2. Worker 从队列取事件 → 处理
3. 处理：当前状态 → 期望状态

本质：生产者-消费者模式（队列）

多个 Controller 并发：
- Deployment Controller
- ReplicaSet Controller
- Node Controller
- ...

每个 Controller 有自己的工作队列
```

---

## 最终地图

```
                数组
                  ↓
                链表
                  ↓
         ┌────栈────队列────┐
         ↓                  ↓
      HashMap              Heap
         ↓                  ↓
         └────→ 树 ←────────┘
                  ↓
                图
    ─────────────────────────
         ↓        ↓        ↓
       LRU    SkipList  B+Tree
         ↓        ↓        ↓
    ─────────────────────────
         ↓        ↓        ↓
      Redis    MySQL    Kafka
         ↓        ↓        ↓
         └────→ Go Runtime ←┘
                  ↓
            Kubernetes
```

---

## 核心洞察

### 1. 系统是数据结构的组合，不是单一结构

```
Redis ZSet = Dict + SkipList
MySQL = B+Tree + BufferPool(LRU) + RedoLog
Kafka = Queue + Segment(Array) + Index(Sparse)
```

### 2. 存储介质决定数据结构

```
内存 → 数组、链表、HashTable（随机访问快）
磁盘 → B+Tree、LSM-Tree（顺序读写快）
SSD → 介于两者之间
```

### 3. 约束不是限制，是简化

```
栈 = 数组/链表 + "只在一端操作"的约束
  → 简化了设计，获得了 LIFO 的语义

队列 = 数组/链表 + "两端操作"的约束
  → 简化了设计，获得了 FIFO 的语义
```

### 4. 工程 = 理论 + 权衡

```
理论上：SkipList ≈ 红黑树（都是 O(log n)）
工程上：Redis 选 SkipList → 实现简单 + 范围查询友好

理论上：红黑树可以替代 B+Tree
工程上：磁盘 IO 是瓶颈 → B+Tree 维度设计更适合
```

---

**下一篇：回到起点，重新审视整个知识体系**
