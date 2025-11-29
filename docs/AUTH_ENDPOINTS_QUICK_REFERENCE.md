# Authentication Endpoints

**Last Updated**: 2025-11-19  
**Base URL**: `http://localhost:8080/api/v1`

---

| Symbol | Meaning                   |
| ------ | ------------------------- |
| 🔓     | Public (no auth required) |
| 🔒     | Protected (JWT required)  |
| 🚦     | Rate limited              |
| ⏰     | Token/session expiry      |

---

## 🆕 NEW ENDPOINTS (15)

### OTP Registration

#### 1. Register with OTP

```http
POST /auth/register-otp 🔓 🚦
Content-Type: application/json

{
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "user_type": "jobseeker",
  "phone": "081234567890"  // optional
}

→ 201 Created
{
  "success": true,
  "message": "Registration successful. Please check your email for OTP verification code.",
  "data": {
    "email": "john@example.com",
    "note": "OTP code is valid for 5 minutes."
  }
}
```

🚦 Rate: 5 requests/hour  
⏰ OTP expires: 5 minutes

---

#### 2. Verify Email OTP

```http
POST /auth/verify-email-otp 🔓 🚦
Content-Type: application/json

{
  "email": "john@example.com",
  "otp_code": "123456"
}

→ 200 OK
{
  "success": true,
  "message": "Email verified successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": { ... }
  }
}
```

🚦 Rate: 10 requests/hour  
⏰ Max 5 attempts per OTP

---

#### 3. Resend OTP

```http
POST /auth/resend-otp 🔓 🚦
Content-Type: application/json

{
  "email": "john@example.com"
}

→ 200 OK
{
  "success": true,
  "message": "OTP code has been resent to your email.",
  "data": {
    "email": "john@example.com",
    "note": "OTP code is valid for 5 minutes."
  }
}
```

🚦 Rate: 3 requests/hour  
⏰ Must wait 60 seconds between requests

---

### Forgot Password with OTP

#### 4. Request Password Reset OTP

```http
POST /auth/forgot-password-otp 🔓 🚦
Content-Type: application/json

{
  "email": "john@example.com"
}

→ 200 OK
{
  "success": true,
  "message": "Password reset OTP has been sent to your email.",
  "data": {
    "email": "john@example.com",
    "note": "OTP code is valid for 5 minutes."
  }
}
```

🚦 Rate: 3 requests/hour per email  
⏰ OTP expires: 5 minutes  
⚠️ **Security**: Silent success if email doesn't exist

---

#### 5. Reset Password with OTP

```http
POST /auth/reset-password-otp 🔓 🚦
Content-Type: application/json

{
  "email": "john@example.com",
  "otp_code": "123456",
  "new_password": "NewSecurePass123!"
}

→ 200 OK
{
  "success": true,
  "message": "Password has been reset successfully. You can now login with your new password.",
  "data": null
}
```

🚦 Rate: 10 requests/hour  
⏰ Max 5 attempts per OTP  
📧 Confirmation email sent automatically

---

### Refresh Token Device Management

#### 6. Login with Remember Me

```http
POST /auth/login-remember 🔓 🚦
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!",
  "remember_me": true,         // optional, default: false
  "device_id": "chrome-001"    // optional, client-generated
}

→ 200 OK
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "a1b2c3d4e5f6g7h8i9j0...",  // 64-byte token
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": { ... }
  }
}
```

🚦 Rate: 10 requests/hour  
⏰ Refresh token expires:

- 30 days (remember_me = false)
- 90 days (remember_me = true)

---

#### 7. Refresh Access Token

```http
POST /auth/refresh 🔓
Content-Type: application/json
Authorization: Bearer <EXPIRED_OR_VALID_JWT>

{
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0..."
}

→ 200 OK
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",  // NEW JWT
    "refresh_token": "b2c3d4e5f6g7h8i9j0k1...",  // NEW refresh token
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

⚠️ **Important**: Old refresh token is REVOKED (token rotation)  
⏰ New JWT expires: 1 hour

---

#### 8. List Active Devices

```http
GET /auth/devices 🔒
Authorization: Bearer <JWT_TOKEN>

→ 200 OK
{
  "success": true,
  "message": "Active devices retrieved successfully",
  "data": {
    "devices": [
      {
        "id": 1,
        "device_name": "chrome-001",
        "device_type": "Chrome on Windows",
        "ip_address": "192.168.1.100",
        "last_used_at": "2025-10-24T10:30:00Z",
        "created_at": "2025-10-24T08:00:00Z",
        "is_current": true
      }
    ],
    "total": 1
  }
}
```

---

#### 9. Revoke Device

```http
POST /auth/devices/revoke 🔒
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>

{
  "device_id": "chrome-001"
}

→ 200 OK
{
  "success": true,
  "message": "Device session revoked successfully",
  "data": null
}
```

---

#### 10. Logout All Devices

```http
POST /auth/logout-all 🔒
Authorization: Bearer <JWT_TOKEN>

→ 200 OK
{
  "success": true,
  "message": "Logged out from all devices successfully",
  "data": null
}
```

**Important**: Revokes ALL refresh tokens for user

---

### Google OAuth

#### 11. Get Google OAuth URL

```http
GET /auth/oauth/google?client=mobile&redirect_uri=myapp://oauth-callback&code_challenge=<S256> 🔓

Query Params:
- `client` (optional): `web` (default) or `mobile`
- `redirect_uri` (optional for web): must exist in `ALLOWED_MOBILE_REDIRECT_URIS` when `client=mobile`
- `post_login_redirect_uri` (optional): Deep-link to receive JWT via fragment (`myapp://oauth#token=...`)
- `code_challenge` + `code_challenge_method` (optional): PKCE S256 for mobile flows

→ 200 OK
{
  "success": true,
  "message": "Google auth URL generated",
  "data": {
    "auth_url": "https://accounts.google.com/o/oauth2/v2/auth?...",
    "state": "J2lQJd9iXo...",
    "expires_in": 300
  }
}
```

**Usage**: Redirect/launch browser to `auth_url`, persist the returned `state`.

---

#### 12. Handle Google Callback

```http
GET /auth/oauth/google/callback?code=4/0AY0e-g6XXX&state=abc123 🔓

→ 200 OK
{
  "success": true,
  "message": "Google authentication successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

⚠️ **Note**: Automatically called by Google after user login  
💡 **Auto-create**: New user created if email doesn't exist  
📱 **Deep-link mode**: when `post_login_redirect_uri` was provided, backend redirects to `myapp://oauth-callback#token=<APP_JWT>` instead of returning JSON.

---

#### 13. Exchange Google Authorization Code (Mobile PKCE)

Used by Flutter once it receives the Google `code` via app deep-link.

```http
POST /auth/oauth/google/exchange 🔓
Content-Type: application/json

{
  "code": "4/0AY0e-g6XXX",
  "code_verifier": "generated-by-app",
  "state": "J2lQJd9iXo...",
  "redirect_uri": "myapp://oauth-callback"
}

→ 200 OK
{
  "success": true,
  "message": "Google authentication successful",
  "data": {
    "access_token": "<APP_JWT>",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

✅ **Flow summary**:
1. Mobile fetches auth URL with PKCE (`code_challenge`) and allowed `redirect_uri`.
2. Google redirects to the app with `code` + original `state`.
3. Mobile POSTs to `/auth/oauth/google/exchange` with `code_verifier`.
4. Backend talks to Google, upserts user, stores refresh token, returns Keerja JWT.

**Optional — preferred mobile deep-link (one-time code)**

When the backend uses a safer one-time code (single-use) in the deep-link instead of embedding the JWT, mobile apps should POST that one-time code to the exchange-one-time endpoint to obtain the app JWT:

```http
POST /auth/oauth/google/exchange-one-time 🔓
Content-Type: application/json

{
  "code": "v1_onetime_AbCdEf..."
}

→ 200 OK
{
  "success": true,
  "message": "Token exchange successful",
  "data": {
    "access_token": "<APP_JWT>",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

Notes:
- One-time codes are single-use and short-lived (default TTL 2 minutes).
- Backend generates the code and redirects mobile apps to `myapp://oauth-callback?code=<ONE_TIME_CODE>`.

---

#### 14. List Connected Providers

```http
GET /auth/oauth/connected 🔒
Authorization: Bearer <JWT_TOKEN>

→ 200 OK
{
  "success": true,
  "message": "Connected providers retrieved successfully",
  "data": [
    {
      "id": 1,
      "provider": "google",
      "provider_id": "105123456789012345678",
      "email": "john@gmail.com",
      "display_name": "John Doe",
      "connected_at": "2025-10-24T08:00:00Z",
      "last_login_at": "2025-10-24T10:30:00Z"
    }
  ]
}
```

---

#### 15. Disconnect OAuth Provider

```http
DELETE /auth/oauth/google 🔒
Authorization: Bearer <JWT_TOKEN>

→ 200 OK
{
  "success": true,
  "message": "OAuth provider disconnected successfully",
  "data": null
}
```

⚠️ **Note**: Replace `google` with provider name (google, facebook, github)

---

## LEGACY ENDPOINTS (8)

### Traditional Registration

```http
POST /auth/register 🔓 🚦
{
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "user_type": "jobseeker"
}
```

**Deprecated**: Use `/auth/register-otp` instead

---

### Traditional Login

```http
POST /auth/login 🔓 🚦
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

⚠️ **Deprecated**: Use `/auth/login-remember` instead

---

### Email Verification (Token-based)

```http
POST /auth/verify-email 🔓
{
  "token": "uuid-token-from-email"
}
```

---

### Forgot Password

```http
POST /auth/forgot-password 🔓 🚦
{
  "email": "john@example.com"
}
```

---

### Reset Password

```http
POST /auth/reset-password 🔓
{
  "token": "reset-token-from-email",
  "new_password": "NewSecurePass123!"
}
```

---

### Resend Verification Email

```http
POST /auth/resend-verification 🔓 🚦
{
  "email": "john@example.com"
}
```

---

### Legacy Refresh Token

```http
POST /auth/refresh-token 🔒
Authorization: Bearer <JWT_TOKEN>
```

**Deprecated**: Use `/auth/refresh` instead

---

### Logout (JWT invalidation)

```http
POST /auth/logout 🔒
Authorization: Bearer <JWT_TOKEN>
```

---

## Common Workflows

### 1. New User Registration Flow (OTP)

```
1. POST /auth/register-otp
   → User receives OTP email (123456)

2. POST /auth/verify-email-otp
   → Returns access_token + user data
   → User now verified and logged in

3. (Optional) POST /auth/login-remember
   → Get refresh_token for persistent session
```

---

### 2. Forgot Password Flow (OTP)

```
1. POST /auth/forgot-password-otp
   → email: "john@example.com"
   → User receives OTP email (123456, valid 5 min)

2. POST /auth/reset-password-otp
   → otp_code: "123456"
   → new_password: "NewSecurePass123!"
   → Returns success message
   → User receives confirmation email

3. Login with new password
   → POST /auth/login-remember
   → Use new password
```

**Security**: Max 3 OTP requests/hour, max 5 attempts per OTP

---

### 3. Login with Persistent Session

```
1. POST /auth/login-remember
   → remember_me: true, device_id: "chrome-001"
   → Returns access_token + refresh_token

2. Save both tokens securely (localStorage/secure cookie)

3. After 1 hour (JWT expires):
   → POST /auth/refresh
   → Get NEW access_token + NEW refresh_token
   → Update stored tokens

4. Repeat step 3 for 90 days (token lifespan)
```

---

### 4. Multi-Device Management

```
1. POST /auth/login-remember (Device A: laptop)
   → device_id: "chrome-laptop-001"

2. POST /auth/login-remember (Device B: mobile)
   → device_id: "firefox-mobile-002"

3. GET /auth/devices
   → See both devices active

4. POST /auth/devices/revoke
   → device_id: "chrome-laptop-001"
   → Logout from laptop only

5. POST /auth/logout-all
   → Logout from ALL devices
```

---

### 5. Social Login Flow (Google)

```
1. GET /auth/oauth/google
   → Returns auth_url

2. Redirect user to auth_url
   → User logs in with Google
   → Google redirects to callback

3. GET /auth/oauth/google/callback (automatic)
   → Returns access_token
   → User logged in

4. GET /auth/oauth/connected
   → See Google is connected

5. DELETE /auth/oauth/google
   → Disconnect Google
```

---

## Error Codes

| Status | Error                 | Description                 |
| ------ | --------------------- | --------------------------- |
| 400    | Bad Request           | Invalid request body/params |
| 401    | Unauthorized          | Invalid credentials/token   |
| 404    | Not Found             | Resource not found          |
| 409    | Conflict              | Email already exists        |
| 429    | Too Many Requests     | Rate limit exceeded         |
| 500    | Internal Server Error | Server error                |

### Common Error Responses

#### Invalid OTP

```json
{
  "success": false,
  "message": "Invalid OTP code",
  "error": "invalid OTP code"
}
```

#### Expired OTP

```json
{
  "success": false,
  "message": "OTP code has expired",
  "error": "OTP code has expired"
}
```

#### Too Many OTP Attempts

```json
{
  "success": false,
  "message": "Too many invalid attempts. Please request a new OTP.",
  "error": "too many invalid attempts"
}
```

#### Password Reset Rate Limit

```json
{
  "success": false,
  "message": "Too many password reset requests. Please try again later.",
  "error": "rate limit exceeded"
}
```

#### Expired Refresh Token

```json
{
  "success": false,
  "message": "Refresh token expired",
  "error": "refresh token expired"
}
```

#### Rate Limit Exceeded

```json
{
  "success": false,
  "message": "Too many OTP requests. Please try again later.",
  "error": "rate limit exceeded"
}
```

---

## Security Best Practices

### Client-Side

1. **Store Tokens Securely**

   - JWT: localStorage/memory (short-lived)
   - Refresh Token: httpOnly cookie or secure storage

2. **Implement Token Rotation**

   - Always use new refresh token after `/auth/refresh`
   - Discard old refresh token immediately

3. **Device ID Generation**

   - Use persistent unique ID per device
   - Example: `UUID.v4()` or `fingerprintjs`

4. **Error Handling**
   - Catch 401 errors → redirect to login
   - Catch 429 errors → show rate limit message

### Server-Side (Already Implemented)

- ✅ SHA256 hashing for OTP + refresh tokens
- ✅ Token rotation on refresh
- ✅ Rate limiting per endpoint
- ✅ Device tracking (IP, user agent)
- ✅ CSRF protection via OAuth state

---

## Full Documentation

- **Testing Guide**: `docs/OTP_REFRESH_OAUTH_TESTING_GUIDE.md`
- **Implementation Summary**: `docs/OTP_REFRESH_OAUTH_IMPLEMENTATION_SUMMARY.md`
- **Refresh Token Details**: `docs/REFRESH_TOKEN_TESTING.md`

---

**Total Endpoints**: 20 (12 new + 8 legacy)  
**Last Updated**: 2025-10-24
