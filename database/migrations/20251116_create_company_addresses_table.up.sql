-- Create company_addresses table
CREATE TABLE IF NOT EXISTS company_addresses (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    full_address TEXT NOT NULL,
    latitude NUMERIC(10,6),
    longitude NUMERIC(10,6),
    province_id BIGINT,
    city_id BIGINT,
    district_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_company_addresses_company_id ON company_addresses(company_id);
CREATE INDEX IF NOT EXISTS idx_company_addresses_deleted_at ON company_addresses(deleted_at);
