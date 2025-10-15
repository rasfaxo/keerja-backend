package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gosimple/slug"
)

func GenerateSlug(text string) string {
	return slug.Make(text)
}

// GenerateSlugSimple generates a slug without uniqueness checking
func GenerateSlugSimple(text string) string {
	return GenerateSlug(text)
}

// GenerateUniqueSlug generates a unique slug by checking against existing slugs
// checkFunc should return true if the slug already exists
func GenerateUniqueSlug(text string, checkFunc func(string) bool) string {
	baseSlug := GenerateSlug(text)

	// If base slug is unique, return it
	if !checkFunc(baseSlug) {
		return baseSlug
	}

	// Try appending numbers until find a unique slug
	counter := 1
	for {
		newSlug := fmt.Sprintf("%s-%d", baseSlug, counter)
		if !checkFunc(newSlug) {
			return newSlug
		}
		counter++

		// Prevent infinite loop
		if counter > 1000 {
			// Fallback: append timestamp
			return fmt.Sprintf("%s-%d", baseSlug, getCurrentTimestamp())
		}
	}
}

// SanitizeSlug sanitizes a slug by removing invalid characters
func SanitizeSlug(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Replace spaces and underscores with hyphens
	text = strings.ReplaceAll(text, " ", "-")
	text = strings.ReplaceAll(text, "_", "-")

	// Remove all characters except alphanumeric and hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	text = reg.ReplaceAllString(text, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	text = reg.ReplaceAllString(text, "-")

	// Trim hyphens from start and end
	text = strings.Trim(text, "-")

	return text
}

func ValidateSlug(slug string) bool {
	if slug == "" {
		return false
	}

	// Slug should only contain lowercase letters, numbers, and hyphens
	// Should not start or end with hyphen
	// Should not contain consecutive hyphens
	match := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`).MatchString(slug)
	return match
}

func getCurrentTimestamp() int64 {
	return 0 // Will be implemented when needed with time.Now().Unix()
}
