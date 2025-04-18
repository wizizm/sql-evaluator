# SQL Evaluator

一个轻量级的SQL WHERE子句评估器，可以将SQL WHERE条件转换为Go逻辑进行执行。

## 特性

- 支持将SQL WHERE子句转换为Go逻辑
- 支持json标签映射
- 支持常见的比较操作符（=, !=, >, <, >=, <=）
- 支持AND和OR逻辑组合
- 支持括号表达式
- 支持布尔值比较
- 支持数值类型比较
- 支持字符串比较
- 支持NULL值处理
- 区分NULL和空字符串
- 支持指针和非指针类型字段
- 支持自动类型转换（int和float64之间）
- 支持LIKE/NOT LIKE模式匹配
- 支持IN/NOT IN值列表比较
- 支持BETWEEN/NOT BETWEEN范围比较

## 安装

```bash
go get github.com/wizizm/sql-evaluator
```

## 使用示例

### 基本用法

```go
package main

import (
    "fmt"
    sqlevaluator "github.com/wizizm/sql-evaluator"
)

// User 示例用户模型（使用指针类型支持NULL值）
type User struct {
    ID       *int     `json:"id"`
    Name     *string  `json:"name"`
    Age      *int     `json:"age"`
    Salary   *float64 `json:"salary"`
    IsActive *bool    `json:"is_active"`
}

// 辅助函数：创建指针类型
func strPtr(s string) *string {
    return &s
}

func intPtr(i int) *int {
    return &i
}

func float64Ptr(f float64) *float64 {
    return &f
}

func boolPtr(b bool) *bool {
    return &b
}

func main() {
    // 创建model实例
    user := &User{
        ID:       intPtr(1),
        Name:     strPtr("张三"),
        Age:      intPtr(25),
        Salary:   float64Ptr(5000.50),
        IsActive: boolPtr(true),
    }

    // 创建评估器
    evaluator := sqlevaluator.NewSQLEvaluator(user)

    // 评估WHERE子句
    result, err := evaluator.EvaluateWhere("name = '张三' AND age > 20")
    if err != nil {
        fmt.Printf("评估失败: %v\n", err)
        return
    }

    fmt.Printf("评估结果: %v\n", result)
}
```

### NULL值处理

```go

// NULL值处理示例
userWithNull := &User{
    ID:       intPtr(1),
    Name:     nil, // NULL值
    Age:      intPtr(30),
    Salary:   nil, // NULL值
    IsActive: boolPtr(true),
}
evaluator := sqlevaluator.NewSQLEvaluator(userWithNull)
// 检查字段是否为NULL
isNull, _ := evaluator.EvaluateWhere("name IS NULL")
fmt.Printf("name IS NULL: %v\n", isNull) // 输出: true

// 空字符串示例
userWithEmpty := &User{
    ID:       intPtr(1),
    Name:     strPtr(""), // 空字符串
    Age:      intPtr(25),
    Salary:   float64Ptr(0.0), // 零值
    IsActive: boolPtr(false), // false值
}
evaluator := sqlevaluator.NewSQLEvaluator(userWithEmpty)
// 检查字段是否为空字符串
isEmpty, _ := evaluator.EvaluateWhere("name = ''")
fmt.Printf("name = '': %v\n", isEmpty) // 输出: true
// 区分NULL和空字符串
isNullEmpty, _ := evaluator.EvaluateWhere("name IS NULL OR name = ''")
fmt.Printf("name IS NULL OR name = '': %v\n", isNullEmpty) // 输出: true
```

### 非指针类型字段

```go
// UserWithNonPtr 示例用户模型（使用非指针类型）
type UserWithNonPtr struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Age      int     `json:"age"`
    Salary   float64 `json:"salary"`
    IsActive bool    `json:"is_active"`
}

// 使用非指针类型model
userNonPtr := &UserWithNonPtr{
    ID:       1,
    Name:     "张三",
    Age:      25,
    Salary:   5000.50,
    IsActive: true,
}

// 创建非指针类型评估器
evaluatorNonPtr := sqlevaluator.NewSQLEvaluator(userNonPtr)

// 评估非指针类型WHERE子句
resultNonPtr, err := evaluatorNonPtr.EvaluateWhere("name = '张三' AND age > 20")
if err != nil {
    fmt.Printf("评估失败: %v\n", err)
    return
}

fmt.Printf("非指针类型评估结果: %v\n", resultNonPtr)
```

## 支持的SQL操作

- 相等比较 (=)
- 不等比较 (!=)
- 大于比较 (>)
- 小于比较 (<)
- 大于等于比较 (>=)
- 小于等于比较 (<=)
- AND 逻辑
- OR 逻辑
- 括号表达式
- 布尔值比较
- LIKE/NOT LIKE 模式匹配
- IN/NOT IN 值列表比较
- IS NULL/IS NOT NULL NULL值检查
- BETWEEN/NOT BETWEEN 范围比较

## NULL值处理

SQL Evaluator 支持对NULL值的处理，并区分NULL和空字符串：

- 使用指针类型（如`*string`）表示可为NULL的字段
- `IS NULL`操作符检查字段是否为NULL
- `IS NOT NULL`操作符检查字段是否不为NULL
- 空字符串(`''`)和NULL是不同的值
- 零值（0、0.0、false）和NULL是不同的值
- 使用辅助函数（strPtr、intPtr等）创建指针类型

## 类型转换

SQL Evaluator 支持以下类型转换：

- int 和 float64 之间的自动转换
- 字符串和数值之间的比较（需要显式转换）
- 布尔值和数值之间的比较（需要显式转换）

## 贡献

欢迎提交 Pull Request 和 Issue！

## 许可证

MIT License 