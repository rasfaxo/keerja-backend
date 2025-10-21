package utils

import (
	"html"
	"regexp"
	"strings"
)

// SanitizeString removes leading/trailing whitespace and escapes HTML characters
func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	s = html.EscapeString(s)
	return s
}

// SanitizeStringPtr sanitizes string pointer
func SanitizeStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	sanitized := SanitizeString(*s)
	return &sanitized
}

// SanitizeHTML removes potentially dangerous HTML tags while keeping safe formatting
func SanitizeHTML(s string) string {
	// Remove script tags
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	s = scriptRegex.ReplaceAllString(s, "")
	
	// Remove iframe tags
	iframeRegex := regexp.MustCompile(`(?i)<iframe[^>]*>.*?</iframe>`)
	s = iframeRegex.ReplaceAllString(s, "")
	
	// Remove onclick, onerror, and other event handlers
	eventRegex := regexp.MustCompile(`(?i)\son\w+\s*=\s*["'][^"']*["']`)
	s = eventRegex.ReplaceAllString(s, "")
	
	// Remove javascript: protocol
	jsProtocolRegex := regexp.MustCompile(`(?i)javascript:`)
	s = jsProtocolRegex.ReplaceAllString(s, "")
	
	return strings.TrimSpace(s)
}

// SanitizeHTMLPtr sanitizes HTML string pointer
func SanitizeHTMLPtr(s *string) *string {
	if s == nil {
		return nil
	}
	sanitized := SanitizeHTML(*s)
	return &sanitized
}

// SanitizeURL validates and sanitizes URL
func SanitizeURL(url string) string {
	url = strings.TrimSpace(url)
	
	// Check for valid protocols
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return ""
	}
	
	// Remove javascript: protocol
	jsProtocolRegex := regexp.MustCompile(`(?i)javascript:`)
	url = jsProtocolRegex.ReplaceAllString(url, "")
	
	return url
}

// SanitizeURLPtr sanitizes URL string pointer
func SanitizeURLPtr(url *string) *string {
	if url == nil || *url == "" {
		return nil
	}
	sanitized := SanitizeURL(*url)
	if sanitized == "" {
		return nil
	}
	return &sanitized
}

// TrimSpacePtr trims whitespace from string pointer
func TrimSpacePtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

// SanitizeStringArray sanitizes array of strings
func SanitizeStringArray(arr []string) []string {
	result := make([]string, 0, len(arr))
	for _, s := range arr {
		sanitized := SanitizeString(s)
		if sanitized != "" {
			result = append(result, sanitized)
		}
	}
	return result
}

// SanitizeMap sanitizes map of strings
func SanitizeMap(m map[string]string) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		key := SanitizeString(k)
		value := SanitizeString(v)
		if key != "" && value != "" {
			result[key] = value
		}
	}
	return result
}

// StripHTMLTags removes all HTML tags from string
func StripHTMLTags(s string) string {
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	return htmlTagRegex.ReplaceAllString(s, "")
}

// ValidateNoSQLInjection checks for common SQL injection patterns
func ValidateNoSQLInjection(s string) bool {
	// Common SQL injection patterns
	sqlPatterns := []string{
		"(?i)(union.*select)",
		"(?i)(insert.*into)",
		"(?i)(delete.*from)",
		"(?i)(drop.*table)",
		"(?i)(update.*set)",
		"(?i)(exec.*sp_)",
		"(?i)(';|\"--)",
	}
	
	for _, pattern := range sqlPatterns {
		matched, _ := regexp.MatchString(pattern, s)
		if matched {
			return false
		}
	}
	return true
}

// ValidateNoXSS checks for common XSS patterns
func ValidateNoXSS(s string) bool {
	// Common XSS patterns
	xssPatterns := []string{
		"(?i)<script",
		"(?i)javascript:",
		"(?i)onerror=",
		"(?i)onload=",
		"(?i)onclick=",
		"(?i)<iframe",
	}
	
	for _, pattern := range xssPatterns {
		matched, _ := regexp.MatchString(pattern, s)
		if matched {
			return false
		}
	}
	return true
}
