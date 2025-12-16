package routes

import (
	"keerja-backend/internal/config"
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/handler/http/admin"
	applicationhandler "keerja-backend/internal/handler/http/application"
	authhandler "keerja-backend/internal/handler/http/auth"
	chathandler "keerja-backend/internal/handler/http/chat"
	companyhandler "keerja-backend/internal/handler/http/company"
	jobhandler "keerja-backend/internal/handler/http/job"
	userhandler "keerja-backend/internal/handler/http/jobseeker"
	"keerja-backend/internal/handler/http/master"
	notificationhandler "keerja-backend/internal/handler/http/notification"
	"keerja-backend/internal/handler/websocket"
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Dependencies holds all handler dependencies
type Dependencies struct {
	Config                 *config.Config
	AuthHandler            *authhandler.AuthHandler
	JobHandler             *jobhandler.JobHandler                 // Job management (10 endpoints)
	ApplicationHandler     *applicationhandler.ApplicationHandler // Application management (21 endpoints)
	AdminJobHandler        *admin.AdminJobHandler                 // Admin moderation & job approval
	AdminMasterDataHandler *admin.AdminMasterDataHandler          // Admin master data CRUD

	// Admin handlers
	AdminAuthHandler    *admin.AdminAuthHandler         // Admin authentication
	AdminCompanyHandler *admin.CompanyHandler           // Company moderation
	AdminAuthMiddleware *middleware.AdminAuthMiddleware // Admin auth middleware

	// User handlers (split by domain for better organization)
	UserProfileHandler    *userhandler.UserProfileHandler    // Profile & preferences (5 endpoints)
	UserEducationHandler  *userhandler.UserEducationHandler  // Education CRUD (4 endpoints)
	UserExperienceHandler *userhandler.UserExperienceHandler // Experience CRUD (4 endpoints)
	UserSkillHandler      *userhandler.UserSkillHandler      // Skills management (3 endpoints)
	UserDocumentHandler   *userhandler.UserDocumentHandler   // Document upload (2 endpoints)
	UserMiscHandler       *userhandler.UserMiscHandler       // Certifications, languages, projects (3 endpoints)

	// Company handlers (split by domain for better organization)
	CompanyBasicHandler        *companyhandler.CompanyBasicHandler        // CRUD operations (7 endpoints)
	CompanyImageHandler        *companyhandler.CompanyImageHandler        // Image upload/delete (4 endpoints)
	CompanyAddressHandler      *companyhandler.CompanyAddressHandler      // Address CRUD (4 endpoints)
	CompanyEmployerHandler     *companyhandler.CompanyEmployerHandler     // Employer profile (2 endpoints)
	CompanyVerificationHandler *companyhandler.CompanyVerificationHandler // Verification (3 endpoints)
	CompanyProfileHandler      *companyhandler.CompanyProfileHandler      // Profile & social features (8 endpoints)
	CompanyReviewHandler       *companyhandler.CompanyReviewHandler       // Review system (5 endpoints)
	CompanyStatsHandler        *companyhandler.CompanyStatsHandler        // Statistics & queries (3 endpoints)
	CompanyInviteHandler       *companyhandler.CompanyInviteHandler       // Employee invitation (5 endpoints)
	// Master data handlers
	SkillsMasterHandler *master.SkillsMasterHandler // Skills master data (8 endpoints)
	MasterDataHandlers  *MasterDataHandlers         // Industry, company size, location (10 endpoints)
	MasterDataHandler   *master.MasterDataHandler   // Job titles & options (Phase 1-4)

	// FCM Notification handlers (Firebase Cloud Messaging)
	DeviceTokenHandler      *notificationhandler.DeviceTokenHandler      // Device token management (6 endpoints)
	PushNotificationHandler *notificationhandler.PushNotificationHandler // Push notifications (5 endpoints)

	// Chat handlers
	ChatHandler      *chathandler.ChatHandler // Chat HTTP handler (6 endpoints)
	WebSocketHub     *websocket.Hub           // WebSocket hub
	WebSocketHandler *websocket.Handler       // WebSocket handler

	// Services (for middlewares)
	CompanyService company.CompanyService
}

// SetupRoutes configures all application routes
// This is the main entry point for route configuration
func SetupRoutes(app *fiber.App, deps *Dependencies) {
	// Initialize auth middleware
	authMw := middleware.NewAuthMiddleware(deps.Config)

	// Initialize permission middleware
	permMw := middleware.NewPermissionMiddleware(deps.CompanyService)

	// Get admin auth middleware from dependencies
	adminAuthMw := deps.AdminAuthMiddleware

	// API v1 group
	api := app.Group("/api/v1")

	// Setup route groups (each in separate file)
	SetupAuthRoutes(api, deps, authMw)               // auth_routes.go
	SetupUserRoutes(api, deps, authMw)               // user_routes.go
	SetupJobRoutes(api, deps, authMw)                // job_routes.go
	SetupApplicationRoutes(api, deps, authMw)        // application_routes.go
	SetupCompanyRoutes(api, deps, authMw, permMw)    // company_routes.go
	SetupAdminAuthRoutes(api, deps, adminAuthMw)     // admin_auth_routes.go
	SetupAdminRoutes(api, deps, adminAuthMw)         // admin_routes.go
	SetupSkillsRoutes(api, deps.SkillsMasterHandler) // skills_routes.go

	// Master data routes (industries, company sizes, locations)
	if deps.MasterDataHandlers != nil {
		SetupMasterDataRoutes(api, deps.MasterDataHandlers) // master_routes.go
	}

	// Job master data routes (job titles & options - Phase 1-4)
	if deps.MasterDataHandler != nil {
		SetupJobMasterDataRoutes(api, deps.MasterDataHandler, authMw) // job_master_data_routes.go
	}

	// FCM Notification routes
	if deps.DeviceTokenHandler != nil {
		SetupDeviceTokenRoutes(api, deps.DeviceTokenHandler, authMw) // device_token_routes.go
	}
	if deps.PushNotificationHandler != nil {
		SetupPushNotificationRoutes(api, deps.PushNotificationHandler, authMw) // push_notification_routes.go
	}

	// Chat routes
	if deps.ChatHandler != nil {
		SetupChatRoutes(api, deps.ChatHandler, authMw) // chat_routes.go
	}

	// WebSocket routes
	if deps.WebSocketHandler != nil {
		SetupWebSocketRoutes(app, deps.WebSocketHandler) // websocket_routes.go
	}
}
