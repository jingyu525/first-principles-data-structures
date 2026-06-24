package array

// DynamicArray 动态数组（Go Slice 的简化实现）
type DynamicArray struct {
	data     []int
	size     int // 当前元素数量
	capacity int // 容量
}

// NewDynamicArray 创建动态数组
func NewDynamicArray(capacity int) *DynamicArray {
	return &DynamicArray{
		data:     make([]int, capacity),
		size:     0,
		capacity: capacity,
	}
}

// Size 返回元素数量
func (a *DynamicArray) Size() int {
	return a.size
}

// Get 随机访问 O(1)
func (a *DynamicArray) Get(index int) (int, bool) {
	if index < 0 || index >= a.size {
		return 0, false
	}
	return a.data[index], true
}

// Set 设置值 O(1)
func (a *DynamicArray) Set(index, val int) bool {
	if index < 0 || index >= a.size {
		return false
	}
	a.data[index] = val
	return true
}

// Append 尾部插入 O(1) 均摊
func (a *DynamicArray) Append(val int) {
	if a.size == a.capacity {
		a.resize(a.capacity * 2)
	}
	a.data[a.size] = val
	a.size++
}

// Insert 中间插入 O(n)
func (a *DynamicArray) Insert(index, val int) bool {
	if index < 0 || index > a.size {
		return false
	}
	if a.size == a.capacity {
		a.resize(a.capacity * 2)
	}
	// 后移
	for i := a.size; i > index; i-- {
		a.data[i] = a.data[i-1]
	}
	a.data[index] = val
	a.size++
	return true
}

// Delete 删除 O(n)
func (a *DynamicArray) Delete(index int) (int, bool) {
	if index < 0 || index >= a.size {
		return 0, false
	}
	val := a.data[index]
	// 前移
	for i := index; i < a.size-1; i++ {
		a.data[i] = a.data[i+1]
	}
	a.size--
	return val, true
}

// Search 线性查找 O(n)
func (a *DynamicArray) Search(val int) int {
	for i := 0; i < a.size; i++ {
		if a.data[i] == val {
			return i
		}
	}
	return -1
}

// BinarySearch 二分查找 O(log n)，要求数组有序
func (a *DynamicArray) BinarySearch(val int) int {
	lo, hi := 0, a.size-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if a.data[mid] == val {
			return mid
		} else if a.data[mid] < val {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return -1
}

// ToSlice 返回底层切片
func (a *DynamicArray) ToSlice() []int {
	result := make([]int, a.size)
	copy(result, a.data[:a.size])
	return result
}

func (a *DynamicArray) resize(newCap int) {
	newData := make([]int, newCap)
	copy(newData, a.data[:a.size])
	a.data = newData
	a.capacity = newCap
}
