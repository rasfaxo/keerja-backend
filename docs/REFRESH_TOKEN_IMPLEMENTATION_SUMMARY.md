# Refresh Token & Device Session Management - Implementation Summary

## Overview

Implementasi lengkap sistem refresh token dengan device session management untuk fitur "Remember Me" dan multi-device login tracking.

---

## Requirements (Completed)

✅ **Persistent Device Session ("Remember Me")**

- Buat refresh token dengan masa hidup lebih panjang untuk perangkat tepercaya
- Default: 30 hari
- Remember me: 90 hari

✅ **Simpan Sesi Per-Device**

- Buat tabel `sessions` / `refresh_tokens` dengan:
  - Device info (name, type, ID)
  - IP address
  - Last used timestamp
  - Revoked flag dengan reason tracking

✅ **Multi-Device Support**

- Max 5 perangkat aktif per user
- Auto-revoke oldest device saat limit exceeded
- Device fingerprinting dari user agent

✅ **Token Security**

- 64-byte cryptographically secure random tokens
- SHA256 hashed before storage (never store plaintext)
- Base64 encoded for transport

✅ **Token Rotation**

- New refresh token issued on each access token refresh
- Old token automatically revoked dengan reason "rotated"

✅ **Session Management**

- List active devices
- Revoke specific device
- Logout all devices
- Auto-cleanup expired/revoked tokens

---

## Files Created/Modified

### 1. Database Migration

**`database/migrations/create_refresh_tokens_table.sql`**

```sql
CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,  -- SHA256 hash
    device_name VARCHAR(255),
    device_type VARCHAR(50),
    device_id VARCHAR(255),
    user_agent TEXT,
    ip_address VARCHAR(45),
    last_used_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP,
    revoked_reason VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 6 indexes for performance
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_revoked ON refresh_tokens(revoked);
CREATE INDEX idx_refresh_tokens_device_id ON refresh_tokens(device_id);
CREATE INDEX idx_refresh_tokens_user_device ON refresh_tokens(user_id, device_id);
```

**Status**: Executed successfully

- Table created with 15 columns
- 8 indexes created (6 regular + 2 unique)
- Foreign key constraint to `users` table

---

### 2. Domain Layer

#### **`internal/domain/auth/entity.go`** (Modified)

Added `RefreshToken` entity:

```go
type RefreshToken struct {
    ID            int64
    UserID        int64
    TokenHash     string    // SHA256 hash (64 chars)
    DeviceName    *string
    DeviceType    *string   // mobile, desktop, tablet, unknown
    DeviceID      *string
    UserAgent     *string
    IPAddress     *string
    LastUsedAt    *time.Time
    ExpiresAt     time.Time
    Revoked       bool
    RevokedAt     *time.Time
    RevokedReason *string
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// Helper Methods
func (rt *RefreshToken) IsExpired() bool
func (rt *RefreshToken) IsValid() bool
func (rt *RefreshToken) Revoke(reason string)
func (rt *RefreshToken) UpdateLastUsed()
```

#### **`internal/domain/auth/repository.go`** (Modified)

Added `RefreshTokenRepository` interface (14 methods):

```go
type RefreshTokenRepository interface {
    Create(ctx context.Context, token *RefreshToken) error
    FindByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
    FindByUserID(ctx context.Context, userID int64) ([]RefreshToken, error)
    FindActiveByUserID(ctx context.Context, userID int64) ([]RefreshToken, error)
    FindByUserAndDevice(ctx context.Context, userID int64, deviceID string) (*RefreshToken, error)
    Update(ctx context.Context, token *RefreshToken) error
    UpdateLastUsed(ctx context.Context, id int64, lastUsedAt time.Time) error
    Revoke(ctx context.Context, id int64, reason string) error
    RevokeAllByUserID(ctx context.Context, userID int64, reason string) error
    RevokeByDeviceID(ctx context.Context, userID int64, deviceID string, reason string) error
    DeleteExpired(ctx context.Context) error
    DeleteRevoked(ctx context.Context, days int) error
    CountActiveByUserID(ctx context.Context, userID int64) (int, error)
}
```

---

### 3. Repository Layer

#### **`internal/repository/postgres/auth_repository.go`** (Modified)

Added `refreshTokenRepository` implementation (~150 lines):

```go
type refreshTokenRepository struct {
    db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) auth.RefreshTokenRepository {
    return &refreshTokenRepository{db: db}
}

// All 14 methods implemented with:
// - Context support
// - Null handling for nullable fields
// - Composite queries (user_id + device_id)
// - Atomic updates
// - Soft delete support
```

**Key Implementation Details:**

- Uses GORM with context support
- Proper null handling for optional fields (`*string`, `*time.Time`)
- Composite index usage (user_id + device_id)
- Batch operations for revocation
- Efficient cleanup queries

---

### 4. Service Layer

#### **`internal/service/refresh_token_service.go`** (Created - 367 lines)

Complete business logic for token lifecycle:

```go
// Configuration
const (
    RefreshTokenLength       = 64   // 64 bytes random
    RefreshTokenExpiryDays   = 30   // Default expiry
    RememberMeExpiryDays     = 90   // Remember me expiry
    MaxActiveTokensPerUser   = 5    // Device limit
    RefreshTokenRotation     = true // Enable rotation
)

type RefreshTokenService struct {
    refreshTokenRepo auth.RefreshTokenRepository
    jwtSecret        string
    jwtDuration      time.Duration
}

// Core Methods
func (s *RefreshTokenService) generateRefreshToken() (string, error)
func (s *RefreshTokenService) hashRefreshToken(token string) string
func (s *RefreshTokenService) parseDeviceType(userAgent string) string
func (s *RefreshTokenService) parseDeviceName(userAgent string) string
func (s *RefreshTokenService) CreateRefreshToken(ctx, userID, deviceInfo, rememberMe) (string, error)
func (s *RefreshTokenService) RefreshAccessToken(ctx, refreshToken, userID, email, userType) (string, string, error)
func (s *RefreshTokenService) RevokeRefreshToken(ctx, token, reason) error
func (s *RefreshTokenService) RevokeAllUserTokens(ctx, userID, reason) error
func (s *RefreshTokenService) RevokeDeviceToken(ctx, userID, deviceID, reason) error
func (s *RefreshTokenService) GetUserDevices(ctx, userID) ([]auth.RefreshToken, error)
func (s *RefreshTokenService) CleanupExpiredTokens(ctx) error
func (s *RefreshTokenService) CleanupRevokedTokens(ctx) error
```

**Key Features:**

1. **Token Generation**: 64 bytes crypto/rand → base64 → SHA256 hash
2. **Device Detection**: Parse user agent for device name/type
3. **Device Limit Enforcement**: Auto-revoke oldest when > 5 devices
4. **Token Rotation**: Issue new token on refresh, revoke old
5. **Validation**: Check hash, expiry, revoked status, user match
6. **Maintenance**: Cleanup expired/revoked tokens

**Custom Errors:**

```go
var (
    ErrRefreshTokenNotFound = errors.New("refresh token not found")
    ErrRefreshTokenExpired  = errors.New("refresh token has expired")
    ErrRefreshTokenRevoked  = errors.New("refresh token has been revoked")
)
```

**Device Type Detection:**

```go
// From user agent string:
mobile   → "Mobile", "Android", "iPhone", "iPad"
desktop  → "Windows", "Macintosh", "Linux"
tablet   → "Tablet"
unknown  → everything else
```

**Device Name Examples:**

- "Chrome on Windows"
- "Safari on iPhone"
- "Firefox on Linux"
- "Unknown Browser"

---

### 5. Handler Layer

#### **`internal/handler/http/auth_handler.go`** (Modified)

Added `refreshTokenService` field and 5 new handler methods:

**Updated Constructor:**

```go
func NewAuthHandler(
    authService *service.AuthService,
    otpService *service.OTPService,
    oauthService *service.OAuthService,
    registrationService *service.RegistrationService,
    refreshTokenService *service.RefreshTokenService, // NEW
    userRepo user.UserRepository,
) *AuthHandler
```

**New Handler Methods:**

1. **`LoginWithRememberMe(c *fiber.Ctx) error`**

   - Login + create refresh token
   - Parse device info from headers (User-Agent, IP)
   - Respect `remember_me` flag (30/90 days)
   - Return both access_token and refresh_token

2. **`RefreshAccessToken(c *fiber.Ctx) error`**

   - Validate refresh token
   - Issue new access token (JWT)
   - Rotate refresh token (issue new, revoke old)
   - Update last_used_at timestamp

3. **`GetActiveDevices(c *fiber.Ctx) error`**

   - List all active device sessions
   - Show device name, type, IP, last used
   - Mark current device (TODO: implement detection)

4. **`RevokeDevice(c *fiber.Ctx) error`**

   - Revoke specific device by device_id
   - Immediate effect on next token use

5. **`LogoutAllDevices(c *fiber.Ctx) error`**
   - Revoke all user tokens
   - Force re-login on all devices

---

### 6. DTO Layer

#### **`internal/dto/request/auth_request.go`** (Modified)

Added 3 new request DTOs:

```go
type LoginWithRememberMeRequest struct {
    Email      string `json:"email" validate:"required,email"`
    Password   string `json:"password" validate:"required,min=6"`
    RememberMe bool   `json:"remember_me"`
    DeviceID   string `json:"device_id"`
}

type RefreshAccessTokenRequest struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
}

type RevokeDeviceRequest struct {
    DeviceID string `json:"device_id" validate:"required"`
}
```

#### **`internal/dto/response/auth_response.go`** (Modified)

Added 2 new response DTOs:

```go
type DeviceInfo struct {
    ID         int64   `json:"id"`
    DeviceName *string `json:"device_name"`
    DeviceType *string `json:"device_type"`
    IPAddress  *string `json:"ip_address"`
    LastUsedAt string  `json:"last_used_at"`
    CreatedAt  string  `json:"created_at"`
    IsCurrent  bool    `json:"is_current"`
}

type DeviceListResponse struct {
    Devices []DeviceInfo `json:"devices"`
    Total   int          `json:"total"`
}
```

**Note**: `TokenResponse` already existed, used for refresh endpoint response.

---

### 7. Routes Layer

#### **`internal/routes/auth_routes.go`** (Modified)

Added 5 new routes:

```go
// Public route (rate limited)
auth.Post("/login-remember",
    middleware.AuthRateLimiter(),
    deps.AuthHandler.LoginWithRememberMe,
)

// Protected routes (require valid access token)
auth.Post("/refresh",
    authMw.AuthRequired(),
    deps.AuthHandler.RefreshAccessToken,
)

auth.Get("/devices",
    authMw.AuthRequired(),
    deps.AuthHandler.GetActiveDevices,
)

auth.Post("/devices/revoke",
    authMw.AuthRequired(),
    deps.AuthHandler.RevokeDevice,
)

auth.Post("/logout-all",
    authMw.AuthRequired(),
    deps.AuthHandler.LogoutAllDevices,
)
```

---

### 8. Main Application

#### **`cmd/main.go`** (Modified)

Updated dependency injection:

```go
// Create refresh token repository
refreshTokenRepo := postgres.NewRefreshTokenRepository(db)

// Create refresh token service
refreshTokenService := service.NewRefreshTokenService(
    refreshTokenRepo,
    cfg.JWTSecret,
    time.Duration(cfg.JWTExpirationHours)*time.Hour,
)

// Initialize AuthHandler with RefreshTokenService
authHandler := http.NewAuthHandler(
    authService,
    otpService,
    oauthService,
    registrationService,
    refreshTokenService,  // NEW parameter
    userRepo,
)
```

---

### 9. Documentation

#### **`docs/REFRESH_TOKEN_TESTING.md`** (Created - ~1000 lines)

Comprehensive testing guide:

- Architecture overview
- Security features
- API endpoint documentation (5 endpoints)
- cURL examples for all scenarios
- Postman collection (ready to import)
- 7 testing scenarios (multi-device, token rotation, revocation, etc.)
- Database verification queries
- Rate limiting tests
- Error codes reference
- Troubleshooting guide
- Monitoring queries
- Maintenance jobs

---

## Security Implementation

### Token Security

1. **Generation**: 64 bytes from `crypto/rand` (cryptographically secure)
2. **Encoding**: Base64 for transport
3. **Hashing**: SHA256 before storage (never store plaintext)
4. **Transport**: HTTPS only (recommended)

### Token Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│  1. LOGIN                                                   │
│  ┌──────────┐                                               │
│  │  User    │ ──Login Request──> Server                    │
│  │          │ <──Tokens─────────  ├─ Generate 64-byte      │
│  │          │                      ├─ Base64 encode        │
│  │          │                      ├─ SHA256 hash          │
│  │          │                      └─ Store hash in DB     │
│  └──────────┘                                               │
│  Receives:                                                  │
│  - access_token (JWT, 24h)                                  │
│  - refresh_token (plaintext, 30/90 days)                    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  2. ACCESS TOKEN EXPIRES (After 24h)                        │
│  ┌──────────┐                                               │
│  │  User    │ ──Refresh Request─> Server                   │
│  │          │    (expired access + refresh token)          │
│  │          │                      ├─ Hash incoming token  │
│  │          │                      ├─ Find in DB           │
│  │          │                      ├─ Validate (expiry,    │
│  │          │                      │   revoked, user)      │
│  │          │                      ├─ Issue new JWT        │
│  │          │                      ├─ Generate new refresh │
│  │          │                      └─ Revoke old refresh   │
│  │          │ <──New Tokens─────   (reason: "rotated")    │
│  └──────────┘                                               │
│  Receives:                                                  │
│  - new access_token (JWT, 24h)                              │
│  - new refresh_token (plaintext, 30/90 days)                │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  3. REVOCATION                                              │
│  ┌──────────┐                                               │
│  │  User    │ ──Logout/Revoke──> Server                    │
│  │          │                     ├─ Set revoked = true    │
│  │          │                     ├─ Set revoked_at        │
│  │          │                     └─ Set revoked_reason    │
│  │          │ <──Success─────────                          │
│  └──────────┘                                               │
│  Next refresh attempt: 401 Unauthorized                     │
└─────────────────────────────────────────────────────────────┘
```

### Device Limit Enforcement

```go
// When creating new token:
activeCount := s.refreshTokenRepo.CountActiveByUserID(ctx, userID)

if activeCount >= MaxActiveTokensPerUser {
    // Find oldest device
    tokens, _ := s.refreshTokenRepo.FindActiveByUserID(ctx, userID)
    oldestToken := tokens[0] // Sorted by created_at ASC

    // Revoke oldest
    s.refreshTokenRepo.Revoke(ctx, oldestToken.ID, "device_limit_exceeded")
}

// Create new token
// ...
```

### Token Rotation Strategy

```go
// On refresh:
1. Validate incoming refresh token
2. Generate new access token (JWT)
3. Generate new refresh token
4. Save new refresh token to DB
5. Revoke old refresh token (reason: "rotated")
6. Return both new tokens to client

// Benefits:
- Prevents replay attacks
- Limits window of compromise
- Tracks token usage lineage
```

---

## Database Schema

### Refresh Tokens Table

```
┌─────────────────┬──────────────┬─────────────┬──────────────┐
│ Column          │ Type         │ Nullable    │ Index        │
├─────────────────┼──────────────┼─────────────┼──────────────┤
│ id              │ BIGSERIAL    │ NO          │ PRIMARY KEY  │
│ user_id         │ BIGINT       │ NO          │ INDEXED + FK │
│ token_hash      │ VARCHAR(64)  │ NO          │ UNIQUE       │
│ device_name     │ VARCHAR(255) │ YES         │              │
│ device_type     │ VARCHAR(50)  │ YES         │              │
│ device_id       │ VARCHAR(255) │ YES         │ INDEXED      │
│ user_agent      │ TEXT         │ YES         │              │
│ ip_address      │ VARCHAR(45)  │ YES         │              │
│ last_used_at    │ TIMESTAMP    │ YES         │              │
│ expires_at      │ TIMESTAMP    │ NO          │ INDEXED      │
│ revoked         │ BOOLEAN      │ NO          │ INDEXED      │
│ revoked_at      │ TIMESTAMP    │ YES         │              │
│ revoked_reason  │ VARCHAR(100) │ YES         │              │
│ created_at      │ TIMESTAMP    │ NO          │              │
│ updated_at      │ TIMESTAMP    │ NO          │              │
└─────────────────┴──────────────┴─────────────┴──────────────┘

Indexes (8 total):
1. refresh_tokens_pkey (id) - PRIMARY
2. refresh_tokens_token_hash_key (token_hash) - UNIQUE
3. idx_refresh_tokens_user_id (user_id)
4. idx_refresh_tokens_token_hash (token_hash)
5. idx_refresh_tokens_expires_at (expires_at)
6. idx_refresh_tokens_revoked (revoked)
7. idx_refresh_tokens_device_id (device_id)
8. idx_refresh_tokens_user_device (user_id, device_id) - COMPOSITE

Foreign Keys:
- user_id → users(id) ON DELETE CASCADE
```

---

## Testing Status

### Build & Compilation

```bash
$ go build -o keerja-api.exe ./cmd
# ✅ Success - No errors
```

### ✅ Database Migration

```bash
$ PGPASSWORD=bekerja123 psql -U bekerja -h localhost -p 5434 -d keerja \
    -f database/migrations/create_refresh_tokens_table.sql

CREATE TABLE
CREATE INDEX (x6)
COMMENT (x13)

# Verification:
   table_name   | column_count
----------------+--------------
 refresh_tokens |           15

# ✅ All 8 indexes created successfully
```

### Manual Testing (Pending)

Need to test:

1. Login with remember me (30 days vs 90 days)
2. Token refresh flow
3. Device listing
4. Device revocation
5. Logout all devices
6. Multi-device limit (max 5)
7. Token rotation
8. Rate limiting

**Next Step**: Start server and run testing scenarios from `REFRESH_TOKEN_TESTING.md`

---

## API Endpoints Summary

| Method | Endpoint                      | Auth                     | Description            |
| ------ | ----------------------------- | ------------------------ | ---------------------- |
| POST   | `/api/v1/auth/login-remember` | ❌ Public (rate limited) | Login with remember me |
| POST   | `/api/v1/auth/refresh`        | ✅ Protected             | Refresh access token   |
| GET    | `/api/v1/auth/devices`        | ✅ Protected             | List active devices    |
| POST   | `/api/v1/auth/devices/revoke` | ✅ Protected             | Revoke specific device |
| POST   | `/api/v1/auth/logout-all`     | ✅ Protected             | Logout all devices     |

**Rate Limiting:**

- Login endpoint: 5 requests per 15 minutes per IP
- Other endpoints: No rate limit (protected by auth)

---

## Configuration

### Constants (Service Layer)

```go
const (
    RefreshTokenLength     = 64    // 64 bytes = 512 bits
    RefreshTokenExpiryDays = 30    // Default: 30 days
    RememberMeExpiryDays   = 90    // Remember me: 90 days
    MaxActiveTokensPerUser = 5     // Max devices per user
    RefreshTokenRotation   = true  // Enable token rotation
)
```

### Environment Variables (No changes required)

```env
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION_HOURS=24  # Access token: 24 hours
```

---

## Monitoring & Maintenance

### Recommended Cron Jobs

1. **Cleanup Expired Tokens** (Daily)

```go
// Delete expired tokens
s.refreshTokenRepo.DeleteExpired(ctx)
```

2. **Cleanup Old Revoked Tokens** (Weekly)

```go
// Delete revoked tokens older than 90 days
s.refreshTokenRepo.DeleteRevoked(ctx, 90)
```

### Monitoring Queries

**Active Sessions per User:**

```sql
SELECT
    u.email,
    COUNT(rt.id) as active_devices,
    MAX(rt.last_used_at) as last_activity
FROM users u
LEFT JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.revoked = false AND rt.expires_at > NOW()
GROUP BY u.id, u.email
ORDER BY active_devices DESC;
```

**Token Distribution:**

```sql
SELECT
    device_type,
    COUNT(*) as total,
    ROUND(100.0 * COUNT(*) / SUM(COUNT(*)) OVER(), 2) as percentage
FROM refresh_tokens
WHERE revoked = false AND expires_at > NOW()
GROUP BY device_type;
```

---

## Checklist

### Implementation

- [x] Create database migration
- [x] Add RefreshToken entity
- [x] Implement RefreshTokenRepository interface
- [x] Implement repository (GORM)
- [x] Create RefreshTokenService
- [x] Implement token generation (crypto/rand)
- [x] Implement token hashing (SHA256)
- [x] Implement device detection
- [x] Implement device limit enforcement
- [x] Implement token rotation
- [x] Add handler methods (5 endpoints)
- [x] Create request DTOs
- [x] Create response DTOs
- [x] Update routes
- [x] Wire dependencies in main.go
- [x] Run database migration
- [x] Build application
- [x] Create testing documentation

### Testing (Pending)

- [ ] Test login with remember me (false)
- [ ] Test login with remember me (true)
- [ ] Test token refresh
- [ ] Test token rotation
- [ ] Test device listing
- [ ] Test device revocation
- [ ] Test logout all devices
- [ ] Test multi-device limit (> 5 devices)
- [ ] Test expired token handling
- [ ] Test revoked token handling
- [ ] Test rate limiting
- [ ] Verify database indexes
- [ ] Run cleanup jobs

### Documentation

- [x] API endpoint documentation
- [x] cURL examples
- [x] Postman collection
- [x] Testing scenarios
- [x] Security best practices
- [x] Troubleshooting guide
- [x] Monitoring queries
- [x] Implementation summary

---

## Completion Status

| Category             | Status      | Progress |
| -------------------- | ----------- | -------- |
| Database Schema      | ✅ Complete | 100%     |
| Domain Layer         | ✅ Complete | 100%     |
| Repository Layer     | ✅ Complete | 100%     |
| Service Layer        | ✅ Complete | 100%     |
| Handler Layer        | ✅ Complete | 100%     |
| DTO Layer            | ✅ Complete | 100%     |
| Routes               | ✅ Complete | 100%     |
| Dependency Injection | ✅ Complete | 100%     |
| Build & Compilation  | ✅ Success  | 100%     |
| Database Migration   | ✅ Executed | 100%     |
| Documentation        | ✅ Complete | 100%     |
| Manual Testing       | 🔄 Pending  | 0%       |

**Overall Progress: 95%** (Implementation complete, testing pending)

---

## Notes

### Token Storage Best Practices

**Client-Side:**

- ✅ GOOD: httpOnly cookies (set by server)
- ✅ GOOD: Secure encrypted storage (mobile apps)
- ❌ BAD: localStorage (vulnerable to XSS)
- ❌ BAD: sessionStorage (still vulnerable)

**Server-Side:**

- ✅ Always hash tokens before storage (SHA256)
- ✅ Never log plaintext tokens
- ✅ Use HTTPS for transport
- ✅ Implement rate limiting
- ✅ Monitor suspicious activity

### Future Enhancements

1. **Current Device Detection**: Mark current device in device list
2. **Device Notifications**: Email when new device logs in
3. **Geolocation**: Add country/city from IP address
4. **Browser Fingerprinting**: Enhanced device tracking
5. **2FA Integration**: Require 2FA for new devices
6. **Risk-Based Auth**: Challenge suspicious logins

---

## Related Documentation

- **Testing Guide**: `docs/REFRESH_TOKEN_TESTING.md`
- **OTP Registration**: `docs/OTP_REGISTRATION_TESTING.md`
- **Authentication**: `docs/AUTHENTICATION.md`
- **API Documentation**: `docs/API_DOCUMENTATION.md`

---

## Support

For issues or questions:

1. Check `REFRESH_TOKEN_TESTING.md` for troubleshooting
2. Verify database migration executed successfully
3. Check application logs for errors
4. Test endpoints with Postman collection

---
