# Admin Domain

## Overview

The Admin domain manages the administrative layer of the Keerja job portal, including admin users, roles, and access control. This domain ensures secure system administration with role-based access control (RBAC) and comprehensive audit tracking.

## Entities (2)

### 1. AdminRole

Administrative role with specific permissions and access levels.

**Fields:**

- `ID` (int64): Primary key
- `RoleName` (string): Unique role name (max 100 chars)
- `RoleDescription` (text): Description of role responsibilities
- `AccessLevel` (int16): Numeric access level (1-10, higher = more access)
- `IsSystemRole` (bool): System-defined roles cannot be modified
- `CreatedBy` (int64): Admin user who created this role
- `CreatedAt`, `UpdatedAt`: Timestamps
- `DeletedAt`: Soft delete support

**Relationships:**

- `Creator` (AdminUser): The admin who created this role
- `Users` ([]AdminUser): All users assigned to this role

**Helper Methods:**

- `IsSuperAdmin()`: Checks if access level >= 9
- `IsAdmin()`: Checks if access level >= 7
- `IsModerator()`: Checks if access level >= 5
- `CanModifyRole(targetRole)`: Checks if this role can modify another

**Access Levels:**

- 1-4: Basic access (viewer, reporter)
- 5-6: Moderator (can manage content)
- 7-8: Admin (can manage users and settings)
- 9-10: Super Admin (full system access)

### 2. AdminUser

Administrative user with role-based permissions.

**Fields:**

- `ID` (int64): Primary key
- `UUID` (uuid): Unique identifier
- `FullName` (string): Admin's full name (max 100 chars)
- `Email` (string): Unique email address
- `Phone` (string): Optional phone number
- `PasswordHash` (text): Bcrypt hashed password
- `RoleID` (int64): Reference to AdminRole
- `Status` (string): Enum: active, inactive, suspended
- `LastLogin` (timestamp): Last login time
- `TwoFactorSecret` (string): Secret for 2FA (TOTP)
- `ProfileImageURL` (text): Profile picture URL
- `CreatedBy` (int64): Admin user who created this account
- `CreatedAt`, `UpdatedAt`: Timestamps
- `DeletedAt`: Soft delete support

**Relationships:**

- `Role` (AdminRole): The role assigned to this admin
- `Creator` (AdminUser): The admin who created this account
- `CreatedUsers` ([]AdminUser): Users created by this admin
- `CreatedRoles` ([]AdminRole): Roles created by this admin

**Helper Methods:**

- `IsActive()`, `IsInactive()`, `IsSuspended()`: Status checks
- `Has2FA()`: Checks if 2FA is enabled
- `UpdateLastLogin()`: Updates last login timestamp
- `GetAccessLevel()`: Returns access level from role
- `IsSuperAdmin()`, `IsAdmin()`, `IsModerator()`: Permission checks
- `CanModifyUser(targetUser)`: Checks if can modify another admin

**Security Features:**

- Two-factor authentication (TOTP)
- Password hashing (bcrypt)
- Role-based access control
- Self-referential creator tracking
- Status management (active/inactive/suspended)

## Repository Layer (2 Interfaces)

### AdminRoleRepository (20 methods)

**Basic CRUD (5):**

- `Create`, `FindByID`, `FindByName`, `Update`, `Delete`

**Listing & Search (4):**

- `List(filter)`: Paginated list with filtering
- `ListActive`: All active roles
- `GetSystemRoles`: System-defined roles only
- `GetNonSystemRoles`: Custom roles only

**Access Level Operations (2):**

- `GetRolesByAccessLevel(min, max)`: Roles in range
- `GetRolesByMinAccessLevel(min)`: Roles with minimum access

**Statistics (3):**

- `Count`: Total roles count
- `CountByAccessLevel`: Count by specific level
- `GetRoleStats`: Comprehensive statistics

### AdminUserRepository (41 methods)

**Basic CRUD (6):**

- `Create`, `FindByID`, `FindByUUID`, `FindByEmail`, `Update`, `Delete`

**Listing & Search (4):**

- `List(filter)`: Paginated list with filtering
- `ListByRole`: Users by role
- `ListByStatus`: Users by status
- `SearchUsers`: Full-text search

**Status Operations (7):**

- `UpdateStatus`, `ActivateUser`, `DeactivateUser`, `SuspendUser`
- `GetActiveUsers`, `GetInactiveUsers`, `GetSuspendedUsers`

**Role Operations (2):**

- `UpdateRole`: Assign role to user
- `GetUsersByRole`: All users with specific role

**Authentication & Security (6):**

- `UpdatePassword`: Change password hash
- `UpdateLastLogin`: Track login time
- `Enable2FA`, `Disable2FA`: Manage 2FA
- `Get2FAUsers`: Users with 2FA enabled

**Profile Operations (2):**

- `UpdateProfile`: Update name, phone, image
- `UpdateProfileImage`: Update profile picture only

**Statistics (6):**

- `Count`, `CountByStatus`, `CountByRole`
- `GetUserStats`: Comprehensive user statistics
- `GetActivityStats(startDate, endDate)`: Activity analytics
- `GetRecentLogins(limit)`: Recent login activity

**Audit & Tracking (2):**

- `GetCreatedUsers(creatorID)`: Users created by admin
- `GetUserActivity(userID, startDate, endDate)`: User activity log

## Service Layer (2 Interfaces)

### AdminRoleService (14 methods)

**Role Management (6):**

- `CreateRole(req)`: Create new role
- `UpdateRole(id, req)`: Update existing role
- `DeleteRole(id)`: Delete role (if no users assigned)
- `GetRole(id)`: Get role details with user count
- `GetRoleByName(name)`: Get role by unique name
- `GetRoles(filter)`: Paginated list with filtering

**System Roles (2):**

- `GetSystemRoles`: List system-defined roles
- `GetNonSystemRoles`: List custom roles

**Access Level Operations (3):**

- `GetRolesByAccessLevel(min, max)`: Roles in range
- `PromoteRole(id, newLevel)`: Increase access level
- `DemoteRole(id, newLevel)`: Decrease access level

**Statistics (2):**

- `GetRoleStats`: Overall role statistics
- `GetRoleUsage(roleID)`: Usage info for specific role

**Validation (3):**

- `ValidateRole(req)`: Validate role data
- `CheckRolePermissions(roleID, requiredLevel)`: Permission check
- `CanModifyRole(actorRoleID, targetRoleID)`: Authorization check

### AdminUserService (50 methods)

**User Management (8):**

- `CreateUser(req)`: Create new admin user
- `UpdateUser(id, req)`: Update user details
- `DeleteUser(id)`: Delete admin user
- `GetUser(id)`, `GetUserByUUID(uuid)`, `GetUserByEmail(email)`: Get user
- `GetUsers(filter)`: Paginated list with filtering
- `SearchUsers(query, page, pageSize)`: Full-text search

**Authentication (6):**

- `Login(req)`: Authenticate and generate tokens
- `Logout(userID)`: Invalidate tokens
- `RefreshToken(refreshToken)`: Get new access token
- `ChangePassword(userID, req)`: Change password
- `ResetPassword(req)`: Reset password with token
- `RequestPasswordReset(email)`: Send reset email

**Two-Factor Authentication (4):**

- `Enable2FA(userID)`: Enable and generate secret
- `Verify2FA(userID, code)`: Verify TOTP code
- `Disable2FA(userID, password)`: Disable 2FA
- `Generate2FAQRCode(userID)`: Generate QR code for setup

**Status Management (7):**

- `ActivateUser(id)`: Set status to active
- `DeactivateUser(id)`: Set status to inactive
- `SuspendUser(id, reason)`: Suspend with reason
- `UnsuspendUser(id)`: Remove suspension
- `GetActiveUsers`, `GetInactiveUsers`, `GetSuspendedUsers`: List by status

**Role Management (3):**

- `AssignRole(userID, roleID)`: Assign role to user
- `RemoveRole(userID)`: Remove role from user
- `GetUsersByRole(roleID)`: List users with role

**Profile Management (3):**

- `UpdateProfile(userID, req)`: Update profile info
- `UpdateProfileImage(userID, imageURL)`: Update profile picture
- `GetProfile(userID)`: Get complete profile with stats

**Statistics & Analytics (4):**

- `GetUserStats`: Overall user statistics
- `GetActivityStats(startDate, endDate)`: Activity analytics
- `GetRecentLogins(limit)`: Recent login activity
- `GetUserActivity(userID, startDate, endDate)`: User activity details

**Audit & Tracking (3):**

- `GetCreatedUsers(creatorID)`: Users created by admin
- `TrackLogin(userID)`: Record login event
- `GetLoginHistory(userID, limit)`: Login history

**Validation & Authorization (4):**

- `ValidateUser(req)`: Validate user data
- `CheckUserPermissions(userID, requiredLevel)`: Permission check
- `CanModifyUser(actorID, targetID)`: Authorization check
- `IsUserActive(userID)`: Check active status

## DTOs

### Request DTOs (8)

1. **CreateRoleRequest**: RoleName, RoleDescription, AccessLevel, IsSystemRole, CreatedBy
2. **UpdateRoleRequest**: RoleName, RoleDescription, AccessLevel
3. **CreateUserRequest**: FullName, Email, Phone, Password, RoleID, Status, ProfileImageURL, CreatedBy
4. **UpdateUserRequest**: FullName, Phone, RoleID, Status, ProfileImageURL
5. **LoginRequest**: Email, Password, TwoFACode
6. **ChangePasswordRequest**: CurrentPassword, NewPassword, ConfirmPassword
7. **ResetPasswordRequest**: Email, ResetToken, NewPassword, ConfirmPassword
8. **UpdateProfileRequest**: FullName, Phone, ProfileImageURL

### Response DTOs (13)

1. **AdminRoleResponse**: Role details with user count
2. **RoleListResponse**: Paginated role list
3. **AdminUserResponse**: User details with role and 2FA status
4. **UserListResponse**: Paginated user list
5. **LoginResponse**: User + access/refresh tokens + 2FA requirement
6. **TokenResponse**: Access/refresh tokens
7. **Enable2FAResponse**: Secret, QR code, backup codes
8. **ProfileResponse**: User profile with stats and recent activity
9. **RoleStatsResponse**: Role statistics with distribution
10. **RoleUsageResponse**: Role usage details with user list
11. **UserStatsResponse**: User statistics with status breakdown
12. **ActivityStatsResponse**: Activity analytics with trends
13. **UserActivityResponse**: User activity details with login history

### Supporting Types (3)

1. **LoginRecord**: UserID, LoginTime, IPAddress, UserAgent
2. **AdminRoleFilter**: Search, AccessLevel, MinLevel, MaxLevel, IsSystemRole, CreatedBy, Page, PageSize, SortBy, SortOrder
3. **AdminUserFilter**: Search, Status, RoleID, MinAccessLevel, Has2FA, CreatedBy, LastLoginAfter, LastLoginBefore, CreatedAfter, CreatedBefore, Page, PageSize, SortBy, SortOrder

## Business Features

### 1. Role-Based Access Control (RBAC)

- Hierarchical access levels (1-10)
- System roles (protected from modification)
- Custom roles (organization-specific)
- Role inheritance and permission checks

### 2. Authentication & Security

- Email/password authentication
- JWT tokens (access + refresh)
- Two-factor authentication (TOTP)
- Password reset flow with tokens
- Last login tracking

### 3. User Management

- Admin user lifecycle (create, update, delete)
- Status management (active, inactive, suspended)
- Profile management with images
- Self-referential creator tracking

### 4. Authorization

- Permission checks based on access level
- Hierarchical modification rules (can only modify lower-level admins)
- Role-based access to features
- Action authorization (CanModifyUser, CanModifyRole)

### 5. Audit & Tracking

- Creator tracking for all entities
- Login history and activity logs
- Activity statistics and analytics
- Recent logins and user activity

### 6. Security Features

- Password hashing (bcrypt)
- Two-factor authentication (TOTP)
- Token-based authentication (JWT)
- Status-based access control
- Suspension with reason tracking

## Technical Features

### 1. Data Persistence

- GORM models with proper relationships
- Soft delete support (DeletedAt)
- UUID support for public identifiers
- Indexed fields for performance

### 2. Validation

- Struct validation tags
- Email validation
- Password strength requirements
- Access level constraints (1-10)
- Status enum validation

### 3. Relationships

- Self-referential (AdminUser.Creator, AdminUser.CreatedUsers)
- One-to-many (AdminRole.Users, AdminUser.CreatedRoles)
- Cascade constraints (SET NULL on delete)

### 4. Filtering & Search

- Comprehensive filter types
- Full-text search
- Date range filtering
- Status and role filtering
- Pagination support

### 5. Statistics

- User statistics by status and role
- Role statistics by access level
- Activity analytics with time series
- Login statistics and trends

## Key Workflows

### Admin Creation Flow

1. Super admin creates new admin user
2. Assign role with appropriate access level
3. Send invitation email with password setup link
4. Admin sets password and enables 2FA
5. Admin can now login and perform authorized actions

### Authentication Flow

1. Admin enters email and password
2. System validates credentials
3. If 2FA enabled, prompt for TOTP code
4. System verifies 2FA code
5. Generate access and refresh tokens
6. Track login time and return tokens

### Role Assignment Flow

1. Check if actor has permission to assign role
2. Validate role exists and is active
3. Check actor's access level > target role's level
4. Assign role to user
5. Update user's effective permissions
6. Log role assignment for audit

### Suspension Flow

1. Admin initiates suspension with reason
2. Check if actor has permission to suspend target
3. Set user status to suspended
4. Invalidate all active tokens
5. Log suspension event
6. Notify user of suspension

## Statistics

### Admin Domain Summary

- **Entities**: 2 (AdminRole, AdminUser)
- **Repository Methods**: 61 total (20 + 41)
- **Service Methods**: 64 total (14 + 50)
- **Request DTOs**: 8
- **Response DTOs**: 13
- **Supporting Types**: 3
- **Total Lines**: ~580 (entity: ~140, repository: ~150, service: ~290)

### Complexity Analysis

- **Authentication**: Login, JWT tokens, 2FA, password reset
- **Authorization**: RBAC with access levels, hierarchical permissions
- **Audit**: Creator tracking, login history, activity logs
- **Security**: Bcrypt, TOTP, JWT, status-based access control
- **Statistics**: User stats, role stats, activity analytics

### Integration Points

- **Application Domain**: AdminUser referenced in ApplicationNote, ApplicationDocument, JobApplicationStage, Interview (verified_by, handled_by, author_id, interviewer_id)
- **Job Domain**: AdminUser referenced in Job moderation
- **User/Company Domains**: AdminUser for verification, moderation

## Next Steps

1. **Implementation Phase**:

   - Implement AdminRoleRepository with GORM
   - Implement AdminUserRepository with GORM
   - Implement AdminRoleService with business logic
   - Implement AdminUserService with authentication logic

2. **Security Implementation**:

   - JWT token generation and validation
   - TOTP 2FA implementation
   - Password reset token generation
   - Token blacklist for logout

3. **API Layer**:

   - Admin authentication endpoints
   - Role management endpoints
   - User management endpoints
   - Profile management endpoints
   - Statistics endpoints

4. **Testing**:

   - Unit tests for repositories
   - Unit tests for services
   - Integration tests for authentication flow
   - Integration tests for authorization checks
   - E2E tests for admin workflows

5. **Documentation**:
   - API documentation (Swagger)
   - Authentication guide
   - Role configuration guide
   - Security best practices
