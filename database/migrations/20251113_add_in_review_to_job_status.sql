-- Add 'in_review' to allowed job status values
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_status_check;
ALTER TABLE jobs ADD CONSTRAINT jobs_status_check CHECK (
  status IN (
    'in_review', 'draft', 'pending_review', 'published', 'closed', 'expired', 'suspended', 'rejected'
  )
);