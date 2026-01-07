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

func CsvRowToMap(headers, row []string) map[string]string {
	result := make(map[string]string)
	for i, header := range headers {
		result[header] = row[i]
	}
	return result
}

func StructToCsvSlice(v interface{}, columns []string) []string {
	row := make([]string, len(columns))
	for i, col := range columns {
		val := GetFieldValue(v, col)
		if val == nil {
			row[i] = ""
			continue
		}
		switch v := val.(type) {
		case []string:
			row[i] = strings.Join(v, ",")
		default:
			row[i] = fmt.Sprintf("%v", v)
		}
	}
	return row
}

// Validation helpers
var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	urlRegex      = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	linkedinRegex = regexp.MustCompile(`^https?://(www\.)?linkedin\.com/.*$`)
	uuidRegex     = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}
	return nil
}

func ValidateURL(url string, fieldName string) error {
	if url == "" {
		return nil // URL is optional
	}
	if !urlRegex.MatchString(url) {
		return fmt.Errorf("invalid %s URL format: %s", fieldName, url)
	}
	return nil
}

func ValidateLinkedInURL(url string) error {
	if url == "" {
		return nil // LinkedIn URL is optional
	}
	if !linkedinRegex.MatchString(url) {
		return fmt.Errorf("invalid LinkedIn URL format: %s (must start with https://linkedin.com/)", url)
	}
	return nil
}

func ValidateUUID(uuid string) error {
	if uuid == "" {
		return fmt.Errorf("UUID is required")
	}
	if !uuidRegex.MatchString(uuid) {
		return fmt.Errorf("invalid UUID format: %s", uuid)
	}
	return nil
}

func ValidateRequiredString(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func ValidateNonNegativeInt64(value int64, fieldName string) error {
	if value < 0 {
		return fmt.Errorf("%s must be non-negative, got: %d", fieldName, value)
	}
	return nil
}

func ValidateStringLength(value, fieldName string, minLength, maxLength int) error {
	if value == "" {
		return nil // Empty strings are handled by ValidateRequiredString
	}
	length := len(strings.TrimSpace(value))
	if minLength > 0 && length < minLength {
		return fmt.Errorf("%s must be at least %d characters long, got %d", fieldName, minLength, length)
	}
	if maxLength > 0 && length > maxLength {
		return fmt.Errorf("%s must be at most %d characters long, got %d", fieldName, maxLength, length)
	}
	return nil
}

func ValidateStringMaxLength(value, fieldName string, maxLength int) error {
	return ValidateStringLength(value, fieldName, 0, maxLength)
}

func ValidateStringMinLength(value, fieldName string, minLength int) error {
	return ValidateStringLength(value, fieldName, minLength, 0)
}

func ValidatePhoneNumber(phone, fieldName string) error {
	if phone == "" {
		return nil // Phone is optional
	}
	// Basic phone validation: should contain digits and optionally start with +
	cleaned := GetCleanedPhoneNumber(phone)
	if len(cleaned) < 10 {
		return fmt.Errorf("%s must be at least 10 digits long", fieldName)
	}
	if len(cleaned) > 20 {
		return fmt.Errorf("%s must be at most 20 characters long", fieldName)
	}
	return nil
}

func ValidateEmailStatus(status string) error {
	if status == "" {
		return nil // Email status is optional
	}
	validStatuses := []string{"verified", "unverified", "invalid", "bounced"}
	statusLower := strings.ToLower(status)
	for _, valid := range validStatuses {
		if statusLower == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid email_status: %s (must be one of: %v)", status, validStatuses)
}

func ValidateSeniority(seniority string) error {
	if seniority == "" {
		return nil // Seniority is optional
	}
	// Common seniority levels (can be extended)
	validSeniorities := []string{"executive", "director", "manager", "senior", "mid", "junior", "entry", "intern"}
	seniorityLower := strings.ToLower(seniority)
	for _, valid := range validSeniorities {
		if seniorityLower == valid {
			return nil
		}
	}
	// Allow custom seniority values but warn if they don't match common ones
	// For now, we'll just validate format (non-empty, reasonable length)
	if len(seniority) > 50 {
		return fmt.Errorf("seniority must be at most 50 characters long")
	}
	return nil
}
