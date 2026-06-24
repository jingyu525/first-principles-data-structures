# 04 栈

## 第一性原理：为什么需要栈？

### 问题

```
数组/链表 → 可以在任意位置操作

但如果我只想在"一端"操作呢？

场景1：函数调用 → 最后调用的函数最先返回
场景2：撤销操作 → 最后做的操作最先撤销
场景3：括号匹配 → 最后出现的左括号匹配最先出现的右括号
```

**规律：后进先出（LIFO）**

```
不需要在中间操作 → 可以简化数据结构 → 更高效
```

这就是栈。

---

## 核心特性

### LIFO (Last In, First Out)

```
Push(1) →  [1]
Push(2) →  [1, 2]
Push(3) →  [1, 2, 3]
Pop()   →  [1, 2]  返回 3
Pop()   →  [1]     返回 2
Pop()   →  []      返回 1
```

### 只在一端操作

```
栈顶(top) → Push/Pop/Peek
栈底(bottom) → 不操作
```

---

## 时间复杂度

| 操作 | 复杂度 |
|------|--------|
| Push | O(1) |
| Pop | O(1) |
| Peek | O(1) |
| 查找 | O(n) |

---

## Go 代码实现

```go
// code/04-stack/stack.go

package stack

// Stack 用切片实现的栈
type Stack[T any] struct {
    items []T
}

// NewStack 创建栈
func NewStack[T any]() *Stack[T] {
    return &Stack[T]{items: make([]T, 0)}
}

// Push 入栈 O(1)
func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

// Pop 出栈 O(1)
func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// Peek 查看栈顶 O(1)
func (s *Stack[T]) Peek() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

// IsEmpty 判空
func (s *Stack[T]) IsEmpty() bool {
    return len(s.items) == 0
}

// Size 大小
func (s *Stack[T]) Size() int {
    return len(s.items)
}
```

---

## 工程案例

### 1. Go Goroutine Stack

```
每个 Goroutine 都有一个栈

初始大小：2KB（Go 1.4+）
按需增长和收缩（分段栈 → 连续栈）

func A() {
    B()  // B 的栈帧压入 A 上方
}

栈布局：
┌─────────────┐ ← 高地址
│  B 的栈帧    │
├─────────────┤
│  A 的栈帧    │
├─────────────┤
│  局部变量    │
├─────────────┤ ← 低地址
│  返回地址    │
└─────────────┘
        栈增长方向 ↓
```

```go
// 递归 → 每次调用压栈
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
    // factorial(5)
    //   → factorial(4)
    //     → factorial(3)
    //       → factorial(2)
    //         → factorial(1) = 1
    //       ← 2*1 = 2
    //     ← 3*2 = 6
    //   ← 4*6 = 24
    // ← 5*24 = 120
}
```

### 2. 表达式计算

```
中缀表达式：3 + 4 × 2
后缀表达式（逆波兰）：3 4 2 × +

计算后缀表达式（用栈）：

读到 3 → Push(3)    栈：[3]
读到 4 → Push(4)    栈：[3, 4]
读到 2 → Push(2)    栈：[3, 4, 2]
读到 × → Pop 2, Pop 4, 计算 4×2=8, Push(8)  栈：[3, 8]
读到 + → Pop 8, Pop 3, 计算 3+8=11, Push(11) 栈：[11]

结果：11
```

### 3. 括号匹配

```go
func isValid(s string) bool {
    stack := NewStack[rune]()
    pairs := map[rune]rune{')': '(', '}': '{', ']': '['}

    for _, ch := range s {
        switch ch {
        case '(', '{', '[':
            stack.Push(ch)
        case ')', '}', ']':
            top, ok := stack.Pop()
            if !ok || top != pairs[ch] {
                return false
            }
        }
    }
    return stack.IsEmpty()
}
```

---

## 练习题

### 括号匹配

详见 [`exercises/04-stack/`](../../exercises/04-stack/)

---

## 小结

```
栈 = LIFO + 一端操作

本质上是对数组/链表的操作约束：
→ 只允许在一端 Push/Pop

这个约束看起来是限制，实际是简化
→ 不需要考虑中间位置
→ 实现更简单
→ 适用特定场景

工程价值：
- Go Goroutine Stack → 函数调用
- 逆波兰表达式 → 表达式计算
- 括号匹配 → 编译器解析

关键洞察：
约束不是限制，是简化。
```

---

**上一篇：[03 链表](../03-linked-list/README.md)**
**下一篇：[05 队列](../05-queue/README.md)**
