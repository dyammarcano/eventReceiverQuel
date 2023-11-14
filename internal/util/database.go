package util

import (
	"errors"
	"reflect"
	"time"
)

func ReflectStruct(obj any) (reflect.Value, reflect.Type, error) {
	valueOf := reflect.ValueOf(obj)
	typeOf := valueOf.Type()

	if valueOf.Kind() != reflect.Ptr || typeOf.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, nil, errors.New("obj must be a pointer to a struct")
	}

	return valueOf.Elem(), typeOf.Elem(), nil
}

func IsEmptyValue(value any) bool {
	switch value := value.(type) {
	case string:
		return value == ""
	case time.Time:
		return value.IsZero()
	default:
		zeroValue := reflect.Zero(reflect.TypeOf(value)).Interface()
		return value == zeroValue
	}
}

func OnlyDigits(s string) string {
	var res string
	for _, c := range s {
		if c >= '0' && c <= '9' {
			res += string(c)
		}
	}
	return res
}
