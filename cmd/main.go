package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/config"
	"keerja-backend/internal/handler/http"
	"keerja-backend/internal/jobs"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/repository/postgres"
	"keerja-backend/internal/routes"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"

	_ "keerja-backend/docs" // Import generated docs

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Keerja Backend API
// @version 1.0
// @description Job platform backend API with authentication, job management, and push notifications
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@keerja.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize logger
	config.InitLogger()
	appLogger := config.GetLogger()
	appLogger.Info("Starting Keerja Backend API...")

	// Load configuration
	cfg := config.LoadConfig()
	appLogger.Info(fmt.Sprintf("Environment: %s", cfg.AppEnv))

	// Initialize database
	appLogger.Info("Initializing database connection...")
	if err := config.InitDatabase(); err != nil {
		appLogger.WithError(err).Fatal("Failed to initialize database")
	}
	defer config.CloseDatabase()
	db := config.GetDB()
	appLogger.Info("✓ Database connected successfully")

	// Initialize FCM (Firebase Cloud Messaging)
	if cfg.FCMEnabled {
		appLogger.Info("Initializing Firebase Cloud Messaging...")
		if err := config.InitFCM(cfg); err != nil {
			appLogger.WithError(err).Error("Failed to initialize FCM - push notifications will be disabled")
		} else {
			appLogger.Info("✓ FCM initialized successfully")
		}
		defer config.CloseFCM()
	} else {
		appLogger.Info("FCM is disabled - push notifications will not be sent")
	}

	// Initialize validator
	utils.InitValidator()

	// Initialize repositories
	appLogger.Info("Initializing repositories...")
	userRepo := postgres.NewUserRepository(db)
	emailRepo := postgres.NewEmailRepository(db)
	companyRepo := postgres.NewCompanyRepository(db)
	jobRepo := postgres.NewJobRepository(db)
	applicationRepo := postgres.NewApplicationRepository(db)
	skillsMasterRepo := postgres.NewSkillsMasterRepository(db)
	oauthRepo := postgres.NewOAuthRepository(db)
	otpCodeRepo := postgres.NewOTPCodeRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db)
	deviceTokenRepo := postgres.NewDeviceTokenRepository(db) // FCM device tokens

	// Initialize services
	appLogger.Info("Initializing services...")
	tokenStore := service.NewInMemoryTokenStore()

	// Initialize email service with config
	emailService := service.NewEmailService(emailRepo, cfg)

	// Initialize FCM push service
	fcmPushService := service.NewFCMPushService(deviceTokenRepo, cfg)

	// Initialize notification service with FCM support
	notificationService := service.NewNotificationService(notificationRepo, fcmPushService, emailService)

	// Initialize upload service
	uploadConfig := service.UploadServiceConfig{
		StorageProvider: "local",
		UploadPath:      "./uploads",
		BaseURL:         "http://localhost:8080", // TODO: Get from config
	}
	uploadService := service.NewUploadService(uploadConfig)

	authServiceConfig := service.AuthServiceConfig{
		JWTSecret:   cfg.JWTSecret,
		JWTDuration: time.Duration(cfg.JWTExpirationHours) * time.Hour,
	}

	// Initialize cache
	appLogger.Info("Initializing cache...")
	var cacheService cache.Cache
	if cfg.CacheEnabled {
		cacheService = cache.NewInMemoryCache(cfg.CacheMaxSize, cfg.CacheCleanupInterval)
		appLogger.Info(fmt.Sprintf("Cache enabled (max size: %d, cleanup interval: %s)",
			cfg.CacheMaxSize, cfg.CacheCleanupInterval))
	} else {
		// Use a no-op cache if caching is disabled
		cacheService = cache.NewInMemoryCache(1, 1*time.Minute) // Minimal cache
		appLogger.Info("Cache disabled")
	}

	// Create auth service with email service
	authService := service.NewAuthService(userRepo, emailService, tokenStore, authServiceConfig)

	// Create OAuth service
	googleConfig := service.OAuthConfig{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURI:  cfg.GoogleRedirectURI,
	}
	oauthService := service.NewOAuthService(
		oauthRepo,
		userRepo,
		googleConfig,
		cfg.JWTSecret,
		time.Duration(cfg.JWTExpirationHours)*time.Hour,
	)

	// Create registration service (for OTP-based registration)
	registrationService := service.NewRegistrationService(
		userRepo,
		otpCodeRepo,
		emailService,
		cfg.JWTSecret,
		time.Duration(cfg.JWTExpirationHours)*time.Hour,
	)

	// Create refresh token service (for remember me)
	refreshTokenService := service.NewRefreshTokenService(
		refreshTokenRepo,
		cfg.JWTSecret,
		time.Duration(cfg.JWTExpirationHours)*time.Hour,
	)

	userService := service.NewUserService(userRepo, uploadService, skillsMasterRepo)
	companyService := service.NewCompanyService(companyRepo, uploadService, cacheService, db)
	jobService := service.NewJobService(jobRepo, companyRepo, userRepo)
	applicationService := service.NewApplicationService(applicationRepo, jobRepo, userRepo, companyRepo, emailService, notificationService)
	skillsMasterService := service.NewSkillsMasterService(skillsMasterRepo)

	// Initialize handlers
	appLogger.Info("Initializing handlers...")
	authHandler := http.NewAuthHandler(authService, oauthService, registrationService, refreshTokenService, userRepo, companyRepo)
	userHandler := http.NewUserHandler(userService)

	// Initialize company handlers (split by domain)
	appLogger.Info("Initializing company handlers...")
	companyBasicHandler := http.NewCompanyBasicHandler(companyService)
	companyProfileHandler := http.NewCompanyProfileHandler(companyService)
	companyReviewHandler := http.NewCompanyReviewHandler(companyService)
	companyStatsHandler := http.NewCompanyStatsHandler(companyService)
	companyInviteHandler := http.NewCompanyInviteHandler(companyService, emailService, userService)

	// Initialize job & application handlers
	appLogger.Info("Initializing job & application handlers...")
	jobHandler := http.NewJobHandler(jobService)
	applicationHandler := http.NewApplicationHandler(applicationService)

	// Initialize master data handlers
	appLogger.Info("Initializing master data handlers...")
	skillsMasterHandler := http.NewSkillsMasterHandler(skillsMasterService)

	// Initialize notification handler
	appLogger.Info("Initializing notification handler...")
	notificationHandler := http.NewNotificationHandler(notificationService)

	// Initialize device token & push notification handlers
	appLogger.Info("📱 Initializing FCM handlers...")
	deviceTokenHandler := http.NewDeviceTokenHandler(deviceTokenRepo, fcmPushService, appLogger)
	pushNotificationHandler := http.NewPushNotificationHandler(fcmPushService, appLogger)
	appLogger.Info("✓ FCM handlers initialized successfully")

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName:               cfg.AppName,
		ServerHeader:          "Keerja",
		StrictRouting:         false,
		CaseSensitive:         false,
		ErrorHandler:          nil, // Will be set by middleware
		DisableStartupMessage: false,
	})

	// Setup global middleware (order matters!)
	appLogger.Info("Setting up middleware...")

	// 1. Panic recovery (must be first)
	app.Use(middleware.RecoverPanic(cfg.AppEnv == "development"))

	// 2. Request logging
	if cfg.AppEnv == "development" {
		app.Use(middleware.DetailedLogger())
	} else {
		app.Use(middleware.RequestLogger())
	}

	// 3. CORS
	app.Use(middleware.CORSConfig(cfg))

	// 4. Security headers
	app.Use(middleware.SecurityHeaders())

	// 5. Rate limiting
	app.Use(middleware.RateLimiter(cfg))

	// 6. Error handler
	app.Use(middleware.ErrorHandler(cfg.AppEnv == "development"))

	// Setup routes
	appLogger.Info("Setting up routes...")
	deps := &routes.Dependencies{
		Config:      cfg,
		AuthHandler: authHandler,
		UserHandler: userHandler,

		// Job & Application handlers
		JobHandler:         jobHandler,
		ApplicationHandler: applicationHandler,

		// Company handlers (split by domain)
		CompanyBasicHandler:   companyBasicHandler,
		CompanyProfileHandler: companyProfileHandler,
		CompanyReviewHandler:  companyReviewHandler,
		CompanyStatsHandler:   companyStatsHandler,
		CompanyInviteHandler:  companyInviteHandler,

		// Master data handlers
		SkillsMasterHandler: skillsMasterHandler,

		// Notification handler
		NotificationHandler: notificationHandler,

		// Device Token & Push Notification handlers
		DeviceTokenHandler:      deviceTokenHandler,
		PushNotificationHandler: pushNotificationHandler,

		// Services (for middlewares)
		CompanyService: companyService,
	}
	routes.SetupRoutes(app, deps)

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	appLogger.Info("✓ Swagger documentation available at /swagger/index.html")

	// 404 handler (must be last)
	app.Use(middleware.NotFoundHandler())

	// ==========================================
	// BACKGROUND JOBS SETUP
	// ==========================================
	appLogger.Info("Setting up background jobs...")

	// Initialize scheduler
	scheduler := jobs.NewScheduler()

	// Register jobs
	invitationExpiryJob := jobs.NewInvitationExpiryJob(companyService)
	if err := scheduler.Register(invitationExpiryJob); err != nil {
		appLogger.WithError(err).Fatal("Failed to register invitation expiry job")
	}

	// Register device token cleanup job
	deviceTokenCleanupConfig := jobs.CleanupConfig{
		InactiveDays:        90,  // Clean tokens inactive for 90+ days
		MaxFailureCount:     10,  // Clean tokens with 10+ failures
		BatchSize:           100, // Process 100 tokens per batch
		EnableInactiveClean: true,
		EnableFailureClean:  true,
	}
	deviceTokenCleanupJob := jobs.NewDeviceTokenCleanupJob(deviceTokenRepo, appLogger, deviceTokenCleanupConfig)
	if err := scheduler.Register(deviceTokenCleanupJob); err != nil {
		appLogger.WithError(err).Fatal("Failed to register device token cleanup job")
	}

	// Start scheduler
	scheduler.Start()

	// Ensure scheduler stops on shutdown
	defer scheduler.Stop()

	// Start server in a goroutine
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, port)

	go func() {
		appLogger.Info(fmt.Sprintf("Server listening on %s", addr))
		appLogger.Info(fmt.Sprintf("API Documentation: http://%s/api/v1/health", addr))
		if err := app.Listen(addr); err != nil {
			appLogger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server gracefully...")

	// Stop background jobs first
	appLogger.Info("Stopping background jobs...")
	scheduler.Stop()

	// Graceful shutdown with timeout
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		appLogger.WithError(err).Error("Server forced to shutdown")
	}

	appLogger.Info("Server stopped successfully")
}
