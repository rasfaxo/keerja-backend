package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"keerja-backend/internal/handler/http/master"
)

// MasterDataHandlers holds all master data handlers
type MasterDataHandlers struct {
	IndustryHandler    *master.IndustryHandler
	CompanySizeHandler *master.CompanySizeHandler
	LocationHandler    *master.LocationHandler
}

// SetupMasterDataRoutes sets up routes for all master data endpoints
// Master data includes: industries, company sizes, and locations (provinces, cities, districts)
// These are public endpoints (no authentication required) with caching and rate limiting
func SetupMasterDataRoutes(api fiber.Router, handlers *MasterDataHandlers) {
	// Master data group - /api/v1/master
	master := api.Group("/master")

	// Apply middleware to all master data routes
	master.Use(setupMasterDataMiddleware())

	// Industry routes
	setupIndustryRoutes(master, handlers.IndustryHandler)

	// Company size routes
	setupCompanySizeRoutes(master, handlers.CompanySizeHandler)

	// Location routes (provinces, cities, districts)
	setupLocationRoutes(master, handlers.LocationHandler)
}

// setupMasterDataMiddleware configures middleware for master data endpoints
func setupMasterDataMiddleware() fiber.Handler {
	// Cache configuration for master data
	// Master data changes infrequently, so we can cache responses
	return cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			// Skip cache for non-GET requests
			return c.Method() != fiber.MethodGet
		},
		Expiration:   5 * time.Minute, // Cache for 5 minutes
		CacheControl: true,            // Add Cache-Control header
		KeyGenerator: func(c *fiber.Ctx) string {
			// Generate cache key from URL and query params
			return c.Path() + "?" + string(c.Request().URI().QueryString())
		},
		Storage: nil, // Use in-memory cache (Fiber default)
	})
}

// setupRateLimiter creates a rate limiter for master data endpoints
func setupRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,             // 100 requests
		Expiration: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			// Rate limit by IP address
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Rate limit exceeded. Please try again later.",
			})
		},
	})
}

// setupIndustryRoutes configures routes for industry master data
func setupIndustryRoutes(master fiber.Router, handler *master.IndustryHandler) {
	industries := master.Group("/industries")

	// Apply rate limiting
	industries.Use(setupRateLimiter())

	// GET /api/v1/master/industries - Get all industries
	// Query params: ?active=true, ?search=tech
	industries.Get("/", handler.GetAllIndustries)

	// GET /api/v1/master/industries/:id - Get industry by ID
	industries.Get("/:id", handler.GetIndustryByID)
}

// setupCompanySizeRoutes configures routes for company size master data
func setupCompanySizeRoutes(master fiber.Router, handler *master.CompanySizeHandler) {
	companySizes := master.Group("/company-sizes")

	// Apply rate limiting
	companySizes.Use(setupRateLimiter())

	// GET /api/v1/master/company-sizes - Get all company sizes
	// Query params: ?active=true
	companySizes.Get("/", handler.GetAllCompanySizes)

	// GET /api/v1/master/company-sizes/:id - Get company size by ID
	companySizes.Get("/:id", handler.GetCompanySizeByID)
}

// setupLocationRoutes configures routes for location hierarchy
// Includes provinces, cities, and districts
func setupLocationRoutes(master fiber.Router, handler *master.LocationHandler) {
	locations := master.Group("/locations")

	// Apply rate limiting
	locations.Use(setupRateLimiter())

	// Province routes
	// GET /api/v1/master/locations/provinces - Get all provinces
	// Query params: ?active=true, ?search=jawa
	locations.Get("/provinces", handler.GetAllProvinces)

	// GET /api/v1/master/locations/provinces/:id - Get province by ID
	locations.Get("/provinces/:id", handler.GetProvinceByID)

	// City routes
	// GET /api/v1/master/locations/cities - Get cities by province
	// Query params: ?province_id=32 (required), ?active=true, ?search=bandung
	locations.Get("/cities", handler.GetCities)

	// GET /api/v1/master/locations/cities/:id - Get city by ID with province
	locations.Get("/cities/:id", handler.GetCityByID)

	// District routes
	// GET /api/v1/master/locations/districts - Get districts by city
	// Query params: ?city_id=150 (required), ?active=true, ?search=batu
	locations.Get("/districts", handler.GetDistricts)

	// GET /api/v1/master/locations/districts/:id - Get district by ID with full hierarchy
	locations.Get("/districts/:id", handler.GetDistrictByID)
}
