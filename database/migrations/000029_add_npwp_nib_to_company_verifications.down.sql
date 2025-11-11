-- Rollback: Remove NPWP and NIB fields from company_verifications table
-- Migration: 000029_add_npwp_nib_to_company_verifications

ALTER TABLE company_verifications
DROP COLUMN IF EXISTS npwp_number,
DROP COLUMN IF EXISTS nib_number;
