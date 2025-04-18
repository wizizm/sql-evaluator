package sqlevaluator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/xwb1989/sqlparser"
)

// SQLEvaluator SQL评估器
type SQLEvaluator struct {
	model interface{}
}

// NewSQLEvaluator 创建新的SQL评估器
func NewSQLEvaluator(model interface{}) *SQLEvaluator {
	return &SQLEvaluator{
		model: model,
	}
}

// EvaluateWhere 评估WHERE子句
func (e *SQLEvaluator) EvaluateWhere(whereClause string) (bool, error) {
	// 解析SQL
	stmt, err := sqlparser.Parse("SELECT * FROM `users` WHERE " + whereClause)
	if err != nil {
		return false, fmt.Errorf("解析SQL失败: %v", err)
	}

	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		return false, fmt.Errorf("不是SELECT语句")
	}

	if selectStmt.Where == nil {
		return true, nil
	}

	return e.evaluateExpr(selectStmt.Where.Expr)
}

// evaluateExpr 评估表达式
func (e *SQLEvaluator) evaluateExpr(expr sqlparser.Expr) (bool, error) {
	switch node := expr.(type) {
	case *sqlparser.ComparisonExpr:
		if node.Operator == sqlparser.InStr || node.Operator == sqlparser.NotInStr {
			return e.evaluateInExpr(node)
		}
		return e.evaluateComparison(node)
	case *sqlparser.AndExpr:
		left, err := e.evaluateExpr(node.Left)
		if err != nil {
			return false, err
		}
		right, err := e.evaluateExpr(node.Right)
		if err != nil {
			return false, err
		}
		return left && right, nil
	case *sqlparser.OrExpr:
		left, err := e.evaluateExpr(node.Left)
		if err != nil {
			return false, err
		}
		right, err := e.evaluateExpr(node.Right)
		if err != nil {
			return false, err
		}
		return left || right, nil
	case *sqlparser.ParenExpr:
		return e.evaluateExpr(node.Expr)
	case *sqlparser.RangeCond:
		return e.evaluateRange(node)
	case *sqlparser.IsExpr:
		return e.evaluateIsExpr(node)
	default:
		return false, fmt.Errorf("不支持的表达式类型: %T", expr)
	}
}

// evaluateComparison 评估比较表达式
func (e *SQLEvaluator) evaluateComparison(expr *sqlparser.ComparisonExpr) (bool, error) {
	// 获取左操作数的字段名
	leftField, err := e.getFieldName(expr.Left)
	if err != nil {
		return false, err
	}

	// 获取左操作数的值
	leftVal, err := e.getFieldValue(leftField)
	if err != nil {
		return false, err
	}

	// 处理IS NULL和IS NOT NULL
	if sqlVal, ok := expr.Right.(*sqlparser.SQLVal); ok && sqlVal.Type == 0 {
		switch expr.Operator {
		case "is null":
			return leftVal == nil, nil
		case "is not null":
			return leftVal != nil, nil
		}
	}

	// 如果左操作数为NULL，且不是IS NULL或IS NOT NULL操作，则返回false
	if leftVal == nil {
		return false, nil
	}

	// 处理IN和NOT IN操作符
	if expr.Operator == sqlparser.InStr || expr.Operator == sqlparser.NotInStr {
		// 获取IN列表的值
		values, err := e.getSQLValues(expr.Right)
		if err != nil {
			return false, err
		}

		// 检查左操作数是否在值列表中
		found := false
		for _, val := range values {
			// 如果值为NULL，跳过
			if val == nil {
				continue
			}

			// 尝试类型转换
			leftConverted, rightConverted, err := convertTypes(leftVal, val)
			if err == nil && reflect.DeepEqual(leftConverted, rightConverted) {
				found = true
				break
			}
		}

		// 对于NOT IN，如果找到匹配项则返回false，否则返回true
		if expr.Operator == sqlparser.NotInStr {
			return !found, nil
		}
		// 对于IN，如果找到匹配项则返回true，否则返回false
		return found, nil
	}

	// 处理BETWEEN和NOT BETWEEN操作符
	if expr.Operator == sqlparser.BetweenStr || expr.Operator == sqlparser.NotBetweenStr {
		// 获取范围值
		values, err := e.getSQLValues(expr.Right)
		if err != nil {
			return false, err
		}

		if len(values) != 2 {
			return false, fmt.Errorf("BETWEEN操作符需要两个值")
		}

		// 确保所有值都是相同类型
		lower := values[0]
		upper := values[1]

		// 如果范围值为NULL，则返回false
		if lower == nil || upper == nil {
			return false, nil
		}

		// 尝试类型转换
		leftLower, lowerConverted, err := convertTypes(leftVal, lower)
		if err != nil {
			return false, err
		}

		leftUpper, upperConverted, err := convertTypes(leftVal, upper)
		if err != nil {
			return false, err
		}

		// 比较值是否在范围内
		fromResult, err := compareValues(leftLower, lowerConverted, func(a, b float64) bool { return a >= b })
		if err != nil {
			return false, err
		}

		toResult, err := compareValues(leftUpper, upperConverted, func(a, b float64) bool { return a <= b })
		if err != nil {
			return false, err
		}

		result := fromResult && toResult
		// 对于NOT BETWEEN，如果不在范围内则返回true，否则返回false
		if expr.Operator == sqlparser.NotBetweenStr {
			return !result, nil
		}
		// 对于BETWEEN，如果在范围内则返回true，否则返回false
		return result, nil
	}

	// 获取右操作数的值
	rightVal, err := e.getValue(expr.Right)
	if err != nil {
		return false, err
	}

	// 如果右操作数为NULL，且不是IS NULL或IS NOT NULL操作，则返回false
	if rightVal == nil {
		return false, nil
	}

	// 尝试类型转换
	leftConverted, rightConverted, err := convertTypes(leftVal, rightVal)
	if err != nil {
		return false, err
	}

	// 根据操作符进行比较
	switch expr.Operator {
	case "=":
		return reflect.DeepEqual(leftConverted, rightConverted), nil
	case "!=", "<>":
		return !reflect.DeepEqual(leftConverted, rightConverted), nil
	case ">":
		return compareValues(leftConverted, rightConverted, func(a, b float64) bool { return a > b })
	case ">=":
		return compareValues(leftConverted, rightConverted, func(a, b float64) bool { return a >= b })
	case "<":
		return compareValues(leftConverted, rightConverted, func(a, b float64) bool { return a < b })
	case "<=":
		return compareValues(leftConverted, rightConverted, func(a, b float64) bool { return a <= b })
	case "like":
		return e.evaluateLike(leftConverted, rightConverted)
	case "not like":
		result, err := e.evaluateLike(leftConverted, rightConverted)
		if err != nil {
			return false, err
		}
		return !result, nil
	default:
		return false, fmt.Errorf("不支持的操作符: %s", expr.Operator)
	}
}

// evaluateLike 评估LIKE操作符
func (e *SQLEvaluator) evaluateLike(left, right interface{}) (bool, error) {
	leftStr, ok := left.(string)
	if !ok {
		return false, fmt.Errorf("LIKE操作符的左侧必须是字符串类型")
	}

	rightStr, ok := right.(string)
	if !ok {
		return false, fmt.Errorf("LIKE操作符的右侧必须是字符串类型")
	}

	// 将SQL LIKE模式转换为正则表达式
	pattern := strings.ReplaceAll(rightStr, "%", ".*")
	pattern = strings.ReplaceAll(pattern, "_", ".")
	pattern = "^" + pattern + "$"

	matched, err := regexp.MatchString(pattern, leftStr)
	if err != nil {
		return false, fmt.Errorf("正则表达式匹配失败: %v", err)
	}

	return matched, nil
}

// evaluateInExpr 评估IN表达式
func (e *SQLEvaluator) evaluateInExpr(expr *sqlparser.ComparisonExpr) (bool, error) {
	// 获取左操作数的字段名
	leftField, err := e.getFieldName(expr.Left)
	if err != nil {
		return false, err
	}

	// 获取左操作数的值
	leftVal, err := e.getFieldValue(leftField)
	if err != nil {
		return false, err
	}

	// 如果左操作数为NULL，则返回false
	if leftVal == nil {
		return false, nil
	}

	// 获取IN列表的值
	values, err := e.getSQLValues(expr.Right)
	if err != nil {
		return false, err
	}

	// 检查左操作数是否在值列表中
	found := false
	for _, val := range values {
		// 如果值为NULL，跳过
		if val == nil {
			continue
		}

		// 尝试类型转换
		leftConverted, rightConverted, err := convertTypes(leftVal, val)
		if err == nil && reflect.DeepEqual(leftConverted, rightConverted) {
			found = true
			break
		}
	}

	// 对于NOT IN，如果找到匹配项则返回false，否则返回true
	if expr.Operator == sqlparser.NotInStr {
		return !found, nil
	}
	// 对于IN，如果找到匹配项则返回true，否则返回false
	return found, nil
}

// evaluateIsExpr 评估IS表达式
func (e *SQLEvaluator) evaluateIsExpr(expr *sqlparser.IsExpr) (bool, error) {
	// 获取左操作数的字段名
	leftField, err := e.getFieldName(expr.Expr)
	if err != nil {
		return false, err
	}

	// 获取左操作数的值
	leftVal, err := e.getFieldValue(leftField)
	if err != nil {
		return false, err
	}

	// 根据操作符返回结果
	switch expr.Operator {
	case "is null":
		return leftVal == nil, nil
	case "is not null":
		return leftVal != nil, nil
	default:
		return false, fmt.Errorf("不支持的IS操作符: %s", expr.Operator)
	}
}

// evaluateBetween 评估BETWEEN操作符
func (e *SQLEvaluator) evaluateBetween(left, right interface{}) (bool, error) {
	values, ok := right.([]interface{})
	if !ok || len(values) != 2 {
		return false, fmt.Errorf("BETWEEN操作符需要两个值")
	}

	lower := values[0]
	upper := values[1]

	// 确保所有值都是相同类型
	switch v := left.(type) {
	case int:
		lowerInt, ok1 := lower.(int)
		upperInt, ok2 := upper.(int)
		if !ok1 || !ok2 {
			return false, fmt.Errorf("BETWEEN操作符的值类型必须一致")
		}
		return v >= lowerInt && v <= upperInt, nil
	case float64:
		lowerFloat, ok1 := lower.(float64)
		upperFloat, ok2 := upper.(float64)
		if !ok1 || !ok2 {
			return false, fmt.Errorf("BETWEEN操作符的值类型必须一致")
		}
		return v >= lowerFloat && v <= upperFloat, nil
	case string:
		lowerStr, ok1 := lower.(string)
		upperStr, ok2 := upper.(string)
		if !ok1 || !ok2 {
			return false, fmt.Errorf("BETWEEN操作符的值类型必须一致")
		}
		return v >= lowerStr && v <= upperStr, nil
	default:
		return false, fmt.Errorf("不支持的BETWEEN值类型: %T", left)
	}
}

// evaluateRange 评估范围条件（BETWEEN）
func (e *SQLEvaluator) evaluateRange(expr *sqlparser.RangeCond) (bool, error) {
	// 获取左操作数的字段名
	leftField, err := e.getFieldName(expr.Left)
	if err != nil {
		return false, err
	}

	// 获取左操作数的值
	leftVal, err := e.getFieldValue(leftField)
	if err != nil {
		return false, err
	}

	// 如果左操作数为NULL，则返回false
	if leftVal == nil {
		return false, nil
	}

	// 获取范围的最小值和最大值
	fromVal, err := e.getValue(expr.From)
	if err != nil {
		return false, err
	}

	toVal, err := e.getValue(expr.To)
	if err != nil {
		return false, err
	}

	// 如果范围值为NULL，则返回false
	if fromVal == nil || toVal == nil {
		return false, nil
	}

	// 尝试类型转换
	leftLower, lowerConverted, err := convertTypes(leftVal, fromVal)
	if err != nil {
		return false, err
	}

	leftUpper, upperConverted, err := convertTypes(leftVal, toVal)
	if err != nil {
		return false, err
	}

	// 比较值是否在范围内
	fromResult, err := compareValues(leftLower, lowerConverted, func(a, b float64) bool { return a >= b })
	if err != nil {
		return false, err
	}

	toResult, err := compareValues(leftUpper, upperConverted, func(a, b float64) bool { return a <= b })
	if err != nil {
		return false, err
	}

	result := fromResult && toResult
	// 对于NOT BETWEEN，如果不在范围内则返回true，否则返回false
	if expr.Operator == sqlparser.NotBetweenStr {
		return !result, nil
	}
	// 对于BETWEEN，如果在范围内则返回true，否则返回false
	return result, nil
}

// getValue 获取表达式的值
func (e *SQLEvaluator) getValue(expr sqlparser.Expr) (interface{}, error) {
	switch node := expr.(type) {
	case *sqlparser.ColName:
		return e.getFieldValue(node.Name.String())
	case *sqlparser.SQLVal:
		switch node.Type {
		case sqlparser.StrVal:
			return string(node.Val), nil
		case sqlparser.IntVal:
			// 将字符串转换为整数
			var val int
			_, err := fmt.Sscanf(string(node.Val), "%d", &val)
			if err != nil {
				return nil, fmt.Errorf("无法解析整数值: %v", err)
			}
			return val, nil
		case sqlparser.FloatVal:
			// 将字符串转换为浮点数
			var val float64
			_, err := fmt.Sscanf(string(node.Val), "%f", &val)
			if err != nil {
				return nil, fmt.Errorf("无法解析浮点数值: %v", err)
			}
			return val, nil
		case sqlparser.ValArg:
			return string(node.Val), nil
		default:
			return nil, fmt.Errorf("不支持的SQL值类型: %v", node.Type)
		}
	case sqlparser.BoolVal:
		return bool(node), nil
	case *sqlparser.NullVal:
		return nil, nil
	default:
		return nil, fmt.Errorf("不支持的表达式类型: %T", expr)
	}
}

// getFieldValue 获取字段值
func (e *SQLEvaluator) getFieldValue(fieldName string) (interface{}, error) {
	// 获取字段的反射值
	modelValue := reflect.ValueOf(e.model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	field := modelValue.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	// 处理指针类型
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return nil, nil
		}
		// 获取指针指向的值
		elemValue := field.Elem()
		switch elemValue.Kind() {
		case reflect.Int:
			return int(elemValue.Int()), nil
		case reflect.Float64:
			return elemValue.Float(), nil
		case reflect.String:
			return elemValue.String(), nil
		case reflect.Bool:
			return elemValue.Bool(), nil
		default:
			return elemValue.Interface(), nil
		}
	}

	return field.Interface(), nil
}

// compareValues 比较两个值
func compareValues(a, b interface{}, compare func(float64, float64) bool) (bool, error) {
	switch v1 := a.(type) {
	case int:
		v2, ok := b.(int)
		if !ok {
			return false, fmt.Errorf("类型不匹配: %T 和 %T", a, b)
		}
		return compare(float64(v1), float64(v2)), nil
	case float64:
		v2, ok := b.(float64)
		if !ok {
			return false, fmt.Errorf("类型不匹配: %T 和 %T", a, b)
		}
		return compare(v1, v2), nil
	case string:
		v2, ok := b.(string)
		if !ok {
			return false, fmt.Errorf("类型不匹配: %T 和 %T", a, b)
		}
		return compare(float64(strings.Compare(v1, v2)), 0), nil
	case bool:
		v2, ok := b.(bool)
		if !ok {
			return false, fmt.Errorf("类型不匹配: %T 和 %T", a, b)
		}
		if v1 == v2 {
			return true, nil
		}
		return false, nil
	default:
		return false, fmt.Errorf("不支持的比较类型: %T", a)
	}
}

// isTypeCompatible 检查两个值是否类型兼容
func isTypeCompatible(a, b interface{}) bool {
	if a == nil || b == nil {
		return true
	}

	switch a.(type) {
	case string:
		_, ok := b.(string)
		return ok
	case int:
		_, ok := b.(int)
		return ok
	case float64:
		_, ok := b.(float64)
		return ok
	case bool:
		_, ok := b.(bool)
		return ok
	default:
		return false
	}
}

// convertTypes 尝试转换类型使其兼容
func convertTypes(a, b interface{}) (interface{}, interface{}, error) {
	// 如果任一值为nil，直接返回
	if a == nil || b == nil {
		return a, b, nil
	}

	// 尝试将int转换为float64
	switch v1 := a.(type) {
	case int:
		switch v2 := b.(type) {
		case float64:
			return float64(v1), v2, nil
		case int:
			return v1, v2, nil
		case string:
			// 尝试将字符串转换为数字
			if f, err := strconv.ParseFloat(v2, 64); err == nil {
				return float64(v1), f, nil
			}
		}
	case float64:
		switch v2 := b.(type) {
		case int:
			return v1, float64(v2), nil
		case float64:
			return v1, v2, nil
		case string:
			// 尝试将字符串转换为数字
			if f, err := strconv.ParseFloat(v2, 64); err == nil {
				return v1, f, nil
			}
		}
	case string:
		switch v2 := b.(type) {
		case string:
			return v1, v2, nil
		case int:
			// 尝试将字符串转换为数字
			if f, err := strconv.ParseFloat(v1, 64); err == nil {
				return f, float64(v2), nil
			}
		case float64:
			// 尝试将字符串转换为数字
			if f, err := strconv.ParseFloat(v1, 64); err == nil {
				return f, v2, nil
			}
		}
	case bool:
		switch v2 := b.(type) {
		case bool:
			return v1, v2, nil
		case string:
			// 尝试将字符串转换为布尔值
			if b, err := strconv.ParseBool(v2); err == nil {
				return v1, b, nil
			}
		}
	}

	// 如果类型相同，直接返回
	if reflect.TypeOf(a) == reflect.TypeOf(b) {
		return a, b, nil
	}

	return a, b, fmt.Errorf("无法转换类型: %T 和 %T", a, b)
}

// getFieldName 从SQL表达式中获取字段名
func (e *SQLEvaluator) getFieldName(expr sqlparser.Expr) (string, error) {
	switch v := expr.(type) {
	case *sqlparser.ColName:
		// 将SQL字段名转换为Go结构体字段名
		sqlName := v.Name.String()
		parts := strings.Split(sqlName, "_")
		for i := range parts {
			parts[i] = strings.Title(parts[i])
		}
		fieldName := strings.Join(parts, "")

		// 检查字段是否存在
		modelValue := reflect.ValueOf(e.model)
		if modelValue.Kind() == reflect.Ptr {
			modelValue = modelValue.Elem()
		}

		// 首先尝试直接使用转换后的字段名
		field := modelValue.FieldByName(fieldName)
		if field.IsValid() {
			return fieldName, nil
		}

		// 如果找不到，尝试使用json标签
		modelType := modelValue.Type()
		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			if tag := field.Tag.Get("json"); tag == sqlName {
				return field.Name, nil
			}
		}

		return "", fmt.Errorf("field %s not found", sqlName)
	default:
		return "", fmt.Errorf("unsupported expression type for field name: %T", expr)
	}
}

// getSQLValues 获取SQL值列表
func (e *SQLEvaluator) getSQLValues(expr sqlparser.Expr) ([]interface{}, error) {
	switch node := expr.(type) {
	case sqlparser.ValTuple:
		values := make([]interface{}, len(node))
		for i, val := range node {
			value, err := e.getValue(val)
			if err != nil {
				return nil, err
			}
			values[i] = value
		}
		return values, nil
	default:
		return nil, fmt.Errorf("不支持的表达式类型: %T", expr)
	}
}
