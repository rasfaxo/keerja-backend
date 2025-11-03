-- Create provinces table
-- This table stores master data for Indonesian provinces (34 provinces)

CREATE TABLE IF NOT EXISTS provinces (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) UNIQUE,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_provinces_active ON provinces(is_active, deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_provinces_name ON provinces(name);
CREATE INDEX idx_provinces_code ON provinces(code) WHERE code IS NOT NULL;

-- Add comment
COMMENT ON TABLE provinces IS 'Master data for Indonesian provinces';
COMMENT ON COLUMN provinces.code IS 'Province code from BPS (Badan Pusat Statistik)';
