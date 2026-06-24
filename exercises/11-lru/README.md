# LRU 练习题

## LRU 缓存

实现 LRUCache 类：

- `LRUCache(capacity int)` 初始化
- `Get(key int) int` 获取值，不存在返回 -1
- `Put(key, value int)` 插入值

要求所有操作 O(1)。

### 验证测试

```go
func TestLRU(t *testing.T) {
    lru := NewLRUCache(2)

    lru.Put(1, 1)          // 缓存: {1=1}
    lru.Put(2, 2)          // 缓存: {1=1, 2=2}
    assert(t, lru.Get(1), 1) // 返回 1，缓存: {2=2, 1=1}（1 移到头部）
    lru.Put(3, 3)          // 淘汰 key=2，缓存: {1=1, 3=3}
    assert(t, lru.Get(2), -1) // 返回 -1
    lru.Put(4, 4)          // 淘汰 key=1，缓存: {3=3, 4=4}
    assert(t, lru.Get(1), -1)
    assert(t, lru.Get(3), 3)
    assert(t, lru.Get(4), 4)
}
```

### 完整实现

见 [`code/11-lru/lru.go`](../../code/11-lru/lru.go)
