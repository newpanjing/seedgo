package pkg

import (
	"errors"
	"reflect"
)

// GetStructFields 获取结构体的所有字段名
// @param obj: 任意结构体实例或结构体指针
// @return: 字段名切片
func GetStructFields(obj any) []string {
	t := reflect.TypeOf(obj)

	// 如果是指针，获取其指向的元素类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 确保是结构体
	if t.Kind() != reflect.Struct {
		return []string{}
	}

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	return fields
}

// SetFieldValue 设置结构体中指定字段的值
// @param obj: 结构体指针 (必须是指针，否则无法修改值)
// @param fieldName: 字段名称
// @param value: 要设置的值
// @return: 如果设置失败返回错误
func SetFieldValue(obj any, fieldName string, value any) error {
	v := reflect.ValueOf(obj)

	// 必须是指针才能修改
	if v.Kind() != reflect.Ptr {
		return errors.New("obj must be a pointer")
	}

	// 获取指针指向的元素
	v = v.Elem()

	// 确保是结构体
	if v.Kind() != reflect.Struct {
		return errors.New("obj must point to a struct")
	}

	// 获取字段
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return errors.New("field not found: " + fieldName)
	}

	// 检查字段是否可设置
	if !field.CanSet() {
		return errors.New("field cannot be set (might be unexported): " + fieldName)
	}

	// 获取要设置的值的反射对象
	val := reflect.ValueOf(value)

	// 类型检查与转换
	if field.Type() != val.Type() {
		// 尝试进行简单的类型转换兼容 (例如 int 到 int64)
		if val.Type().ConvertibleTo(field.Type()) {
			val = val.Convert(field.Type())
		} else {
			return errors.New("value type mismatch for field: " + fieldName)
		}
	}

	field.Set(val)
	return nil
}

// GetFieldValue 获取结构体中指定字段的值
// @param obj: 任意结构体实例或结构体指针
// @param fieldName: 字段名称
// @return: 字段的值，如果出错返回 error
func GetFieldValue(obj any, fieldName string) (any, error) {
	v := reflect.ValueOf(obj)

	// 如果是指针，获取其指向的元素
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 确保是结构体
	if v.Kind() != reflect.Struct {
		return nil, errors.New("obj must be a struct or pointer to struct")
	}

	// 获取字段
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, errors.New("field not found: " + fieldName)
	}

	return field.Interface(), nil
}
