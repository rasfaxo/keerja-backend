-- Migration: Add Job Master Data Foreign Keys
-- Purpose: Link jobs table to master data tables (job_types, work_policies, etc.)
-- Created: 2025-11-04

-- Add foreign key columns to jobs table
ALTER TABLE jobs 
ADD COLUMN IF NOT EXISTS job_title_id BIGINT,
ADD COLUMN IF NOT EXISTS job_type_id BIGINT,
ADD COLUMN IF NOT EXISTS work_policy_id BIGINT,
ADD COLUMN IF NOT EXISTS education_level_id BIGINT,
ADD COLUMN IF NOT EXISTS experience_level_id BIGINT,
ADD COLUMN IF NOT EXISTS gender_preference_id BIGINT;

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_jobs_job_title_id ON jobs(job_title_id);
CREATE INDEX IF NOT EXISTS idx_jobs_job_type_id ON jobs(job_type_id);
CREATE INDEX IF NOT EXISTS idx_jobs_work_policy_id ON jobs(work_policy_id);
CREATE INDEX IF NOT EXISTS idx_jobs_education_level_id ON jobs(education_level_id);
CREATE INDEX IF NOT EXISTS idx_jobs_experience_level_id ON jobs(experience_level_id);
CREATE INDEX IF NOT EXISTS idx_jobs_gender_preference_id ON jobs(gender_preference_id);

-- Add foreign key constraints (nullable for flexibility and backward compatibility)
-- Using DO block because PostgreSQL doesn't support IF NOT EXISTS for ADD CONSTRAINT
DO $$
BEGIN
    -- FK to job_titles
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jobs_job_title') THEN
        ALTER TABLE jobs
        ADD CONSTRAINT fk_jobs_job_title 
            FOREIGN KEY (job_title_id) 
            REFERENCES job_titles(id) 
            ON DELETE SET NULL;
    END IF;

    -- FK to job_types
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jobs_job_type') THEN
        ALTER TABLE jobs
        ADD CONSTRAINT fk_jobs_job_type 
            FOREIGN KEY (job_type_id) 
            REFERENCES job_types(id) 
            ON DELETE SET NULL;
    END IF;

    -- FK to work_policies
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jobs_work_policy') THEN
        ALTER TABLE jobs
        ADD CONSTRAINT fk_jobs_work_policy 
            FOREIGN KEY (work_policy_id) 
            REFERENCES work_policies(id) 
            ON DELETE SET NULL;
    END IF;

    -- FK to education_levels
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jobs_education_level') THEN
        ALTER TABLE jobs
        ADD CONSTRAINT fk_jobs_education_level 
            FOREIGN KEY (education_level_id) 
            REFERENCES education_levels(id) 
            ON DELETE SET NULL;
    END IF;

    -- FK to experience_levels
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jobs_experience_level') THEN
        ALTER TABLE jobs
        ADD CONSTRAINT fk_jobs_experience_level 
            FOREIGN KEY (experience_level_id) 
            REFERENCES experience_levels(id) 
            ON DELETE SET NULL;
    END IF;

    -- FK to gender_preferences
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jobs_gender_preference') THEN
        ALTER TABLE jobs
        ADD CONSTRAINT fk_jobs_gender_preference 
            FOREIGN KEY (gender_preference_id) 
            REFERENCES gender_preferences(id) 
            ON DELETE SET NULL;
    END IF;
END $$;

-- Add column comments for documentation
COMMENT ON COLUMN jobs.job_title_id IS 'FK to job_titles master data - standardized job title';
COMMENT ON COLUMN jobs.job_type_id IS 'FK to job_types master data (full_time, part_time, contract, internship, freelance)';
COMMENT ON COLUMN jobs.work_policy_id IS 'FK to work_policies master data (onsite, remote, hybrid)';
COMMENT ON COLUMN jobs.education_level_id IS 'FK to education_levels master data (minimum education requirement)';
COMMENT ON COLUMN jobs.experience_level_id IS 'FK to experience_levels master data (entry, junior, mid, senior, expert, lead)';
COMMENT ON COLUMN jobs.gender_preference_id IS 'FK to gender_preferences master data (male, female, any)';

-- Migrate existing data: Map employment_type string to job_type_id
UPDATE jobs SET job_type_id = (
    SELECT id FROM job_types WHERE code = 
    CASE 
        WHEN LOWER(jobs.employment_type) LIKE '%full%time%' THEN 'full_time'
        WHEN LOWER(jobs.employment_type) LIKE '%part%time%' THEN 'part_time'
        WHEN LOWER(jobs.employment_type) LIKE '%contract%' THEN 'contract'
        WHEN LOWER(jobs.employment_type) LIKE '%internship%' OR LOWER(jobs.employment_type) LIKE '%intern%' OR LOWER(jobs.employment_type) LIKE '%magang%' THEN 'internship'
        WHEN LOWER(jobs.employment_type) LIKE '%freelance%' THEN 'freelance'
    END
) WHERE employment_type IS NOT NULL AND job_type_id IS NULL;

-- Migrate existing data: Map job_level/experience to experience_level_id
UPDATE jobs SET experience_level_id = (
    SELECT id FROM experience_levels WHERE code = 
    CASE 
        WHEN LOWER(jobs.job_level) LIKE '%intern%' OR jobs.experience_min = 0 THEN 'entry'
        WHEN LOWER(jobs.job_level) LIKE '%entry%' OR LOWER(jobs.job_level) LIKE '%junior%' OR jobs.experience_min BETWEEN 1 AND 2 THEN 'junior'
        WHEN LOWER(jobs.job_level) LIKE '%mid%' OR jobs.experience_min BETWEEN 3 AND 4 THEN 'mid'
        WHEN LOWER(jobs.job_level) LIKE '%senior%' OR jobs.experience_min BETWEEN 5 AND 7 THEN 'senior'
        WHEN LOWER(jobs.job_level) LIKE '%manager%' OR LOWER(jobs.job_level) LIKE '%lead%' OR jobs.experience_min >= 8 THEN 'expert'
        WHEN LOWER(jobs.job_level) LIKE '%director%' OR jobs.experience_min >= 10 THEN 'lead'
        ELSE 'mid' -- default fallback
    END
) WHERE job_level IS NOT NULL AND experience_level_id IS NULL;

-- Migrate existing data: Map education_level string to education_level_id
UPDATE jobs SET education_level_id = (
    SELECT id FROM education_levels WHERE code = 
    CASE 
        WHEN LOWER(jobs.education_level) LIKE '%sma%' OR LOWER(jobs.education_level) LIKE '%high school%' THEN 'sma'
        WHEN LOWER(jobs.education_level) LIKE '%d3%' OR LOWER(jobs.education_level) LIKE '%diploma%' THEN 'd3'
        WHEN LOWER(jobs.education_level) LIKE '%s1%' OR LOWER(jobs.education_level) LIKE '%bachelor%' OR LOWER(jobs.education_level) LIKE '%sarjana%' THEN 's1'
        WHEN LOWER(jobs.education_level) LIKE '%s2%' OR LOWER(jobs.education_level) LIKE '%master%' THEN 's2'
        WHEN LOWER(jobs.education_level) LIKE '%s3%' OR LOWER(jobs.education_level) LIKE '%doctor%' OR LOWER(jobs.education_level) LIKE '%phd%' THEN 's3'
        ELSE 'sma' -- default fallback
    END
) WHERE education_level IS NOT NULL AND education_level_id IS NULL;

-- Set default work_policy_id to 'onsite' for jobs without remote_option
UPDATE jobs SET work_policy_id = (
    SELECT id FROM work_policies WHERE code = 'onsite'
) WHERE remote_option = false AND work_policy_id IS NULL;

-- Set default work_policy_id to 'remote' for jobs with remote_option
UPDATE jobs SET work_policy_id = (
    SELECT id FROM work_policies WHERE code = 'remote'
) WHERE remote_option = true AND work_policy_id IS NULL;

-- Set default gender_preference_id to 'any' if not specified
UPDATE jobs SET gender_preference_id = (
    SELECT id FROM gender_preferences WHERE code = 'any'
) WHERE gender_preference_id IS NULL;
