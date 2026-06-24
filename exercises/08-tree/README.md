# 树 练习题

## 1. 二叉树的最大深度

### 思路

DFS 递归或 BFS 层序遍历。

### DFS 实现

```go
func maxDepth(root *TreeNode) int {
    if root == nil {
        return 0
    }
    left := maxDepth(root.Left)
    right := maxDepth(root.Right)
    if left > right {
        return left + 1
    }
    return right + 1
}
```

### BFS 实现

```go
func maxDepthBFS(root *TreeNode) int {
    if root == nil {
        return 0
    }
    depth := 0
    queue := []*TreeNode{root}
    for len(queue) > 0 {
        depth++
        size := len(queue)
        for i := 0; i < size; i++ {
            node := queue[0]
            queue = queue[1:]
            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
    }
    return depth
}
```

## 2. 验证二叉搜索树

### 思路

利用 BST 性质：中序遍历严格递增。

### 实现

```go
func isValidBST(root *TreeNode) bool {
    var prev *int
    return inorderBST(root, &prev)
}

func inorderBST(node *TreeNode, prev **int) bool {
    if node == nil {
        return true
    }
    if !inorderBST(node.Left, prev) {
        return false
    }
    if *prev != nil && node.Val <= **prev {
        return false
    }
    *prev = &node.Val
    return inorderBST(node.Right, prev)
}
```

## 3. 二叉树层序遍历

### 思路

BFS 队列，记录每层大小。

### 实现

```go
func levelOrder(root *TreeNode) [][]int {
    if root == nil {
        return nil
    }
    var result [][]int
    queue := []*TreeNode{root}
    for len(queue) > 0 {
        levelSize := len(queue)
        var level []int
        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:]
            level = append(level, node.Val)
            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        result = append(result, level)
    }
    return result
}
```
