-- Migration: Add job_subcategory_id column to jobs table

BEGIN;

ALTER TABLE jobs
    ADD COLUMN IF NOT EXISTS job_subcategory_id bigint;

-- Add foreign key constraint to job_subcategories
ALTER TABLE jobs
    ADD CONSTRAINT fk_jobs_job_subcategory_id FOREIGN KEY (job_subcategory_id) REFERENCES job_subcategories(id) ON DELETE SET NULL;

COMMIT;
