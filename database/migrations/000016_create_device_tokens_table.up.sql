-- ============================================================================
-- Device Tokens Table for FCM Push Notifications
-- Reference: https://firebase.google.com/docs/cloud-messaging/manage-tokens
-- ============================================================================

-- Create device_tokens table
CREATE TABLE IF NOT EXISTS device_tokens (
    id BIGSERIAL PRIMARY KEY,
    
    -- User relationship
    user_id BIGINT NOT NULL,
    
    -- FCM Token (max 4096 chars per FCM documentation)
    token VARCHAR(4096) NOT NULL,
    
    -- Platform: android, ios, web
    platform VARCHAR(20) NOT NULL CHECK (platform IN ('android', 'ios', 'web')),
    
    -- Device information (optional, for analytics)
    device_info JSONB DEFAULT '{}',
    
    -- Token status
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    
    -- Last used for sending notification
    last_used_at TIMESTAMP,
    
    -- Failure tracking
    failure_count INTEGER DEFAULT 0 NOT NULL,
    last_failure_at TIMESTAMP,
    failure_reason TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Foreign key constraint
    CONSTRAINT fk_device_tokens_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE,
    
    -- Unique constraint: one token per user
    CONSTRAINT unique_user_token UNIQUE (user_id, token)
);

-- ============================================================================
-- Indexes for Performance Optimization
-- ============================================================================

-- Index for fast lookup by user (most common query)
CREATE INDEX IF NOT EXISTS idx_device_tokens_user_id 
    ON device_tokens(user_id) WHERE is_active = TRUE;

-- Index for fast lookup by token (for validation)
CREATE INDEX IF NOT EXISTS idx_device_tokens_token 
    ON device_tokens(token) WHERE is_active = TRUE;

-- Index for cleanup job (inactive tokens)
CREATE INDEX IF NOT EXISTS idx_device_tokens_inactive 
    ON device_tokens(is_active, last_used_at);

-- Index for platform-specific queries
CREATE INDEX IF NOT EXISTS idx_device_tokens_platform 
    ON device_tokens(platform, is_active);

-- Index for failure tracking
CREATE INDEX IF NOT EXISTS idx_device_tokens_failure 
    ON device_tokens(failure_count, last_failure_at) 
    WHERE failure_count > 0;

-- Composite index for user + platform queries
CREATE INDEX IF NOT EXISTS idx_device_tokens_user_platform 
    ON device_tokens(user_id, platform) WHERE is_active = TRUE;

-- ============================================================================
-- Trigger for updated_at
-- ============================================================================

CREATE OR REPLACE FUNCTION update_device_tokens_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_device_tokens_updated_at
    BEFORE UPDATE ON device_tokens
    FOR EACH ROW
    EXECUTE FUNCTION update_device_tokens_updated_at();

-- ============================================================================
-- Comments for Documentation
-- ============================================================================

COMMENT ON TABLE device_tokens IS 'Stores FCM device registration tokens for push notifications';
COMMENT ON COLUMN device_tokens.token IS 'FCM registration token (max 4096 chars per Firebase docs)';
COMMENT ON COLUMN device_tokens.platform IS 'Device platform: android, ios, or web';
COMMENT ON COLUMN device_tokens.device_info IS 'Device metadata stored as JSON (model, OS version, app version)';
COMMENT ON COLUMN device_tokens.failure_count IS 'Number of consecutive failures (auto-deactivate after threshold)';
COMMENT ON COLUMN device_tokens.last_used_at IS 'Last time token was used to send a notification';
COMMENT ON COLUMN device_tokens.failure_reason IS 'Reason for last failure (from FCM error response)';
