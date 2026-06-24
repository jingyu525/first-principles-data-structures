package queue

// CircularQueue O(1) Dequeue 的循环队列
type CircularQueue[T any] struct {
	items []T
	head  int
	tail  int
	size  int
	cap   int
}

// NewCircularQueue 创建循环队列
func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
	return &CircularQueue[T]{
		items: make([]T, capacity),
		cap:   capacity,
	}
}

// Enqueue 入队 O(1)
func (q *CircularQueue[T]) Enqueue(item T) bool {
	if q.size == q.cap {
		return false
	}
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.cap
	q.size++
	return true
}

// Dequeue 出队 O(1)
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

// Peek 查看队首 O(1)
func (q *CircularQueue[T]) Peek() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}
	return q.items[q.head], true
}

// IsFull 判满
func (q *CircularQueue[T]) IsFull() bool {
	return q.size == q.cap
}

// IsEmpty 判空
func (q *CircularQueue[T]) IsEmpty() bool {
	return q.size == 0
}

// Size 大小
func (q *CircularQueue[T]) Size() int {
	return q.size
}
