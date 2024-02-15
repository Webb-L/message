package model

import (
	"database/sql/driver"
	"errors"
	"strings"
)

// StringArray 是一个自定义类型，表示字符串数组。用于MYSQL不支持字符串数组。
type StringArray []string

// Scan 实现了 sql.Scanner 接口，用于将数据库中的原始数据转换为 StringArray 类型。
func (sa *StringArray) Scan(src interface{}) error {
	var source string
	switch src := src.(type) {
	case []byte:
		source = string(src)
	case string:
		source = src
	default:
		return errors.New("incompatible type for SenderIDs")
	}

	*sa = strings.Split(source, ",")
	return nil
}

// Value 实现了 driver.Valuer 接口，用于将 StringArray 类型转换为数据库中的原始数据。
func (sa StringArray) Value() (driver.Value, error) {
	return strings.Join(sa, ","), nil
}
