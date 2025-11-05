-- Rollback: Restore employment_type and job_level check constraints
-- Purpose: Rollback migration 000028

-- Restore job_level check constraint
ALTER TABLE jobs 
ADD CONSTRAINT jobs_job_level_check 
CHECK ((job_level::text = ANY (ARRAY['Internship'::character varying, 'Entry Level'::character varying, 'Mid Level'::character varying, 'Senior Level'::character varying, 'Manager'::character varying, 'Director'::character varying]::text[])));

-- Restore employment_type check constraint
ALTER TABLE jobs 
ADD CONSTRAINT jobs_employment_type_check 
CHECK ((employment_type::text = ANY (ARRAY['Full-Time'::character varying, 'Part-Time'::character varying, 'Contract'::character varying, 'Internship'::character varying, 'Freelance'::character varying]::text[])));
