package service

import (
	"fmt"
	"time"
)

// CacheKeys contains all cache key constants and generation functions
const (
	// Master data cache keys with TTL documentation
	JobOptions         = "master:job_options"         // TTL: 3 days (static)
	JobCategories      = "master:job_categories"      // TTL: 7 days (static)
	CategoriesTree     = "master:job_categories_tree" // TTL: 7 days (static)
	SkillsCachePattern = "master:skills:page_%d_%d"   // TTL: 6 hours
)

// GenerateSkillsCacheKey creates a cache key for paginated skills
func GenerateSkillsCacheKey(page, limit int) string {
	return fmt.Sprintf("master:skills:page_%d_%d", page, limit)
}

// GenerateCompanyAddressesCacheKey creates a cache key for user's company addresses
func GenerateCompanyAddressesCacheKey(userID int64) string {
	return fmt.Sprintf("master:company_addresses:%d", userID)
}

// CacheTTLConfig defines TTL values for tiered caching strategy
const (
	StaticDataTTL         = 7 * 24 * time.Hour // 7 days for static admin-managed data (categories)
	RarelyChangingTTL     = 3 * 24 * time.Hour // 3 days for rarely changing data (job types, policies)
	FrequentlyChangingTTL = 6 * time.Hour      // 6 hours for data that changes frequently (skills)
	UserSpecificTTL       = 1 * time.Hour      // 1 hour for user-specific data (company addresses)
)
