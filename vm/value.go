package vm

import (
	"fmt"
	"strconv"
)

type Value struct {
	rawValue interface{}
	// 数字、文本、其他
}

func (value Value) RawValue() (result interface{}) {
	return value.rawValue
}

func (value Value) AsInt() (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
		}
	}()

	int64Value, err := strconv.ParseInt(fmt.Sprintf("%d", value.rawValue), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(int64Value), nil
}

func (value Value) AsInt64() (result int64, err error) {
	defer func() {
		if r := recover();r!=nil{
			err = fmt.Errorf("%+v", r)
		}
	}()

	return strconv.ParseInt(fmt.Sprintf("%d", value.rawValue), 10, 64)
}

func (value Value) AsString() (result string, err error) {
	defer func() {
		r := recover()
		err = fmt.Errorf("%+v", r)
	}()

	if str, ok := value.rawValue.(string); ok {
		return str, nil
	}

	return fmt.Sprintf("%s", value.rawValue), nil
}









