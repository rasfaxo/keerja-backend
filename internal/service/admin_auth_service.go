package service

import (
	"context"
	"errors"
	"time"

	"keerja-backend/internal/config"
	"keerja-backend/internal/domain/admin"
	"keerja-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// Common errors for admin auth
var (
	ErrAdminInvalidCredentials = errors.New("invalid email or password")
	ErrAdminNotFound           = errors.New("admin user not found")
	ErrAdminNotActive          = errors.New("admin account is not active")
	ErrAdminSuspended          = errors.New("admin account is suspended")
	ErrAdminInvalidPassword    = errors.New("current password is incorrect")
	ErrAdminInvalidToken       = errors.New("invalid refresh token")
	ErrAdminTokenExpired       = errors.New("refresh token expired")
)

// AdminAuthService handles admin authentication operations
type AdminAuthService struct {
	adminUserRepo admin.AdminUserRepository
	adminRoleRepo admin.AdminRoleRepository
	config        *config.Config
}

// NewAdminAuthService creates new admin auth service
func NewAdminAuthService(
	adminUserRepo admin.AdminUserRepository,
	adminRoleRepo admin.AdminRoleRepository,
	cfg *config.Config,
) *AdminAuthService {
	return &AdminAuthService{
		adminUserRepo: adminUserRepo,
		adminRoleRepo: adminRoleRepo,
		config:        cfg,
	}
}

// Login authenticates admin user
func (s *AdminAuthService) Login(ctx context.Context, email, password string) (*admin.AdminUser, string, string, error) {
	// Find admin by email
	adminUser, err := s.adminUserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", ErrAdminInvalidCredentials
	}

	// Check if admin exists
	if adminUser == nil {
		return nil, "", "", ErrAdminInvalidCredentials
	}

	// Check account status
	if adminUser.Status == "inactive" {
		return nil, "", "", ErrAdminNotActive
	}
	if adminUser.Status == "suspended" {
		return nil, "", "", ErrAdminSuspended
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", ErrAdminInvalidCredentials
	}

	// Load role information
	if adminUser.RoleID != nil {
		role, err := s.adminRoleRepo.FindByID(ctx, *adminUser.RoleID)
		if err == nil && role != nil {
			adminUser.Role = role
		}
	}

	// Generate access token (1 hour expiry)
	accessToken, err := utils.GenerateAdminToken(
		adminUser.ID,
		adminUser.Email,
		adminUser.RoleID,
		adminUser.GetAccessLevel(),
		s.config.JWTSecret,
		time.Hour*1, // 1 hour
	)
	if err != nil {
		return nil, "", "", err
	}

	// Generate refresh token (7 days expiry)
	refreshToken, err := utils.GenerateAdminToken(
		adminUser.ID,
		adminUser.Email,
		adminUser.RoleID,
		adminUser.GetAccessLevel(),
		s.config.JWTSecret,
		time.Hour*24*7, // 7 days
	)
	if err != nil {
		return nil, "", "", err
	}

	// Update last login
	if err := s.adminUserRepo.UpdateLastLogin(ctx, adminUser.ID); err != nil {
		// Log error but don't fail login
		println("Failed to update last login:", err.Error())
	}

	// Reload admin to get updated last_login
	updatedAdmin, err := s.adminUserRepo.FindByID(ctx, adminUser.ID)
	if err == nil && updatedAdmin != nil {
		adminUser = updatedAdmin
		// Reload role again
		if adminUser.RoleID != nil {
			role, err := s.adminRoleRepo.FindByID(ctx, *adminUser.RoleID)
			if err == nil && role != nil {
				adminUser.Role = role
			}
		}
	}

	return adminUser, accessToken, refreshToken, nil
}

// Logout invalidates admin session
func (s *AdminAuthService) Logout(ctx context.Context, adminID int64) error {
	// For now, JWT tokens are stateless
	// In production, you might want to:
	// 1. Add token to blacklist/redis
	// 2. Clear session from database
	// 3. Revoke refresh tokens

	// Verify admin exists
	adminUser, err := s.adminUserRepo.FindByID(ctx, adminID)
	if err != nil {
		return ErrAdminNotFound
	}
	if adminUser == nil {
		return ErrAdminNotFound
	}

	return nil
}

// RefreshToken generates new access token from refresh token
func (s *AdminAuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Validate refresh token
	claims, err := utils.ValidateAdminToken(refreshToken, s.config.JWTSecret)
	if err != nil {
		if err == utils.ErrExpiredToken {
			return "", "", ErrAdminTokenExpired
		}
		return "", "", ErrAdminInvalidToken
	}

	// Verify admin still exists and is active
	adminUser, err := s.adminUserRepo.FindByID(ctx, claims.AdminID)
	if err != nil || adminUser == nil {
		return "", "", ErrAdminNotFound
	}

	if adminUser.Status != "active" {
		return "", "", ErrAdminNotActive
	}

	// Load role information
	if adminUser.RoleID != nil {
		role, err := s.adminRoleRepo.FindByID(ctx, *adminUser.RoleID)
		if err == nil && role != nil {
			adminUser.Role = role
		}
	}

	// Generate new access token (1 hour)
	newAccessToken, err := utils.GenerateAdminToken(
		adminUser.ID,
		adminUser.Email,
		adminUser.RoleID,
		adminUser.GetAccessLevel(),
		s.config.JWTSecret,
		time.Hour*1,
	)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token (7 days)
	newRefreshToken, err := utils.GenerateAdminToken(
		adminUser.ID,
		adminUser.Email,
		adminUser.RoleID,
		adminUser.GetAccessLevel(),
		s.config.JWTSecret,
		time.Hour*24*7,
	)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// GetCurrentProfile gets admin profile by ID
func (s *AdminAuthService) GetCurrentProfile(ctx context.Context, adminID int64) (*admin.AdminUser, error) {
	adminUser, err := s.adminUserRepo.FindByID(ctx, adminID)
	if err != nil {
		return nil, ErrAdminNotFound
	}
	if adminUser == nil {
		return nil, ErrAdminNotFound
	}

	// Load role information
	if adminUser.RoleID != nil {
		role, err := s.adminRoleRepo.FindByID(ctx, *adminUser.RoleID)
		if err == nil && role != nil {
			adminUser.Role = role
		}
	}

	return adminUser, nil
}

// ChangePassword changes admin password
func (s *AdminAuthService) ChangePassword(ctx context.Context, adminID int64, currentPassword, newPassword string) error {
	// Get admin user
	adminUser, err := s.adminUserRepo.FindByID(ctx, adminID)
	if err != nil || adminUser == nil {
		return ErrAdminNotFound
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrAdminInvalidPassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password using repository method
	return s.adminUserRepo.UpdatePassword(ctx, adminID, string(hashedPassword))
}

// ValidateSession checks if admin session is valid
func (s *AdminAuthService) ValidateSession(ctx context.Context, adminID int64) (bool, error) {
	adminUser, err := s.adminUserRepo.FindByID(ctx, adminID)
	if err != nil || adminUser == nil {
		return false, ErrAdminNotFound
	}

	// Check if admin is active
	if adminUser.Status != "active" {
		return false, ErrAdminNotActive
	}

	return true, nil
}

// InvalidateAllSessions invalidates all admin sessions
func (s *AdminAuthService) InvalidateAllSessions(ctx context.Context, adminID int64) error {
	// For JWT-based auth, this would typically:
	// 1. Increment a version number in database
	// 2. Add all tokens to blacklist
	// 3. Clear session store

	// Verify admin exists
	_, err := s.adminUserRepo.FindByID(ctx, adminID)
	if err != nil {
		return ErrAdminNotFound
	}

	// In production, implement token blacklist or session versioning
	return nil
}
