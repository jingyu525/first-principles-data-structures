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

// NewBST 创建 BST
func NewBST() *BST {
	return &BST{}
}

// Search 查找 O(log n) 平均，O(n) 最坏
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
		} else if val > cur.Val {
			if cur.Right == nil {
				cur.Right = &TreeNode{Val: val}
				return
			}
			cur = cur.Right
		} else {
			return // 已存在，不重复插入
		}
	}
}

// Inorder 中序遍历（输出有序序列）
func (t *BST) Inorder() []int {
	var result []int
	inorder(t.root, &result)
	return result
}

func inorder(node *TreeNode, result *[]int) {
	if node == nil {
		return
	}
	inorder(node.Left, result)
	*result = append(*result, node.Val)
	inorder(node.Right, result)
}

// Preorder 前序遍历
func (t *BST) Preorder() []int {
	var result []int
	preorder(t.root, &result)
	return result
}

func preorder(node *TreeNode, result *[]int) {
	if node == nil {
		return
	}
	*result = append(*result, node.Val)
	preorder(node.Left, result)
	preorder(node.Right, result)
}

// LevelOrder 层序遍历 (BFS)
func (t *BST) LevelOrder() [][]int {
	if t.root == nil {
		return nil
	}
	var result [][]int
	queue := []*TreeNode{t.root}
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
