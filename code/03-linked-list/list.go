package list

// ListNode 单链表节点
type ListNode struct {
	Val  int
	Next *ListNode
}

// LinkedList 单链表
type LinkedList struct {
	head *ListNode
	size int
}

// NewLinkedList 创建链表
func NewLinkedList() *LinkedList {
	return &LinkedList{head: nil, size: 0}
}

// Size 返回大小
func (l *LinkedList) Size() int {
	return l.size
}

// InsertHead 头部插入 O(1)
func (l *LinkedList) InsertHead(val int) {
	node := &ListNode{Val: val, Next: l.head}
	l.head = node
	l.size++
}

// InsertTail 尾部插入 O(n)
func (l *LinkedList) InsertTail(val int) {
	node := &ListNode{Val: val}
	if l.head == nil {
		l.head = node
		l.size++
		return
	}
	cur := l.head
	for cur.Next != nil {
		cur = cur.Next
	}
	cur.Next = node
	l.size++
}

// Delete 删除指定值的第一个节点 O(n)
func (l *LinkedList) Delete(val int) bool {
	if l.head == nil {
		return false
	}
	if l.head.Val == val {
		l.head = l.head.Next
		l.size--
		return true
	}
	cur := l.head
	for cur.Next != nil {
		if cur.Next.Val == val {
			cur.Next = cur.Next.Next
			l.size--
			return true
		}
		cur = cur.Next
	}
	return false
}

// Search 查找 O(n)
func (l *LinkedList) Search(val int) int {
	cur := l.head
	idx := 0
	for cur != nil {
		if cur.Val == val {
			return idx
		}
		cur = cur.Next
		idx++
	}
	return -1
}

// Get 获取第 k 个元素 O(k)
func (l *LinkedList) Get(k int) (int, bool) {
	if k < 0 || k >= l.size {
		return 0, false
	}
	cur := l.head
	for i := 0; i < k; i++ {
		cur = cur.Next
	}
	return cur.Val, true
}

// ToSlice 转为切片
func (l *LinkedList) ToSlice() []int {
	result := make([]int, 0, l.size)
	cur := l.head
	for cur != nil {
		result = append(result, cur.Val)
		cur = cur.Next
	}
	return result
}

// Reverse 反转链表
func (l *LinkedList) Reverse() {
	var prev *ListNode
	cur := l.head
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	l.head = prev
}

// HasCycle 环检测（Floyd 判圈算法）
func (l *LinkedList) HasCycle() bool {
	slow, fast := l.head, l.head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return true
		}
	}
	return false
}
