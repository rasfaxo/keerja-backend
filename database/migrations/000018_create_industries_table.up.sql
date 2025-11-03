-- Create industries table
-- This table stores master data for company industries (e.g., Technology, Healthcare, etc.)

CREATE TABLE IF NOT EXISTS industries (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    icon_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true NOT NULL,
    display_order INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_industries_active ON industries(is_active, deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_industries_slug ON industries(slug);
CREATE INDEX idx_industries_display_order ON industries(display_order);

-- Add comment
COMMENT ON TABLE industries IS 'Master data for company industries';
COMMENT ON COLUMN industries.slug IS 'URL-friendly version of industry name';
COMMENT ON COLUMN industries.display_order IS 'Order for displaying in UI dropdowns';
