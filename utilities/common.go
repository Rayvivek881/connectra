package utilities

import (
	"cmp"
	"fmt"
	"reflect"
	"strings"
)

func GetFieldValue(v interface{}, fieldName string) string {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return ""
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			jsonTag = strings.Split(jsonTag, ",")[0]
			if jsonTag == fieldName {
				fieldValue := rv.Field(i)
				return fmt.Sprintf("%v", fieldValue.Interface())
			}
		}
		if strings.EqualFold(field.Name, fieldName) {
			fieldValue := rv.Field(i)
			return fmt.Sprintf("%v", fieldValue.Interface())
		}
	}
	return ""
}

func InlineIf(condition bool, trueValue any, falseValue any) any {
	if condition {
		return trueValue
	}
	return falseValue
}

func Max[T cmp.Ordered](a, b T) T {
	return InlineIf(a > b, a, b).(T)
}

func Min[T cmp.Ordered](a, b T) T {
	return InlineIf(a < b, a, b).(T)
}
