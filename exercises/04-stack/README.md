# 栈练习题

## 1. 有效括号

判断字符串中的括号是否有效匹配。

### 思路

遍历字符串，遇到左括号入栈，右括号与栈顶匹配。

### 实现

```go
func isValid(s string) bool {
    stack := []rune{}
    pairs := map[rune]rune{
        ')': '(',
        '}': '{',
        ']': '[',
    }

    for _, ch := range s {
        switch ch {
        case '(', '{', '[':
            stack = append(stack, ch)
        case ')', '}', ']':
            if len(stack) == 0 || stack[len(stack)-1] != pairs[ch] {
                return false
            }
            stack = stack[:len(stack)-1]
        }
    }
    return len(stack) == 0
}
```

## 2. 逆波兰表达式求值

计算逆波兰（后缀）表达式的值。

### 思路

遇到数字入栈，遇到操作符弹出两个计算后入栈。

### 实现

```go
func evalRPN(tokens []string) int {
    stack := []int{}
    for _, t := range tokens {
        switch t {
        case "+", "-", "*", "/":
            b, a := stack[len(stack)-1], stack[len(stack)-2]
            stack = stack[:len(stack)-2]
            switch t {
            case "+": stack = append(stack, a+b)
            case "-": stack = append(stack, a-b)
            case "*": stack = append(stack, a*b)
            case "/": stack = append(stack, a/b)
            }
        default:
            v, _ := strconv.Atoi(t)
            stack = append(stack, v)
        }
    }
    return stack[0]
}
```

## 3. 最小栈

设计一个栈支持 push、pop、top 和 getMin（O(1)）。

### 思路

维护两个栈：一个存数据，一个存当前最小值。

### 实现

```go
type MinStack struct {
    data []int
    mins []int
}

func (s *MinStack) Push(val int) {
    s.data = append(s.data, val)
    if len(s.mins) == 0 || val <= s.mins[len(s.mins)-1] {
        s.mins = append(s.mins, val)
    }
}

func (s *MinStack) Pop() {
    if len(s.data) == 0 {
        return
    }
    if s.data[len(s.data)-1] == s.mins[len(s.mins)-1] {
        s.mins = s.mins[:len(s.mins)-1]
    }
    s.data = s.data[:len(s.data)-1]
}

func (s *MinStack) Top() int {
    return s.data[len(s.data)-1]
}

func (s *MinStack) GetMin() int {
    return s.mins[len(s.mins)-1]
}
```
