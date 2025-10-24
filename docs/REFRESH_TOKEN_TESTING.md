# Refresh Token & Device Session Management - Testing Guide

## Overview

Sistem refresh token memungkinkan user untuk:

- **Remember Me**: Login dengan persistent session (30/90 hari)
- **Multi-Device**: Kelola hingga 5 perangkat aktif secara bersamaan
- **Device Tracking**: Monitor perangkat yang login (nama, tipe, IP, last used)
- **Session Revocation**: Logout dari device tertentu atau semua device

## Architecture

```
┌─────────────────┐
│  Client Device  │
│  (Browser/App)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Auth Handler   │  ← LoginWithRememberMe()
│                 │  ← RefreshAccessToken()
│                 │  ← GetActiveDevices()
│                 │  ← RevokeDevice()
│                 │  ← LogoutAllDevices()
└────────┬────────┘
         │
         ▼
┌──────────────────────────┐
│  RefreshTokenService     │
│  • Generate 64-byte token│
│  • SHA256 hashing        │
│  • Device detection      │
│  • Token rotation        │
│  • Max 5 devices/user    │
│  • Expiry: 30/90 days    │
└────────┬─────────────────┘
         │
         ▼
┌─────────────────────────┐
│  RefreshTokenRepository │
│  (GORM + PostgreSQL)    │
└─────────────────────────┘
```

## Security Features

1. **Token Security**

   - 64 bytes cryptographically secure random
   - SHA256 hashed before storage (never store plaintext)
   - Base64 encoded for transport

2. **Device Limit**

   - Maximum 5 active devices per user
   - Oldest device auto-revoked when limit exceeded

3. **Token Rotation**

   - New refresh token issued on each access token refresh
   - Old token automatically revoked

4. **Expiry Management**

   - Default: 30 days
   - Remember me: 90 days
   - Auto-cleanup of expired tokens

5. **Revocation**
   - Single device logout
   - All devices logout
   - Immediate effect on next token use

## API Endpoints

### 1. Login with Remember Me

**POST** `/api/v1/auth/login-remember`

Login dan dapatkan refresh token untuk persistent session.

**Request Body:**

```json
{
  "email": "john.doe@example.com",
  "password": "SecurePass123!",
  "remember_me": true,
  "device_id": "browser-chrome-windows-xyz123"
}
```

**Fields:**

- `email` (required): User email
- `password` (required): User password
- `remember_me` (optional, default: false):
  - `false` = 30 days expiry
  - `true` = 90 days expiry
- `device_id` (optional): Unique device identifier untuk multi-device tracking

**Success Response (200):**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "ZjQzMmRlZmE4YjQ3MzJkZjg5YWJjZGVm...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "full_name": "John Doe",
      "email": "john.doe@example.com",
      "phone": "628123456789",
      "user_type": "jobseeker",
      "is_verified": true,
      "status": "active"
    }
  }
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login-remember \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123!",
    "remember_me": true,
    "device_id": "browser-chrome-windows"
  }'
```

---

### 2. Refresh Access Token

**POST** `/api/v1/auth/refresh`

Dapatkan access token baru menggunakan refresh token.

**Headers:**

```
Authorization: Bearer <current_access_token>
```

**Request Body:**

```json
{
  "refresh_token": "ZjQzMmRlZmE4YjQ3MzJkZjg5YWJjZGVm..."
}
```

**Success Response (200):**

```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "YzEyM2RlZmE4YjQ3MzJkZjg5YWJjZGVm...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

**Notes:**

- Old refresh token will be **revoked** (token rotation)
- New refresh token will be issued
- Last used timestamp updated
- Access token expires in 24 hours

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "refresh_token": "ZjQzMmRlZmE4YjQ3MzJkZjg5YWJjZGVm..."
  }'
```

**Error Responses:**

1. **Invalid Token (401)**

```json
{
  "success": false,
  "message": "Invalid refresh token",
  "error": "token not found or invalid"
}
```

2. **Expired Token (401)**

```json
{
  "success": false,
  "message": "Refresh token expired",
  "error": "token has expired"
}
```

3. **Revoked Token (401)**

```json
{
  "success": false,
  "message": "Refresh token has been revoked",
  "error": "token was revoked"
}
```

---

### 3. Get Active Devices

**GET** `/api/v1/auth/devices`

Dapatkan list semua perangkat yang aktif login.

**Headers:**

```
Authorization: Bearer <access_token>
```

**Success Response (200):**

```json
{
  "success": true,
  "message": "Active devices retrieved successfully",
  "data": {
    "devices": [
      {
        "id": 1,
        "device_name": "Chrome on Windows",
        "device_type": "desktop",
        "ip_address": "192.168.1.100",
        "last_used_at": "2024-01-15 14:30:00",
        "created_at": "2024-01-01 10:00:00",
        "is_current": true
      },
      {
        "id": 2,
        "device_name": "Safari on iPhone",
        "device_type": "mobile",
        "ip_address": "192.168.1.101",
        "last_used_at": "2024-01-14 08:15:00",
        "created_at": "2024-01-05 12:30:00",
        "is_current": false
      }
    ],
    "total": 2
  }
}
```

**Device Types:**

- `mobile`: Smartphone/tablet
- `desktop`: Desktop/laptop browser
- `tablet`: Tablet device
- `unknown`: Unidentified device

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/auth/devices \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

### 4. Revoke Device Session

**POST** `/api/v1/auth/devices/revoke`

Logout dari perangkat tertentu (revoke refresh token).

**Headers:**

```
Authorization: Bearer <access_token>
```

**Request Body:**

```json
{
  "device_id": "browser-chrome-windows"
}
```

**Success Response (200):**

```json
{
  "success": true,
  "message": "Device session revoked successfully",
  "data": null
}
```

**Use Cases:**

- User melihat device yang tidak dikenal
- User ingin logout dari device tertentu tanpa akses fisik
- Security: Revoke compromised device

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/devices/revoke \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "device_id": "browser-chrome-windows"
  }'
```

---

### 5. Logout All Devices

**POST** `/api/v1/auth/logout-all`

Logout dari **semua** perangkat sekaligus.

**Headers:**

```
Authorization: Bearer <access_token>
```

**Success Response (200):**

```json
{
  "success": true,
  "message": "Logged out from all devices successfully",
  "data": null
}
```

**Use Cases:**

- Password changed → force re-login semua device
- Security breach → immediate revocation
- User wants fresh start

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout-all \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## Testing Scenarios

### Scenario 1: Login with Remember Me (30 Days)

```bash
# 1. Login without remember me (default: 30 days)
curl -X POST http://localhost:8080/api/v1/auth/login-remember \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123!",
    "remember_me": false,
    "device_id": "laptop-chrome"
  }'

# Expected: refresh_token valid for 30 days
```

### Scenario 2: Login with Remember Me (90 Days)

```bash
# 1. Login with remember me enabled
curl -X POST http://localhost:8080/api/v1/auth/login-remember \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123!",
    "remember_me": true,
    "device_id": "phone-safari"
  }'

# Expected: refresh_token valid for 90 days
```

### Scenario 3: Multi-Device Login

```bash
# Login from 3 different devices
for device in laptop-chrome desktop-firefox phone-safari; do
  curl -X POST http://localhost:8080/api/v1/auth/login-remember \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"john.doe@example.com\",
      \"password\": \"SecurePass123!\",
      \"remember_me\": true,
      \"device_id\": \"$device\"
    }"
done

# Check active devices
curl -X GET http://localhost:8080/api/v1/auth/devices \
  -H "Authorization: Bearer <access_token>"

# Expected: 3 devices listed
```

### Scenario 4: Token Refresh Flow

```bash
# 1. Save tokens from login
ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
REFRESH_TOKEN="ZjQzMmRlZmE4YjQ3MzJkZjg5YWJjZGVm..."

# 2. Wait for access token to expire (24 hours) OR test immediately

# 3. Refresh token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }"

# Expected:
# - New access_token
# - New refresh_token
# - Old refresh_token revoked
```

### Scenario 5: Device Limit (Max 5)

```bash
# Login from 6 devices
for i in {1..6}; do
  curl -X POST http://localhost:8080/api/v1/auth/login-remember \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"john.doe@example.com\",
      \"password\": \"SecurePass123!\",
      \"remember_me\": true,
      \"device_id\": \"device-$i\"
    }"
done

# Check devices
curl -X GET http://localhost:8080/api/v1/auth/devices \
  -H "Authorization: Bearer <access_token>"

# Expected: Only 5 devices (device-2 to device-6)
# device-1 automatically revoked (oldest)
```

### Scenario 6: Revoke Specific Device

```bash
# 1. Get device list
RESPONSE=$(curl -X GET http://localhost:8080/api/v1/auth/devices \
  -H "Authorization: Bearer <access_token>")

echo $RESPONSE

# 2. Revoke device
curl -X POST http://localhost:8080/api/v1/auth/devices/revoke \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "device_id": "laptop-chrome"
  }'

# 3. Try to use revoked token (should fail)
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "refresh_token": "<revoked_refresh_token>"
  }'

# Expected: 401 Unauthorized - "Refresh token has been revoked"
```

### Scenario 7: Logout All Devices

```bash
# 1. Login from multiple devices
# ... (see Scenario 3)

# 2. Logout all
curl -X POST http://localhost:8080/api/v1/auth/logout-all \
  -H "Authorization: Bearer <access_token>"

# 3. Try to use any refresh token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "refresh_token": "<any_refresh_token>"
  }'

# Expected: 401 Unauthorized - All tokens revoked
```

---

## Database Verification

### Check Refresh Tokens Table

```sql
-- Connect to database
PGPASSWORD=bekerja123 psql -U bekerja -h localhost -p 5434 -d keerja

-- View all refresh tokens
SELECT
    id,
    user_id,
    LEFT(token_hash, 20) || '...' as token_hash,
    device_name,
    device_type,
    ip_address,
    last_used_at,
    expires_at,
    revoked,
    created_at
FROM refresh_tokens
ORDER BY created_at DESC;

-- Count active tokens per user
SELECT
    user_id,
    COUNT(*) as active_tokens
FROM refresh_tokens
WHERE revoked = false
    AND expires_at > NOW()
GROUP BY user_id;

-- View token expiry distribution
SELECT
    device_type,
    COUNT(*) as total,
    COUNT(CASE WHEN revoked = false THEN 1 END) as active,
    COUNT(CASE WHEN revoked = true THEN 1 END) as revoked,
    COUNT(CASE WHEN expires_at < NOW() THEN 1 END) as expired
FROM refresh_tokens
GROUP BY device_type;
```

### Verify Token Rotation

```sql
-- Check if old token was revoked after refresh
SELECT
    token_hash,
    revoked,
    revoked_at,
    revoked_reason
FROM refresh_tokens
WHERE user_id = 1
ORDER BY created_at DESC
LIMIT 5;

-- Expected:
-- Latest token: revoked = false
-- Previous tokens: revoked = true, revoked_reason = 'rotated'
```

### Check Device Limits

```sql
-- Verify max 5 devices per user
SELECT
    user_id,
    COUNT(*) as active_devices
FROM refresh_tokens
WHERE revoked = false
    AND expires_at > NOW()
GROUP BY user_id
HAVING COUNT(*) > 5;

-- Expected: 0 rows (no user should have > 5 active devices)
```

---

## Rate Limiting

Endpoint `POST /auth/login-remember` menggunakan rate limiting:

- **Max**: 5 requests per 15 minutes per IP
- **Middleware**: `AuthRateLimiter()`

Test rate limit:

```bash
# Rapid requests (should block after 5)
for i in {1..10}; do
  echo "Request $i:"
  curl -X POST http://localhost:8080/api/v1/auth/login-remember \
    -H "Content-Type: application/json" \
    -d '{
      "email": "test@example.com",
      "password": "wrong"
    }'
  echo
done

# Expected:
# Requests 1-5: Normal responses
# Requests 6-10: 429 Too Many Requests
```

---

## Error Codes

| Status | Message                           | Cause                        | Solution                             |
| ------ | --------------------------------- | ---------------------------- | ------------------------------------ |
| 400    | Invalid request body              | Malformed JSON               | Check JSON syntax                    |
| 400    | Validation failed                 | Missing required fields      | Provide email/password/refresh_token |
| 401    | Invalid email or password         | Wrong credentials            | Check credentials                    |
| 401    | Email not verified                | Email not verified           | Verify email first                   |
| 401    | Invalid refresh token             | Token not found              | Login again                          |
| 401    | Refresh token expired             | Token expired (30/90 days)   | Login again                          |
| 401    | Refresh token has been revoked    | Token revoked manually       | Login again                          |
| 401    | Unauthorized                      | Missing/invalid access token | Provide valid access token           |
| 429    | Too many requests                 | Rate limit exceeded          | Wait 15 minutes                      |
| 500    | Failed to create refresh token    | Database error               | Check logs                           |
| 500    | Failed to refresh token           | Service error                | Check logs                           |
| 500    | Failed to get devices             | Database error               | Check logs                           |
| 500    | Failed to revoke device           | Database error               | Check logs                           |
| 500    | Failed to logout from all devices | Database error               | Check logs                           |

---

## Postman Collection

Import this collection for quick testing:

```json
{
  "info": {
    "name": "Keerja - Refresh Token API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Login with Remember Me",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"john.doe@example.com\",\n  \"password\": \"SecurePass123!\",\n  \"remember_me\": true,\n  \"device_id\": \"{{$randomUUID}}\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/v1/auth/login-remember",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "auth", "login-remember"]
        }
      }
    },
    {
      "name": "Refresh Access Token",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          },
          {
            "key": "Authorization",
            "value": "Bearer {{access_token}}"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"refresh_token\": \"{{refresh_token}}\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/v1/auth/refresh",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "auth", "refresh"]
        }
      }
    },
    {
      "name": "Get Active Devices",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{access_token}}"
          }
        ],
        "url": {
          "raw": "{{base_url}}/api/v1/auth/devices",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "auth", "devices"]
        }
      }
    },
    {
      "name": "Revoke Device",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          },
          {
            "key": "Authorization",
            "value": "Bearer {{access_token}}"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"device_id\": \"browser-chrome-windows\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/v1/auth/devices/revoke",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "auth", "devices", "revoke"]
        }
      }
    },
    {
      "name": "Logout All Devices",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{access_token}}"
          }
        ],
        "url": {
          "raw": "{{base_url}}/api/v1/auth/logout-all",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "auth", "logout-all"]
        }
      }
    }
  ],
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080"
    },
    {
      "key": "access_token",
      "value": ""
    },
    {
      "key": "refresh_token",
      "value": ""
    }
  ]
}
```

**Setup Postman:**

1. Import collection
2. Set environment variable `base_url` = `http://localhost:8080`
3. After login, copy `access_token` and `refresh_token` to environment variables
4. Use `{{access_token}}` and `{{refresh_token}}` in requests

---

## Security Best Practices

### Client-Side Implementation

1. **Store Tokens Securely**

   ```javascript
   // ❌ BAD: localStorage (vulnerable to XSS)
   localStorage.setItem("refresh_token", token);

   // ✅ GOOD: httpOnly cookie (set by server)
   // Or secure encrypted storage
   ```

2. **Auto-Refresh Before Expiry**

   ```javascript
   // Refresh 5 minutes before expiry
   const refreshTime = tokenExpiry - 5 * 60 * 1000;
   setTimeout(() => {
     refreshAccessToken();
   }, refreshTime);
   ```

3. **Handle Token Rotation**
   ```javascript
   async function refreshAccessToken() {
     const response = await fetch("/api/v1/auth/refresh", {
       method: "POST",
       headers: {
         Authorization: `Bearer ${accessToken}`,
         "Content-Type": "application/json",
       },
       body: JSON.stringify({
         refresh_token: refreshToken,
       }),
     });

     const data = await response.json();

     // Update both tokens
     accessToken = data.data.access_token;
     refreshToken = data.data.refresh_token; // New token!
   }
   ```

### Server-Side Security

1. **Token Hashing**: SHA256 before database storage
2. **Rate Limiting**: Prevent brute force attacks
3. **Device Limits**: Max 5 devices per user
4. **Token Rotation**: Prevent replay attacks
5. **Auto Cleanup**: Remove expired/revoked tokens

---

## Troubleshooting

### Issue: "Refresh token not found"

**Cause**: Token tidak ada di database atau sudah dihapus

**Solutions:**

1. Check if token was revoked: `SELECT * FROM refresh_tokens WHERE token_hash = '<hash>'`
2. Check if token expired
3. Login ulang untuk mendapatkan token baru

### Issue: "Too many active devices"

**Cause**: User sudah mencapai limit 5 devices

**Solutions:**

1. Revoke device lama: `POST /auth/devices/revoke`
2. Logout all: `POST /auth/logout-all`
3. Oldest device akan otomatis di-revoke saat login baru

### Issue: Token rotation not working

**Cause**: Old token masih bisa digunakan setelah refresh

**Check:**

```sql
SELECT token_hash, revoked, revoked_reason
FROM refresh_tokens
WHERE user_id = 1
ORDER BY created_at DESC
LIMIT 5;
```

**Expected**: Previous token `revoked = true, revoked_reason = 'rotated'`

---

## Maintenance Jobs

### Cleanup Expired Tokens

Run periodic cleanup (recommended: daily cron job):

```go
// In service layer
func (s *RefreshTokenService) CleanupExpiredTokens(ctx context.Context) error {
    return s.refreshTokenRepo.DeleteExpired(ctx)
}
```

**SQL equivalent:**

```sql
DELETE FROM refresh_tokens
WHERE expires_at < NOW();
```

### Cleanup Old Revoked Tokens

Remove revoked tokens older than 90 days:

```go
func (s *RefreshTokenService) CleanupRevokedTokens(ctx context.Context) error {
    return s.refreshTokenRepo.DeleteRevoked(ctx, 90)
}
```

**SQL equivalent:**

```sql
DELETE FROM refresh_tokens
WHERE revoked = true
    AND revoked_at < NOW() - INTERVAL '90 days';
```

---

## Monitoring Queries

### Active Sessions per User

```sql
SELECT
    u.id,
    u.email,
    COUNT(rt.id) as active_devices,
    MAX(rt.last_used_at) as last_activity
FROM users u
LEFT JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.revoked = false
    AND rt.expires_at > NOW()
GROUP BY u.id, u.email
ORDER BY active_devices DESC;
```

### Token Usage Statistics

```sql
SELECT
    DATE(created_at) as date,
    COUNT(*) as tokens_created,
    COUNT(CASE WHEN revoked THEN 1 END) as tokens_revoked,
    COUNT(CASE WHEN expires_at < NOW() THEN 1 END) as tokens_expired
FROM refresh_tokens
WHERE created_at >= NOW() - INTERVAL '30 days'
GROUP BY DATE(created_at)
ORDER BY date DESC;
```

### Device Type Distribution

```sql
SELECT
    device_type,
    COUNT(*) as total,
    ROUND(100.0 * COUNT(*) / SUM(COUNT(*)) OVER(), 2) as percentage
FROM refresh_tokens
WHERE revoked = false
    AND expires_at > NOW()
GROUP BY device_type
ORDER BY total DESC;
```

---

## Next Steps

1. **Testing**: Test all endpoints dengan Postman atau cURL
2. **Frontend Integration**: Implementasi auto-refresh dan token storage
3. **Monitoring**: Setup dashboard untuk tracking active sessions
4. **Security Audit**: Review token lifecycle dan revocation flow
5. **Documentation**: Update API docs dengan refresh token endpoints

---

## Summary

✅ **Implemented Features:**

- Login with remember me (30/90 days)
- Refresh access token dengan rotation
- Multi-device session management
- Device tracking (name, type, IP, last used)
- Single device revocation
- All devices logout
- Device limit enforcement (max 5)
- Auto-cleanup expired/revoked tokens

✅ **Security Features:**

- 64-byte cryptographically secure tokens
- SHA256 hashing (no plaintext storage)
- Token rotation on refresh
- Rate limiting on login
- Device fingerprinting
- Revocation tracking with reasons

✅ **Database:**

- `refresh_tokens` table created
- 8 indexes for performance
- Foreign key constraint to `users`
- Automated cleanup queries
