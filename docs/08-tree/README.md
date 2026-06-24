# 08 树

## 第一性原理：为什么需要树？

### 问题

```
数组/链表 → 线性结构 → 一维关系
堆 → 完全二叉树 → 但只是为"最值"设计的

实际世界中的关系：
- 公司组织架构 → 层级关系
- 文件系统 → 目录层级
- HTML DOM → 嵌套结构

这些都需要：层级结构 + 快速查找
```

这就是树。

---

## 树的定义

```
树是一种「一对多」的层级数据结构

术语：
- 根节点 (Root)：树的起点
- 父节点 (Parent)
- 子节点 (Child)
- 叶节点 (Leaf)：没有子节点
- 高度 (Height)：从根到最远叶子的边数
- 深度 (Depth)：从根到该节点的边数
```

```
        1          ← 根节点，深度 0
      / | \
     2  3  4       ← 深度 1
    / \    |
   5   6   7       ← 深度 2，叶节点
```

---

## 二叉树 → 二叉搜索树 (BST)

### BST 特性

```
- 左子树所有节点 < 根节点
- 右子树所有节点 > 根节点
- 左右子树也是 BST

       8
     /   \
    3     10
   / \      \
  1   6      14
     / \    /
    4   7  13
```

**查找一个值**：类似二分查找
```
找 7：8 > 7 → 左 / 3 < 7 → 右 / 6 < 7 → 右 → 找到！
找 9：8 < 9 → 右 / 10 > 9 → 左 → nil → 不存在
```

### BST 的问题：退化

```
插入顺序：1, 2, 3, 4, 5, 6

      1
       \
        2
         \
          3
           \
            4
             \
              5
               \
                6

这棵树退化成链表 → 查找 O(n)
```

### 解决方案：自平衡树

```
AVL 树：
- 严格平衡：左右子树高度差 ≤ 1
- 插入/删除后旋转修复
- 查找快，但插入/删除旋转多

红黑树：
- 宽松平衡：最长路径 ≤ 2 × 最短路径
- 插入/删除旋转少
- 工程中更常用
```

---

## 时间复杂度

| 操作 | BST (平均) | BST (最坏) | AVL | 红黑树 |
|------|-----------|-----------|-----|--------|
| 查找 | O(log n) | O(n) | O(log n) | O(log n) |
| 插入 | O(log n) | O(n) | O(log n) | O(log n) |
| 删除 | O(log n) | O(n) | O(log n) | O(log n) |

---

## Go 代码实现

```go
// code/08-tree/bst.go

package tree

// TreeNode BST 节点
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

// BST 二叉搜索树
type BST struct {
    root *TreeNode
}

// Search 查找 O(log n) 平均
func (t *BST) Search(val int) *TreeNode {
    cur := t.root
    for cur != nil {
        if cur.Val == val {
            return cur
        } else if val < cur.Val {
            cur = cur.Left
        } else {
            cur = cur.Right
        }
    }
    return nil
}

// Insert 插入 O(log n) 平均
func (t *BST) Insert(val int) {
    if t.root == nil {
        t.root = &TreeNode{Val: val}
        return
    }
    cur := t.root
    for {
        if val < cur.Val {
            if cur.Left == nil {
                cur.Left = &TreeNode{Val: val}
                return
            }
            cur = cur.Left
        } else {
            if cur.Right == nil {
                cur.Right = &TreeNode{Val: val}
                return
            }
            cur = cur.Right
        }
    }
}
```

### DFS 与 BFS

```go
// DFS 深度优先遍历

// 前序：根 → 左 → 右
func preorder(root *TreeNode) []int {
    if root == nil {
        return nil
    }
    var res []int
    res = append(res, root.Val)
    res = append(res, preorder(root.Left)...)
    res = append(res, preorder(root.Right)...)
    return res
}

// 中序：左 → 根 → 右（BST 中序 = 有序序列）
func inorder(root *TreeNode) []int {
    if root == nil {
        return nil
    }
    var res []int
    res = append(res, inorder(root.Left)...)
    res = append(res, root.Val)
    res = append(res, inorder(root.Right)...)
    return res
}

// BFS 广度优先（层序）— 用队列
func levelOrder(root *TreeNode) [][]int {
    if root == nil {
        return nil
    }
    var res [][]int
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
        res = append(res, level)
    }
    return res
}
```

---

## 工程案例

### 1. 组织架构

```
CEO
├── CTO
│   ├── 后端组
│   └── 前端组
└── CFO
    ├── 财务
    └── 法务

这是树 → 层级结构的天然建模
```

### 2. 文件系统

```
/
├── home/
│   ├── user/
│   │   ├── documents/
│   │   └── downloads/
│   └── guest/
├── etc/
│   └── nginx/
└── var/
    └── log/

每个目录是一个节点
子目录是子节点
文件是叶子节点
```

### 3. MySQL 索引为什么不用 BST？

```
树在内存中很好，但在磁盘上呢？

问题：磁盘 IO 是瓶颈
- BS T → 每个节点一个值，一个磁盘页可能只存一个节点 → 浪费
- 树太高 → 磁盘 IO 次数多

解决方案：B+Tree → 一个节点存多个值
详细见 [10 B+Tree](../10-bplus-tree/README.md)
```

---

## 练习题

### DFS / BFS

详见 [`exercises/08-tree/`](../../exercises/08-tree/)

---

## 小结

```
树 = 层级结构 + 一对多关系

二叉搜索树 → O(log n) 查找（平衡时）
退化问题 → AVL/红黑树解决

树的局限：
- BST 在磁盘上效率低 → 引出 B+Tree
- BST 不适合范围查找 → 引出 B+Tree 的叶子链表

工程价值：
- 文件系统 → 树形目录
- 组织架构 → 层级建模
- DFS/BFS → 遍历骨架（几乎所有树算法的基础）
```

---

**上一篇：[07 Heap](../07-heap/README.md)**
**下一篇：[09 SkipList](../09-skiplist/README.md)**
