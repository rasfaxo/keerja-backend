-- Migration: Initial Database Schema
-- Description: Rollback for baseline schema from Keerja production database
-- Direction: down

-- Drop all tables in correct order to respect foreign key constraints
DROP TABLE IF EXISTS application_documents CASCADE;
DROP TABLE IF EXISTS application_notes CASCADE;
DROP TABLE IF EXISTS admin_users CASCADE;
DROP TABLE IF EXISTS admin_roles CASCADE;
DROP TABLE IF EXISTS interviews CASCADE;
DROP TABLE IF EXISTS job_application_stages CASCADE;
DROP TABLE IF EXISTS job_applications CASCADE;
DROP TABLE IF EXISTS job_requirements CASCADE;
DROP TABLE IF EXISTS job_skills CASCADE;
DROP TABLE IF EXISTS job_benefits CASCADE;
DROP TABLE IF EXISTS job_locations CASCADE;
DROP TABLE IF EXISTS jobs CASCADE;
DROP TABLE IF EXISTS job_subcategories CASCADE;
DROP TABLE IF EXISTS job_titles CASCADE;
DROP TABLE IF EXISTS job_types CASCADE;
DROP TABLE IF EXISTS job_categories CASCADE;
DROP TABLE IF EXISTS company_followers CASCADE;
DROP TABLE IF EXISTS company_documents CASCADE;
DROP TABLE IF EXISTS company_employees CASCADE;
DROP TABLE IF EXISTS company_invitations CASCADE;
DROP TABLE IF EXISTS company_industries CASCADE;
DROP TABLE IF EXISTS company_verifications CASCADE;
DROP TABLE IF EXISTS company_reviews CASCADE;
DROP TABLE IF EXISTS company_profiles CASCADE;
DROP TABLE IF EXISTS company_addresses CASCADE;
DROP TABLE IF EXISTS companies CASCADE;
DROP TABLE IF EXISTS company_sizes CASCADE;
DROP TABLE IF EXISTS employer_users CASCADE;
DROP TABLE IF EXISTS email_logs CASCADE;
DROP TABLE IF EXISTS push_notification_logs CASCADE;
DROP TABLE IF EXISTS notification_preferences CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS device_tokens CASCADE;
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS otp_codes CASCADE;
DROP TABLE IF EXISTS oauth_providers CASCADE;
DROP TABLE IF EXISTS user_skills CASCADE;
DROP TABLE IF EXISTS user_projects CASCADE;
DROP TABLE IF EXISTS user_preferences CASCADE;
DROP TABLE IF EXISTS user_languages CASCADE;
DROP TABLE IF EXISTS user_experiences CASCADE;
DROP TABLE IF EXISTS user_educations CASCADE;
DROP TABLE IF EXISTS user_certifications CASCADE;
DROP TABLE IF EXISTS user_documents CASCADE;
DROP TABLE IF EXISTS user_profiles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS skills_master CASCADE;
DROP TABLE IF EXISTS benefits_master CASCADE;
DROP TABLE IF EXISTS gender_preferences CASCADE;
DROP TABLE IF EXISTS experience_levels CASCADE;
DROP TABLE IF EXISTS education_levels CASCADE;
DROP TABLE IF EXISTS industries CASCADE;
DROP TABLE IF EXISTS work_policies CASCADE;
DROP TABLE IF EXISTS districts CASCADE;
DROP TABLE IF EXISTS cities CASCADE;
DROP TABLE IF EXISTS provinces CASCADE;

-- Drop all functions and triggers
DROP TRIGGER IF EXISTS trigger_update_company_invitations_updated_at ON company_invitations CASCADE;
DROP TRIGGER IF EXISTS trigger_device_tokens_updated_at ON device_tokens CASCADE;
DROP FUNCTION IF EXISTS update_timestamp() CASCADE;
DROP FUNCTION IF EXISTS update_company_invitations_updated_at() CASCADE;
