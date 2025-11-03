-- Create cities table
-- This table stores master data for Indonesian cities/regencies (Kota/Kabupaten)

CREATE TABLE IF NOT EXISTS cities (
    id BIGSERIAL PRIMARY KEY,
    province_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'Kota' or 'Kabupaten'
    code VARCHAR(10),
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_cities_province FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE RESTRICT
);

-- Create indexes
CREATE INDEX idx_cities_province ON cities(province_id, is_active, deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_cities_name ON cities(name);
CREATE INDEX idx_cities_type ON cities(type);
CREATE INDEX idx_cities_code ON cities(code) WHERE code IS NOT NULL;

-- Add comment
COMMENT ON TABLE cities IS 'Master data for Indonesian cities and regencies';
COMMENT ON COLUMN cities.type IS 'City type: "Kota" (city) or "Kabupaten" (regency)';
COMMENT ON COLUMN cities.code IS 'City code from BPS (Badan Pusat Statistik)';
