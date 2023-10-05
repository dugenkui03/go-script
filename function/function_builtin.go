package function

import "fmt"

// Add todo 字符串相加标识字符串拼接，所以这里还应该是 interface{}
func Add(args ...int64) int64 {
	result := int64(0)
	for _, arg := range args {
		result = result + arg
	}

	return result
}


func Int64(val interface{}) (int64, error) {
	if val == nil {
		return 0, nil
	}

	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, fmt.Errorf("无法将类型 %T 转换为 int64", val)
	}
}


func Multiplication(args ...int64) int64 {
	// 参数为0，返回1，参考 x^0 = 1
	result := int64(1)
	for _, arg := range args {
		result = result * arg
	}
	return result
}
