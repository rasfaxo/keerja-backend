-- Migration: Refactor jobs table - Add requirement fields and remove deprecated columns
-- Version: 000031
-- Description:
--   1. Add age requirement fields (min_age, max_age)
--   2. Add salary display preference field
--   3. Add company address FK for proper location relation
--   4. Remove deprecated fields (job_level, employment_type, education_level)
--   5. Add proper indexes and constraints

-- ============================================================================
-- STEP 1: Add new columns
-- ============================================================================

-- Add age requirement fields
ALTER TABLE jobs
ADD COLUMN IF NOT EXISTS min_age INTEGER,
ADD COLUMN IF NOT EXISTS max_age INTEGER;

-- Add salary display preference
ALTER TABLE jobs
ADD COLUMN IF NOT EXISTS salary_display VARCHAR(20) DEFAULT 'range';

-- Add company address FK (for future company_addresses table)
ALTER TABLE jobs
ADD COLUMN IF NOT EXISTS company_address_id BIGINT;

-- ============================================================================
-- STEP 2: Add constraints
-- ============================================================================

-- Age constraints
ALTER TABLE jobs
ADD CONSTRAINT jobs_min_age_check CHECK (min_age IS NULL OR (min_age >= 17 AND min_age <= 100));

ALTER TABLE jobs
ADD CONSTRAINT jobs_max_age_check CHECK (max_age IS NULL OR (max_age >= 17 AND max_age <= 100));

ALTER TABLE jobs
ADD CONSTRAINT jobs_age_range_check CHECK (min_age IS NULL OR max_age IS NULL OR min_age <= max_age);

-- Salary display constraint
ALTER TABLE jobs
ADD CONSTRAINT jobs_salary_display_check CHECK (
    salary_display IN ('range', 'starting_from', 'up_to', 'hidden')
);

-- ============================================================================
-- STEP 3: Add indexes
-- ============================================================================

-- Index for age filtering
CREATE INDEX IF NOT EXISTS idx_jobs_age_range ON jobs(min_age, max_age)
WHERE min_age IS NOT NULL OR max_age IS NOT NULL;

-- Index for company address FK
CREATE INDEX IF NOT EXISTS idx_jobs_company_address ON jobs(company_address_id)
WHERE company_address_id IS NOT NULL;

-- ============================================================================
-- STEP 4: Add comments
-- ============================================================================

COMMENT ON COLUMN jobs.min_age IS 'Minimum age requirement for job applicants (17-100)';
COMMENT ON COLUMN jobs.max_age IS 'Maximum age requirement for job applicants (17-100)';
COMMENT ON COLUMN jobs.salary_display IS 'How salary should be displayed: range (show min-max), starting_from (show min only), up_to (show max only), hidden (hide salary)';
COMMENT ON COLUMN jobs.company_address_id IS 'FK to company_addresses table - specific work location for this job';

-- ============================================================================
-- STEP 5: Remove deprecated columns
-- ============================================================================

-- Check if deprecated columns have data
DO $$
BEGIN
    -- Log warning if deprecated columns contain data
    IF EXISTS (SELECT 1 FROM jobs WHERE job_level IS NOT NULL AND job_level != '') THEN
        RAISE NOTICE 'WARNING: job_level column contains data. It will be dropped. Data should be migrated to experience_level_id.';
    END IF;

    IF EXISTS (SELECT 1 FROM jobs WHERE employment_type IS NOT NULL AND employment_type != '') THEN
        RAISE NOTICE 'WARNING: employment_type column contains data. It will be dropped. Data should be migrated to job_type_id.';
    END IF;

    IF EXISTS (SELECT 1 FROM jobs WHERE education_level IS NOT NULL AND education_level != '') THEN
        RAISE NOTICE 'WARNING: education_level column contains data. It will be dropped. Data should be migrated to education_level_id.';
    END IF;
END $$;

-- Drop deprecated columns
ALTER TABLE jobs
DROP COLUMN IF EXISTS job_level,
DROP COLUMN IF EXISTS employment_type,
DROP COLUMN IF EXISTS education_level;

-- ============================================================================
-- STEP 6: Update status constraint to include new statuses
-- ============================================================================

-- Drop old constraint if exists
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_status_check;

-- Add updated status constraint
ALTER TABLE jobs
ADD CONSTRAINT jobs_status_check CHECK (
    status IN (
        'draft',
        'pending_review',
        'published',
        'closed',
        'expired',
        'suspended',
        'rejected'
    )
);

-- ============================================================================
-- STEP 7: Add foreign key for company_address_id (when table exists)
-- ============================================================================

-- Note: This FK will be added when company_addresses table is created
-- For now, just add a comment
COMMENT ON COLUMN jobs.company_address_id IS
'FK to company_addresses table (to be added when company_addresses table is created). For now, this can reference companies.id as a temporary measure.';

-- ============================================================================
-- COMPLETED
-- ============================================================================

-- Log completion
DO $$
BEGIN
    RAISE NOTICE 'Migration 000031 completed successfully';
    RAISE NOTICE 'Added fields: min_age, max_age, salary_display, company_address_id';
    RAISE NOTICE 'Removed fields: job_level, employment_type, education_level';
    RAISE NOTICE 'Updated: status constraint, added indexes and constraints';
END $$;
