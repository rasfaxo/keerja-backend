-- Alter companies table to add master data relations
-- This migration adds foreign keys to master data tables and new fields

-- Add new columns
ALTER TABLE companies ADD COLUMN IF NOT EXISTS industry_id BIGINT;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS company_size_id BIGINT;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS province_id BIGINT;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS city_id BIGINT;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS district_id BIGINT;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS full_address TEXT;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS description TEXT;

-- Create foreign key constraints
ALTER TABLE companies 
    ADD CONSTRAINT fk_companies_industry 
    FOREIGN KEY (industry_id) 
    REFERENCES industries(id) 
    ON DELETE RESTRICT;

ALTER TABLE companies 
    ADD CONSTRAINT fk_companies_company_size 
    FOREIGN KEY (company_size_id) 
    REFERENCES company_sizes(id) 
    ON DELETE RESTRICT;

ALTER TABLE companies 
    ADD CONSTRAINT fk_companies_province 
    FOREIGN KEY (province_id) 
    REFERENCES provinces(id) 
    ON DELETE RESTRICT;

ALTER TABLE companies 
    ADD CONSTRAINT fk_companies_city 
    FOREIGN KEY (city_id) 
    REFERENCES cities(id) 
    ON DELETE RESTRICT;

ALTER TABLE companies 
    ADD CONSTRAINT fk_companies_district 
    FOREIGN KEY (district_id) 
    REFERENCES districts(id) 
    ON DELETE RESTRICT;

-- Create indexes for foreign keys
CREATE INDEX IF NOT EXISTS idx_companies_industry ON companies(industry_id);
CREATE INDEX IF NOT EXISTS idx_companies_size ON companies(company_size_id);
CREATE INDEX IF NOT EXISTS idx_companies_province ON companies(province_id);
CREATE INDEX IF NOT EXISTS idx_companies_city ON companies(city_id);
CREATE INDEX IF NOT EXISTS idx_companies_district ON companies(district_id);

-- Add comments
COMMENT ON COLUMN companies.industry_id IS 'Foreign key to industries master table';
COMMENT ON COLUMN companies.company_size_id IS 'Foreign key to company_sizes master table';
COMMENT ON COLUMN companies.province_id IS 'Foreign key to provinces master table (derived from district)';
COMMENT ON COLUMN companies.city_id IS 'Foreign key to cities master table (derived from district)';
COMMENT ON COLUMN companies.district_id IS 'Foreign key to districts master table (replaces old city/province fields)';
COMMENT ON COLUMN companies.full_address IS 'Complete office address (replaces old address field for new entries)';
COMMENT ON COLUMN companies.description IS 'Company description (replaces old about field for new entries)';

-- Note: Old fields (industry, size_category, city, province, address, about) are kept for backward compatibility
-- but should be deprecated in favor of the new structure
