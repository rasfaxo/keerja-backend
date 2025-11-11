-- Remove social media fields from companies table
ALTER TABLE companies
DROP COLUMN IF EXISTS instagram_url,
DROP COLUMN IF EXISTS facebook_url,
DROP COLUMN IF EXISTS linkedin_url,
DROP COLUMN IF EXISTS twitter_url,
DROP COLUMN IF EXISTS short_description;
