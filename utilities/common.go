package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"vivek-ray/constants"

	"github.com/google/uuid"
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

func GetCleanedString(strValue string) string {
	strValue = strings.TrimFunc(strValue, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	spaceRegex := regexp.MustCompile(`\s+`)
	return spaceRegex.ReplaceAllString(strValue, " ")
}

func GenerateUUID5(value string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(value)).String()
}

func StringToInt64(value string) int64 {
	if value == "" {
		return 0
	}
	normalized := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, value)

	if normalized == "" {
		return 0
	}

	result, err := strconv.ParseInt(normalized, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

func StringToFloat64(value string) float64 {
	if value == "" {
		return 0
	}
	normalized := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == '.' {
			return r
		}
		return -1
	}, value)

	if normalized == "" {
		return 0
	}

	result, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0
	}
	return result
}

func SplitAndTrim(value, sep string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(strings.ToLower(part))
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
