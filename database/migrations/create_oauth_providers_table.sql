-- Migration: Create oauth_providers table for social login
-- Created: 2025-10-24

CREATE TABLE IF NOT EXISTS oauth_providers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,                   -- 'google', 'facebook', 'github', etc.
    provider_user_id TEXT NOT NULL,           -- user ID from OAuth provider
    email TEXT,                               -- email from provider
    name TEXT,                                -- name from provider
    avatar_url TEXT,                          -- profile picture URL
    access_token TEXT,                        -- OAuth access token (encrypted)
    refresh_token TEXT,                       -- OAuth refresh token (encrypted)
    token_expiry TIMESTAMPTZ,                 -- when access token expires
    raw_data JSONB,                           -- full OAuth response for reference
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    
    -- Ensure one provider account per user
    UNIQUE(provider, provider_user_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_oauth_providers_user_id ON oauth_providers(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_providers_provider ON oauth_providers(provider);
CREATE INDEX IF NOT EXISTS idx_oauth_providers_email ON oauth_providers(email);

COMMENT ON TABLE oauth_providers IS 'Stores OAuth provider connections for social login';
COMMENT ON COLUMN oauth_providers.provider IS 'OAuth provider name (google, facebook, etc.)';
COMMENT ON COLUMN oauth_providers.provider_user_id IS 'User ID from the OAuth provider';
COMMENT ON COLUMN oauth_providers.raw_data IS 'Full OAuth user profile data';
