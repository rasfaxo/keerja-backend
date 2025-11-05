-- Rollback migration: Remove Job Master Data Foreign Keys
-- Purpose: Revert jobs table to previous state without master data FKs

-- Drop foreign key constraints
ALTER TABLE jobs
DROP CONSTRAINT IF EXISTS fk_jobs_job_title,
DROP CONSTRAINT IF EXISTS fk_jobs_job_type,
DROP CONSTRAINT IF EXISTS fk_jobs_work_policy,
DROP CONSTRAINT IF EXISTS fk_jobs_education_level,
DROP CONSTRAINT IF EXISTS fk_jobs_experience_level,
DROP CONSTRAINT IF EXISTS fk_jobs_gender_preference;

-- Drop indexes
DROP INDEX IF EXISTS idx_jobs_job_title_id;
DROP INDEX IF EXISTS idx_jobs_job_type_id;
DROP INDEX IF EXISTS idx_jobs_work_policy_id;
DROP INDEX IF EXISTS idx_jobs_education_level_id;
DROP INDEX IF EXISTS idx_jobs_experience_level_id;
DROP INDEX IF EXISTS idx_jobs_gender_preference_id;

-- Drop columns
ALTER TABLE jobs 
DROP COLUMN IF EXISTS job_title_id,
DROP COLUMN IF EXISTS job_type_id,
DROP COLUMN IF EXISTS work_policy_id,
DROP COLUMN IF EXISTS education_level_id,
DROP COLUMN IF EXISTS experience_level_id,
DROP COLUMN IF EXISTS gender_preference_id;
