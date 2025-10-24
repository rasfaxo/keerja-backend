# OTP Registration Flow - Testing Guide

## Overview

Implementasi OTP verification untuk proses registrasi user baru 

## Flow Diagram

```
User Registration → Send OTP → Verify Email → Auto Login
      ↓                ↓            ↓             ↓
   Create User    Email Sent    Set Verified  Return JWT
 (unverified)    (5 min TTL)   (is_verified)   Token
```

## Endpoints

### 1. Register User with OTP

**POST** `/api/v1/auth/register-otp`

**Request Body:**

```json
{
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "phone": "08123456789",
  "user_type": "jobseeker"
}
```

**Response (201):**

```json
{
  "success": true,
  "message": "Registration successful! Please check your email for OTP verification code.",
  "data": {
    "email": "john@example.com",
    "note": "OTP code has been sent to your email. Valid for 5 minutes."
  }
}
```

**What Happens:**

- Validates input (email format, password strength, etc.)
- Checks if email already exists
- Hashes password with bcrypt
- Creates user with `is_verified = false`
- Generates 6-digit OTP code
- Hashes OTP with SHA256 (email + code)
- Saves OTP to `otp_codes` table (expires in 5 minutes)
- Sends professional email with OTP code

### 2. Verify Email with OTP

**POST** `/api/v1/auth/verify-email-otp`

**Request Body:**

```json
{
  "email": "john@example.com",
  "otp_code": "123456"
}
```

**Response (200):**

```json
{
  "success": true,
  "message": "Email verified successfully! You are now logged in.",
  "data": {
    "user": {
      "id": 1,
      "uuid": "123e4567-e89b-12d3-a456-426614174000",
      "full_name": "John Doe",
      "email": "john@example.com",
      "user_type": "jobseeker",
      "is_verified": true,
      "status": "active"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

**What Happens:**

- Finds user by email
- Checks if already verified
- Finds latest OTP for user
- Validates OTP not expired (5 min window)
- Validates OTP not already used
- Checks max attempts not exceeded (5 attempts)
- Verifies OTP hash matches
- Marks OTP as used
- Updates user `is_verified = true`
- Generates JWT token for auto-login
- Returns token + user data

### 3. Resend OTP

**POST** `/api/v1/auth/resend-otp`

**Request Body:**

```json
{
  "email": "john@example.com"
}
```

**Response (200):**

```json
{
  "success": true,
  "message": "OTP code has been resent to your email.",
  "data": {
    "email": "john@example.com",
    "note": "OTP code is valid for 5 minutes."
  }
}
```

**What Happens:**

- Finds user by email
- Checks if already verified
- Rate limiting: max 3 OTP per hour
- Resend window: must wait 60 seconds
- Generates new OTP
- Sends new email

## Security Features

### 1. **OTP Hashing (SHA256)**

```go
// Never store plaintext OTP
hash := sha256.Sum256([]byte(email + "|" + otpCode))
```

### 2. **Rate Limiting**

- Max 3 OTP requests per hour per user
- 60-second resend window
- Prevents brute force attacks

### 3. **Attempt Limiting**

- Max 5 failed verification attempts
- Auto-lockout after max attempts
- Requires new OTP request

### 4. **Time-based Expiry**

- OTP expires in 5 minutes
- Automatic cleanup of expired OTPs
- Prevents replay attacks

### 5. **One-time Use Enforcement**

- OTP marked as `is_used = true` after verification
- Cannot reuse same OTP
- Prevents token reuse attacks

### 6. **Input Validation**

- Email format validation
- Password strength (min 8 chars)
- OTP format (6-digit numeric)
- User type validation (jobseeker/employer)

## Database Schema

### Table: `otp_codes`

```sql
CREATE TABLE otp_codes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    otp_hash TEXT NOT NULL,              -- SHA256 hash
    type VARCHAR(50) NOT NULL,           -- 'email_verification'
    expired_at TIMESTAMPTZ NOT NULL,     -- 5 minutes from creation
    is_used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMPTZ,
    attempts INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Indexes
CREATE INDEX idx_otp_codes_user_id ON otp_codes(user_id);
CREATE INDEX idx_otp_codes_type ON otp_codes(type);
CREATE INDEX idx_otp_codes_expired_at ON otp_codes(expired_at);
CREATE INDEX idx_otp_codes_user_type ON otp_codes(user_id, type);
```

## Testing with cURL

### 1. Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register-otp \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "phone": "08123456789",
    "user_type": "jobseeker"
  }'
```

### 2. Check Email

Open MailHog: http://localhost:8025

Look for email with 6-digit OTP code.

### 3. Verify OTP

```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-email-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "otp_code": "123456"
  }'
```

### 4. Resend OTP (if needed)

```bash
curl -X POST http://localhost:8080/api/v1/auth/resend-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

## Error Codes & Responses

### Invalid OTP (401)

```json
{
  "success": false,
  "message": "Invalid OTP code",
  "errors": "invalid OTP code"
}
```

### OTP Expired (401)

```json
{
  "success": false,
  "message": "OTP code has expired. Please request a new one.",
  "errors": "OTP code has expired"
}
```

### Too Many Attempts (429)

```json
{
  "success": false,
  "message": "Too many failed attempts. Please request a new OTP.",
  "errors": "too many failed OTP verification attempts"
}
```

### Rate Limit Exceeded (429)

```json
{
  "success": false,
  "message": "Too many OTP requests. Please try again later.",
  "errors": "too many OTP requests, please try again later"
}
```

### Resend Too Soon (429)

```json
{
  "success": false,
  "message": "Please wait 60 seconds before requesting a new OTP.",
  "errors": "please wait before requesting a new OTP"
}
```

### Email Already Exists (409)

```json
{
  "success": false,
  "message": "Email already exists",
  "errors": "email already exists"
}
```

### User Already Verified (400)

```json
{
  "success": false,
  "message": "Email already verified",
  "errors": "user email already verified"
}
```

## Email Template

**Subject:** Verifikasi Email Registrasi - Keerja

**Body Preview:**

```
Selamat Datang di Keerja!

Halo John Doe,

Terima kasih telah mendaftar di Keerja. Untuk menyelesaikan
registrasi, silakan verifikasi email Anda dengan memasukkan
kode OTP berikut:

┌─────────────────┐
│    123456       │
└─────────────────┘

Kode OTP ini akan kadaluarsa dalam 5 menit.

⚠️ Perhatian Keamanan:
• Jangan bagikan kode OTP ini kepada siapapun
• Tim Keerja tidak akan pernah meminta kode OTP Anda
• Kode ini hanya valid untuk satu kali penggunaan
```

## Configuration

Required environment variables:

```env
# Email/SMTP (already configured)
SMTP_HOST=mailhog
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=noreply@keerja.com

# JWT (already configured)
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION_HOURS=24
```

## Production Checklist

- [x] OTP hashing with SHA256
- [x] Rate limiting (3 requests/hour)
- [x] Attempt limiting (max 5 attempts)
- [x] Time-based expiry (5 minutes)
- [x] One-time use enforcement
- [x] Email validation
- [x] Password hashing (bcrypt)
- [x] Professional email template
- [x] Proper error handling
- [x] Input sanitization
- [x] Database indexes for performance
- [x] Foreign key constraints
- [x] Cleanup expired OTPs (scheduled job)
- [x] Auto-login after verification (JWT)

## Architecture Overview

```
┌─────────────┐
│   Handler   │ → Validate, Sanitize, Map DTOs
└──────┬──────┘
       ↓
┌─────────────┐
│   Service   │ → Business Logic, OTP Generation/Verification
└──────┬──────┘
       ↓
┌─────────────┐
│ Repository  │ → Database Operations (GORM)
└──────┬──────┘
       ↓
┌─────────────┐
│  Database   │ → PostgreSQL (otp_codes, users)
└─────────────┘
```

### Layer Responsibilities

1. **Handler Layer** (`auth_handler.go`)

   - Parse & validate HTTP requests
   - Sanitize user input
   - Map DTOs to domain models
   - Handle HTTP status codes
   - Return formatted responses

2. **Service Layer** (`registration_service.go`)

   - Generate OTP codes
   - Hash OTP with SHA256
   - Validate OTP attempts & expiry
   - Rate limiting logic
   - Send emails via EmailService
   - Update user verification status
   - Generate JWT tokens

3. **Repository Layer** (`auth_repository.go`)

   - CRUD operations for otp_codes
   - Database queries with GORM
   - Transaction management
   - Index utilization

4. **Domain Layer** (`auth/entity.go`, `auth/repository.go`)
   - Define entities (OTPCode)
   - Define repository interfaces
   - Business rules & validation
   
---

**Next Steps:**

1. Test semua endpoint dengan Postman/cURL
2. Verifikasi email diterima di MailHog (http://localhost:8025)
3. Test error cases (invalid OTP, expired, too many attempts, etc.)
4. Monitor logs untuk debugging
5. Setup cron job untuk cleanup expired OTPs
