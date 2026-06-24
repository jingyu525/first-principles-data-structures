package chash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// HashFunc hash 函数类型
type HashFunc func(data []byte) uint32

// ConsistentHash 一致性哈希
type ConsistentHash struct {
	hashFunc HashFunc
	replicas int            // 虚拟节点倍数
	hashRing []int          // 排序的 hash 环
	hashMap  map[int]string // hash → 节点名
}

// New 创建一致性哈希
func New(replicas int, fn HashFunc) *ConsistentHash {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}
	return &ConsistentHash{
		hashFunc: fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

// AddNodes 添加节点
func (ch *ConsistentHash) AddNodes(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < ch.replicas; i++ {
			virtualKey := node + "#" + strconv.Itoa(i)
			hash := int(ch.hashFunc([]byte(virtualKey)))
			ch.hashRing = append(ch.hashRing, hash)
			ch.hashMap[hash] = node
		}
	}
	sort.Ints(ch.hashRing)
}

// RemoveNode 删除节点
func (ch *ConsistentHash) RemoveNode(node string) {
	for i := 0; i < ch.replicas; i++ {
		virtualKey := node + "#" + strconv.Itoa(i)
		hash := int(ch.hashFunc([]byte(virtualKey)))
		delete(ch.hashMap, hash)
	}
	ch.hashRing = make([]int, 0, len(ch.hashMap))
	for h := range ch.hashMap {
		ch.hashRing = append(ch.hashRing, h)
	}
	sort.Ints(ch.hashRing)
}

// GetNode 根据 key 获取节点
func (ch *ConsistentHash) GetNode(key string) string {
	if len(ch.hashRing) == 0 {
		return ""
	}
	hash := int(ch.hashFunc([]byte(key)))
	idx := sort.Search(len(ch.hashRing), func(i int) bool {
		return ch.hashRing[i] >= hash
	})
	if idx == len(ch.hashRing) {
		idx = 0
	}
	return ch.hashMap[ch.hashRing[idx]]
}
