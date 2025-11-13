package config

// Environment variable constants
const (
	// Server
	EnvServerPort = "SERVER_PORT"
	EnvServerHost = "SERVER_HOST"
	EnvAppEnv     = "APP_ENV"
	EnvAppName    = "APP_NAME"

	// Database
	EnvDBHost     = "DB_HOST"
	EnvDBPort     = "DB_PORT"
	EnvDBUser     = "DB_USER"
	EnvDBPassword = "DB_PASSWORD"
	EnvDBName     = "DB_NAME"
	EnvDBSSLMode  = "DB_SSLMODE"

	// JWT
	EnvJWTSecret                = "JWT_SECRET"
	EnvJWTExpirationHours       = "JWT_EXPIRATION_HOURS"
	EnvJWTRefreshExpirationDays = "JWT_REFRESH_EXPIRATION_DAYS"

	// Email
	EnvSMTPHost     = "SMTP_HOST"
	EnvSMTPPort     = "SMTP_PORT"
	EnvSMTPUsername = "SMTP_USERNAME"
	EnvSMTPPassword = "SMTP_PASSWORD"
	EnvSMTPFrom     = "SMTP_FROM"

	// Storage
	EnvStorageProvider = "STORAGE_PROVIDER"
	EnvAWSRegion       = "AWS_REGION"
	EnvAWSBucket       = "AWS_BUCKET"
	EnvAWSAccessKey    = "AWS_ACCESS_KEY"
	EnvAWSSecretKey    = "AWS_SECRET_KEY"
	EnvCloudinaryURL   = "CLOUDINARY_URL"
	EnvUploadPath      = "UPLOAD_PATH"

	// Redis
	EnvRedisHost     = "REDIS_HOST"
	EnvRedisPort     = "REDIS_PORT"
	EnvRedisPassword = "REDIS_PASSWORD"
	EnvRedisDB       = "REDIS_DB"

	// Rate Limiting
	EnvRateLimitEnabled       = "RATE_LIMIT_ENABLED"
	EnvRateLimitMax           = "RATE_LIMIT_MAX"
	EnvRateLimitWindowSeconds = "RATE_LIMIT_WINDOW_SECONDS"

	// CORS
	EnvAllowedOrigins = "ALLOWED_ORIGINS"

	// Pagination
	EnvDefaultPageSize = "DEFAULT_PAGE_SIZE"
	EnvMaxPageSize     = "MAX_PAGE_SIZE"
)

// Default values
const (
	DefaultServerPort = "8080"
	DefaultServerHost = "0.0.0.0"
	DefaultAppEnv     = "development"
	DefaultAppName    = "Keerja API"

	DefaultDBHost    = "localhost"
	DefaultDBPort    = "5432"
	DefaultDBUser    = "bekerja"
	DefaultDBName    = "keerja"
	DefaultDBSSLMode = "disable"

	DefaultJWTExpirationHours       = 24
	DefaultJWTRefreshExpirationDays = 7

	DefaultSMTPPort = "587"
	DefaultSMTPFrom = "noreply@keerja.com"

	DefaultStorageProvider = "local"
	DefaultUploadPath      = "./uploads"

	DefaultRedisHost = "localhost"
	DefaultRedisPort = "6379"
	DefaultRedisDB   = 0

	DefaultRateLimitMax           = 100
	DefaultRateLimitWindowSeconds = 60

	DefaultPageSize    = 10
	DefaultMaxPageSize = 100
)

// Application environments
const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

// Storage providers
const (
	StorageLocal      = "local"
	StorageS3         = "s3"
	StorageCloudinary = "cloudinary"
)
