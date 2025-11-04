-- Create company_sizes table
-- This table stores master data for company size categories (e.g., 1-10 employees, 11-50 employees, etc.)

CREATE TABLE IF NOT EXISTS company_sizes (
    id BIGSERIAL PRIMARY KEY,
    label VARCHAR(100) NOT NULL,
    min_employees INT NOT NULL,
    max_employees INT,
    display_order INT DEFAULT 0 NOT NULL,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_company_sizes_active ON company_sizes(is_active, deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_company_sizes_display_order ON company_sizes(display_order);

-- Add comment
COMMENT ON TABLE company_sizes IS 'Master data for company size categories';
COMMENT ON COLUMN company_sizes.label IS 'Display label (e.g., "1 - 10 karyawan")';
COMMENT ON COLUMN company_sizes.min_employees IS 'Minimum number of employees in this range';
COMMENT ON COLUMN company_sizes.max_employees IS 'Maximum number of employees (NULL for unlimited)';
