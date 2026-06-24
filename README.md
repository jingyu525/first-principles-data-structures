# 从第一性原理理解数据结构与算法

> 不是刷题，而是理解 Go Runtime、Redis、MySQL、Kafka、Kubernetes 背后的设计思想。

---

## 项目地图

```
数组
 ↓
链表
 ↓
栈
 ↓
队列
 ↓
HashMap
 ↓
Heap
 ↓
Tree
 ↓
Graph
─────────────────
 ↓
LRU
 ↓
SkipList
 ↓
B+Tree
 ↓
HashRing
─────────────────
 ↓
Redis
 ↓
MySQL
 ↓
Kafka
 ↓
Go Runtime
 ↓
Kubernetes
```

---

## 为什么做这个项目

市面上有无数 LeetCode 题解、算法导论笔记、数据结构教程。

但缺少的是：**从工程视角理解数据结构**。

本项目的目标是：

| 传统教程 | 本项目 |
|---------|-------|
| 学数组为了刷题 | 学数组为了理解 Kafka 缓冲区 |
| 学链表为了通过面试 | 学链表为了理解 Redis QuickList |
| 学跳表为了应付笔试 | 学跳表为了看懂 Redis ZSet 源码 |
| 学 B+Tree 背定义 | 学 B+Tree 理解 MySQL InnoDB 为什么快 |
| 学图算法背模板 | 学图算法分析 Kubernetes 服务依赖 |

---

## 目录

### 第一部分：第一性原理
- [01 计算机的本质：存储与处理](docs/01-first-principle/README.md)

### 第二部分：基础数据结构
- [02 数组](docs/02-array/README.md)
- [03 链表](docs/03-linked-list/README.md)
- [04 栈](docs/04-stack/README.md)
- [05 队列](docs/05-queue/README.md)

### 第三部分：核心结构
- [06 HashMap](docs/06-hashmap/README.md)
- [07 Heap（堆）](docs/07-heap/README.md)

### 第四部分：树
- [08 树](docs/08-tree/README.md)
- [09 SkipList（跳表）](docs/09-skiplist/README.md)
- [10 B+Tree](docs/10-bplus-tree/README.md)

### 第五部分：工程组合结构
- [11 LRU 缓存](docs/11-lru/README.md)
- [12 一致性哈希](docs/12-consistent-hash/README.md)

### 第六部分：图
- [13 图](docs/13-graph/README.md)

### 第七部分：系统设计
- [14 从数据结构看系统设计](docs/14-system-design/README.md)

---

## 代码

所有 Go 代码示例在 [`code/`](code/) 目录下。

---

## 练习题

在 [`exercises/`](exercises/) 目录下，每题包含题目描述、思路分析和 Go 实现。

---

## 里程碑

| 版本 | 时间 | 内容 |
|------|------|------|
| V1 | 30天 | 数组、链表、栈、队列、HashMap、Heap |
| V2 | 60天 | Tree、SkipList、B+Tree、LRU、HashRing |
| V3 | 90天 | Redis 源码、MySQL 源码、Go Runtime |

---

## License

MIT
