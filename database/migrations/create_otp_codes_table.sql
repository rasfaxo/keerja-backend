-- Migration: Create otp_codes table for email verification during registration
-- Created: 2025-10-24
-- Purpose: Store OTP codes for user email verification and other authentication purposes

-- Drop table if exists (for development only)
-- DROP TABLE IF EXISTS otp_codes CASCADE;

CREATE TABLE IF NOT EXISTS otp_codes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    otp_hash TEXT NOT NULL,                    -- SHA256 hash of OTP code
    type VARCHAR(50) NOT NULL,                 -- 'email_verification', 'password_reset', etc.
    expired_at TIMESTAMPTZ NOT NULL,           -- expiration timestamp
    is_used BOOLEAN DEFAULT FALSE NOT NULL,    -- whether OTP has been used
    used_at TIMESTAMPTZ,                       -- when OTP was used
    attempts INT DEFAULT 0 NOT NULL,           -- failed verification attempts
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    
    -- Foreign key constraint
    CONSTRAINT fk_otp_codes_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_otp_codes_user_id ON otp_codes(user_id);
CREATE INDEX IF NOT EXISTS idx_otp_codes_type ON otp_codes(type);
CREATE INDEX IF NOT EXISTS idx_otp_codes_expired_at ON otp_codes(expired_at);
CREATE INDEX IF NOT EXISTS idx_otp_codes_user_type ON otp_codes(user_id, type);
CREATE INDEX IF NOT EXISTS idx_otp_codes_is_used ON otp_codes(is_used);

-- Comments for documentation
COMMENT ON TABLE otp_codes IS 'Stores OTP verification codes for user authentication';
COMMENT ON COLUMN otp_codes.user_id IS 'Reference to users table';
COMMENT ON COLUMN otp_codes.otp_hash IS 'SHA256 hash of OTP code (never store plaintext)';
COMMENT ON COLUMN otp_codes.type IS 'Purpose of OTP: email_verification, password_reset';
COMMENT ON COLUMN otp_codes.expired_at IS 'When the OTP code expires (typically 5 minutes)';
COMMENT ON COLUMN otp_codes.is_used IS 'Whether OTP has been successfully used';
COMMENT ON COLUMN otp_codes.used_at IS 'Timestamp when OTP was used';
COMMENT ON COLUMN otp_codes.attempts IS 'Number of failed verification attempts';

-- Verification query
SELECT 
    'otp_codes' AS table_name,
    COUNT(*) AS column_count
FROM information_schema.columns
WHERE table_name = 'otp_codes'
AND table_schema = 'public';

-- Show indexes
SELECT 
    indexname,
    tablename,
    indexdef
FROM pg_indexes
WHERE tablename = 'otp_codes'
ORDER BY indexname;
