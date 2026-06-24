# 图 练习题

## 1. 课程表（拓扑排序）

给定课程数量 n 和先修关系 `prerequisites`（`[a,b]` 表示修 a 必须先修 b），判断能否完成所有课程（是否有环）。

### 思路

Kahn 算法（BFS 拓扑排序）：计算入度 → 入度为 0 入队 → 逐个出队减入度。

### 实现

```go
func canFinish(numCourses int, prerequisites [][]int) bool {
    // 构建图和入度
    graph := make([][]int, numCourses)
    indegree := make([]int, numCourses)
    for _, p := range prerequisites {
        a, b := p[0], p[1]
        graph[b] = append(graph[b], a) // b → a
        indegree[a]++
    }

    // 入度为 0 的入队
    queue := []int{}
    for i := 0; i < numCourses; i++ {
        if indegree[i] == 0 {
            queue = append(queue, i)
        }
    }

    count := 0
    for len(queue) > 0 {
        course := queue[0]
        queue = queue[1:]
        count++

        for _, next := range graph[course] {
            indegree[next]--
            if indegree[next] == 0 {
                queue = append(queue, next)
            }
        }
    }

    return count == numCourses
}
```

## 2. 岛屿数量（图的 DFS/BFS）

网格中 '1' 是陆地，'0' 是水。计算岛屿数量。

### 思路

DFS：遍历每个格子，遇到 '1' 就 DFS 淹没（标记为 '0'），岛屿数 +1。

### 实现

```go
func numIslands(grid [][]byte) int {
    if len(grid) == 0 {
        return 0
    }
    count := 0
    for i := range grid {
        for j := range grid[i] {
            if grid[i][j] == '1' {
                count++
                dfs(grid, i, j)
            }
        }
    }
    return count
}

func dfs(grid [][]byte, i, j int) {
    if i < 0 || i >= len(grid) || j < 0 || j >= len(grid[0]) || grid[i][j] == '0' {
        return
    }
    grid[i][j] = '0'
    dfs(grid, i+1, j)
    dfs(grid, i-1, j)
    dfs(grid, i, j+1)
    dfs(grid, i, j-1)
}
```
