# HashMap 练习题

## 1. Two Sum

给定数组和 target，找两个数之和等于 target 的索引。

### 思路

遍历时用 HashMap 存 `值 → 索引`，O(n)。

### 实现

```go
func twoSum(nums []int, target int) []int {
    m := make(map[int]int)
    for i, v := range nums {
        if j, ok := m[target-v]; ok {
            return []int{j, i}
        }
        m[v] = i
    }
    return nil
}
```

## 2. 无重复字符的最长子串

### 思路

滑动窗口 + HashMap 记录字符上次出现位置。

### 实现

```go
func lengthOfLongestSubstring(s string) int {
    lastPos := make(map[byte]int)
    start, maxLen := 0, 0

    for i := range s {
        if pos, ok := lastPos[s[i]]; ok && pos >= start {
            start = pos + 1
        }
        lastPos[s[i]] = i
        if i-start+1 > maxLen {
            maxLen = i - start + 1
        }
    }
    return maxLen
}
```

## 3. 字母异位词分组

### 思路

对每个字符串排序作为 key，用 HashMap 分组。

### 实现

```go
func groupAnagrams(strs []string) [][]string {
    groups := make(map[string][]string)
    for _, s := range strs {
        key := sortString(s)
        groups[key] = append(groups[key], s)
    }
    result := make([][]string, 0, len(groups))
    for _, g := range groups {
        result = append(result, g)
    }
    return result
}

func sortString(s string) string {
    b := []byte(s)
    sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
    return string(b)
}
```
