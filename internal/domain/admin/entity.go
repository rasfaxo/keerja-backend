package admin

import (
	"time"

	"gorm.io/gorm"
)

// AdminRole represents an administrative role with specific permissions
// Maps to: admin_roles table
type AdminRole struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleName        string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"role_name" validate:"required,min=3,max=100"`
	RoleDescription string    `gorm:"type:text" json:"role_description,omitempty"`
	AccessLevel     int16     `gorm:"type:smallint;default:5" json:"access_level" validate:"min=1,max=10"`
	IsSystemRole    bool      `gorm:"default:false" json:"is_system_role"`
	CreatedBy       *int64    `gorm:"index" json:"created_by,omitempty"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Creator *AdminUser  `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"creator,omitempty"`
	Users   []AdminUser `gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL" json:"users,omitempty"`
}

// TableName specifies the table name for AdminRole
func (AdminRole) TableName() string {
	return "admin_roles"
}

// IsSuperAdmin checks if the role has highest access level
func (r *AdminRole) IsSuperAdmin() bool {
	return r.AccessLevel >= 9
}

// IsAdmin checks if the role has admin level access
func (r *AdminRole) IsAdmin() bool {
	return r.AccessLevel >= 7
}

// IsModerator checks if the role has moderator level access
func (r *AdminRole) IsModerator() bool {
	return r.AccessLevel >= 5
}

// CanModifyRole checks if this role can modify another role based on access level
func (r *AdminRole) CanModifyRole(targetRole *AdminRole) bool {
	return r.AccessLevel > targetRole.AccessLevel
}

// AdminUser represents an administrative user in the system
// Maps to: admin_users table
type AdminUser struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID            string         `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex" json:"uuid"`
	FullName        string         `gorm:"type:varchar(100);not null" json:"full_name" validate:"required,min=2,max=100"`
	Email           string         `gorm:"type:varchar(150);not null;uniqueIndex" json:"email" validate:"required,email,max=150"`
	Phone           string         `gorm:"type:varchar(20)" json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	PasswordHash    string         `gorm:"type:text;not null" json:"-"`
	RoleID          *int64         `gorm:"index" json:"role_id,omitempty"`
	Status          string         `gorm:"type:varchar(20);default:'active';index" json:"status" validate:"oneof=active inactive suspended"`
	LastLogin       *time.Time     `json:"last_login,omitempty"`
	TwoFactorSecret string         `gorm:"type:varchar(100)" json:"-"`
	ProfileImageURL string         `gorm:"type:text" json:"profile_image_url,omitempty"`
	CreatedBy       *int64         `gorm:"index" json:"created_by,omitempty"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Role         *AdminRole  `gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL" json:"role,omitempty"`
	Creator      *AdminUser  `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"creator,omitempty"`
	CreatedUsers []AdminUser `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"created_users,omitempty"`
	CreatedRoles []AdminRole `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"created_roles,omitempty"`
}

// TableName specifies the table name for AdminUser
func (AdminUser) TableName() string {
	return "admin_users"
}

// IsActive checks if the admin user status is active
func (a *AdminUser) IsActive() bool {
	return a.Status == "active"
}

// IsInactive checks if the admin user status is inactive
func (a *AdminUser) IsInactive() bool {
	return a.Status == "inactive"
}

// IsSuspended checks if the admin user status is suspended
func (a *AdminUser) IsSuspended() bool {
	return a.Status == "suspended"
}

// Has2FA checks if two-factor authentication is enabled
func (a *AdminUser) Has2FA() bool {
	return a.TwoFactorSecret != ""
}

// UpdateLastLogin updates the last login timestamp
func (a *AdminUser) UpdateLastLogin() {
	now := time.Now()
	a.LastLogin = &now
}

// GetAccessLevel returns the access level from the role
func (a *AdminUser) GetAccessLevel() int16 {
	if a.Role != nil {
		return a.Role.AccessLevel
	}
	return 0
}

// IsSuperAdmin checks if the user has super admin role
func (a *AdminUser) IsSuperAdmin() bool {
	return a.Role != nil && a.Role.IsSuperAdmin()
}

// IsAdmin checks if the user has admin role
func (a *AdminUser) IsAdmin() bool {
	return a.Role != nil && a.Role.IsAdmin()
}

// IsModerator checks if the user has moderator role
func (a *AdminUser) IsModerator() bool {
	return a.Role != nil && a.Role.IsModerator()
}

// CanModifyUser checks if this admin can modify another admin user
func (a *AdminUser) CanModifyUser(targetUser *AdminUser) bool {
	if a.Role == nil || targetUser.Role == nil {
		return false
	}
	return a.Role.AccessLevel > targetUser.Role.AccessLevel
}
