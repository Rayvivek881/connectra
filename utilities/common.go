package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"vivek-ray/constants"
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

func AddToBuffer(buf *bytes.Buffer, data any) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	buf.Write(dataBytes)
	buf.WriteByte('\n')
	return nil
}

func InlineIf(condition bool, trueValue any, falseValue any) any {
	if condition {
		return trueValue
	}
	return falseValue
}

func ValidatePageSize(limit int) error {
	if limit > constants.MaxPageSize {
		return constants.PageSizeExceededError
	}
	return nil
}

func ValidateElasticPagination(page, limit int) error {
	if limit > constants.MaxPageSize {
		return constants.PageSizeExceededError
	}
	if page > constants.MaxElasticPageNumber {
		return constants.PageNumberExceededError
	}
	return nil
}

func UniqueStringSlice(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
