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

var spaceRegex = regexp.MustCompile(`\s+`) // regex to match one or more spaces

func GetFieldValue(v interface{}, fieldName string) any {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			jsonTag = strings.Split(jsonTag, ",")[0]
			if jsonTag == fieldName {
				return rv.Field(i).Interface()
			}
		}
		if strings.EqualFold(field.Name, fieldName) {
			return rv.Field(i).Interface()
		}
	}
	return nil
}

func ToStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		result := make([]string, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i).Interface()
			if s, ok := elem.(string); ok {
				result = append(result, s)
			} else {
				result = append(result, fmt.Sprintf("%v", elem))
			}
		}
		return result
	}

	return []string{fmt.Sprintf("%v", v)}
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
	seen := make(map[string]struct{})
	result := make([]string, 0)
	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func GetCleanedString(strValue string) string {
	strValue = strings.TrimFunc(strValue, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

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

func SplitAndTrim(value, sep string) []string {
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

func GetCleanedPhoneNumber(phoneNumber string) string {
	var result strings.Builder
	for i, r := range phoneNumber {
		if r == '+' && i == 0 {
			result.WriteRune(r)
		} else if r >= '0' && r <= '9' {
			result.WriteRune(r)
		}
	}

	cleaned := result.String()
	if !strings.HasPrefix(cleaned, "+") && cleaned != "" {
		cleaned = "+" + cleaned
	}
	return cleaned
}
