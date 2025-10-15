package utils

// StringPtr returns a pointer to the provided string value
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to the provided bool value
func BoolPtr(b bool) *bool {
	return &b
}

// Int64Ptr returns a pointer to the provided int64 value
func Int64Ptr(i int64) *int64 {
	return &i
}

// IntPtr returns a pointer to the provided int value
func IntPtr(i int) *int {
	return &i
}

// Float64Ptr returns a pointer to the provided float64 value
func Float64Ptr(f float64) *float64 {
	return &f
}

// Int16Ptr returns a pointer to the provided int16 value
func Int16Ptr(i int16) *int16 {
	return &i
}

// StringValue returns the value of a string pointer or empty string if nil
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BoolValue returns the value of a bool pointer or false if nil
func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Int64Value returns the value of an int64 pointer or 0 if nil
func Int64Value(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// Float64Value returns the value of a float64 pointer or 0.0 if nil
func Float64Value(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}
