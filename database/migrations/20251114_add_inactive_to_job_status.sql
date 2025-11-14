-- Migration: Add 'inactive' to jobs.status check constraint
-- Date: 2025-11-14
-- Purpose: Allow application code to set job status to 'inactive' without DB constraint failure

ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_status_check;

ALTER TABLE jobs ADD CONSTRAINT jobs_status_check CHECK (
  status IN (
    'in_review',
    'pending_review',
    'draft',
    'published',
    'closed',
    'expired',
    'suspended',
    'rejected',
    'inactive'
  )
);

-- Note: Apply this migration to dev/test DB before enabling the 'inactive' status in production.
