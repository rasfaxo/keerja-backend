
-- Dumped from database version 17.6
-- Dumped by pg_dump version 17.6
-- Compatible with PostgreSQL 15+ (removed transaction_timeout)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

-- ================================================================
--  NUCLEAR OPTION: RESET TOTAL (Development Only)
-- ================================================================
DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;
SET search_path TO public;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- ================================================================

CREATE FUNCTION public.update_company_invitations_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$

BEGIN

    NEW.updated_at = NOW();

    RETURN NEW;

END;

$$;


CREATE FUNCTION public.update_device_tokens_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$

BEGIN

    NEW.updated_at = CURRENT_TIMESTAMP;

    RETURN NEW;

END;

$$;


SET default_tablespace = '';

SET default_table_access_method = heap;


CREATE TABLE public.admin_roles (
    id bigint NOT NULL,
    role_name character varying(100) NOT NULL,
    role_description text,
    access_level smallint DEFAULT 5,
    is_system_role boolean DEFAULT false,
    created_by bigint,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT admin_roles_access_level_check CHECK (((access_level >= 1) AND (access_level <= 10)))
);


CREATE SEQUENCE public.admin_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.admin_roles_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.admin_roles_id_seq1 OWNED BY public.admin_roles.id;


--

--

CREATE TABLE public.admin_users (
    id bigint NOT NULL,
    uuid uuid DEFAULT gen_random_uuid(),
    full_name character varying(100) NOT NULL,
    email character varying(150) NOT NULL,
    phone character varying(20),
    password_hash text NOT NULL,
    role_id bigint,
    status character varying(20) DEFAULT 'active'::character varying,
    last_login timestamp without time zone,
    two_factor_secret character varying(100),
    profile_image_url text,
    created_by bigint,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    deleted_at timestamp without time zone,
    CONSTRAINT admin_users_status_check CHECK (((status)::text = ANY (ARRAY[('active'::character varying)::text, ('inactive'::character varying)::text, ('suspended'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.admin_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.admin_users_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.admin_users_id_seq1 OWNED BY public.admin_users.id;


--

--

CREATE TABLE public.application_documents (
    id bigint NOT NULL,
    application_id bigint NOT NULL,
    user_id bigint NOT NULL,
    document_type character varying(50) DEFAULT 'cv'::character varying,
    file_name character varying(255),
    file_url text NOT NULL,
    file_type character varying(50),
    file_size bigint,
    uploaded_at timestamp without time zone DEFAULT now(),
    is_verified boolean DEFAULT false,
    verified_by bigint,
    verified_at timestamp without time zone,
    notes text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT application_documents_document_type_check CHECK (((document_type)::text = ANY (ARRAY[('cv'::character varying)::text, ('cover_letter'::character varying)::text, ('portfolio'::character varying)::text, ('certificate'::character varying)::text, ('transcript'::character varying)::text, ('other'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.application_documents_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.application_documents_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.application_documents_id_seq1 OWNED BY public.application_documents.id;


--

--

CREATE TABLE public.application_notes (
    id bigint NOT NULL,
    application_id bigint NOT NULL,
    stage_id bigint,
    author_id bigint NOT NULL,
    note_type character varying(30) DEFAULT 'internal'::character varying,
    note_text text NOT NULL,
    visibility character varying(20) DEFAULT 'internal'::character varying,
    sentiment character varying(20) DEFAULT 'neutral'::character varying,
    is_pinned boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT application_notes_note_type_check CHECK (((note_type)::text = ANY (ARRAY[('evaluation'::character varying)::text, ('feedback'::character varying)::text, ('reminder'::character varying)::text, ('internal'::character varying)::text]))),
    CONSTRAINT application_notes_sentiment_check CHECK (((sentiment)::text = ANY (ARRAY[('positive'::character varying)::text, ('neutral'::character varying)::text, ('negative'::character varying)::text]))),
    CONSTRAINT application_notes_visibility_check CHECK (((visibility)::text = ANY (ARRAY[('internal'::character varying)::text, ('public'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.application_notes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.application_notes_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.application_notes_id_seq1 OWNED BY public.application_notes.id;


--

--

CREATE TABLE public.benefits_master (
    id bigint NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(150) NOT NULL,
    category character varying(50) DEFAULT 'other'::character varying,
    description text,
    icon character varying(100),
    is_active boolean DEFAULT true,
    popularity_score numeric(5,2) DEFAULT 0.00,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT benefits_master_category_check CHECK (((category)::text = ANY (ARRAY[('financial'::character varying)::text, ('health'::character varying)::text, ('career'::character varying)::text, ('lifestyle'::character varying)::text, ('flexibility'::character varying)::text, ('other'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.benefits_master_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.benefits_master_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.benefits_master_id_seq1 OWNED BY public.benefits_master.id;


--

--

CREATE TABLE public.cities (
    id bigint NOT NULL,
    province_id bigint NOT NULL,
    name character varying(255) NOT NULL,
    type character varying(50) NOT NULL,
    code character varying(10),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);


--

--

COMMENT ON TABLE public.cities IS 'Master data for Indonesian cities and regencies';


--

--

COMMENT ON COLUMN public.cities.type IS 'City type: "Kota" (city) or "Kabupaten" (regency)';


--

--

COMMENT ON COLUMN public.cities.code IS 'City code from BPS (Badan Pusat Statistik)';


--

--

CREATE SEQUENCE public.cities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.cities_id_seq OWNED BY public.cities.id;


--

--

CREATE TABLE public.companies (
    id bigint NOT NULL,
    uuid uuid DEFAULT gen_random_uuid(),
    company_name character varying(200) NOT NULL,
    slug character varying(200) NOT NULL,
    legal_name character varying(200),
    registration_number character varying(100),
    industry character varying(100),
    company_type character varying(50),
    size_category character varying(50),
    website_url text,
    email_domain character varying(100),
    phone character varying(30),
    address text,
    city character varying(100),
    province character varying(100),
    country character varying(100) DEFAULT 'Indonesia'::character varying,
    postal_code character varying(10),
    latitude numeric(10,6),
    longitude numeric(10,6),
    logo_url text,
    banner_url text,
    about text,
    culture text,
    benefits text[],
    verified boolean DEFAULT false,
    verified_at timestamp without time zone,
    verified_by bigint,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    industry_id bigint,
    company_size_id bigint,
    district_id bigint,
    full_address text,
    description text,
    province_id bigint,
    city_id bigint,
    instagram_url text,
    facebook_url text,
    linkedin_url text,
    twitter_url text,
    short_description text,
    CONSTRAINT companies_company_type_check CHECK (((company_type)::text = ANY (ARRAY[('private'::character varying)::text, ('public'::character varying)::text, ('startup'::character varying)::text, ('ngo'::character varying)::text, ('government'::character varying)::text]))),
    CONSTRAINT companies_size_category_check CHECK (((size_category)::text = ANY (ARRAY[('1-10'::character varying)::text, ('11-50'::character varying)::text, ('51-200'::character varying)::text, ('201-1000'::character varying)::text, ('1000+'::character varying)::text])))
);


--

--

COMMENT ON COLUMN public.companies.industry_id IS 'Foreign key to industries master table';


--

--

COMMENT ON COLUMN public.companies.company_size_id IS 'Foreign key to company_sizes master table';


--

--

COMMENT ON COLUMN public.companies.district_id IS 'Foreign key to districts master table (replaces old city/province fields)';


--

--

COMMENT ON COLUMN public.companies.full_address IS 'Complete office address (replaces old address field for new entries)';


--

--

COMMENT ON COLUMN public.companies.description IS 'Company description (replaces old about field for new entries)';


--

--

COMMENT ON COLUMN public.companies.province_id IS 'Foreign key to provinces master table (derived from district)';


--

--

COMMENT ON COLUMN public.companies.city_id IS 'Foreign key to cities master table (derived from district)';


--

--

COMMENT ON COLUMN public.companies.instagram_url IS 'Instagram profile URL';


--

--

COMMENT ON COLUMN public.companies.facebook_url IS 'Facebook page URL';


--

--

COMMENT ON COLUMN public.companies.linkedin_url IS 'LinkedIn company page URL';


--

--

COMMENT ON COLUMN public.companies.twitter_url IS 'Twitter profile URL';


--

--

COMMENT ON COLUMN public.companies.short_description IS 'Short description (Singkat) - brief company overview';


--

--

CREATE SEQUENCE public.companies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.companies_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.companies_id_seq1 OWNED BY public.companies.id;


--

--

CREATE TABLE public.company_addresses (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    full_address text NOT NULL,
    latitude numeric(10,6),
    longitude numeric(10,6),
    province_id bigint,
    city_id bigint,
    district_id bigint,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone
);


--

--

CREATE SEQUENCE public.company_addresses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.company_addresses_id_seq OWNED BY public.company_addresses.id;


--

--

CREATE TABLE public.company_documents (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    uploaded_by bigint,
    document_type character varying(50) NOT NULL,
    document_number character varying(100),
    document_name character varying(150),
    file_path text NOT NULL,
    issue_date date,
    expiry_date date,
    status character varying(20) DEFAULT 'pending'::character varying,
    verified_by bigint,
    verified_at timestamp without time zone,
    rejection_reason text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT company_documents_document_type_check CHECK (((document_type)::text = ANY (ARRAY[('SIUP'::character varying)::text, ('NPWP'::character varying)::text, ('NIB'::character varying)::text, ('AKTA'::character varying)::text, ('TDP'::character varying)::text, ('ISO'::character varying)::text, ('SERTIFIKAT'::character varying)::text, ('LAINNYA'::character varying)::text]))),
    CONSTRAINT company_documents_status_check CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('approved'::character varying)::text, ('rejected'::character varying)::text, ('expired'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.company_documents_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.company_documents_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.company_documents_id_seq1 OWNED BY public.company_documents.id;


--

--

CREATE TABLE public.company_employees (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    user_id bigint,
    full_name character varying(150),
    job_title character varying(100),
    department character varying(100),
    employment_type character varying(30) DEFAULT 'permanent'::character varying,
    employment_status character varying(30) DEFAULT 'active'::character varying,
    join_date date,
    end_date date,
    salary_range_min numeric(12,2),
    salary_range_max numeric(12,2),
    added_by bigint,
    note text,
    is_visible_public boolean DEFAULT false,
    verified boolean DEFAULT false,
    verified_at timestamp without time zone,
    verified_by bigint,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT company_employees_employment_status_check CHECK (((employment_status)::text = ANY (ARRAY[('active'::character varying)::text, ('resigned'::character varying)::text, ('terminated'::character varying)::text, ('on_leave'::character varying)::text]))),
    CONSTRAINT company_employees_employment_type_check CHECK (((employment_type)::text = ANY (ARRAY[('permanent'::character varying)::text, ('contract'::character varying)::text, ('intern'::character varying)::text, ('freelance'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.company_employees_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.company_employees_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.company_employees_id_seq1 OWNED BY public.company_employees.id;


--

--

CREATE TABLE public.company_followers (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    user_id bigint NOT NULL,
    followed_at timestamp without time zone DEFAULT now(),
    unfollowed_at timestamp without time zone,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

CREATE SEQUENCE public.company_followers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.company_followers_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.company_followers_id_seq1 OWNED BY public.company_followers.id;


--

--

CREATE TABLE public.company_industries (
    id bigint NOT NULL,
    code character varying(20) NOT NULL,
    name character varying(150) NOT NULL,
    description text,
    parent_id bigint,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

CREATE SEQUENCE public.company_industries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.company_industries_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.company_industries_id_seq1 OWNED BY public.company_industries.id;


--

--

CREATE TABLE public.company_invitations (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    email character varying(150) NOT NULL,
    full_name character varying(150) NOT NULL,
    "position" character varying(100),
    role character varying(30) DEFAULT 'recruiter'::character varying NOT NULL,
    token character varying(64) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    invited_by bigint NOT NULL,
    accepted_by bigint,
    accepted_at timestamp without time zone,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT company_invitations_role_check CHECK (((role)::text = ANY ((ARRAY['admin'::character varying, 'recruiter'::character varying, 'viewer'::character varying])::text[]))),
    CONSTRAINT company_invitations_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'accepted'::character varying, 'rejected'::character varying, 'expired'::character varying])::text[])))
);


--

--

COMMENT ON TABLE public.company_invitations IS 'Stores company employee invitation records with token-based acceptance system';


--

--

COMMENT ON COLUMN public.company_invitations.token IS 'Unique invitation token valid for 7 days';


--

--

COMMENT ON COLUMN public.company_invitations.status IS 'Invitation status: pending, accepted, rejected, or expired';


--

--

COMMENT ON COLUMN public.company_invitations.expires_at IS 'Token expiration timestamp (7 days from creation)';


--

--

CREATE SEQUENCE public.company_invitations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.company_invitations_id_seq OWNED BY public.company_invitations.id;


--

--

CREATE TABLE public.company_profiles (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    tagline character varying(200),
    short_description text,
    long_description text,
    mission text,
    vision text,
    "values" text[],
    culture text,
    work_environment text,
    gallery_urls text[],
    video_url text,
    awards text[],
    social_links jsonb,
    hiring_tagline character varying(200),
    seo_title character varying(200),
    seo_keywords text[],
    seo_description text,
    status character varying(20) DEFAULT 'draft'::character varying,
    verified boolean DEFAULT false,
    verified_at timestamp without time zone,
    verified_by bigint,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT company_profiles_status_check CHECK (((status)::text = ANY (ARRAY[('draft'::character varying)::text, ('published'::character varying)::text, ('suspended'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.company_profiles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.company_profiles_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.company_profiles_id_seq1 OWNED BY public.company_profiles.id;


--

--

CREATE TABLE public.company_reviews (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    user_id bigint,
    reviewer_type character varying(30),
    position_title character varying(100),
    employment_period character varying(50),
    rating_overall numeric(2,1),
    rating_culture numeric(2,1),
    rating_worklife numeric(2,1),
    rating_salary numeric(2,1),
    rating_management numeric(2,1),
    pros text,
    cons text,
    advice_to_management text,
    is_anonymous boolean DEFAULT true,
    recommend_to_friend boolean DEFAULT true,
    status character varying(20) DEFAULT 'pending'::character varying,
    moderated_by bigint,
    moderated_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT company_reviews_rating_overall_check CHECK (((rating_overall >= (0)::numeric) AND (rating_overall <= (5)::numeric))),
    CONSTRAINT company_reviews_reviewer_type_check CHECK (((reviewer_type)::text = ANY (ARRAY[('employee'::character varying)::text, ('ex-employee'::character varying)::text, ('applicant'::character varying)::text]))),
    CONSTRAINT company_reviews_status_check CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('approved'::character varying)::text, ('rejected'::character varying)::text, ('hidden'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.company_reviews_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.company_reviews_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.company_reviews_id_seq1 OWNED BY public.company_reviews.id;


--

--

CREATE TABLE public.company_sizes (
    id bigint NOT NULL,
    label character varying(100) NOT NULL,
    min_employees integer NOT NULL,
    max_employees integer,
    display_order integer DEFAULT 0 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);


--

--

COMMENT ON TABLE public.company_sizes IS 'Master data for company size categories';


--

--

COMMENT ON COLUMN public.company_sizes.label IS 'Display label (e.g., "1 - 10 karyawan")';


--

--

COMMENT ON COLUMN public.company_sizes.min_employees IS 'Minimum number of employees in this range';


--

--

COMMENT ON COLUMN public.company_sizes.max_employees IS 'Maximum number of employees (NULL for unlimited)';


--

--

CREATE SEQUENCE public.company_sizes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.company_sizes_id_seq OWNED BY public.company_sizes.id;


--

--

CREATE TABLE public.company_verifications (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    requested_by bigint,
    reviewed_by bigint,
    reviewed_at timestamp without time zone,
    status character varying(20) DEFAULT 'pending'::character varying,
    verification_score numeric(5,2) DEFAULT 0.00,
    verification_notes text,
    rejection_reason text,
    verification_expiry date,
    badge_granted boolean DEFAULT false,
    auto_expired boolean DEFAULT false,
    last_checked timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    npwp_number character varying(50) DEFAULT ''::character varying NOT NULL,
    nib_number character varying(50),
    CONSTRAINT company_verifications_status_check CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('under_review'::character varying)::text, ('verified'::character varying)::text, ('rejected'::character varying)::text, ('blacklisted'::character varying)::text, ('expired'::character varying)::text])))
);


--

--

COMMENT ON COLUMN public.company_verifications.npwp_number IS 'Nomor NPWP Perusahaan (Required for verification)';


--

--

COMMENT ON COLUMN public.company_verifications.nib_number IS 'Nomor Induk Berusaha 13 digit (Optional)';


--

--

CREATE SEQUENCE public.company_verifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.company_verifications_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.company_verifications_id_seq1 OWNED BY public.company_verifications.id;


--

--

CREATE TABLE public.device_tokens (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    token character varying(4096) NOT NULL,
    platform character varying(20) NOT NULL,
    device_info jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true NOT NULL,
    last_used_at timestamp without time zone,
    failure_count integer DEFAULT 0 NOT NULL,
    last_failure_at timestamp without time zone,
    failure_reason text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT device_tokens_platform_check CHECK (((platform)::text = ANY ((ARRAY['android'::character varying, 'ios'::character varying, 'web'::character varying])::text[])))
);


--

--

COMMENT ON TABLE public.device_tokens IS 'Stores FCM device registration tokens for push notifications';


--

--

COMMENT ON COLUMN public.device_tokens.token IS 'FCM registration token (max 4096 chars per Firebase docs)';


--

--

COMMENT ON COLUMN public.device_tokens.platform IS 'Device platform: android, ios, or web';


--

--

COMMENT ON COLUMN public.device_tokens.device_info IS 'Device metadata stored as JSON (model, OS version, app version)';


--

--

COMMENT ON COLUMN public.device_tokens.last_used_at IS 'Last time token was used to send a notification';


--

--

COMMENT ON COLUMN public.device_tokens.failure_count IS 'Number of consecutive failures (auto-deactivate after threshold)';


--

--

COMMENT ON COLUMN public.device_tokens.failure_reason IS 'Reason for last failure (from FCM error response)';


--

--

CREATE SEQUENCE public.device_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.device_tokens_id_seq OWNED BY public.device_tokens.id;


--

--

CREATE TABLE public.districts (
    id bigint NOT NULL,
    city_id bigint NOT NULL,
    name character varying(255) NOT NULL,
    code character varying(10),
    postal_code character varying(10),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);


--

--

COMMENT ON TABLE public.districts IS 'Master data for Indonesian districts (Kecamatan)';


--

--

COMMENT ON COLUMN public.districts.code IS 'District code from BPS (Badan Pusat Statistik)';


--

--

COMMENT ON COLUMN public.districts.postal_code IS 'Postal code for this district';


--

--

CREATE SEQUENCE public.districts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.districts_id_seq OWNED BY public.districts.id;


--

--

CREATE TABLE public.education_levels (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50) NOT NULL,
    description text,
    "order" integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

COMMENT ON TABLE public.education_levels IS 'Master data for education levels (SMA, D3, S1, S2, S3)';


--

--

CREATE SEQUENCE public.education_levels_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.education_levels_id_seq OWNED BY public.education_levels.id;


--

--

CREATE TABLE public.email_logs (
    id bigint NOT NULL,
    recipient character varying(255) NOT NULL,
    subject character varying(500) NOT NULL,
    body text,
    template character varying(100),
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    provider character varying(50),
    sent_at timestamp without time zone,
    failure_reason text,
    metadata jsonb,
    retry_count integer DEFAULT 0,
    max_retries integer DEFAULT 3,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

COMMENT ON TABLE public.email_logs IS 'Logs of all emails sent by the system';


--

--

CREATE SEQUENCE public.email_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.email_logs_id_seq OWNED BY public.email_logs.id;


--

--

CREATE TABLE public.employer_users (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    company_id bigint NOT NULL,
    role character varying(30) DEFAULT 'recruiter'::character varying,
    position_title character varying(100),
    department character varying(100),
    email_company character varying(150),
    phone_company character varying(30),
    is_verified boolean DEFAULT false,
    verified_at timestamp without time zone,
    verified_by bigint,
    is_active boolean DEFAULT true,
    last_login timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT employer_users_role_check CHECK (((role)::text = ANY (ARRAY[('owner'::character varying)::text, ('admin'::character varying)::text, ('recruiter'::character varying)::text, ('viewer'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.employer_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.employer_users_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.employer_users_id_seq1 OWNED BY public.employer_users.id;


--

--

CREATE TABLE public.experience_levels (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50) NOT NULL,
    description text,
    min_years integer DEFAULT 0 NOT NULL,
    max_years integer,
    "order" integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

COMMENT ON TABLE public.experience_levels IS 'Master data for experience levels (Fresh Graduate, 1-3 years, 3-5 years, etc.)';


--

--

COMMENT ON COLUMN public.experience_levels.min_years IS 'Minimum years of experience';


--

--

COMMENT ON COLUMN public.experience_levels.max_years IS 'Maximum years of experience (NULL = unlimited)';


--

--

CREATE SEQUENCE public.experience_levels_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.experience_levels_id_seq OWNED BY public.experience_levels.id;


--

--

CREATE TABLE public.gender_preferences (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50) NOT NULL,
    description text,
    "order" integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

COMMENT ON TABLE public.gender_preferences IS 'Master data for gender preferences (Male, Female, Any)';


--

--

CREATE SEQUENCE public.gender_preferences_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.gender_preferences_id_seq OWNED BY public.gender_preferences.id;


--

--

CREATE TABLE public.industries (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    description text,
    icon_url character varying(500),
    is_active boolean DEFAULT true NOT NULL,
    display_order integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);


--

--

COMMENT ON TABLE public.industries IS 'Master data for company industries';


--

--

COMMENT ON COLUMN public.industries.slug IS 'URL-friendly version of industry name';


--

--

COMMENT ON COLUMN public.industries.display_order IS 'Order for displaying in UI dropdowns';


--

--

CREATE SEQUENCE public.industries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.industries_id_seq OWNED BY public.industries.id;


--

--

CREATE TABLE public.interviews (
    id bigint NOT NULL,
    application_id bigint NOT NULL,
    stage_id bigint,
    interviewer_id bigint,
    scheduled_at timestamp without time zone NOT NULL,
    ended_at timestamp without time zone,
    interview_type character varying(20) DEFAULT 'online'::character varying,
    meeting_link text,
    location text,
    status character varying(20) DEFAULT 'scheduled'::character varying,
    overall_score numeric(4,2),
    technical_score numeric(4,2),
    communication_score numeric(4,2),
    personality_score numeric(4,2),
    remarks text,
    feedback_summary text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT interviews_interview_type_check CHECK (((interview_type)::text = ANY (ARRAY[('online'::character varying)::text, ('onsite'::character varying)::text, ('hybrid'::character varying)::text]))),
    CONSTRAINT interviews_overall_score_check CHECK (((overall_score >= (0)::numeric) AND (overall_score <= (100)::numeric))),
    CONSTRAINT interviews_status_check CHECK (((status)::text = ANY (ARRAY[('scheduled'::character varying)::text, ('completed'::character varying)::text, ('rescheduled'::character varying)::text, ('cancelled'::character varying)::text, ('no_show'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.interviews_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.interviews_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.interviews_id_seq1 OWNED BY public.interviews.id;


--

--

CREATE TABLE public.job_application_stages (
    id bigint NOT NULL,
    application_id bigint NOT NULL,
    stage_name character varying(50) NOT NULL,
    description text,
    handled_by bigint,
    started_at timestamp without time zone DEFAULT now(),
    completed_at timestamp without time zone,
    duration interval GENERATED ALWAYS AS ((completed_at - started_at)) STORED,
    notes text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT job_application_stages_stage_name_check CHECK (((stage_name)::text = ANY (ARRAY[('applied'::character varying)::text, ('screening'::character varying)::text, ('shortlisted'::character varying)::text, ('interview'::character varying)::text, ('offered'::character varying)::text, ('hired'::character varying)::text, ('rejected'::character varying)::text, ('withdrawn'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.job_application_stages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.job_application_stages_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.job_application_stages_id_seq1 OWNED BY public.job_application_stages.id;


--

--

CREATE TABLE public.job_applications (
    id bigint NOT NULL,
    job_id bigint NOT NULL,
    user_id bigint NOT NULL,
    company_id bigint,
    applied_at timestamp without time zone DEFAULT now(),
    status character varying(30) DEFAULT 'applied'::character varying,
    source character varying(50) DEFAULT 'keerja_portal'::character varying,
    match_score numeric(5,2) DEFAULT 0.00,
    notes text,
    viewed_by_employer boolean DEFAULT false,
    is_bookmarked boolean DEFAULT false,
    resume_url text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT job_applications_status_check CHECK (((status)::text = ANY (ARRAY[('applied'::character varying)::text, ('screening'::character varying)::text, ('shortlisted'::character varying)::text, ('interview'::character varying)::text, ('offered'::character varying)::text, ('hired'::character varying)::text, ('rejected'::character varying)::text, ('withdrawn'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.job_applications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.job_applications_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.job_applications_id_seq1 OWNED BY public.job_applications.id;


--

--

CREATE TABLE public.job_benefits (
    id bigint NOT NULL,
    job_id bigint NOT NULL,
    benefit_id bigint,
    benefit_name character varying(150) NOT NULL,
    description text,
    is_highlight boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

CREATE SEQUENCE public.job_benefits_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.job_benefits_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.job_benefits_id_seq1 OWNED BY public.job_benefits.id;


--

--

CREATE TABLE public.job_categories (
    id bigint NOT NULL,
    parent_id bigint,
    code character varying(30) NOT NULL,
    name character varying(150) NOT NULL,
    description text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

CREATE SEQUENCE public.job_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.job_categories_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.job_categories_id_seq1 OWNED BY public.job_categories.id;


--

--

CREATE TABLE public.job_locations (
    id bigint NOT NULL,
    job_id bigint NOT NULL,
    company_id bigint,
    location_type character varying(20) DEFAULT 'onsite'::character varying,
    address text,
    city character varying(100),
    province character varying(100),
    postal_code character varying(20),
    country character varying(100) DEFAULT 'Indonesia'::character varying,
    latitude numeric(10,6),
    longitude numeric(10,6),
    google_place_id character varying(100),
    map_url text,
    is_primary boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT job_locations_location_type_check CHECK (((location_type)::text = ANY (ARRAY[('onsite'::character varying)::text, ('hybrid'::character varying)::text, ('remote'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.job_locations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.job_locations_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.job_locations_id_seq1 OWNED BY public.job_locations.id;


--

--

CREATE TABLE public.job_requirements (
    id bigint NOT NULL,
    job_id bigint NOT NULL,
    requirement_type character varying(50) DEFAULT 'other'::character varying,
    requirement_text text NOT NULL,
    skill_id bigint,
    min_experience smallint,
    max_experience smallint,
    education_level character varying(50),
    language character varying(50),
    is_mandatory boolean DEFAULT true,
    priority smallint DEFAULT 1,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT job_requirements_requirement_type_check CHECK (((requirement_type)::text = ANY (ARRAY[('education'::character varying)::text, ('experience'::character varying)::text, ('skill'::character varying)::text, ('language'::character varying)::text, ('certification'::character varying)::text, ('other'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.job_requirements_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.job_requirements_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.job_requirements_id_seq1 OWNED BY public.job_requirements.id;


--

--

CREATE TABLE public.job_skills (
    id bigint NOT NULL,
    job_id bigint NOT NULL,
    skill_id bigint NOT NULL,
    importance_level character varying(20) DEFAULT 'required'::character varying,
    weight numeric(3,2) DEFAULT 1.00,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT job_skills_importance_level_check CHECK (((importance_level)::text = ANY (ARRAY[('required'::character varying)::text, ('preferred'::character varying)::text, ('optional'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.job_skills_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.job_skills_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.job_skills_id_seq1 OWNED BY public.job_skills.id;


--

--

CREATE TABLE public.job_subcategories (
    id bigint NOT NULL,
    category_id bigint NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(150) NOT NULL,
    description text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

CREATE SEQUENCE public.job_subcategories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.job_subcategories_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.job_subcategories_id_seq1 OWNED BY public.job_subcategories.id;


--

--

CREATE TABLE public.job_titles (
    id bigint NOT NULL,
    name character varying(200) NOT NULL,
    normalized_name character varying(220) NOT NULL,
    description text,
    recommended_category_id bigint,
    popularity_score integer DEFAULT 0,
    search_count integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

COMMENT ON TABLE public.job_titles IS 'Master data for job titles (e.g., Software Engineer, Data Analyst)';


--

--

COMMENT ON COLUMN public.job_titles.recommended_category_id IS 'Recommended job category for this title';


--

--

COMMENT ON COLUMN public.job_titles.popularity_score IS 'Popularity score for ranking (higher = more popular)';


--

--

COMMENT ON COLUMN public.job_titles.search_count IS 'Number of times searched by users';


--

--

CREATE SEQUENCE public.job_titles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.job_titles_id_seq OWNED BY public.job_titles.id;


--

--

CREATE TABLE public.job_types (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50) NOT NULL,
    description text,
    "order" integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

COMMENT ON TABLE public.job_types IS 'Master data for job types (Full-Time, Part-Time, Internship, Freelance, Contract)';


--

--

COMMENT ON COLUMN public.job_types.code IS 'Unique code for programmatic reference';


--

--

CREATE SEQUENCE public.job_types_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.job_types_id_seq OWNED BY public.job_types.id;


--

--

CREATE TABLE public.jobs (
    id bigint NOT NULL,
    uuid uuid DEFAULT gen_random_uuid(),
    company_id bigint NOT NULL,
    employer_user_id bigint,
    category_id bigint,
    title character varying(200) NOT NULL,
    slug character varying(220),
    description text NOT NULL,
    requirements text,
    responsibilities text,
    location character varying(150),
    city character varying(100),
    province character varying(100),
    remote_option boolean DEFAULT false,
    salary_min numeric(12,2),
    salary_max numeric(12,2),
    currency character varying(10) DEFAULT 'IDR'::character varying,
    experience_min smallint,
    experience_max smallint,
    total_hires smallint DEFAULT 1,
    status character varying(20) DEFAULT 'draft'::character varying,
    views_count bigint DEFAULT 0,
    applications_count bigint DEFAULT 0,
    published_at timestamp without time zone,
    expired_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    job_title_id bigint,
    job_type_id bigint,
    work_policy_id bigint,
    education_level_id bigint,
    experience_level_id bigint,
    gender_preference_id bigint,
    min_age integer,
    max_age integer,
    salary_display character varying(20) DEFAULT 'range'::character varying,
    company_address_id bigint,
    job_subcategory_id bigint,
    CONSTRAINT jobs_age_range_check CHECK (((min_age IS NULL) OR (max_age IS NULL) OR (min_age <= max_age))),
    CONSTRAINT jobs_max_age_check CHECK (((max_age IS NULL) OR ((max_age >= 17) AND (max_age <= 100)))),
    CONSTRAINT jobs_min_age_check CHECK (((min_age IS NULL) OR ((min_age >= 17) AND (min_age <= 100)))),
    CONSTRAINT jobs_salary_display_check CHECK (((salary_display)::text = ANY ((ARRAY['range'::character varying, 'starting_from'::character varying, 'up_to'::character varying, 'hidden'::character varying])::text[]))),
    CONSTRAINT jobs_status_check CHECK (((status)::text = ANY ((ARRAY['in_review'::character varying, 'pending_review'::character varying, 'draft'::character varying, 'published'::character varying, 'closed'::character varying, 'expired'::character varying, 'suspended'::character varying, 'rejected'::character varying, 'inactive'::character varying])::text[])))
);


--

--

COMMENT ON COLUMN public.jobs.job_title_id IS 'FK to job_titles master data - standardized job title';


--

--

COMMENT ON COLUMN public.jobs.job_type_id IS 'FK to job_types master data (full_time, part_time, contract, internship, freelance)';


--

--

COMMENT ON COLUMN public.jobs.work_policy_id IS 'FK to work_policies master data (onsite, remote, hybrid)';


--

--

COMMENT ON COLUMN public.jobs.education_level_id IS 'FK to education_levels master data (minimum education requirement)';


--

--

COMMENT ON COLUMN public.jobs.experience_level_id IS 'FK to experience_levels master data (entry, junior, mid, senior, expert, lead)';


--

--

COMMENT ON COLUMN public.jobs.gender_preference_id IS 'FK to gender_preferences master data (male, female, any)';


--

--

COMMENT ON COLUMN public.jobs.min_age IS 'Minimum age requirement for job applicants (17-100)';


--

--

COMMENT ON COLUMN public.jobs.max_age IS 'Maximum age requirement for job applicants (17-100)';


--

--

COMMENT ON COLUMN public.jobs.salary_display IS 'How salary should be displayed: range (show min-max), starting_from (show min only), up_to (show max only), hidden (hide salary)';


--

--

COMMENT ON COLUMN public.jobs.company_address_id IS 'FK to company_addresses table (to be added when company_addresses table is created). For now, this can reference companies.id as a temporary measure.';


--

--

CREATE SEQUENCE public.jobs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.jobs_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.jobs_id_seq1 OWNED BY public.jobs.id;


--

--

CREATE TABLE public.notification_preferences (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    email_enabled boolean DEFAULT true,
    push_enabled boolean DEFAULT true,
    sms_enabled boolean DEFAULT false,
    job_applications_enabled boolean DEFAULT true,
    interview_enabled boolean DEFAULT true,
    status_updates_enabled boolean DEFAULT true,
    job_recommendations_enabled boolean DEFAULT true,
    company_updates_enabled boolean DEFAULT true,
    marketing_enabled boolean DEFAULT false,
    weekly_digest_enabled boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

COMMENT ON TABLE public.notification_preferences IS 'Stores user notification preferences and settings';


--

--

CREATE SEQUENCE public.notification_preferences_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.notification_preferences_id_seq OWNED BY public.notification_preferences.id;


--

--

CREATE TABLE public.notifications (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    type character varying(50) NOT NULL,
    title character varying(255) NOT NULL,
    message text NOT NULL,
    data json,
    is_read boolean DEFAULT false,
    read_at timestamp without time zone,
    priority character varying(20) DEFAULT 'normal'::character varying,
    category character varying(50) NOT NULL,
    action_url character varying(500),
    icon character varying(100),
    sender_id bigint,
    related_id bigint,
    related_type character varying(50),
    expires_at timestamp without time zone,
    is_sent boolean DEFAULT false,
    sent_at timestamp without time zone,
    channel character varying(50),
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

COMMENT ON TABLE public.notifications IS 'Stores all user notifications';


--

--

COMMENT ON COLUMN public.notifications.type IS 'Type of notification: job_application, interview, status_update, etc.';


--

--

COMMENT ON COLUMN public.notifications.data IS 'Additional metadata stored as JSON';


--

--

COMMENT ON COLUMN public.notifications.priority IS 'Priority level: low, normal, high, urgent';


--

--

COMMENT ON COLUMN public.notifications.category IS 'Category: application, job, account, system';


--

--

COMMENT ON COLUMN public.notifications.related_type IS 'Type of related entity: job, application, interview, etc.';


--

--

COMMENT ON COLUMN public.notifications.channel IS 'Delivery channel: in_app, email, push, sms';


--

--

CREATE SEQUENCE public.notifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.notifications_id_seq OWNED BY public.notifications.id;


--

--

CREATE TABLE public.oauth_providers (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    provider text NOT NULL,
    provider_user_id text NOT NULL,
    email text,
    name text,
    avatar_url text,
    access_token text,
    refresh_token text,
    token_expiry timestamp with time zone,
    raw_data jsonb,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--

--

COMMENT ON TABLE public.oauth_providers IS 'Stores OAuth provider connections for social login';


--

--

COMMENT ON COLUMN public.oauth_providers.provider IS 'OAuth provider name (google, facebook, etc.)';


--

--

COMMENT ON COLUMN public.oauth_providers.provider_user_id IS 'User ID from the OAuth provider';


--

--

COMMENT ON COLUMN public.oauth_providers.raw_data IS 'Full OAuth user profile data';


--

--

CREATE SEQUENCE public.oauth_providers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.oauth_providers_id_seq OWNED BY public.oauth_providers.id;


--

--

CREATE TABLE public.otp_codes (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    otp_hash text NOT NULL,
    type character varying(50) NOT NULL,
    expired_at timestamp with time zone NOT NULL,
    is_used boolean DEFAULT false NOT NULL,
    used_at timestamp with time zone,
    attempts integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--

--

COMMENT ON TABLE public.otp_codes IS 'Stores OTP verification codes for user authentication';


--

--

COMMENT ON COLUMN public.otp_codes.user_id IS 'Reference to users table';


--

--

COMMENT ON COLUMN public.otp_codes.otp_hash IS 'SHA256 hash of OTP code (never store plaintext)';


--

--

COMMENT ON COLUMN public.otp_codes.type IS 'Purpose of OTP: email_verification, password_reset';


--

--

COMMENT ON COLUMN public.otp_codes.expired_at IS 'When the OTP code expires (typically 5 minutes)';


--

--

COMMENT ON COLUMN public.otp_codes.is_used IS 'Whether OTP has been successfully used';


--

--

COMMENT ON COLUMN public.otp_codes.used_at IS 'Timestamp when OTP was used';


--

--

COMMENT ON COLUMN public.otp_codes.attempts IS 'Number of failed verification attempts';


--

--

CREATE SEQUENCE public.otp_codes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.otp_codes_id_seq OWNED BY public.otp_codes.id;


--

--

CREATE TABLE public.provinces (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    code character varying(10),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);


--

--

COMMENT ON TABLE public.provinces IS 'Master data for Indonesian provinces';


--

--

COMMENT ON COLUMN public.provinces.code IS 'Province code from BPS (Badan Pusat Statistik)';


--

--

CREATE SEQUENCE public.provinces_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.provinces_id_seq OWNED BY public.provinces.id;


--

--

CREATE TABLE public.push_notification_logs (
    id bigint NOT NULL,
    notification_id bigint,
    device_token_id bigint,
    user_id bigint NOT NULL,
    fcm_message_id character varying(255),
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    error_code character varying(100),
    error_message text,
    fcm_response jsonb,
    sent_at timestamp without time zone,
    delivered_at timestamp without time zone,
    clicked_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT push_notification_logs_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'sent'::character varying, 'delivered'::character varying, 'failed'::character varying, 'clicked'::character varying])::text[])))
);


--

--

COMMENT ON TABLE public.push_notification_logs IS 'Tracks push notification delivery for analytics and debugging';


--

--

COMMENT ON COLUMN public.push_notification_logs.notification_id IS 'Optional reference to notifications table if it exists';


--

--

COMMENT ON COLUMN public.push_notification_logs.fcm_message_id IS 'Unique message ID returned by Firebase';


--

--

COMMENT ON COLUMN public.push_notification_logs.status IS 'Delivery status: pending, sent, delivered, failed, clicked';


--

--

COMMENT ON COLUMN public.push_notification_logs.error_code IS 'FCM error code (e.g., InvalidRegistration, NotRegistered)';


--

--

COMMENT ON COLUMN public.push_notification_logs.fcm_response IS 'Full response from FCM API (for debugging)';


--

--

CREATE SEQUENCE public.push_notification_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.push_notification_logs_id_seq OWNED BY public.push_notification_logs.id;


--

--

CREATE TABLE public.refresh_tokens (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    token_hash text NOT NULL,
    device_name character varying(255),
    device_type character varying(50),
    device_id character varying(255),
    user_agent text,
    ip_address character varying(45),
    last_used_at timestamp with time zone DEFAULT now(),
    expires_at timestamp with time zone NOT NULL,
    revoked boolean DEFAULT false NOT NULL,
    revoked_at timestamp with time zone,
    revoked_reason character varying(255),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--

--

COMMENT ON TABLE public.refresh_tokens IS 'Stores refresh tokens for persistent device sessions (Remember Me)';


--

--

COMMENT ON COLUMN public.refresh_tokens.user_id IS 'Reference to users table';


--

--

COMMENT ON COLUMN public.refresh_tokens.token_hash IS 'SHA256 hash of refresh token (never store plaintext)';


--

--

COMMENT ON COLUMN public.refresh_tokens.device_name IS 'User-friendly device name (e.g., "iPhone 13", "Chrome on Windows")';


--

--

COMMENT ON COLUMN public.refresh_tokens.device_type IS 'Type of device: mobile, desktop, tablet, unknown';


--

--

COMMENT ON COLUMN public.refresh_tokens.device_id IS 'Unique identifier for device (fingerprint)';


--

--

COMMENT ON COLUMN public.refresh_tokens.user_agent IS 'Full user agent string for device identification';


--

--

COMMENT ON COLUMN public.refresh_tokens.ip_address IS 'IP address when token was created/last used';


--

--

COMMENT ON COLUMN public.refresh_tokens.last_used_at IS 'Last time this token was used to refresh access token';


--

--

COMMENT ON COLUMN public.refresh_tokens.expires_at IS 'When refresh token expires (typically 30 days)';


--

--

COMMENT ON COLUMN public.refresh_tokens.revoked IS 'Whether token has been manually revoked';


--

--

COMMENT ON COLUMN public.refresh_tokens.revoked_at IS 'Timestamp when token was revoked';


--

--

COMMENT ON COLUMN public.refresh_tokens.revoked_reason IS 'Reason for revocation (logout, security, suspicious activity)';


--

--

CREATE SEQUENCE public.refresh_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;


--

--

CREATE TABLE public.skills_master (
    id bigint NOT NULL,
    code character varying(50),
    name character varying(150) NOT NULL,
    normalized_name character varying(150),
    category_id bigint,
    description text,
    skill_type character varying(30) DEFAULT 'technical'::character varying,
    difficulty_level character varying(20) DEFAULT 'intermediate'::character varying,
    popularity_score numeric(5,2) DEFAULT 0.00,
    aliases text[],
    parent_id bigint,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT skills_master_difficulty_level_check CHECK (((difficulty_level)::text = ANY (ARRAY[('beginner'::character varying)::text, ('intermediate'::character varying)::text, ('advanced'::character varying)::text]))),
    CONSTRAINT skills_master_skill_type_check CHECK (((skill_type)::text = ANY (ARRAY[('technical'::character varying)::text, ('soft'::character varying)::text, ('language'::character varying)::text, ('tool'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.skills_master_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.skills_master_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.skills_master_id_seq1 OWNED BY public.skills_master.id;


--

--

CREATE TABLE public.user_certifications (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    certification_name character varying(150) NOT NULL,
    issuing_organization character varying(150) NOT NULL,
    issue_date date,
    expiration_date date,
    credential_id character varying(100),
    credential_url text,
    description text,
    verified boolean DEFAULT false,
    verification_date timestamp without time zone,
    file_url text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

CREATE SEQUENCE public.user_certifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.user_certifications_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.user_certifications_id_seq1 OWNED BY public.user_certifications.id;


--

--

CREATE TABLE public.user_documents (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    document_type character varying(50),
    document_name character varying(150) NOT NULL,
    file_url text NOT NULL,
    file_size bigint,
    mime_type character varying(100),
    description text,
    uploaded_at timestamp without time zone DEFAULT now(),
    verified boolean DEFAULT false,
    verified_at timestamp without time zone,
    verified_by bigint,
    is_active boolean DEFAULT true,
    checksum character varying(100),
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT user_documents_document_type_check CHECK (((document_type)::text = ANY (ARRAY[('resume'::character varying)::text, ('id_card'::character varying)::text, ('certificate'::character varying)::text, ('portfolio'::character varying)::text, ('transcript'::character varying)::text, ('others'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.user_documents_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.user_documents_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.user_documents_id_seq1 OWNED BY public.user_documents.id;


--

--

CREATE TABLE public.user_educations (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    institution_name character varying(150) NOT NULL,
    major character varying(100),
    degree_level character varying(50),
    start_year integer,
    end_year integer,
    gpa numeric(3,2),
    activities text,
    description text,
    is_current boolean DEFAULT false,
    verified boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT user_educations_degree_level_check CHECK (((degree_level)::text = ANY (ARRAY[('SMA'::character varying)::text, ('D1'::character varying)::text, ('D2'::character varying)::text, ('D3'::character varying)::text, ('S1'::character varying)::text, ('S2'::character varying)::text, ('S3'::character varying)::text, ('Other'::character varying)::text]))),
    CONSTRAINT user_educations_end_year_check CHECK ((end_year >= 1950)),
    CONSTRAINT user_educations_start_year_check CHECK ((start_year >= 1950))
);


--

--

CREATE SEQUENCE public.user_educations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.user_educations_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.user_educations_id_seq1 OWNED BY public.user_educations.id;


--

--

CREATE TABLE public.user_experiences (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    company_name character varying(150) NOT NULL,
    position_title character varying(150) NOT NULL,
    industry character varying(100),
    employment_type character varying(30),
    start_date date NOT NULL,
    end_date date,
    is_current boolean DEFAULT false,
    description text,
    achievements text,
    location_city character varying(100),
    location_country character varying(100) DEFAULT 'Indonesia'::character varying,
    verified boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT user_experiences_employment_type_check CHECK (((employment_type)::text = ANY (ARRAY[('full-time'::character varying)::text, ('part-time'::character varying)::text, ('internship'::character varying)::text, ('freelance'::character varying)::text, ('contract'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.user_experiences_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.user_experiences_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.user_experiences_id_seq1 OWNED BY public.user_experiences.id;


--

--

CREATE TABLE public.user_languages (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    language_name character varying(100) NOT NULL,
    proficiency_level character varying(50),
    certification_name character varying(100),
    certification_score character varying(50),
    certification_date date,
    verified boolean DEFAULT false,
    is_active boolean DEFAULT true,
    notes text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT user_languages_proficiency_level_check CHECK (((proficiency_level)::text = ANY (ARRAY[('basic'::character varying)::text, ('intermediate'::character varying)::text, ('advanced'::character varying)::text, ('fluent'::character varying)::text, ('native'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.user_languages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.user_languages_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.user_languages_id_seq1 OWNED BY public.user_languages.id;


--

--

CREATE TABLE public.user_preferences (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    language_preference character varying(10),
    theme_preference character varying(10) DEFAULT 'light'::character varying,
    preferred_job_type character varying(50) DEFAULT 'onsite'::character varying,
    preferred_industry character varying(100),
    preferred_location character varying(100),
    preferred_salary_min numeric(12,2),
    preferred_salary_max numeric(12,2),
    email_notifications boolean DEFAULT true,
    sms_notifications boolean DEFAULT false,
    push_notifications boolean DEFAULT true,
    email_marketing boolean DEFAULT false,
    profile_visibility character varying(20) DEFAULT 'public'::character varying,
    show_online_status boolean DEFAULT true,
    allow_direct_messages boolean DEFAULT true,
    data_sharing_consent boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT user_preferences_profile_visibility_check CHECK (((profile_visibility)::text = ANY (ARRAY[('public'::character varying)::text, ('private'::character varying)::text, ('recruiter-only'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.user_preferences_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.user_preferences_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.user_preferences_id_seq1 OWNED BY public.user_preferences.id;


--

--

CREATE TABLE public.user_profiles (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    headline character varying(150),
    bio text,
    gender character varying(10),
    birth_date date,
    location_city character varying(100),
    location_country character varying(100),
    desired_position character varying(150),
    desired_salary_min numeric(12,2),
    desired_salary_max numeric(12,2),
    experience_level character varying(50),
    industry_interest character varying(100),
    availability_status character varying(50) DEFAULT 'open'::character varying,
    profile_visibility boolean DEFAULT true,
    slug character varying(100),
    avatar_url text,
    cover_url text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    nationality character varying(100),
    address text,
    location_state character varying(100),
    postal_code character varying(10),
    linkedin_url character varying(255),
    portfolio_url character varying(255),
    github_url character varying(255),
    province_id bigint,
    city_id bigint,
    district_id bigint,
    CONSTRAINT user_profiles_experience_level_check CHECK (((experience_level)::text = ANY (ARRAY[('internship'::character varying)::text, ('junior'::character varying)::text, ('mid'::character varying)::text, ('senior'::character varying)::text, ('lead'::character varying)::text]))),
    CONSTRAINT user_profiles_gender_check CHECK (((gender)::text = ANY (ARRAY[('male'::character varying)::text, ('female'::character varying)::text, ('other'::character varying)::text])))
);


--

--

COMMENT ON COLUMN public.user_profiles.nationality IS 'User nationality';


--

--

COMMENT ON COLUMN public.user_profiles.address IS 'User full address';


--

--

COMMENT ON COLUMN public.user_profiles.location_state IS 'User location state/province';


--

--

COMMENT ON COLUMN public.user_profiles.postal_code IS 'User postal/zip code';


--

--

COMMENT ON COLUMN public.user_profiles.linkedin_url IS 'User LinkedIn profile URL';


--

--

COMMENT ON COLUMN public.user_profiles.portfolio_url IS 'User portfolio/website URL';


--

--

COMMENT ON COLUMN public.user_profiles.github_url IS 'User GitHub profile URL';


--

--

COMMENT ON COLUMN public.user_profiles.province_id IS 'Foreign key to provinces master table (replaces location_state)';


--

--

COMMENT ON COLUMN public.user_profiles.city_id IS 'Foreign key to cities master table (replaces location_city)';


--

--

COMMENT ON COLUMN public.user_profiles.district_id IS 'Foreign key to districts master table (more granular location)';


--

--

CREATE SEQUENCE public.user_profiles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.user_profiles_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.user_profiles_id_seq1 OWNED BY public.user_profiles.id;


--

--

CREATE TABLE public.user_projects (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    project_title character varying(150) NOT NULL,
    role_in_project character varying(100),
    project_type character varying(50),
    description text,
    industry character varying(100),
    start_date date,
    end_date date,
    is_current boolean DEFAULT false,
    project_url text,
    repo_url text,
    media_urls text[],
    skills_used text[],
    collaborators text[],
    verified boolean DEFAULT false,
    featured boolean DEFAULT false,
    visibility character varying(20) DEFAULT 'public'::character varying,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT user_projects_project_type_check CHECK (((project_type)::text = ANY (ARRAY[('personal'::character varying)::text, ('freelance'::character varying)::text, ('company'::character varying)::text, ('academic'::character varying)::text, ('community'::character varying)::text]))),
    CONSTRAINT user_projects_visibility_check CHECK (((visibility)::text = ANY (ARRAY[('public'::character varying)::text, ('private'::character varying)::text, ('limited'::character varying)::text])))
);


--

--

CREATE SEQUENCE public.user_projects_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.user_projects_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.user_projects_id_seq1 OWNED BY public.user_projects.id;


--

--

CREATE TABLE public.user_skills (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    skill_name character varying(100) NOT NULL,
    skill_level character varying(20),
    years_experience integer,
    last_used_at date,
    verified boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT user_skills_skill_level_check CHECK (((skill_level)::text = ANY (ARRAY[('beginner'::character varying)::text, ('intermediate'::character varying)::text, ('advanced'::character varying)::text, ('expert'::character varying)::text]))),
    CONSTRAINT user_skills_years_experience_check CHECK ((years_experience >= 0))
);


--

--

CREATE SEQUENCE public.user_skills_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.user_skills_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.user_skills_id_seq1 OWNED BY public.user_skills.id;


--

--

CREATE TABLE public.users (
    id bigint NOT NULL,
    uuid uuid DEFAULT gen_random_uuid(),
    full_name character varying(150) NOT NULL,
    email character varying(150) NOT NULL,
    phone character varying(20),
    password_hash text NOT NULL,
    user_type character varying(20),
    is_verified boolean DEFAULT false,
    status character varying(20) DEFAULT 'active'::character varying,
    last_login timestamp without time zone,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    has_company boolean DEFAULT false NOT NULL,
    company_id bigint,
    CONSTRAINT users_user_type_check CHECK (((user_type)::text = ANY (ARRAY[('jobseeker'::character varying)::text, ('employer'::character varying)::text, ('admin'::character varying)::text])))
);


--

--

COMMENT ON COLUMN public.users.has_company IS 'Flag indicating if user has created a company';


--

--

COMMENT ON COLUMN public.users.company_id IS 'Foreign key to the company created by this user';


--

--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

CREATE SEQUENCE public.users_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.users_id_seq1 OWNED BY public.users.id;


--

--

CREATE TABLE public.work_policies (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50) NOT NULL,
    description text,
    icon_url character varying(500),
    "order" integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--

--

COMMENT ON TABLE public.work_policies IS 'Master data for work policies (On-site, Remote, Hybrid)';


--

--

COMMENT ON COLUMN public.work_policies.icon_url IS 'URL to icon representing this work policy';


--

--

CREATE SEQUENCE public.work_policies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--

--

ALTER SEQUENCE public.work_policies_id_seq OWNED BY public.work_policies.id;


--

--

ALTER TABLE ONLY public.admin_roles ALTER COLUMN id SET DEFAULT nextval('public.admin_roles_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.admin_users ALTER COLUMN id SET DEFAULT nextval('public.admin_users_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.application_documents ALTER COLUMN id SET DEFAULT nextval('public.application_documents_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.application_notes ALTER COLUMN id SET DEFAULT nextval('public.application_notes_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.benefits_master ALTER COLUMN id SET DEFAULT nextval('public.benefits_master_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.cities ALTER COLUMN id SET DEFAULT nextval('public.cities_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.companies ALTER COLUMN id SET DEFAULT nextval('public.companies_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.company_addresses ALTER COLUMN id SET DEFAULT nextval('public.company_addresses_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.company_documents ALTER COLUMN id SET DEFAULT nextval('public.company_documents_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.company_employees ALTER COLUMN id SET DEFAULT nextval('public.company_employees_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.company_followers ALTER COLUMN id SET DEFAULT nextval('public.company_followers_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.company_industries ALTER COLUMN id SET DEFAULT nextval('public.company_industries_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.company_invitations ALTER COLUMN id SET DEFAULT nextval('public.company_invitations_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.company_profiles ALTER COLUMN id SET DEFAULT nextval('public.company_profiles_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.company_reviews ALTER COLUMN id SET DEFAULT nextval('public.company_reviews_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.company_sizes ALTER COLUMN id SET DEFAULT nextval('public.company_sizes_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.company_verifications ALTER COLUMN id SET DEFAULT nextval('public.company_verifications_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.device_tokens ALTER COLUMN id SET DEFAULT nextval('public.device_tokens_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.districts ALTER COLUMN id SET DEFAULT nextval('public.districts_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.education_levels ALTER COLUMN id SET DEFAULT nextval('public.education_levels_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.email_logs ALTER COLUMN id SET DEFAULT nextval('public.email_logs_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.employer_users ALTER COLUMN id SET DEFAULT nextval('public.employer_users_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.experience_levels ALTER COLUMN id SET DEFAULT nextval('public.experience_levels_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.gender_preferences ALTER COLUMN id SET DEFAULT nextval('public.gender_preferences_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.industries ALTER COLUMN id SET DEFAULT nextval('public.industries_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.interviews ALTER COLUMN id SET DEFAULT nextval('public.interviews_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.job_application_stages ALTER COLUMN id SET DEFAULT nextval('public.job_application_stages_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.job_applications ALTER COLUMN id SET DEFAULT nextval('public.job_applications_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.job_benefits ALTER COLUMN id SET DEFAULT nextval('public.job_benefits_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.job_categories ALTER COLUMN id SET DEFAULT nextval('public.job_categories_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.job_locations ALTER COLUMN id SET DEFAULT nextval('public.job_locations_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.job_requirements ALTER COLUMN id SET DEFAULT nextval('public.job_requirements_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.job_skills ALTER COLUMN id SET DEFAULT nextval('public.job_skills_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.job_subcategories ALTER COLUMN id SET DEFAULT nextval('public.job_subcategories_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.job_titles ALTER COLUMN id SET DEFAULT nextval('public.job_titles_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.job_types ALTER COLUMN id SET DEFAULT nextval('public.job_types_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.jobs ALTER COLUMN id SET DEFAULT nextval('public.jobs_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.notification_preferences ALTER COLUMN id SET DEFAULT nextval('public.notification_preferences_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.notifications ALTER COLUMN id SET DEFAULT nextval('public.notifications_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.oauth_providers ALTER COLUMN id SET DEFAULT nextval('public.oauth_providers_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.otp_codes ALTER COLUMN id SET DEFAULT nextval('public.otp_codes_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.provinces ALTER COLUMN id SET DEFAULT nextval('public.provinces_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.push_notification_logs ALTER COLUMN id SET DEFAULT nextval('public.push_notification_logs_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.skills_master ALTER COLUMN id SET DEFAULT nextval('public.skills_master_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.user_certifications ALTER COLUMN id SET DEFAULT nextval('public.user_certifications_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.user_documents ALTER COLUMN id SET DEFAULT nextval('public.user_documents_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.user_educations ALTER COLUMN id SET DEFAULT nextval('public.user_educations_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.user_experiences ALTER COLUMN id SET DEFAULT nextval('public.user_experiences_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.user_languages ALTER COLUMN id SET DEFAULT nextval('public.user_languages_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.user_preferences ALTER COLUMN id SET DEFAULT nextval('public.user_preferences_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.user_profiles ALTER COLUMN id SET DEFAULT nextval('public.user_profiles_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.user_projects ALTER COLUMN id SET DEFAULT nextval('public.user_projects_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.user_skills ALTER COLUMN id SET DEFAULT nextval('public.user_skills_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq1'::regclass);


--

--

ALTER TABLE ONLY public.work_policies ALTER COLUMN id SET DEFAULT nextval('public.work_policies_id_seq'::regclass);


--

--

ALTER TABLE ONLY public.admin_roles
    ADD CONSTRAINT admin_roles_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.admin_roles
    ADD CONSTRAINT admin_roles_role_name_key UNIQUE (role_name);


--

--

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_email_key UNIQUE (email);


--

--

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.application_documents
    ADD CONSTRAINT application_documents_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.application_notes
    ADD CONSTRAINT application_notes_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.benefits_master
    ADD CONSTRAINT benefits_master_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.benefits_master
    ADD CONSTRAINT benefits_master_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.benefits_master
    ADD CONSTRAINT benefits_master_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.cities
    ADD CONSTRAINT cities_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_slug_key UNIQUE (slug);


--

--

ALTER TABLE ONLY public.company_addresses
    ADD CONSTRAINT company_addresses_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.company_documents
    ADD CONSTRAINT company_documents_company_id_document_type_document_number_key UNIQUE (company_id, document_type, document_number);


--

--

ALTER TABLE ONLY public.company_documents
    ADD CONSTRAINT company_documents_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.company_employees
    ADD CONSTRAINT company_employees_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.company_followers
    ADD CONSTRAINT company_followers_company_id_user_id_key UNIQUE (company_id, user_id);


--

--

ALTER TABLE ONLY public.company_followers
    ADD CONSTRAINT company_followers_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.company_industries
    ADD CONSTRAINT company_industries_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.company_industries
    ADD CONSTRAINT company_industries_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.company_industries
    ADD CONSTRAINT company_industries_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.company_invitations
    ADD CONSTRAINT company_invitations_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.company_invitations
    ADD CONSTRAINT company_invitations_token_key UNIQUE (token);


--

--

ALTER TABLE ONLY public.company_profiles
    ADD CONSTRAINT company_profiles_company_id_key UNIQUE (company_id);


--

--

ALTER TABLE ONLY public.company_profiles
    ADD CONSTRAINT company_profiles_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.company_reviews
    ADD CONSTRAINT company_reviews_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.company_sizes
    ADD CONSTRAINT company_sizes_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.company_verifications
    ADD CONSTRAINT company_verifications_company_id_key UNIQUE (company_id);


--

--

ALTER TABLE ONLY public.company_verifications
    ADD CONSTRAINT company_verifications_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.device_tokens
    ADD CONSTRAINT device_tokens_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.districts
    ADD CONSTRAINT districts_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.education_levels
    ADD CONSTRAINT education_levels_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.education_levels
    ADD CONSTRAINT education_levels_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.education_levels
    ADD CONSTRAINT education_levels_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.email_logs
    ADD CONSTRAINT email_logs_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.employer_users
    ADD CONSTRAINT employer_users_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.employer_users
    ADD CONSTRAINT employer_users_user_id_company_id_key UNIQUE (user_id, company_id);


--

--

ALTER TABLE ONLY public.experience_levels
    ADD CONSTRAINT experience_levels_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.experience_levels
    ADD CONSTRAINT experience_levels_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.experience_levels
    ADD CONSTRAINT experience_levels_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.gender_preferences
    ADD CONSTRAINT gender_preferences_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.gender_preferences
    ADD CONSTRAINT gender_preferences_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.gender_preferences
    ADD CONSTRAINT gender_preferences_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.industries
    ADD CONSTRAINT industries_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.industries
    ADD CONSTRAINT industries_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.industries
    ADD CONSTRAINT industries_slug_key UNIQUE (slug);


--

--

ALTER TABLE ONLY public.interviews
    ADD CONSTRAINT interviews_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.job_application_stages
    ADD CONSTRAINT job_application_stages_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.job_applications
    ADD CONSTRAINT job_applications_job_id_user_id_key UNIQUE (job_id, user_id);


--

--

ALTER TABLE ONLY public.job_applications
    ADD CONSTRAINT job_applications_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.job_benefits
    ADD CONSTRAINT job_benefits_job_id_benefit_name_key UNIQUE (job_id, benefit_name);


--

--

ALTER TABLE ONLY public.job_benefits
    ADD CONSTRAINT job_benefits_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.job_categories
    ADD CONSTRAINT job_categories_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.job_categories
    ADD CONSTRAINT job_categories_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.job_categories
    ADD CONSTRAINT job_categories_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.job_locations
    ADD CONSTRAINT job_locations_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.job_requirements
    ADD CONSTRAINT job_requirements_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.job_skills
    ADD CONSTRAINT job_skills_job_id_skill_id_key UNIQUE (job_id, skill_id);


--

--

ALTER TABLE ONLY public.job_skills
    ADD CONSTRAINT job_skills_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.job_subcategories
    ADD CONSTRAINT job_subcategories_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.job_subcategories
    ADD CONSTRAINT job_subcategories_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.job_subcategories
    ADD CONSTRAINT job_subcategories_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.job_titles
    ADD CONSTRAINT job_titles_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.job_titles
    ADD CONSTRAINT job_titles_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.job_titles
    ADD CONSTRAINT job_titles_slug_key UNIQUE (normalized_name);


--

--

ALTER TABLE ONLY public.job_types
    ADD CONSTRAINT job_types_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.job_types
    ADD CONSTRAINT job_types_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.job_types
    ADD CONSTRAINT job_types_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_slug_key UNIQUE (slug);


--

--

ALTER TABLE ONLY public.notification_preferences
    ADD CONSTRAINT notification_preferences_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.notification_preferences
    ADD CONSTRAINT notification_preferences_user_id_key UNIQUE (user_id);


--

--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.oauth_providers
    ADD CONSTRAINT oauth_providers_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.oauth_providers
    ADD CONSTRAINT oauth_providers_provider_provider_user_id_key UNIQUE (provider, provider_user_id);


--

--

ALTER TABLE ONLY public.otp_codes
    ADD CONSTRAINT otp_codes_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.provinces
    ADD CONSTRAINT provinces_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.provinces
    ADD CONSTRAINT provinces_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.push_notification_logs
    ADD CONSTRAINT push_notification_logs_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_hash_key UNIQUE (token_hash);


--

--

ALTER TABLE ONLY public.skills_master
    ADD CONSTRAINT skills_master_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.skills_master
    ADD CONSTRAINT skills_master_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.skills_master
    ADD CONSTRAINT skills_master_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.device_tokens
    ADD CONSTRAINT unique_user_token UNIQUE (user_id, token);


--

--

ALTER TABLE ONLY public.user_certifications
    ADD CONSTRAINT user_certifications_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.user_documents
    ADD CONSTRAINT user_documents_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.user_educations
    ADD CONSTRAINT user_educations_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.user_experiences
    ADD CONSTRAINT user_experiences_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.user_languages
    ADD CONSTRAINT user_languages_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.user_preferences
    ADD CONSTRAINT user_preferences_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.user_preferences
    ADD CONSTRAINT user_preferences_user_id_key UNIQUE (user_id);


--

--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_slug_key UNIQUE (slug);


--

--

ALTER TABLE ONLY public.user_projects
    ADD CONSTRAINT user_projects_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.user_skills
    ADD CONSTRAINT user_skills_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--

--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--

--

ALTER TABLE ONLY public.work_policies
    ADD CONSTRAINT work_policies_code_key UNIQUE (code);


--

--

ALTER TABLE ONLY public.work_policies
    ADD CONSTRAINT work_policies_name_key UNIQUE (name);


--

--

ALTER TABLE ONLY public.work_policies
    ADD CONSTRAINT work_policies_pkey PRIMARY KEY (id);


--

--

CREATE INDEX idx_admin_roles_access_level ON public.admin_roles USING btree (access_level);


--

--

CREATE INDEX idx_admin_users_status ON public.admin_users USING btree (status);


--

--

CREATE INDEX idx_application_documents_verified ON public.application_documents USING btree (is_verified);


--

--

CREATE INDEX idx_application_notes_stage_id ON public.application_notes USING btree (stage_id);


--

--

CREATE INDEX idx_benefits_master_popularity ON public.benefits_master USING btree (popularity_score DESC);


--

--

CREATE INDEX idx_cities_code ON public.cities USING btree (code) WHERE (code IS NOT NULL);


--

--

CREATE INDEX idx_cities_name ON public.cities USING btree (name);


--

--

CREATE INDEX idx_cities_province ON public.cities USING btree (province_id, is_active, deleted_at) WHERE (deleted_at IS NULL);


--

--

CREATE INDEX idx_cities_type ON public.cities USING btree (type);


--

--

CREATE INDEX idx_companies_city ON public.companies USING btree (city_id);


--

--

CREATE INDEX idx_companies_district ON public.companies USING btree (district_id);


--

--

CREATE INDEX idx_companies_industry ON public.companies USING btree (industry_id);


--

--

CREATE INDEX idx_companies_province ON public.companies USING btree (province_id);


--

--

CREATE INDEX idx_companies_size ON public.companies USING btree (company_size_id);


--

--

CREATE INDEX idx_companies_verified ON public.companies USING btree (verified);


--

--

CREATE INDEX idx_company_addresses_company_id ON public.company_addresses USING btree (company_id);


--

--

CREATE INDEX idx_company_addresses_deleted_at ON public.company_addresses USING btree (deleted_at);


--

--

CREATE INDEX idx_company_documents_expiry ON public.company_documents USING btree (expiry_date);


--

--

CREATE INDEX idx_company_employees_type ON public.company_employees USING btree (employment_type);


--

--

CREATE INDEX idx_company_followers_active ON public.company_followers USING btree (is_active);


--

--

CREATE INDEX idx_company_industries_active ON public.company_industries USING btree (is_active);


--

--

CREATE INDEX idx_company_invitations_company_id ON public.company_invitations USING btree (company_id);


--

--

CREATE INDEX idx_company_invitations_email ON public.company_invitations USING btree (email);


--

--

CREATE INDEX idx_company_invitations_expires_at ON public.company_invitations USING btree (expires_at);


--

--

CREATE INDEX idx_company_invitations_status ON public.company_invitations USING btree (status);


--

--

CREATE INDEX idx_company_invitations_token ON public.company_invitations USING btree (token);


--

--

CREATE INDEX idx_company_profiles_status ON public.company_profiles USING btree (status);


--

--

CREATE INDEX idx_company_reviews_status ON public.company_reviews USING btree (status);


--

--

CREATE INDEX idx_company_sizes_active ON public.company_sizes USING btree (is_active, deleted_at) WHERE (deleted_at IS NULL);


--

--

CREATE INDEX idx_company_sizes_display_order ON public.company_sizes USING btree (display_order);


--

--

CREATE INDEX idx_company_verifications_expiry ON public.company_verifications USING btree (verification_expiry);


--

--

CREATE INDEX idx_device_tokens_failure ON public.device_tokens USING btree (failure_count, last_failure_at) WHERE (failure_count > 0);


--

--

CREATE INDEX idx_device_tokens_inactive ON public.device_tokens USING btree (is_active, last_used_at);


--

--

CREATE INDEX idx_device_tokens_platform ON public.device_tokens USING btree (platform, is_active);


--

--

CREATE INDEX idx_device_tokens_token ON public.device_tokens USING btree (token) WHERE (is_active = true);


--

--

CREATE INDEX idx_device_tokens_user_id ON public.device_tokens USING btree (user_id) WHERE (is_active = true);


--

--

CREATE INDEX idx_device_tokens_user_platform ON public.device_tokens USING btree (user_id, platform) WHERE (is_active = true);


--

--

CREATE INDEX idx_districts_city ON public.districts USING btree (city_id, is_active, deleted_at) WHERE (deleted_at IS NULL);


--

--

CREATE INDEX idx_districts_code ON public.districts USING btree (code) WHERE (code IS NOT NULL);


--

--

CREATE INDEX idx_districts_name ON public.districts USING btree (name);


--

--

CREATE INDEX idx_districts_postal_code ON public.districts USING btree (postal_code) WHERE (postal_code IS NOT NULL);


--

--

CREATE INDEX idx_education_levels_code ON public.education_levels USING btree (code);


--

--

CREATE INDEX idx_education_levels_display_order ON public.education_levels USING btree ("order");


--

--

CREATE INDEX idx_email_logs_created_at ON public.email_logs USING btree (created_at);


--

--

CREATE INDEX idx_email_logs_recipient ON public.email_logs USING btree (recipient);


--

--

CREATE INDEX idx_email_logs_status ON public.email_logs USING btree (status);


--

--

CREATE INDEX idx_email_logs_template ON public.email_logs USING btree (template);


--

--

CREATE INDEX idx_employer_users_active ON public.employer_users USING btree (is_active);


--

--

CREATE INDEX idx_experience_levels_code ON public.experience_levels USING btree (code);


--

--

CREATE INDEX idx_experience_levels_display_order ON public.experience_levels USING btree ("order");


--

--

CREATE INDEX idx_experience_levels_min_years ON public.experience_levels USING btree (min_years);


--

--

CREATE INDEX idx_gender_preferences_code ON public.gender_preferences USING btree (code);


--

--

CREATE INDEX idx_gender_preferences_display_order ON public.gender_preferences USING btree ("order");


--

--

CREATE INDEX idx_industries_active ON public.industries USING btree (is_active, deleted_at) WHERE (deleted_at IS NULL);


--

--

CREATE INDEX idx_industries_display_order ON public.industries USING btree (display_order);


--

--

CREATE INDEX idx_industries_slug ON public.industries USING btree (slug);


--

--

CREATE INDEX idx_interviews_date ON public.interviews USING btree (scheduled_at);


--

--

CREATE INDEX idx_job_application_stages_started_at ON public.job_application_stages USING btree (started_at);


--

--

CREATE INDEX idx_job_applications_company_id ON public.job_applications USING btree (company_id);


--

--

CREATE INDEX idx_job_benefits_highlight ON public.job_benefits USING btree (is_highlight);


--

--

CREATE INDEX idx_job_categories_active ON public.job_categories USING btree (is_active);


--

--

CREATE INDEX idx_job_locations_geo ON public.job_locations USING gist (point((longitude)::double precision, (latitude)::double precision));


--

--

CREATE INDEX idx_job_requirements_skill_id ON public.job_requirements USING btree (skill_id);


--

--

CREATE INDEX idx_job_skills_importance ON public.job_skills USING btree (importance_level);


--

--

CREATE INDEX idx_job_subcategories_active ON public.job_subcategories USING btree (is_active);


--

--

CREATE INDEX idx_job_titles_is_active ON public.job_titles USING btree (is_active);


--

--

CREATE INDEX idx_job_titles_name ON public.job_titles USING btree (name);


--

--

CREATE INDEX idx_job_titles_popularity ON public.job_titles USING btree (popularity_score DESC);


--

--

CREATE INDEX idx_job_titles_search_count ON public.job_titles USING btree (search_count DESC);


--

--

CREATE INDEX idx_job_titles_slug ON public.job_titles USING btree (normalized_name);


--

--

CREATE INDEX idx_job_types_code ON public.job_types USING btree (code);


--

--

CREATE INDEX idx_job_types_display_order ON public.job_types USING btree ("order");


--

--

CREATE INDEX idx_jobs_age_range ON public.jobs USING btree (min_age, max_age) WHERE ((min_age IS NOT NULL) OR (max_age IS NOT NULL));


--

--

CREATE INDEX idx_jobs_company_address ON public.jobs USING btree (company_address_id) WHERE (company_address_id IS NOT NULL);


--

--

CREATE INDEX idx_jobs_education_level_id ON public.jobs USING btree (education_level_id);


--

--

CREATE INDEX idx_jobs_experience_level_id ON public.jobs USING btree (experience_level_id);


--

--

CREATE INDEX idx_jobs_gender_preference_id ON public.jobs USING btree (gender_preference_id);


--

--

CREATE INDEX idx_jobs_job_title_id ON public.jobs USING btree (job_title_id);


--

--

CREATE INDEX idx_jobs_job_type_id ON public.jobs USING btree (job_type_id);


--

--

CREATE INDEX idx_jobs_status ON public.jobs USING btree (status);


--

--

CREATE INDEX idx_jobs_work_policy_id ON public.jobs USING btree (work_policy_id);


--

--

CREATE UNIQUE INDEX idx_notification_preferences_user_id ON public.notification_preferences USING btree (user_id);


--

--

CREATE INDEX idx_notifications_category ON public.notifications USING btree (category);


--

--

CREATE INDEX idx_notifications_created_at ON public.notifications USING btree (created_at);


--

--

CREATE INDEX idx_notifications_is_read ON public.notifications USING btree (is_read);


--

--

CREATE INDEX idx_notifications_related_id ON public.notifications USING btree (related_id);


--

--

CREATE INDEX idx_notifications_sender_id ON public.notifications USING btree (sender_id);


--

--

CREATE INDEX idx_notifications_type ON public.notifications USING btree (type);


--

--

CREATE INDEX idx_notifications_user_category ON public.notifications USING btree (user_id, category);


--

--

CREATE INDEX idx_notifications_user_id ON public.notifications USING btree (user_id);


--

--

CREATE INDEX idx_notifications_user_read ON public.notifications USING btree (user_id, is_read);


--

--

CREATE INDEX idx_oauth_providers_email ON public.oauth_providers USING btree (email);


--

--

CREATE INDEX idx_oauth_providers_provider ON public.oauth_providers USING btree (provider);


--

--

CREATE INDEX idx_oauth_providers_user_id ON public.oauth_providers USING btree (user_id);


--

--

CREATE INDEX idx_otp_codes_expired_at ON public.otp_codes USING btree (expired_at);


--

--

CREATE INDEX idx_otp_codes_is_used ON public.otp_codes USING btree (is_used);


--

--

CREATE INDEX idx_otp_codes_type ON public.otp_codes USING btree (type);


--

--

CREATE INDEX idx_otp_codes_user_id ON public.otp_codes USING btree (user_id);


--

--

CREATE INDEX idx_otp_codes_user_type ON public.otp_codes USING btree (user_id, type);


--

--

CREATE INDEX idx_provinces_active ON public.provinces USING btree (is_active, deleted_at) WHERE (deleted_at IS NULL);


--

--

CREATE INDEX idx_provinces_code ON public.provinces USING btree (code) WHERE (code IS NOT NULL);


--

--

CREATE INDEX idx_provinces_name ON public.provinces USING btree (name);


--

--

CREATE INDEX idx_push_logs_device_token_id ON public.push_notification_logs USING btree (device_token_id, created_at DESC) WHERE (device_token_id IS NOT NULL);


--

--

CREATE INDEX idx_push_logs_failed ON public.push_notification_logs USING btree (user_id, status, created_at DESC) WHERE ((status)::text = 'failed'::text);


--

--

CREATE INDEX idx_push_logs_fcm_message_id ON public.push_notification_logs USING btree (fcm_message_id) WHERE (fcm_message_id IS NOT NULL);


--

--

CREATE INDEX idx_push_logs_notification_id ON public.push_notification_logs USING btree (notification_id) WHERE (notification_id IS NOT NULL);


--

--

CREATE INDEX idx_push_logs_status ON public.push_notification_logs USING btree (status, created_at DESC);


--

--

CREATE INDEX idx_push_logs_user_id ON public.push_notification_logs USING btree (user_id, created_at DESC);


--

--

CREATE INDEX idx_refresh_tokens_device_id ON public.refresh_tokens USING btree (device_id);


--

--

CREATE INDEX idx_refresh_tokens_expires_at ON public.refresh_tokens USING btree (expires_at);


--

--

CREATE INDEX idx_refresh_tokens_revoked ON public.refresh_tokens USING btree (revoked);


--

--

CREATE INDEX idx_refresh_tokens_token_hash ON public.refresh_tokens USING btree (token_hash);


--

--

CREATE INDEX idx_refresh_tokens_user_device ON public.refresh_tokens USING btree (user_id, device_id);


--

--

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);


--

--

CREATE INDEX idx_skills_master_popularity ON public.skills_master USING btree (popularity_score DESC);


--

--

CREATE INDEX idx_user_certifications_verified ON public.user_certifications USING btree (verified);


--

--

CREATE INDEX idx_user_educations_degree ON public.user_educations USING btree (degree_level);


--

--

CREATE INDEX idx_user_experiences_employment_type ON public.user_experiences USING btree (employment_type);


--

--

CREATE INDEX idx_user_languages_verified ON public.user_languages USING btree (verified);


--

--

CREATE INDEX idx_user_preferences_jobtype ON public.user_preferences USING btree (preferred_job_type);


--

--

CREATE INDEX idx_user_profiles_city ON public.user_profiles USING btree (city_id);


--

--

CREATE INDEX idx_user_profiles_district ON public.user_profiles USING btree (district_id);


--

--

CREATE INDEX idx_user_profiles_experience_level ON public.user_profiles USING btree (experience_level);


--

--

CREATE INDEX idx_user_profiles_location_state ON public.user_profiles USING btree (location_state);


--

--

CREATE INDEX idx_user_profiles_nationality ON public.user_profiles USING btree (nationality);


--

--

CREATE INDEX idx_user_profiles_province ON public.user_profiles USING btree (province_id);


--

--

CREATE INDEX idx_user_projects_visibility ON public.user_projects USING btree (visibility);


--

--

CREATE INDEX idx_user_skills_level ON public.user_skills USING btree (skill_level);


--

--

CREATE INDEX idx_users_company ON public.users USING btree (company_id);


--

--

CREATE INDEX idx_users_created_at ON public.users USING btree (created_at DESC);


--

--

CREATE INDEX idx_users_has_company ON public.users USING btree (has_company) WHERE (has_company = true);


--

--

CREATE INDEX idx_work_policies_code ON public.work_policies USING btree (code);


--

--

CREATE INDEX idx_work_policies_display_order ON public.work_policies USING btree ("order");


--

--

CREATE TRIGGER trigger_device_tokens_updated_at BEFORE UPDATE ON public.device_tokens FOR EACH ROW EXECUTE FUNCTION public.update_device_tokens_updated_at();


--

--

CREATE TRIGGER trigger_update_company_invitations_updated_at BEFORE UPDATE ON public.company_invitations FOR EACH ROW EXECUTE FUNCTION public.update_company_invitations_updated_at();


--

--

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.admin_users(id);


--

--

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.admin_roles(id);


--

--

ALTER TABLE ONLY public.application_documents
    ADD CONSTRAINT application_documents_application_id_fkey FOREIGN KEY (application_id) REFERENCES public.job_applications(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.application_documents
    ADD CONSTRAINT application_documents_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.application_documents
    ADD CONSTRAINT application_documents_verified_by_fkey FOREIGN KEY (verified_by) REFERENCES public.admin_users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.application_notes
    ADD CONSTRAINT application_notes_application_id_fkey FOREIGN KEY (application_id) REFERENCES public.job_applications(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.application_notes
    ADD CONSTRAINT application_notes_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.admin_users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.application_notes
    ADD CONSTRAINT application_notes_stage_id_fkey FOREIGN KEY (stage_id) REFERENCES public.job_application_stages(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_verified_by_fkey FOREIGN KEY (verified_by) REFERENCES public.admin_users(id);


--

--

ALTER TABLE ONLY public.company_addresses
    ADD CONSTRAINT company_addresses_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.company_documents
    ADD CONSTRAINT company_documents_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.company_documents
    ADD CONSTRAINT company_documents_uploaded_by_fkey FOREIGN KEY (uploaded_by) REFERENCES public.employer_users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.company_documents
    ADD CONSTRAINT company_documents_verified_by_fkey FOREIGN KEY (verified_by) REFERENCES public.admin_users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.company_employees
    ADD CONSTRAINT company_employees_added_by_fkey FOREIGN KEY (added_by) REFERENCES public.employer_users(id);


--

--

ALTER TABLE ONLY public.company_employees
    ADD CONSTRAINT company_employees_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.company_employees
    ADD CONSTRAINT company_employees_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.company_employees
    ADD CONSTRAINT company_employees_verified_by_fkey FOREIGN KEY (verified_by) REFERENCES public.admin_users(id);


--

--

ALTER TABLE ONLY public.company_followers
    ADD CONSTRAINT company_followers_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.company_followers
    ADD CONSTRAINT company_followers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.company_industries
    ADD CONSTRAINT company_industries_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.company_industries(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.company_invitations
    ADD CONSTRAINT company_invitations_accepted_by_fkey FOREIGN KEY (accepted_by) REFERENCES public.users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.company_invitations
    ADD CONSTRAINT company_invitations_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.company_invitations
    ADD CONSTRAINT company_invitations_invited_by_fkey FOREIGN KEY (invited_by) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.company_profiles
    ADD CONSTRAINT company_profiles_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.company_profiles
    ADD CONSTRAINT company_profiles_verified_by_fkey FOREIGN KEY (verified_by) REFERENCES public.admin_users(id);


--

--

ALTER TABLE ONLY public.company_reviews
    ADD CONSTRAINT company_reviews_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.company_reviews
    ADD CONSTRAINT company_reviews_moderated_by_fkey FOREIGN KEY (moderated_by) REFERENCES public.admin_users(id);


--

--

ALTER TABLE ONLY public.company_reviews
    ADD CONSTRAINT company_reviews_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.company_verifications
    ADD CONSTRAINT company_verifications_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.company_verifications
    ADD CONSTRAINT company_verifications_requested_by_fkey FOREIGN KEY (requested_by) REFERENCES public.employer_users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.company_verifications
    ADD CONSTRAINT company_verifications_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES public.admin_users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.employer_users
    ADD CONSTRAINT employer_users_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.employer_users
    ADD CONSTRAINT employer_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.employer_users
    ADD CONSTRAINT employer_users_verified_by_fkey FOREIGN KEY (verified_by) REFERENCES public.admin_users(id);


--

--

ALTER TABLE ONLY public.cities
    ADD CONSTRAINT fk_cities_province FOREIGN KEY (province_id) REFERENCES public.provinces(id) ON DELETE RESTRICT;


--

--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT fk_companies_city FOREIGN KEY (city_id) REFERENCES public.cities(id) ON DELETE RESTRICT;


--

--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT fk_companies_company_size FOREIGN KEY (company_size_id) REFERENCES public.company_sizes(id) ON DELETE RESTRICT;


--

--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT fk_companies_district FOREIGN KEY (district_id) REFERENCES public.districts(id) ON DELETE RESTRICT;


--

--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT fk_companies_industry FOREIGN KEY (industry_id) REFERENCES public.industries(id) ON DELETE RESTRICT;


--

--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT fk_companies_province FOREIGN KEY (province_id) REFERENCES public.provinces(id) ON DELETE RESTRICT;


--

--

ALTER TABLE ONLY public.device_tokens
    ADD CONSTRAINT fk_device_tokens_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.districts
    ADD CONSTRAINT fk_districts_city FOREIGN KEY (city_id) REFERENCES public.cities(id) ON DELETE RESTRICT;


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT fk_jobs_education_level FOREIGN KEY (education_level_id) REFERENCES public.education_levels(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT fk_jobs_experience_level FOREIGN KEY (experience_level_id) REFERENCES public.experience_levels(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT fk_jobs_gender_preference FOREIGN KEY (gender_preference_id) REFERENCES public.gender_preferences(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT fk_jobs_job_subcategory_id FOREIGN KEY (job_subcategory_id) REFERENCES public.job_subcategories(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT fk_jobs_job_title FOREIGN KEY (job_title_id) REFERENCES public.job_titles(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT fk_jobs_job_type FOREIGN KEY (job_type_id) REFERENCES public.job_types(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT fk_jobs_work_policy FOREIGN KEY (work_policy_id) REFERENCES public.work_policies(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.otp_codes
    ADD CONSTRAINT fk_otp_codes_user_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.push_notification_logs
    ADD CONSTRAINT fk_push_logs_device_token FOREIGN KEY (device_token_id) REFERENCES public.device_tokens(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.push_notification_logs
    ADD CONSTRAINT fk_push_logs_notification FOREIGN KEY (notification_id) REFERENCES public.notifications(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.push_notification_logs
    ADD CONSTRAINT fk_push_logs_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT fk_refresh_tokens_user_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT fk_user_profiles_city FOREIGN KEY (city_id) REFERENCES public.cities(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT fk_user_profiles_district FOREIGN KEY (district_id) REFERENCES public.districts(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT fk_user_profiles_province FOREIGN KEY (province_id) REFERENCES public.provinces(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_company FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.interviews
    ADD CONSTRAINT interviews_application_id_fkey FOREIGN KEY (application_id) REFERENCES public.job_applications(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.interviews
    ADD CONSTRAINT interviews_interviewer_id_fkey FOREIGN KEY (interviewer_id) REFERENCES public.admin_users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.interviews
    ADD CONSTRAINT interviews_stage_id_fkey FOREIGN KEY (stage_id) REFERENCES public.job_application_stages(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.job_application_stages
    ADD CONSTRAINT job_application_stages_application_id_fkey FOREIGN KEY (application_id) REFERENCES public.job_applications(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.job_application_stages
    ADD CONSTRAINT job_application_stages_handled_by_fkey FOREIGN KEY (handled_by) REFERENCES public.admin_users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.job_applications
    ADD CONSTRAINT job_applications_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.job_applications
    ADD CONSTRAINT job_applications_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.job_applications
    ADD CONSTRAINT job_applications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.job_benefits
    ADD CONSTRAINT job_benefits_benefit_id_fkey FOREIGN KEY (benefit_id) REFERENCES public.benefits_master(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.job_benefits
    ADD CONSTRAINT job_benefits_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.job_categories
    ADD CONSTRAINT job_categories_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.job_categories(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.job_locations
    ADD CONSTRAINT job_locations_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.job_locations
    ADD CONSTRAINT job_locations_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.job_requirements
    ADD CONSTRAINT job_requirements_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.job_requirements
    ADD CONSTRAINT job_requirements_skill_id_fkey FOREIGN KEY (skill_id) REFERENCES public.skills_master(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.job_skills
    ADD CONSTRAINT job_skills_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.job_skills
    ADD CONSTRAINT job_skills_skill_id_fkey FOREIGN KEY (skill_id) REFERENCES public.skills_master(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.job_subcategories
    ADD CONSTRAINT job_subcategories_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.job_categories(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.job_titles
    ADD CONSTRAINT job_titles_recommended_category_id_fkey FOREIGN KEY (recommended_category_id) REFERENCES public.job_categories(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.job_categories(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_employer_user_id_fkey FOREIGN KEY (employer_user_id) REFERENCES public.employer_users(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.notification_preferences
    ADD CONSTRAINT notification_preferences_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.oauth_providers
    ADD CONSTRAINT oauth_providers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.skills_master
    ADD CONSTRAINT skills_master_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.job_categories(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.skills_master
    ADD CONSTRAINT skills_master_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.skills_master(id) ON DELETE SET NULL;


--

--

ALTER TABLE ONLY public.user_certifications
    ADD CONSTRAINT user_certifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.user_documents
    ADD CONSTRAINT user_documents_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.user_documents
    ADD CONSTRAINT user_documents_verified_by_fkey FOREIGN KEY (verified_by) REFERENCES public.admin_users(id);


--

--

ALTER TABLE ONLY public.user_educations
    ADD CONSTRAINT user_educations_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.user_experiences
    ADD CONSTRAINT user_experiences_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.user_languages
    ADD CONSTRAINT user_languages_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.user_preferences
    ADD CONSTRAINT user_preferences_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.user_projects
    ADD CONSTRAINT user_projects_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

ALTER TABLE ONLY public.user_skills
    ADD CONSTRAINT user_skills_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--

--

-- Create conversations table
CREATE TABLE IF NOT EXISTS public.conversations (
    id bigserial PRIMARY KEY,
    uuid uuid DEFAULT gen_random_uuid() UNIQUE NOT NULL,
    last_message_at timestamp,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    deleted_at timestamp
);

-- Create chat_participants table
CREATE TABLE IF NOT EXISTS public.chat_participants (
    id bigserial PRIMARY KEY,
    conversation_id bigint NOT NULL,
    user_id bigint NOT NULL,
    is_archived boolean DEFAULT false NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    CONSTRAINT chat_participants_conversation_id_fkey 
        FOREIGN KEY (conversation_id) 
        REFERENCES public.conversations(id) 
        ON DELETE CASCADE,
    CONSTRAINT chat_participants_user_id_fkey 
        FOREIGN KEY (user_id) 
        REFERENCES public.users(id) 
        ON DELETE CASCADE,
    CONSTRAINT chat_participants_conversation_user_unique 
        UNIQUE (conversation_id, user_id)
);

-- Create messages table
CREATE TABLE IF NOT EXISTS public.messages (
    id bigserial PRIMARY KEY,
    uuid uuid DEFAULT gen_random_uuid() UNIQUE NOT NULL,
    conversation_id bigint NOT NULL,
    sender_id bigint NOT NULL,
    content varchar(5000) NOT NULL,
    is_read boolean DEFAULT false NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    deleted_at timestamp,
    CONSTRAINT messages_conversation_id_fkey 
        FOREIGN KEY (conversation_id) 
        REFERENCES public.conversations(id) 
        ON DELETE CASCADE,
    CONSTRAINT messages_sender_id_fkey 
        FOREIGN KEY (sender_id) 
        REFERENCES public.users(id) 
        ON DELETE CASCADE
);

-- Create indexes for conversations
CREATE INDEX IF NOT EXISTS idx_conversations_deleted_at ON public.conversations USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_conversations_last_message_at ON public.conversations USING btree (last_message_at DESC);

-- Create indexes for chat_participants
CREATE INDEX IF NOT EXISTS idx_chat_participants_user_id ON public.chat_participants USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_chat_participants_conversation_id ON public.chat_participants USING btree (conversation_id);

-- Create indexes for messages
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON public.messages USING btree (conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON public.messages USING btree (sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON public.messages USING btree (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_deleted_at ON public.messages USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_created ON public.messages USING btree (conversation_id, created_at DESC);


--
-- PostgreSQL database dump complete
--
