# 数组练习题

## 1. 反转数组

原地反转一个整数数组，要求 O(1) 额外空间。

### 思路

双指针从两端向中间交换。

### 实现

```go
func reverse(arr []int) {
    for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
        arr[i], arr[j] = arr[j], arr[i]
    }
}
```

## 2. 二分查找

在有序数组中查找目标值，返回索引，不存在返回 -1。

### 思路

每次取中间值比较，缩小一半搜索范围。

### 实现

```go
func binarySearch(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if nums[mid] == target {
            return mid
        } else if nums[mid] < target {
            lo = mid + 1
        } else {
            hi = mid - 1
        }
    }
    return -1
}
```

## 3. 合并两个有序数组

合并两个有序数组，结果保持有序。

### 思路

从后往前填充，避免移动元素。

### 实现

```go
func merge(nums1 []int, m int, nums2 []int, n int) {
    i, j, k := m-1, n-1, m+n-1
    for j >= 0 {
        if i >= 0 && nums1[i] > nums2[j] {
            nums1[k] = nums1[i]
            i--
        } else {
            nums1[k] = nums2[j]
            j--
        }
        k--
    }
}
```
