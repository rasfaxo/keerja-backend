-- Migration Rollback: Drop Job Master Data Tables
-- Created: 2025-11-04

DROP TABLE IF EXISTS gender_preferences CASCADE;
DROP TABLE IF EXISTS experience_levels CASCADE;
DROP TABLE IF EXISTS education_levels CASCADE;
DROP TABLE IF EXISTS work_policies CASCADE;
DROP TABLE IF EXISTS job_types CASCADE;
DROP TABLE IF EXISTS job_titles CASCADE;
