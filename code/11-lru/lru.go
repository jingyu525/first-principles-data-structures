package lru

// Node 双向链表节点
type Node struct {
	key, value int
	prev, next *Node
}

// LRUCache LRU 缓存
type LRUCache struct {
	capacity int
	cache    map[int]*Node
	head     *Node // 哨兵（最近使用）
	tail     *Node // 哨兵（最久未使用）
}

// NewLRUCache 创建 LRU
func NewLRUCache(capacity int) *LRUCache {
	l := &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*Node),
		head:     &Node{},
		tail:     &Node{},
	}
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

// Get 获取值 O(1)
func (l *LRUCache) Get(key int) int {
	if node, ok := l.cache[key]; ok {
		l.moveToHead(node)
		return node.value
	}
	return -1
}

// Put 插入值 O(1)
func (l *LRUCache) Put(key, value int) {
	if node, ok := l.cache[key]; ok {
		node.value = value
		l.moveToHead(node)
		return
	}

	newNode := &Node{key: key, value: value}
	l.cache[key] = newNode
	l.addToHead(newNode)

	if len(l.cache) > l.capacity {
		removed := l.removeTail()
		delete(l.cache, removed.key)
	}
}

// Size 当前缓存大小
func (l *LRUCache) Size() int {
	return len(l.cache)
}

// addToHead 添加到头部（最近使用）
func (l *LRUCache) addToHead(node *Node) {
	node.prev = l.head
	node.next = l.head.next
	l.head.next.prev = node
	l.head.next = node
}

// removeNode 删除节点
func (l *LRUCache) removeNode(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// moveToHead 移到头部
func (l *LRUCache) moveToHead(node *Node) {
	l.removeNode(node)
	l.addToHead(node)
}

// removeTail 删除尾部节点（最久未使用）
func (l *LRUCache) removeTail() *Node {
	node := l.tail.prev
	l.removeNode(node)
	return node
}
