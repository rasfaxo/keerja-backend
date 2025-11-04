-- Create districts table
-- This table stores master data for Indonesian districts (Kecamatan)

CREATE TABLE IF NOT EXISTS districts (
    id BIGSERIAL PRIMARY KEY,
    city_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10),
    postal_code VARCHAR(10),
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_districts_city FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE RESTRICT
);

-- Create indexes
CREATE INDEX idx_districts_city ON districts(city_id, is_active, deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_districts_name ON districts(name);
CREATE INDEX idx_districts_code ON districts(code) WHERE code IS NOT NULL;
CREATE INDEX idx_districts_postal_code ON districts(postal_code) WHERE postal_code IS NOT NULL;

-- Add comment
COMMENT ON TABLE districts IS 'Master data for Indonesian districts (Kecamatan)';
COMMENT ON COLUMN districts.code IS 'District code from BPS (Badan Pusat Statistik)';
COMMENT ON COLUMN districts.postal_code IS 'Postal code for this district';
