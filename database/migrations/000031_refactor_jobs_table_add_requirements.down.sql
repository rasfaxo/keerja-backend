-- Migration Rollback: Refactor jobs table - Revert requirement fields
-- Version: 000031 DOWN
-- Description:
--   1. Restore deprecated fields (job_level, employment_type, education_level)
--   2. Remove new requirement fields (min_age, max_age)
--   3. Remove salary display preference field
--   4. Remove company address FK
--   5. Remove indexes and constraints

-- ============================================================================
-- STEP 1: Restore deprecated columns
-- ============================================================================

-- Restore job_level column
ALTER TABLE jobs
ADD COLUMN IF NOT EXISTS job_level VARCHAR(50);

-- Restore employment_type column
ALTER TABLE jobs
ADD COLUMN IF NOT EXISTS employment_type VARCHAR(30);

-- Restore education_level column
ALTER TABLE jobs
ADD COLUMN IF NOT EXISTS education_level VARCHAR(50);

-- ============================================================================
-- STEP 2: Restore old constraints
-- ============================================================================

-- Restore old status constraint (remove new statuses)
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_status_check;

ALTER TABLE jobs
ADD CONSTRAINT jobs_status_check CHECK (
    status IN ('draft', 'published', 'closed', 'expired', 'suspended')
);

-- ============================================================================
-- STEP 3: Restore comments for deprecated fields
-- ============================================================================

COMMENT ON COLUMN jobs.job_level IS 'DEPRECATED: Use experience_level_id instead (foreign key to experience_levels master data)';
COMMENT ON COLUMN jobs.employment_type IS 'DEPRECATED: Use job_type_id instead (foreign key to job_types master data)';
COMMENT ON COLUMN jobs.education_level IS 'DEPRECATED: Use education_level_id instead (foreign key to education_levels master data)';

-- ============================================================================
-- STEP 4: Drop new indexes
-- ============================================================================

DROP INDEX IF EXISTS idx_jobs_age_range;
DROP INDEX IF EXISTS idx_jobs_company_address;

-- ============================================================================
-- STEP 5: Drop new constraints
-- ============================================================================

ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_min_age_check;
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_max_age_check;
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_age_range_check;
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_salary_display_check;

-- ============================================================================
-- STEP 6: Drop new columns
-- ============================================================================

ALTER TABLE jobs
DROP COLUMN IF EXISTS min_age,
DROP COLUMN IF EXISTS max_age,
DROP COLUMN IF EXISTS salary_display,
DROP COLUMN IF EXISTS company_address_id;

-- ============================================================================
-- COMPLETED
-- ============================================================================

DO $$
BEGIN
    RAISE NOTICE 'Migration 000031 rollback completed successfully';
    RAISE NOTICE 'Restored fields: job_level, employment_type, education_level';
    RAISE NOTICE 'Removed fields: min_age, max_age, salary_display, company_address_id';
    RAISE NOTICE 'Reverted: status constraint, removed indexes and constraints';
END $$;
