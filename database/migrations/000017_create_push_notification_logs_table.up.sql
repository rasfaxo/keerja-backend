-- ============================================================================
-- Push Notification Logs Table
-- Purpose: Track FCM push notification delivery for analytics and debugging
-- Dependencies: Requires notifications table and device_tokens table
-- ============================================================================

CREATE TABLE IF NOT EXISTS push_notification_logs (
    id BIGSERIAL PRIMARY KEY,
    
    -- Relationships
    notification_id BIGINT,
    device_token_id BIGINT,
    user_id BIGINT NOT NULL,
    
    -- FCM Message ID (returned by Firebase)
    fcm_message_id VARCHAR(255),
    
    -- Status: pending, sent, delivered, failed, clicked
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'delivered', 'failed', 'clicked')),
    
    -- Error information
    error_code VARCHAR(100),
    error_message TEXT,
    
    -- Response from FCM
    fcm_response JSONB,
    
    -- Delivery timing
    sent_at TIMESTAMP,
    delivered_at TIMESTAMP,
    clicked_at TIMESTAMP,
    
    -- Timestamp
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Foreign key constraints
    CONSTRAINT fk_push_logs_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_push_logs_device_token FOREIGN KEY (device_token_id) 
        REFERENCES device_tokens(id) ON DELETE SET NULL
);

-- ============================================================================
-- Indexes for Performance
-- ============================================================================

-- Index for user queries
CREATE INDEX IF NOT EXISTS idx_push_logs_user_id 
    ON push_notification_logs(user_id, created_at DESC);

-- Index for status filtering
CREATE INDEX IF NOT EXISTS idx_push_logs_status 
    ON push_notification_logs(status, created_at DESC);

-- Index for FCM message lookup
CREATE INDEX IF NOT EXISTS idx_push_logs_fcm_message_id 
    ON push_notification_logs(fcm_message_id) 
    WHERE fcm_message_id IS NOT NULL;

-- Index for device token analytics
CREATE INDEX IF NOT EXISTS idx_push_logs_device_token_id 
    ON push_notification_logs(device_token_id, created_at DESC)
    WHERE device_token_id IS NOT NULL;

-- Index for notification tracking
CREATE INDEX IF NOT EXISTS idx_push_logs_notification_id 
    ON push_notification_logs(notification_id)
    WHERE notification_id IS NOT NULL;

-- Index for failed notifications
CREATE INDEX IF NOT EXISTS idx_push_logs_failed 
    ON push_notification_logs(user_id, status, created_at DESC)
    WHERE status = 'failed';

-- ============================================================================
-- Comments for Documentation
-- ============================================================================

COMMENT ON TABLE push_notification_logs IS 'Tracks push notification delivery for analytics and debugging';
COMMENT ON COLUMN push_notification_logs.fcm_message_id IS 'Unique message ID returned by Firebase';
COMMENT ON COLUMN push_notification_logs.status IS 'Delivery status: pending, sent, delivered, failed, clicked';
COMMENT ON COLUMN push_notification_logs.fcm_response IS 'Full response from FCM API (for debugging)';
COMMENT ON COLUMN push_notification_logs.error_code IS 'FCM error code (e.g., InvalidRegistration, NotRegistered)';
COMMENT ON COLUMN push_notification_logs.notification_id IS 'Optional reference to notifications table if it exists';
