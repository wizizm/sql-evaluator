package sqlevaluator

import (
	"testing"
)

// User 示例用户模型
type User struct {
	ID       *int     `json:"id"`
	Name     *string  `json:"name"`
	Age      *int     `json:"age"`
	Salary   *float64 `json:"salary"`
	IsActive *bool    `json:"is_active"`
}

// UserWithNonPtr 示例用户模型（非指针字段）
type UserWithNonPtr struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	Salary   float64 `json:"salary"`
	IsActive bool    `json:"is_active"`
}

// strPtr 返回字符串的指针
func strPtr(s string) *string {
	return &s
}

// intPtr 返回整数的指针
func intPtr(i int) *int {
	return &i
}

// float64Ptr 返回float64的指针
func float64Ptr(f float64) *float64 {
	return &f
}

// boolPtr 返回bool的指针
func boolPtr(b bool) *bool {
	return &b
}

func TestSQLEvaluator(t *testing.T) {
	tests := []struct {
		name        string
		model       interface{}
		whereClause string
		want        bool
		wantErr     bool
	}{
		{
			name: "简单相等比较",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name = '张三' AND age = 25",
			want:        true,
			wantErr:     false,
		},
		{
			name: "数值比较",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(30),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "age > 25 AND salary >= 5000",
			want:        true,
			wantErr:     false,
		},
		{
			name: "布尔值比较",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "is_active = true",
			want:        true,
			wantErr:     false,
		},
		{
			name: "复杂条件",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("李四"),
				Age:      intPtr(35),
				Salary:   float64Ptr(8000.75),
				IsActive: boolPtr(true),
			},
			whereClause: "(name = '李四' AND age > 30) OR (salary > 7000 AND is_active = true)",
			want:        true,
			wantErr:     false,
		},
		{
			name: "LIKE操作符测试",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三丰"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name LIKE '张%'",
			want:        true,
			wantErr:     false,
		},
		{
			name: "NOT LIKE操作符测试",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三丰"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name NOT LIKE '李%'",
			want:        true,
			wantErr:     false,
		},
		{
			name: "IN操作符测试",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "age IN (20, 25, 30)",
			want:        true,
			wantErr:     false,
		},
		{
			name: "NOT IN操作符测试",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "age NOT IN (20, 30, 35)",
			want:        true,
			wantErr:     false,
		},
		{
			name: "IS NULL操作符测试",
			model: &User{
				ID:       intPtr(1),
				Name:     nil,
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name IS NULL",
			want:        true,
			wantErr:     false,
		},
		{
			name: "IS NOT NULL操作符测试",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name IS NOT NULL",
			want:        true,
			wantErr:     false,
		},
		{
			name: "空字符串测试",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr(""),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name = ''",
			want:        true,
			wantErr:     false,
		},
		{
			name: "空字符串不是NULL测试",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr(""),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name IS NULL",
			want:        false,
			wantErr:     false,
		},
		{
			name: "NULL不是空字符串测试",
			model: &User{
				ID:       intPtr(1),
				Name:     nil,
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name = ''",
			want:        false,
			wantErr:     false,
		},
		{
			name: "BETWEEN操作符测试-数值",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "age BETWEEN 20 AND 30",
			want:        true,
			wantErr:     false,
		},
		{
			name: "NOT BETWEEN操作符测试-数值",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(35),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "age NOT BETWEEN 20 AND 30",
			want:        true,
			wantErr:     false,
		},
		{
			name: "BETWEEN操作符测试-字符串",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name BETWEEN '张' AND '李'",
			want:        true,
			wantErr:     false,
		},
		{
			name: "组合条件测试",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三丰"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "(name LIKE '张%' AND age BETWEEN 20 AND 30) OR (salary IN (5000.50, 6000) AND is_active = true)",
			want:        true,
			wantErr:     false,
		},
		{
			name: "错误条件测试-类型不匹配",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name > 100",
			want:        false,
			wantErr:     true,
		},
		{
			name: "错误条件测试-无效操作符",
			model: &User{
				ID:       intPtr(1),
				Name:     strPtr("张三"),
				Age:      intPtr(25),
				Salary:   float64Ptr(5000.50),
				IsActive: boolPtr(true),
			},
			whereClause: "name CONTAINS '张'",
			want:        false,
			wantErr:     true,
		},
		{
			name: "所有字段NULL测试",
			model: &User{
				ID:       nil,
				Name:     nil,
				Age:      nil,
				Salary:   nil,
				IsActive: nil,
			},
			whereClause: "id IS NULL AND name IS NULL AND age IS NULL AND salary IS NULL AND is_active IS NULL",
			want:        true,
			wantErr:     false,
		},
		{
			name: "所有字段空值测试",
			model: &User{
				ID:       intPtr(0),
				Name:     strPtr(""),
				Age:      intPtr(0),
				Salary:   float64Ptr(0.0),
				IsActive: boolPtr(false),
			},
			whereClause: "id = 0 AND name = '' AND age = 0 AND salary = 0.0 AND is_active = false",
			want:        true,
			wantErr:     false,
		},
		{
			name: "混合NULL和非NULL值测试",
			model: &User{
				ID:       intPtr(1),
				Name:     nil,
				Age:      intPtr(30),
				Salary:   nil,
				IsActive: boolPtr(true),
			},
			whereClause: "id = 1 AND name IS NULL AND age = 30 AND salary IS NULL AND is_active = true",
			want:        true,
			wantErr:     false,
		},
		{
			name: "NULL值比较返回false测试",
			model: &User{
				ID:       nil,
				Name:     nil,
				Age:      nil,
				Salary:   nil,
				IsActive: nil,
			},
			whereClause: "id = 1 OR name = '张三' OR age > 20 OR salary < 5000 OR is_active = true",
			want:        false,
			wantErr:     false,
		},
		{
			name: "NULL值与NULL比较测试",
			model: &User{
				ID:       nil,
				Name:     nil,
				Age:      nil,
				Salary:   nil,
				IsActive: nil,
			},
			whereClause: "id = NULL OR name = NULL OR age = NULL OR salary = NULL OR is_active = NULL",
			want:        false,
			wantErr:     false,
		},
		{
			name: "空值与NULL比较测试",
			model: &User{
				ID:       intPtr(0),
				Name:     strPtr(""),
				Age:      intPtr(0),
				Salary:   float64Ptr(0.0),
				IsActive: boolPtr(false),
			},
			whereClause: "id = NULL OR name = NULL OR age = NULL OR salary = NULL OR is_active = NULL",
			want:        false,
			wantErr:     false,
		},
		{
			name: "复杂NULL条件测试",
			model: &User{
				ID:       intPtr(1),
				Name:     nil,
				Age:      intPtr(30),
				Salary:   nil,
				IsActive: boolPtr(true),
			},
			whereClause: "(id = 1 AND name IS NULL) OR (age = 30 AND salary IS NULL) OR (is_active = true AND name IS NULL)",
			want:        true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewSQLEvaluator(tt.model)
			got, err := evaluator.EvaluateWhere(tt.whereClause)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateWhere() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSQLEvaluatorNullValues 测试所有类型字段的NULL值处理
func TestSQLEvaluatorNullValues(t *testing.T) {
	// 创建包含NULL值的测试数据
	userWithNulls := &User{
		ID:       nil,
		Name:     nil,
		Age:      nil,
		Salary:   nil,
		IsActive: nil,
	}

	// 创建评估器
	evaluator := NewSQLEvaluator(userWithNulls)

	// NULL值测试用例
	nullTestCases := []struct {
		name     string
		where    string
		expected bool
	}{
		{
			name:     "ID IS NULL",
			where:    "id IS NULL",
			expected: true,
		},
		{
			name:     "Name IS NULL",
			where:    "name IS NULL",
			expected: true,
		},
		{
			name:     "Age IS NULL",
			where:    "age IS NULL",
			expected: true,
		},
		{
			name:     "Salary IS NULL",
			where:    "salary IS NULL",
			expected: true,
		},
		{
			name:     "IsActive IS NULL",
			where:    "is_active IS NULL",
			expected: true,
		},
		{
			name:     "ID IS NOT NULL",
			where:    "id IS NOT NULL",
			expected: false,
		},
		{
			name:     "Name IS NOT NULL",
			where:    "name IS NOT NULL",
			expected: false,
		},
		{
			name:     "Age IS NOT NULL",
			where:    "age IS NOT NULL",
			expected: false,
		},
		{
			name:     "Salary IS NOT NULL",
			where:    "salary IS NOT NULL",
			expected: false,
		},
		{
			name:     "IsActive IS NOT NULL",
			where:    "is_active IS NOT NULL",
			expected: false,
		},
		{
			name:     "NULL值比较返回false",
			where:    "id = 1",
			expected: false,
		},
		{
			name:     "NULL值与NULL比较",
			where:    "id = NULL",
			expected: false,
		},
		{
			name:     "多个字段IS NULL",
			where:    "id IS NULL AND name IS NULL AND age IS NULL",
			expected: true,
		},
		{
			name:     "混合NULL和非NULL条件",
			where:    "id IS NULL AND name = '张三'",
			expected: false,
		},
	}

	for _, tt := range nullTestCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluator.EvaluateWhere(tt.where)

			if err != nil {
				t.Errorf("EvaluateWhere() error = %v", err)
				return
			}

			if got != tt.expected {
				t.Errorf("EvaluateWhere() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestSQLEvaluatorEmptyValues 测试空值（非NULL）的处理
func TestSQLEvaluatorEmptyValues(t *testing.T) {
	// 创建包含空值的测试数据
	userWithEmptyValues := &User{
		ID:       intPtr(0),
		Name:     strPtr(""),
		Age:      intPtr(0),
		Salary:   float64Ptr(0.0),
		IsActive: boolPtr(false),
	}

	// 创建评估器
	evaluator := NewSQLEvaluator(userWithEmptyValues)

	// 空值测试用例
	emptyTestCases := []struct {
		name     string
		where    string
		expected bool
	}{
		{
			name:     "ID = 0",
			where:    "id = 0",
			expected: true,
		},
		{
			name:     "Name = ''",
			where:    "name = ''",
			expected: true,
		},
		{
			name:     "Age = 0",
			where:    "age = 0",
			expected: true,
		},
		{
			name:     "Salary = 0.0",
			where:    "salary = 0.0",
			expected: true,
		},
		{
			name:     "IsActive = false",
			where:    "is_active = false",
			expected: true,
		},
		{
			name:     "空值IS NULL返回false",
			where:    "id IS NULL",
			expected: false,
		},
		{
			name:     "空字符串IS NULL返回false",
			where:    "name IS NULL",
			expected: false,
		},
		{
			name:     "空值IS NOT NULL返回true",
			where:    "id IS NOT NULL",
			expected: true,
		},
		{
			name:     "空字符串IS NOT NULL返回true",
			where:    "name IS NOT NULL",
			expected: true,
		},
		{
			name:     "空值与NULL比较",
			where:    "id = NULL",
			expected: false,
		},
		{
			name:     "空字符串与NULL比较",
			where:    "name = NULL",
			expected: false,
		},
	}

	for _, tt := range emptyTestCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluator.EvaluateWhere(tt.where)

			if err != nil {
				t.Errorf("EvaluateWhere() error = %v", err)
				return
			}

			if got != tt.expected {
				t.Errorf("EvaluateWhere() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestSQLEvaluatorMixedValues 测试混合NULL和非NULL值的处理
func TestSQLEvaluatorMixedValues(t *testing.T) {
	// 创建混合NULL和非NULL值的测试数据
	userWithMixedValues := &User{
		ID:       intPtr(1),
		Name:     nil,
		Age:      intPtr(30),
		Salary:   nil,
		IsActive: boolPtr(true),
	}

	// 创建评估器
	evaluator := NewSQLEvaluator(userWithMixedValues)

	// 混合值测试用例
	mixedTestCases := []struct {
		name     string
		where    string
		expected bool
	}{
		{
			name:     "非NULL字段比较",
			where:    "id = 1",
			expected: true,
		},
		{
			name:     "NULL字段IS NULL",
			where:    "name IS NULL",
			expected: true,
		},
		{
			name:     "非NULL字段IS NULL",
			where:    "id IS NULL",
			expected: false,
		},
		{
			name:     "NULL字段IS NOT NULL",
			where:    "name IS NOT NULL",
			expected: false,
		},
		{
			name:     "非NULL字段IS NOT NULL",
			where:    "id IS NOT NULL",
			expected: true,
		},
		{
			name:     "NULL和非NULL字段组合",
			where:    "id = 1 AND name IS NULL",
			expected: true,
		},
		{
			name:     "多个NULL字段组合",
			where:    "name IS NULL AND salary IS NULL",
			expected: true,
		},
		{
			name:     "多个非NULL字段组合",
			where:    "id = 1 AND age = 30 AND is_active = true",
			expected: true,
		},
		{
			name:     "NULL和非NULL字段OR组合",
			where:    "name IS NULL OR id = 2",
			expected: true,
		},
		{
			name:     "复杂NULL条件",
			where:    "(id = 1 AND name IS NULL) OR (age = 30 AND salary IS NULL)",
			expected: true,
		},
	}

	for _, tt := range mixedTestCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluator.EvaluateWhere(tt.where)

			if err != nil {
				t.Errorf("EvaluateWhere() error = %v", err)
				return
			}

			if got != tt.expected {
				t.Errorf("EvaluateWhere() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestSQLEvaluatorNonPtr 测试非指针类型字段的SQL评估
func TestSQLEvaluatorNonPtr(t *testing.T) {
	tests := []struct {
		name        string
		model       interface{}
		whereClause string
		want        bool
		wantErr     bool
	}{
		{
			name: "非指针字段-简单相等比较",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "张三",
				Age:      25,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "name = '张三' AND age = 25",
			want:        true,
			wantErr:     false,
		},
		{
			name: "非指针字段-数值比较",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "张三",
				Age:      30,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "age > 25 AND salary >= 5000",
			want:        true,
			wantErr:     false,
		},
		{
			name: "非指针字段-布尔值比较",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "张三",
				Age:      25,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "is_active = true",
			want:        true,
			wantErr:     false,
		},
		{
			name: "非指针字段-LIKE操作符",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "张三丰",
				Age:      25,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "name LIKE '张%'",
			want:        true,
			wantErr:     false,
		},
		{
			name: "非指针字段-IN操作符",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "张三",
				Age:      25,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "age IN (20, 25, 30)",
			want:        true,
			wantErr:     false,
		},
		{
			name: "非指针字段-NOT IN操作符",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "张三",
				Age:      25,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "age NOT IN (20, 30, 35)",
			want:        true,
			wantErr:     false,
		},
		{
			name: "非指针字段-BETWEEN操作符",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "张三",
				Age:      25,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "age BETWEEN 20 AND 30",
			want:        true,
			wantErr:     false,
		},
		{
			name: "非指针字段-NOT BETWEEN操作符",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "张三",
				Age:      35,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "age NOT BETWEEN 20 AND 30",
			want:        true,
			wantErr:     false,
		},
		{
			name: "非指针字段-空字符串比较",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "",
				Age:      25,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "name = ''",
			want:        true,
			wantErr:     false,
		},
		{
			name: "非指针字段-空字符串比较",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "",
				Age:      25,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "name IS NULL",
			want:        false,
			wantErr:     false,
		},
		{
			name: "非指针字段-复杂条件",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "张三丰",
				Age:      25,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "(name LIKE '张%' AND age BETWEEN 20 AND 30) OR (salary IN (5000.50, 6000) AND is_active = true)",
			want:        true,
			wantErr:     false,
		},
		{
			name: "非指针字段-类型不匹配错误",
			model: &UserWithNonPtr{
				ID:       1,
				Name:     "张三",
				Age:      25,
				Salary:   5000.50,
				IsActive: true,
			},
			whereClause: "name > 100",
			want:        false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := NewSQLEvaluator(tt.model)
			got, err := evaluator.EvaluateWhere(tt.whereClause)

			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateWhere() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("EvaluateWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}
