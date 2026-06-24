package queue

// Queue 用切片实现的泛型队列
type Queue[T any] struct {
	items []T
}

// NewQueue 创建队列
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: make([]T, 0)}
}

// Enqueue 入队 O(1) 均摊
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue 出队 O(n)
// 注意：切片实现会 copy 剩余元素，生产环境建议用循环队列
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek 查看队首 O(1)
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

// IsEmpty 判空
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size 大小
func (q *Queue[T]) Size() int {
	return len(q.items)
}
