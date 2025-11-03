-- Alter users table to add company relation
-- This migration adds tracking for users who have created companies

-- Add new columns
ALTER TABLE users ADD COLUMN IF NOT EXISTS has_company BOOLEAN DEFAULT false NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS company_id BIGINT;

-- Create foreign key constraint
ALTER TABLE users 
    ADD CONSTRAINT fk_users_company 
    FOREIGN KEY (company_id) 
    REFERENCES companies(id) 
    ON DELETE SET NULL;

-- Create index for foreign key
CREATE INDEX IF NOT EXISTS idx_users_company ON users(company_id);

-- Create index for has_company for quick filtering
CREATE INDEX IF NOT EXISTS idx_users_has_company ON users(has_company) WHERE has_company = true;

-- Add comments
COMMENT ON COLUMN users.has_company IS 'Flag indicating if user has created a company';
COMMENT ON COLUMN users.company_id IS 'Foreign key to the company created by this user';

-- Note: This helps quickly identify employer users and their associated companies
-- When a user creates a company:
-- 1. has_company is set to true
-- 2. company_id is set to the created company ID
-- 3. user_type should be updated to 'employer'
