# 队列练习题

## 1. 用两个栈实现队列

### 思路

一个栈负责入队（Push），一个栈负责出队（Pop）。
出队栈空时，把入队栈的所有元素倒入出队栈。

### 实现

```go
type MyQueue struct {
    in  []int // 入队栈
    out []int // 出队栈
}

func (q *MyQueue) Push(x int) {
    q.in = append(q.in, x)
}

func (q *MyQueue) Pop() int {
    q.ensure()
    v := q.out[len(q.out)-1]
    q.out = q.out[:len(q.out)-1]
    return v
}

func (q *MyQueue) Peek() int {
    q.ensure()
    return q.out[len(q.out)-1]
}

func (q *MyQueue) Empty() bool {
    return len(q.in) == 0 && len(q.out) == 0
}

func (q *MyQueue) ensure() {
    if len(q.out) == 0 {
        for len(q.in) > 0 {
            q.out = append(q.out, q.in[len(q.in)-1])
            q.in = q.in[:len(q.in)-1]
        }
    }
}
```

均摊 O(1)：每个元素最多入栈两次、出栈两次。

## 2. 滑动窗口最大值

给定数组和窗口大小 k，返回每个窗口的最大值。

### 思路

用双端队列（deque）维护窗口内的「可能最大值」索引。

### 实现

```go
func maxSlidingWindow(nums []int, k int) []int {
    result := make([]int, 0, len(nums)-k+1)
    deque := make([]int, 0) // 存储索引

    for i, v := range nums {
        // 移除窗口外的元素
        for len(deque) > 0 && deque[0] <= i-k {
            deque = deque[1:]
        }
        // 移除比当前元素小的元素（它们不会成为最大值）
        for len(deque) > 0 && nums[deque[len(deque)-1]] < v {
            deque = deque[:len(deque)-1]
        }
        deque = append(deque, i)

        if i >= k-1 {
            result = append(result, nums[deque[0]])
        }
    }
    return result
}
```
