-- Add NPWP and NIB fields to company_verifications table
-- Migration: 000029_add_npwp_nib_to_company_verifications

ALTER TABLE company_verifications
ADD COLUMN IF NOT EXISTS npwp_number VARCHAR(50) NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS nib_number VARCHAR(50);

-- Add comment for documentation
COMMENT ON COLUMN company_verifications.npwp_number IS 'Nomor NPWP Perusahaan (Required for verification)';
COMMENT ON COLUMN company_verifications.nib_number IS 'Nomor Induk Berusaha 13 digit (Optional)';
