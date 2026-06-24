package list

// DoubleNode 双向链表节点
type DoubleNode struct {
	Val  int
	Prev *DoubleNode
	Next *DoubleNode
}

// DoubleLinkedList 双向链表（带 head/tail 哨兵）
type DoubleLinkedList struct {
	head *DoubleNode
	tail *DoubleNode
	size int
}

// NewDoubleLinkedList 创建双向链表
func NewDoubleLinkedList() *DoubleLinkedList {
	return &DoubleLinkedList{head: nil, tail: nil, size: 0}
}

// Size 返回大小
func (l *DoubleLinkedList) Size() int {
	return l.size
}

// InsertHead 头部插入 O(1)
func (l *DoubleLinkedList) InsertHead(val int) {
	node := &DoubleNode{Val: val}
	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		node.Next = l.head
		l.head.Prev = node
		l.head = node
	}
	l.size++
}

// InsertTail O(1)
func (l *DoubleLinkedList) InsertTail(val int) {
	node := &DoubleNode{Val: val}
	if l.tail == nil {
		l.head = node
		l.tail = node
	} else {
		l.tail.Next = node
		node.Prev = l.tail
		l.tail = node
	}
	l.size++
}

// Delete O(1) — 有前驱指针
func (l *DoubleLinkedList) Delete(node *DoubleNode) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		l.head = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		l.tail = node.Prev
	}
	l.size--
}

// DeleteVal 按值删除 O(n)
func (l *DoubleLinkedList) DeleteVal(val int) bool {
	cur := l.head
	for cur != nil {
		if cur.Val == val {
			l.Delete(cur)
			return true
		}
		cur = cur.Next
	}
	return false
}

// ToSlice 正向遍历转为切片
func (l *DoubleLinkedList) ToSlice() []int {
	result := make([]int, 0, l.size)
	cur := l.head
	for cur != nil {
		result = append(result, cur.Val)
		cur = cur.Next
	}
	return result
}

// ToSliceReverse 反向遍历转为切片
func (l *DoubleLinkedList) ToSliceReverse() []int {
	result := make([]int, 0, l.size)
	cur := l.tail
	for cur != nil {
		result = append(result, cur.Val)
		cur = cur.Prev
	}
	return result
}
