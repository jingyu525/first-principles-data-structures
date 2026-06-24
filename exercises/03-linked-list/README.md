# 链表练习题

## 1. 反转链表

### 思路

用三个指针：prev, cur, next。逐个反转指针方向。

### 实现

```go
type ListNode struct {
    Val  int
    Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
    var prev *ListNode
    cur := head
    for cur != nil {
        next := cur.Next
        cur.Next = prev
        prev = cur
        cur = next
    }
    return prev
}
```

## 2. 链表环检测

### 思路

快慢指针。慢指针每次走 1 步，快指针每次走 2 步。
如果有环，快慢指针会相遇。

### 实现

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

## 3. 链表环入口

### 思路

找到相遇点后，head 和 slow 同时各走一步，相遇处就是环入口。（数学可证）

### 实现

```go
func detectCycle(head *ListNode) *ListNode {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            // 有环，找入口
            for head != slow {
                head = head.Next
                slow = slow.Next
            }
            return head
        }
    }
    return nil
}
```
