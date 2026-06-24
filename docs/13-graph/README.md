# 13 图

## 第一性原理：为什么需要图？

### 问题

```
树 → 一对多层级关系
图 → 多对多网状关系

现实世界：
- 社交网络：人和人的关系
- 地图导航：城市之间的道路
- 微服务：服务之间的调用关系
- 包依赖：模块之间的依赖关系

这些都是「网状」的，不是「层级」的
```

这就是图。

---

## 图的定义

```
图 G = (V, E)

V (Vertex)：顶点集合
E (Edge)：边集合

    1 ──→ 2
    │     │
    ↓     ↓
    3 ←── 4
```

### 图的基本概念

| 概念 | 说明 |
|------|------|
| 有向图 | 边有方向 (A → B) |
| 无向图 | 边无方向 (A — B) |
| 权重 | 边上的数值 |
| 度 | 连接某顶点的边数 |
| 入度 | 指向该顶点的边数（有向图） |
| 出度 | 从该顶点出发的边数（有向图） |
| 路径 | 顶点序列 v₁→v₂→...→vₖ |
| 环 | 起点 = 终点的路径 |

---

## 图的表示

### 邻接矩阵

```
   1  2  3  4
1  0  1  0  1
2  1  0  1  0
3  0  1  0  1
4  1  0  1  0

存在边 → 1，不存在 → 0

优点：判断两点是否连通 O(1)
缺点：O(V²) 空间，稀疏图浪费
```

### 邻接表

```go
// 邻接表：每个顶点维护它连接到的顶点列表
graph := map[int][]int{
    1: {2, 4},
    2: {1, 3},
    3: {2, 4},
    4: {1, 3},
}

// 优点：O(V+E) 空间，稀疏图高效
// 缺点：判断两点连通 O(degree)
```

---

## 核心算法

### DFS（深度优先搜索）

```go
func dfs(graph map[int][]int, start int, visited map[int]bool) {
    visited[start] = true
    fmt.Println(start)

    for _, neighbor := range graph[start] {
        if !visited[neighbor] {
            dfs(graph, neighbor, visited)
        }
    }
}
```

```
DFS 的路径：

    1
   / \
  2   4
  |
  3

DFS(1): 1 → 2 → 3 → 4

用途：
- 检测环
- 拓扑排序
- 连通分量
- 路径搜索（不保证最短）
```

### BFS（广度优先搜索）

```go
func bfs(graph map[int][]int, start int) {
    visited := map[int]bool{start: true}
    queue := []int{start}

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        fmt.Println(node)

        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
}
```

```
BFS 的路径：

    1
   / \
  2   4
  |
  3

BFS(1): 1 → 2 → 4 → 3

用途：
- 最短路径（无权图）
- 层级遍历
- 社交网络"几度人脉"
```

### 拓扑排序

```
问题：任务有依赖关系，如何安排执行顺序？

A → B → C
↓
D → E

A 必须在 B 之前，A 必须在 D 之前...

拓扑排序 BFS（Kahn 算法）：
1. 计算所有节点的入度
2. 入度为 0 的节点入队
3. 出队 → 其邻居入度 -1 → 新入度为 0 的入队
4. 重复直到队列空

结果：A, B, D, C, E（一种可能的顺序）
```

```go
func topologicalSort(n int, edges [][]int) []int {
    // 构建图和入度
    graph := make(map[int][]int)
    indegree := make([]int, n)
    for _, e := range edges {
        from, to := e[0], e[1]
        graph[from] = append(graph[from], to)
        indegree[to]++
    }

    // 入度为 0 的入队
    queue := []int{}
    for i := 0; i < n; i++ {
        if indegree[i] == 0 {
            queue = append(queue, i)
        }
    }

    var result []int
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        result = append(result, node)

        for _, neighbor := range graph[node] {
            indegree[neighbor]--
            if indegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }

    if len(result) != n {
        return nil // 有环，无法拓扑排序
    }
    return result
}
```

---

## 工程案例

### 1. 微服务调用链

```
服务 A → 服务 B → 服务 D
  │         │
  ↓         ↓
服务 C    服务 E

这种图的意义：
- 依赖分析：A 挂了影响谁？
- 调用链追踪：一次请求经过了哪些服务？
- 循环依赖检测：A → B → C → A（危险！）
```

### 2. Kubernetes 依赖分析

```
Pod1 → Service1 → Pod2
   │                 │
   ↓                 ↓
ConfigMap1       PVC → PV

K8S 资源之间的依赖关系本质是图
- Pod 依赖 ConfigMap、Secret、PVC
- Service 依赖 Pod（通过 Label Selector）
- Deployment 管理 Pod 的生命周期
```

### 3. 包依赖管理

```
go.mod:
  my-app → lib-a v1.0
         → lib-b v2.0
         → lib-c v1.5

  lib-a → lib-d v1.0
  lib-b → lib-d v2.0  ← 冲突！

这是有向图
拓扑排序 → 确定编译顺序
循环依赖 → 编译失败
版本冲突 → 需要解决
```

---

## 练习题

### 图的遍历和拓扑排序

详见 [`exercises/13-graph/`](../../exercises/13-graph/)

---

## 小结

```
图 = 顶点 + 边 → 多对多关系

核心算法：
- DFS → 路径搜索、环检测、连通分量
- BFS → 最短路径（无权图）、层级遍历
- 拓扑排序 → 依赖排序、死锁检测

工程价值：
- 微服务 → 调用链分析
- Kubernetes → 资源依赖图
- 包管理 → 依赖解析
- 社交网络 → 关系图

关键洞察：
图是所有数据结构中最通用的。
树是图的特例（无环连通图）。
链表是图的特例（每个节点出度 ≤ 1）。
```

---

**上一篇：[12 一致性哈希](../12-consistent-hash/README.md)**
**下一篇：[14 系统设计](../14-system-design/README.md)**
