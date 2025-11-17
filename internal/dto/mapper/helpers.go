package mapper

import "fmt"

// Helper functions for mapper package

// StringPtr converts string to *string
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// PtrToString converts *string to string
func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Int64Ptr converts int64 to *int64
func Int64Ptr(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

// PtrToInt64 converts *int64 to int64
func PtrToInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// Float64Ptr converts float64 to *float64
func Float64Ptr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

// PtrToFloat64 converts *float64 to float64
func PtrToFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

// BoolPtr converts bool to *bool
func BoolPtr(b bool) *bool {
	return &b
}

// PtrToBool converts *bool to bool
func PtrToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Float64ToString converts float64 to string with 2 decimal places
func Float64ToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}
