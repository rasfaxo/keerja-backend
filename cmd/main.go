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
	"keerja-backend/internal/handler/http/admin"
	applicationhandler "keerja-backend/internal/handler/http/application"
	authhandler "keerja-backend/internal/handler/http/auth"
	companyhandler "keerja-backend/internal/handler/http/company"
	"keerja-backend/internal/handler/http/health"
	jobhandler "keerja-backend/internal/handler/http/job"
	userhandler "keerja-backend/internal/handler/http/jobseeker"
	"keerja-backend/internal/handler/http/master"
	notificationhandler "keerja-backend/internal/handler/http/notification"
	"keerja-backend/internal/jobs"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/repository/postgres"
	"keerja-backend/internal/routes"
	"keerja-backend/internal/service"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

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
	appLogger.Info(" Database connected successfully")

	// Initialize Redis
	appLogger.Info("Initializing Redis connection...")
	redisClient, err := config.InitRedis(cfg)
	if err != nil {
		appLogger.WithError(err).Fatal("Failed to initialize Redis")
	}
	defer func() {
		if cerr := config.CloseRedis(); cerr != nil {
			appLogger.WithError(cerr).Warn("Failed to close Redis connection")
		}
	}()

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

	// Admin repositories
	adminUserRepo := postgres.NewAdminUserRepository(db)
	adminRoleRepo := postgres.NewAdminRoleRepository(db)

	// FCM Notification repository
	deviceTokenRepo := postgres.NewDeviceTokenRepository(db)

	// Master data repositories
	appLogger.Info("Initializing master data repositories...")
	industryRepo := postgres.NewIndustryRepository(db)
	companySizeRepo := postgres.NewCompanySizeRepository(db)
	provinceRepo := postgres.NewProvinceRepository(db)
	cityRepo := postgres.NewCityRepository(db)
	districtRepo := postgres.NewDistrictRepository(db)

	// Job master data repositories
	jobTitleRepo := postgres.NewJobTitleRepository(db)
	jobOptionsRepo := postgres.NewJobOptionsRepository(db)
	appLogger.Info("✓ Master data repositories initialized")

	// Initialize services
	appLogger.Info("Initializing services...")
	tokenStore := service.NewInMemoryTokenStore()

	// Initialize email service with config
	emailService := service.NewEmailService(emailRepo, cfg)

	// Initialize FCM service (Firebase Cloud Messaging)
	appLogger.Info("Initializing Firebase Cloud Messaging (FCM)...")
	fcmService := service.NewFCMPushService(deviceTokenRepo, cfg)
	if config.IsFCMEnabled() {
		appLogger.Info("FCM service initialized successfully")
	} else {
		appLogger.Warn("FCM service disabled (set FCM_ENABLED=true to enable)")
	}

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
	stateStore := service.NewRedisOAuthStateStore(redisClient)
	oauthService := service.NewOAuthService(
		oauthRepo,
		userRepo,
		googleConfig,
		cfg.JWTSecret,
		time.Duration(cfg.JWTExpirationHours)*time.Hour,
		stateStore,
		cfg.AllowedMobileRedirectURIs,
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

	// Admin services
	appLogger.Info("Initializing admin services...")
	adminAuthService := service.NewAdminAuthService(adminUserRepo, adminRoleRepo, cfg)
	adminCompanyService := service.NewAdminCompanyService(companyRepo, jobRepo, emailService, cacheService)
	appLogger.Info("✓ Admin services initialized")

	// Master data services
	appLogger.Info("Initializing master data services...")
	industryService := service.NewIndustryService(industryRepo, cacheService)
	companySizeService := service.NewCompanySizeService(companySizeRepo, cacheService)
	provinceService := service.NewProvinceService(provinceRepo, cacheService)
	cityService := service.NewCityService(cityRepo, provinceRepo, cacheService)
	districtService := service.NewDistrictService(districtRepo, cityRepo, provinceRepo, cacheService)

	// Job master data services
	jobTitleService := service.NewJobTitleService(jobTitleRepo)
	jobOptionsService := service.NewJobOptionsService(jobOptionsRepo, cacheService)
	appLogger.Info("✓ Master data services initialized")

	// Company service
	companyService := service.NewCompanyService(
		companyRepo,
		uploadService,
		cacheService,
		db,
		industryService,
		companySizeService,
		districtService,
		jobRepo,
		userService,
		userRepo,
	)

	jobService := service.NewJobService(
		jobRepo,
		companyRepo,
		userRepo,
		jobOptionsRepo,
		jobTitleRepo,
		industryService,
		districtService,
	)

	// Admin job service (orchestrates admin operations on jobs)
	adminJobService := service.NewAdminJobService(jobRepo)

	applicationService := service.NewApplicationService(applicationRepo, jobRepo, userRepo, companyRepo, emailService, nil) // notificationService disabled temporarily
	skillsMasterService := service.NewSkillsMasterService(skillsMasterRepo)

	// Initialize handlers
	appLogger.Info("Initializing handlers...")
	authHandler := authhandler.NewAuthHandler(authService, oauthService, registrationService, refreshTokenService, userRepo, companyRepo)

	// Initialize user handlers (split by domain)
	appLogger.Info("Initializing user handlers...")
	userProfileHandler := userhandler.NewUserProfileHandler(userService)
	userEducationHandler := userhandler.NewUserEducationHandler(userService)
	userExperienceHandler := userhandler.NewUserExperienceHandler(userService)
	userSkillHandler := userhandler.NewUserSkillHandler(userService)
	userDocumentHandler := userhandler.NewUserDocumentHandler(userService)
	userMiscHandler := userhandler.NewUserMiscHandler(userService)

	// Initialize company handlers (split by domain)
	appLogger.Info("Initializing company handlers...")
	companyBasicHandler := companyhandler.NewCompanyBasicHandler(
		companyService,
		industryRepo,
		companySizeRepo,
		provinceRepo,
		cityRepo,
		districtRepo,
	)
	companyImageHandler := companyhandler.NewCompanyImageHandler(companyService)
	companyAddressHandler := companyhandler.NewCompanyAddressHandler(
		companyService,
		provinceRepo,
		cityRepo,
		districtRepo,
	)
	companyEmployerHandler := companyhandler.NewCompanyEmployerHandler(
		companyService,
		userService,
		provinceRepo,
		cityRepo,
		districtRepo,
	)
	companyVerificationHandler := companyhandler.NewCompanyVerificationHandler(companyService)
	companyProfileHandler := companyhandler.NewCompanyProfileHandler(companyService)
	companyReviewHandler := companyhandler.NewCompanyReviewHandler(companyService)
	companyStatsHandler := companyhandler.NewCompanyStatsHandler(companyService)
	companyInviteHandler := companyhandler.NewCompanyInviteHandler(companyService, emailService, userService)

	// Initialize job & application handlers
	appLogger.Info("Initializing job & application handlers...")
	jobHandler := jobhandler.NewJobHandler(jobService, companyService, jobOptionsService, skillsMasterService)
	applicationHandler := applicationhandler.NewApplicationHandler(applicationService)

	// Initialize admin handlers
	appLogger.Info("Initializing admin handlers...")
	adminJobHandler := admin.NewAdminJobHandler(adminJobService)
	adminAuthHandler := admin.NewAdminAuthHandler(adminAuthService)
	adminCompanyHandler := admin.NewCompanyHandler(adminCompanyService)

	// Initialize admin master data services
	appLogger.Info("Initializing admin master data services...")
	adminIndustryService := service.NewAdminIndustryService(industryService, industryRepo, db, cacheService)
	adminCompanySizeService := service.NewAdminCompanySizeService(companySizeService, companySizeRepo, db, cacheService)
	adminProvinceService := service.NewAdminProvinceService(provinceService, provinceRepo, db, cacheService)
	adminCityService := service.NewAdminCityService(cityService, cityRepo, db, cacheService)
	adminDistrictService := service.NewAdminDistrictService(districtService, districtRepo, db, cacheService)
	adminJobTypeService := service.NewAdminJobTypeService(jobOptionsService, jobOptionsRepo, db, cacheService)
	appLogger.Info("✓ Admin master data services initialized")

	// Initialize admin master data handler
	adminMasterDataHandler := admin.NewAdminMasterDataHandler(
		adminProvinceService,
		adminCityService,
		adminDistrictService,
		adminIndustryService,
		adminCompanySizeService,
		adminJobTypeService,
	)

	// Initialize master data handlers
	appLogger.Info("Initializing master data handlers...")
	skillsMasterHandler := master.NewSkillsMasterHandler(skillsMasterService)

	// Initialize master data handlers (industries, company sizes, locations)
	industryHandler := master.NewIndustryHandler(industryService)
	companySizeHandler := master.NewCompanySizeHandler(companySizeService)
	locationHandler := master.NewLocationHandler(provinceService, cityService, districtService)

	// Initialize job master data handler (job titles & options)
	masterDataHandler := master.NewMasterDataHandler(jobTitleService, jobOptionsService, jobService, companyService, skillsMasterService)

	masterDataHandlers := &routes.MasterDataHandlers{
		IndustryHandler:    industryHandler,
		CompanySizeHandler: companySizeHandler,
		LocationHandler:    locationHandler,
	}
	appLogger.Info("✓ Master data handlers initialized")

	// Initialize FCM notification handlers
	appLogger.Info("Initializing FCM notification handlers...")
	deviceTokenHandler := notificationhandler.NewDeviceTokenHandler(deviceTokenRepo, fcmService, appLogger)
	pushNotificationHandler := notificationhandler.NewPushNotificationHandler(fcmService, appLogger)
	appLogger.Info("FCM handlers initialized successfully")

	// Initialize health check handler
	appLogger.Info("Initializing health check handler...")
	healthHandler := health.NewHealthHandler(db, redisClient, cfg.AppVersion)
	appLogger.Info("✓ Health check handler initialized")

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

	// Setup health check routes (before auth middleware)
	routes.SetupHealthRoutes(app, healthHandler)

	adminAuthMw := middleware.NewAdminAuthMiddleware(cfg, adminUserRepo)
	deps := &routes.Dependencies{
		Config:      cfg,
		AuthHandler: authHandler,

		// User handlers (split by domain)
		UserProfileHandler:    userProfileHandler,
		UserEducationHandler:  userEducationHandler,
		UserExperienceHandler: userExperienceHandler,
		UserSkillHandler:      userSkillHandler,
		UserDocumentHandler:   userDocumentHandler,
		UserMiscHandler:       userMiscHandler,

		// Admin handlers
		AdminAuthHandler:    adminAuthHandler,
		AdminCompanyHandler: adminCompanyHandler,
		AdminAuthMiddleware: adminAuthMw,

		// Job & Application handlers
		JobHandler:             jobHandler,
		ApplicationHandler:     applicationHandler,
		AdminJobHandler:        adminJobHandler,
		AdminMasterDataHandler: adminMasterDataHandler,

		// Company handlers (split by domain)
		CompanyBasicHandler:        companyBasicHandler,
		CompanyImageHandler:        companyImageHandler,
		CompanyAddressHandler:      companyAddressHandler,
		CompanyEmployerHandler:     companyEmployerHandler,
		CompanyVerificationHandler: companyVerificationHandler,
		CompanyProfileHandler:      companyProfileHandler,
		CompanyReviewHandler:       companyReviewHandler,
		CompanyStatsHandler:        companyStatsHandler,
		CompanyInviteHandler:       companyInviteHandler,

		// Master data handlers
		SkillsMasterHandler: skillsMasterHandler,
		MasterDataHandlers:  masterDataHandlers,
		MasterDataHandler:   masterDataHandler,

		// FCM Notification handlers
		DeviceTokenHandler:      deviceTokenHandler,
		PushNotificationHandler: pushNotificationHandler,

		// Services (for middlewares)
		CompanyService: companyService,
	}
	routes.SetupRoutes(app, deps)

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
