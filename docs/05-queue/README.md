# 05 队列

## 第一性原理：为什么需要队列？

### 问题

```
栈是 LIFO → 后进先出

但生活中还有一种模式：

排队 → 先到先服务

系统场景：
- 消息处理 → 先收到的消息先处理
- 请求处理 → 先到的请求先响应
- 任务调度 → 先提交的任务先执行
```

**规律：先进先出（FIFO）**

这就是队列。

---

## 核心特性

### FIFO (First In, First Out)

```
Enqueue(1) →  [1]
Enqueue(2) →  [1, 2]
Enqueue(3) →  [1, 2, 3]
Dequeue()  →  [2, 3]   返回 1
Dequeue()  →  [3]      返回 2
Dequeue()  →  []       返回 3
```

---

## 时间复杂度

| 操作 | 数组队列 | 链表队列 | 循环队列 |
|------|---------|---------|---------|
| Enqueue | O(1)/O(n)* | O(1) | O(1) |
| Dequeue | O(n)** | O(1) | O(1) |
| Peek | O(1) | O(1) | O(1) |

\* 数组满时需要扩容
\*\* 普通数组 Dequeue 需要移动所有元素

---

## Go 代码实现

### 普通队列（切片）

```go
// code/05-queue/queue.go

package queue

type Queue[T any] struct {
    items []T
}

func NewQueue[T any]() *Queue[T] {
    return &Queue[T]{items: make([]T, 0)}
}

func (q *Queue[T]) Enqueue(item T) {
    q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
    if len(q.items) == 0 {
        var zero T
        return zero, false
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item, true
}
```

> 注意：切片实现的 Dequeue 是 O(n)，因为 `items[1:]` 需要拷贝。

### 循环队列

```go
// code/05-queue/circular_queue.go

// CircularQueue 循环队列 — 真正的 O(1) Dequeue
type CircularQueue[T any] struct {
    items []T
    head  int
    tail  int
    size  int
    cap   int
}

func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
    return &CircularQueue[T]{
        items: make([]T, capacity),
        cap:   capacity,
    }
}

func (q *CircularQueue[T]) Enqueue(item T) bool {
    if q.size == q.cap {
        return false // 已满
    }
    q.items[q.tail] = item
    q.tail = (q.tail + 1) % q.cap
    q.size++
    return true
}

func (q *CircularQueue[T]) Dequeue() (T, bool) {
    if q.size == 0 {
        var zero T
        return zero, false
    }
    item := q.items[q.head]
    q.head = (q.head + 1) % q.cap
    q.size--
    return item, true
}
```

```
循环队列示意图 (cap=4):

Enqueue(1,2,3)     Dequeue() 后    Enqueue(4)
┌─┬─┬─┬─┐        ┌─┬─┬─┬─┐        ┌─┬─┬─┬─┐
│1│2│3│ │        │ │2│3│ │        │4│2│3│ │
└─┴─┴─┴─┘        └─┴─┴─┴─┘        └─┴─┴─┴─┘
 h   t            h t              t h
```

### 阻塞队列（Go Channel）

```go
// Go Channel 就是线程安全的阻塞队列

ch := make(chan int, 10) // 容量为 10 的有缓冲 channel

ch <- 1    // Enqueue（队列满时阻塞）
v := <-ch  // Dequeue（队列空时阻塞）

// 非阻塞操作
select {
case ch <- v:
    // 发送成功
default:
    // 队列满
}
```

---

## 工程案例

### 1. Go Channel

```
Go 的并发哲学：通过通信共享内存

Channel = 线程安全的阻塞队列

底层结构 (runtime/chan.go)：
type hchan struct {
    qcount   uint           // 队列中元素数量
    dataqsiz uint           // 循环队列大小
    buf      unsafe.Pointer // 指向循环队列的指针
    elemsize uint16
    closed   uint32
    sendx    uint   // 发送索引
    recvx    uint   // 接收索引
    recvq    waitq  // 等待接收的 goroutine 队列
    sendq    waitq  // 等待发送的 goroutine 队列
    lock     mutex
}
```

```
Channel 就是：循环队列 + 等待队列 + 锁

当缓冲区满时：发送者加入 sendq，阻塞等待
当缓冲区空时：接收者加入 recvq，阻塞等待
```

### 2. Kafka

```
Kafka = 分布式消息队列

核心抽象：队列

但 Kafka 的队列和普通队列不同：

普通队列：Dequeue 后消息消失
Kafka：消费后消息保留（基于 offset）

Topic/Partition 的本质：
→ 一个有序的、持久化的消息序列
→ 可以看作磁盘上的循环队列
```

```
分区内的消息：
offset:   0    1    2    3    4    5
        ┌────┬────┬────┬────┬────┬────┐
        │msg0│msg1│msg2│msg3│msg4│msg5│
        └────┴────┴────┴────┴────┴────┘
                  ↑              ↑
               consumer       producer
```

### 3. RabbitMQ / 其他消息队列

```
所有消息队列的共性：

1. 队列抽象 → FIFO
2. 持久化 → 磁盘存储
3. 消费确认 → 可靠消费
4. 路由 → 消息分发策略

本质都是：在队列的基础上添加可靠性、持久性、分布式能力
```

---

## 队列变体

| 类型 | 特点 | 典型应用 |
|------|------|---------|
| 普通队列 | FIFO | 基本消息传递 |
| 循环队列 | 固定大小环形缓冲区 | Go Channel、Kafka 缓冲区 |
| 阻塞队列 | 空阻塞、满阻塞 | Go Channel |
| 优先队列 | 按优先级出队 | K8S Scheduler |
| 延迟队列 | 延迟消费 | 定时任务 |

---

## 练习题

### 用两个栈实现队列

详见 [`exercises/05-queue/`](../../exercises/05-queue/)

---

## 小结

```
队列 = FIFO + 两端操作

核心变体：
- 循环队列 → 解决 Dequeue O(n) 的问题
- 阻塞队列 → 解决生产者-消费者同步问题

工程价值：
- Go Channel → 并发通信的基石
- Kafka → 分布式消息队列
- 消息队列 → 系统解耦
```

---

**上一篇：[04 栈](../04-stack/README.md)**
**下一篇：[06 HashMap](../06-hashmap/README.md)**
