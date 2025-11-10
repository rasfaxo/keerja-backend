-- Rollback: Drop job_subcategories and job_categories
-- Created: 2025-11-07

DROP TABLE IF EXISTS job_subcategories CASCADE;
DROP TABLE IF EXISTS job_categories CASCADE;
