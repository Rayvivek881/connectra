package utilities

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "valid email with plus",
			email:   "user+tag@example.com",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "invalid email - no @",
			email:   "invalidemail.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no domain",
			email:   "user@",
			wantErr: true,
		},
		{
			name:    "invalid email - no TLD",
			email:   "user@example",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		fieldName string
		wantErr   bool
	}{
		{
			name:      "valid HTTP URL",
			url:       "http://example.com",
			fieldName: "website",
			wantErr:   false,
		},
		{
			name:      "valid HTTPS URL",
			url:       "https://example.com/path",
			fieldName: "website",
			wantErr:   false,
		},
		{
			name:      "empty URL (optional)",
			url:       "",
			fieldName: "website",
			wantErr:   false,
		},
		{
			name:      "invalid URL - no protocol",
			url:       "example.com",
			fieldName: "website",
			wantErr:   true,
		},
		{
			name:      "invalid URL - spaces",
			url:       "https://example.com/path with spaces",
			fieldName: "website",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLinkedInURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid LinkedIn URL",
			url:     "https://linkedin.com/in/username",
			wantErr: false,
		},
		{
			name:    "valid LinkedIn URL with www",
			url:     "https://www.linkedin.com/company/companyname",
			wantErr: false,
		},
		{
			name:    "valid LinkedIn URL HTTP",
			url:     "http://linkedin.com/in/username",
			wantErr: false,
		},
		{
			name:    "empty URL (optional)",
			url:     "",
			wantErr: false,
		},
		{
			name:    "invalid LinkedIn URL - wrong domain",
			url:     "https://example.com/in/username",
			wantErr: true,
		},
		{
			name:    "invalid LinkedIn URL - no protocol",
			url:     "linkedin.com/in/username",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLinkedInURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLinkedInURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name    string
		uuid    string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			uuid:    "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "valid UUID uppercase",
			uuid:    "550E8400-E29B-41D4-A716-446655440000",
			wantErr: true, // Should be lowercase
		},
		{
			name:    "empty UUID",
			uuid:    "",
			wantErr: true,
		},
		{
			name:    "invalid UUID - wrong format",
			uuid:    "not-a-uuid",
			wantErr: true,
		},
		{
			name:    "invalid UUID - too short",
			uuid:    "550e8400-e29b-41d4-a716",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUUID(tt.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRequiredString(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		wantErr   bool
	}{
		{
			name:      "non-empty string",
			value:     "test value",
			fieldName: "name",
			wantErr:   false,
		},
		{
			name:      "empty string",
			value:     "",
			fieldName: "name",
			wantErr:   true,
		},
		{
			name:      "whitespace only",
			value:     "   ",
			fieldName: "name",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequiredString(tt.value, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequiredString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNonNegativeInt64(t *testing.T) {
	tests := []struct {
		name      string
		value     int64
		fieldName string
		wantErr   bool
	}{
		{
			name:      "zero value",
			value:     0,
			fieldName: "count",
			wantErr:   false,
		},
		{
			name:      "positive value",
			value:     100,
			fieldName: "count",
			wantErr:   false,
		},
		{
			name:      "negative value",
			value:     -1,
			fieldName: "count",
			wantErr:   true,
		},
		{
			name:      "large positive value",
			value:     999999999,
			fieldName: "count",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonNegativeInt64(tt.value, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNonNegativeInt64() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateStringLength(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		minLength int
		maxLength int
		wantErr   bool
	}{
		{
			name:      "valid length",
			value:     "test",
			fieldName: "name",
			minLength: 1,
			maxLength: 10,
			wantErr:   false,
		},
		{
			name:      "too short",
			value:     "a",
			fieldName: "name",
			minLength: 2,
			maxLength: 10,
			wantErr:   true,
		},
		{
			name:      "too long",
			value:     "this is a very long string",
			fieldName: "name",
			minLength: 1,
			maxLength: 10,
			wantErr:   true,
		},
		{
			name:      "exact min length",
			value:     "a",
			fieldName: "name",
			minLength: 1,
			maxLength: 10,
			wantErr:   false,
		},
		{
			name:      "exact max length",
			value:     "1234567890",
			fieldName: "name",
			minLength: 1,
			maxLength: 10,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStringLength(tt.value, tt.fieldName, tt.minLength, tt.maxLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStringLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name      string
		phone     string
		fieldName string
		wantErr   bool
	}{
		{
			name:      "valid phone number",
			phone:     "+1234567890",
			fieldName: "mobile_phone",
			wantErr:   false,
		},
		{
			name:      "valid phone with dashes",
			phone:     "+1-234-567-8900",
			fieldName: "mobile_phone",
			wantErr:   false,
		},
		{
			name:      "valid phone with spaces",
			phone:     "+1 234 567 8900",
			fieldName: "mobile_phone",
			wantErr:   false,
		},
		{
			name:      "empty phone (optional)",
			phone:     "",
			fieldName: "mobile_phone",
			wantErr:   false,
		},
		{
			name:      "invalid phone - too short",
			phone:     "123",
			fieldName: "mobile_phone",
			wantErr:   true,
		},
		{
			name:      "invalid phone - letters",
			phone:     "abc123",
			fieldName: "mobile_phone",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneNumber(tt.phone, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCleanedPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "phone with dashes",
			input:    "+1-234-567-8900",
			expected: "+12345678900",
		},
		{
			name:     "phone with spaces",
			input:    "+1 234 567 8900",
			expected: "+12345678900",
		},
		{
			name:     "phone with parentheses",
			input:    "+1 (234) 567-8900",
			expected: "+12345678900",
		},
		{
			name:     "already clean",
			input:    "+12345678900",
			expected: "+12345678900",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCleanedPhoneNumber(tt.input)
			if result != tt.expected {
				t.Errorf("GetCleanedPhoneNumber() = %v, want %v", result, tt.expected)
			}
		})
	}
}
