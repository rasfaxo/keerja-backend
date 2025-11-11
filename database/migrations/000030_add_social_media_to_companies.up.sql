-- Add social media fields to companies table
ALTER TABLE companies
ADD COLUMN IF NOT EXISTS instagram_url TEXT,
ADD COLUMN IF NOT EXISTS facebook_url TEXT,
ADD COLUMN IF NOT EXISTS linkedin_url TEXT,
ADD COLUMN IF NOT EXISTS twitter_url TEXT,
ADD COLUMN IF NOT EXISTS short_description TEXT;

-- Add comment for documentation
COMMENT ON COLUMN companies.instagram_url IS 'Instagram profile URL';
COMMENT ON COLUMN companies.facebook_url IS 'Facebook page URL';
COMMENT ON COLUMN companies.linkedin_url IS 'LinkedIn company page URL';
COMMENT ON COLUMN companies.twitter_url IS 'Twitter profile URL';
COMMENT ON COLUMN companies.short_description IS 'Short description (Singkat) - brief company overview';
