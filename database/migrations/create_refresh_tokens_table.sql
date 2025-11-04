-- Migration: Create refresh_tokens table for persistent device sessions
-- Created: 2025-10-24
-- Purpose: Store refresh tokens for "Remember Me" functionality with device tracking

-- Drop table if exists (for development only)
-- DROP TABLE IF EXISTS refresh_tokens CASCADE;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,           -- SHA256 hash of refresh token
    device_name VARCHAR(255),                  -- User-friendly device name
    device_type VARCHAR(50),                   -- 'mobile', 'desktop', 'tablet', 'unknown'
    device_id VARCHAR(255),                    -- Unique device identifier
    user_agent TEXT,                           -- Full user agent string
    ip_address VARCHAR(45),                    -- IPv4 or IPv6
    last_used_at TIMESTAMPTZ DEFAULT now(),    -- Track last usage
    expires_at TIMESTAMPTZ NOT NULL,           -- Token expiration (30 days default)
    revoked BOOLEAN DEFAULT FALSE NOT NULL,    -- Manual revocation flag
    revoked_at TIMESTAMPTZ,                    -- When token was revoked
    revoked_reason VARCHAR(255),               -- Reason for revocation
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    
    -- Foreign key constraint
    CONSTRAINT fk_refresh_tokens_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked ON refresh_tokens(revoked);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_device_id ON refresh_tokens(device_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_device ON refresh_tokens(user_id, device_id);

-- Comments for documentation
COMMENT ON TABLE refresh_tokens IS 'Stores refresh tokens for persistent device sessions (Remember Me)';
COMMENT ON COLUMN refresh_tokens.user_id IS 'Reference to users table';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA256 hash of refresh token (never store plaintext)';
COMMENT ON COLUMN refresh_tokens.device_name IS 'User-friendly device name (e.g., "iPhone 13", "Chrome on Windows")';
COMMENT ON COLUMN refresh_tokens.device_type IS 'Type of device: mobile, desktop, tablet, unknown';
COMMENT ON COLUMN refresh_tokens.device_id IS 'Unique identifier for device (fingerprint)';
COMMENT ON COLUMN refresh_tokens.user_agent IS 'Full user agent string for device identification';
COMMENT ON COLUMN refresh_tokens.ip_address IS 'IP address when token was created/last used';
COMMENT ON COLUMN refresh_tokens.last_used_at IS 'Last time this token was used to refresh access token';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'When refresh token expires (typically 30 days)';
COMMENT ON COLUMN refresh_tokens.revoked IS 'Whether token has been manually revoked';
COMMENT ON COLUMN refresh_tokens.revoked_at IS 'Timestamp when token was revoked';
COMMENT ON COLUMN refresh_tokens.revoked_reason IS 'Reason for revocation (logout, security, suspicious activity)';

-- Verification query
SELECT 
    'refresh_tokens' AS table_name,
    COUNT(*) AS column_count
FROM information_schema.columns
WHERE table_name = 'refresh_tokens'
AND table_schema = 'public';

-- Show indexes
SELECT 
    indexname,
    tablename,
    indexdef
FROM pg_indexes
WHERE tablename = 'refresh_tokens'
ORDER BY indexname;
