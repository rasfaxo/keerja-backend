-- Migration: Remove employment_type and job_level check constraints (use master data only)
-- Purpose: Allow jobs table to work with master data IDs only, remove legacy string constraints
-- Created: 2025-11-04

-- Drop the employment_type check constraint
ALTER TABLE jobs 
DROP CONSTRAINT IF EXISTS jobs_employment_type_check;

-- Drop the job_level check constraint
ALTER TABLE jobs 
DROP CONSTRAINT IF EXISTS jobs_job_level_check;

-- Make employment_type and job_level nullable (they're not used anymore, master data IDs are used instead)
ALTER TABLE jobs 
ALTER COLUMN employment_type DROP NOT NULL,
ALTER COLUMN job_level DROP NOT NULL;

COMMENT ON COLUMN jobs.employment_type IS 'DEPRECATED: Use job_type_id instead (foreign key to job_types master data)';
COMMENT ON COLUMN jobs.job_level IS 'DEPRECATED: Use experience_level_id instead (foreign key to experience_levels master data)';
