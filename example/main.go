package main

import (
	"fmt"

	"github.com/wizizm/sql-evaluator/pkg/sqlevaluator"
)

// User 示例用户模型
type User struct {
	ID       *int     `xorm:"id"`
	Name     *string  `xorm:"name"`
	Age      *int     `xorm:"age"`
	Salary   *float64 `xorm:"salary"`
	IsActive *bool    `xorm:"is_active"`
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
