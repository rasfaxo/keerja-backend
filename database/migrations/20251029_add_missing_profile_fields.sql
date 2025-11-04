-- Migration: Add missing fields to user_profiles table
-- Date: 2025-10-29
-- Description: Add nationality, address, location_state, postal_code, linkedin_url, portfolio_url, github_url

-- Add missing columns to user_profiles
ALTER TABLE public.user_profiles 
ADD COLUMN IF NOT EXISTS nationality varchar(100),
ADD COLUMN IF NOT EXISTS address text,
ADD COLUMN IF NOT EXISTS location_state varchar(100),
ADD COLUMN IF NOT EXISTS postal_code varchar(10),
ADD COLUMN IF NOT EXISTS linkedin_url varchar(255),
ADD COLUMN IF NOT EXISTS portfolio_url varchar(255),
ADD COLUMN IF NOT EXISTS github_url varchar(255);

-- Add indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_user_profiles_location_state ON public.user_profiles USING btree (location_state);
CREATE INDEX IF NOT EXISTS idx_user_profiles_nationality ON public.user_profiles USING btree (nationality);

-- Add comments
COMMENT ON COLUMN public.user_profiles.nationality IS 'User nationality';
COMMENT ON COLUMN public.user_profiles.address IS 'User full address';
COMMENT ON COLUMN public.user_profiles.location_state IS 'User location state/province';
COMMENT ON COLUMN public.user_profiles.postal_code IS 'User postal/zip code';
COMMENT ON COLUMN public.user_profiles.linkedin_url IS 'User LinkedIn profile URL';
COMMENT ON COLUMN public.user_profiles.portfolio_url IS 'User portfolio/website URL';
COMMENT ON COLUMN public.user_profiles.github_url IS 'User GitHub profile URL';
