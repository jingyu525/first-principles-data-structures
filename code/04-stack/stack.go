package stack

// Stack 用切片实现的泛型栈
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
