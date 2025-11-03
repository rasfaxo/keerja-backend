-- Alter user_profiles table to add master data relations
-- This migration adds foreign keys to master data tables (provinces, cities, districts)

-- Add new columns for location master data
ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS province_id BIGINT;
ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS city_id BIGINT;
ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS district_id BIGINT;

-- Create foreign key constraints
ALTER TABLE user_profiles 
    ADD CONSTRAINT fk_user_profiles_province 
    FOREIGN KEY (province_id) 
    REFERENCES provinces(id) 
    ON DELETE SET NULL;

ALTER TABLE user_profiles 
    ADD CONSTRAINT fk_user_profiles_city 
    FOREIGN KEY (city_id) 
    REFERENCES cities(id) 
    ON DELETE SET NULL;

ALTER TABLE user_profiles 
    ADD CONSTRAINT fk_user_profiles_district 
    FOREIGN KEY (district_id) 
    REFERENCES districts(id) 
    ON DELETE SET NULL;

-- Create indexes for foreign keys
CREATE INDEX IF NOT EXISTS idx_user_profiles_province ON user_profiles(province_id);
CREATE INDEX IF NOT EXISTS idx_user_profiles_city ON user_profiles(city_id);
CREATE INDEX IF NOT EXISTS idx_user_profiles_district ON user_profiles(district_id);

-- Add comments
COMMENT ON COLUMN user_profiles.province_id IS 'Foreign key to provinces master table (replaces location_state)';
COMMENT ON COLUMN user_profiles.city_id IS 'Foreign key to cities master table (replaces location_city)';
COMMENT ON COLUMN user_profiles.district_id IS 'Foreign key to districts master table (more granular location)';

-- Note: Old fields (location_city, location_state, location_country) are kept for backward compatibility
-- but should be deprecated in favor of the new structure
