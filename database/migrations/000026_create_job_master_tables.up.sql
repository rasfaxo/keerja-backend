-- Migration: Create Job Master Data Tables (Phase 1-4)
-- Created: 2025-11-04
-- Purpose: Create master data tables for job posting workflow

-- ============================================
-- Table: job_titles
-- Purpose: Master data for job titles with search functionality
-- ============================================
CREATE TABLE IF NOT EXISTS job_titles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL UNIQUE,
    normalized_name VARCHAR(200),
    recommended_category_id BIGINT REFERENCES job_categories(id) ON DELETE SET NULL,
    popularity_score NUMERIC(5,2) DEFAULT 0.00,
    search_count BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_job_titles_name ON job_titles(name);
CREATE INDEX idx_job_titles_normalized ON job_titles(normalized_name);
CREATE INDEX idx_job_titles_popularity ON job_titles(popularity_score DESC);
CREATE INDEX idx_job_titles_search_count ON job_titles(search_count DESC);
CREATE INDEX idx_job_titles_is_active ON job_titles(is_active);

COMMENT ON TABLE job_titles IS 'Master data for job titles (e.g., Software Engineer, Data Analyst)';
COMMENT ON COLUMN job_titles.normalized_name IS 'Normalized name for search and matching';
COMMENT ON COLUMN job_titles.recommended_category_id IS 'Recommended job category for this title';
COMMENT ON COLUMN job_titles.popularity_score IS 'Popularity score for ranking (higher = more popular)';
COMMENT ON COLUMN job_titles.search_count IS 'Number of times searched by users';

-- ============================================
-- Table: job_types
-- Purpose: Master data for employment types
-- ============================================
CREATE TABLE IF NOT EXISTS job_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    "order" INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_job_types_code ON job_types(code);
CREATE INDEX idx_job_types_order ON job_types("order");

COMMENT ON TABLE job_types IS 'Master data for job types (Full-Time, Part-Time, Internship, Freelance, Contract)';
COMMENT ON COLUMN job_types.code IS 'Unique code for programmatic reference';

-- ============================================
-- Table: work_policies
-- Purpose: Master data for work location policies
-- ============================================
CREATE TABLE IF NOT EXISTS work_policies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    "order" INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_work_policies_code ON work_policies(code);
CREATE INDEX idx_work_policies_order ON work_policies("order");

COMMENT ON TABLE work_policies IS 'Master data for work policies (On-site, Remote, Hybrid)';

-- ============================================
-- Table: education_levels
-- Purpose: Master data for minimum education requirements
-- ============================================
CREATE TABLE IF NOT EXISTS education_levels (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    "order" INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_education_levels_code ON education_levels(code);
CREATE INDEX idx_education_levels_order ON education_levels("order");

COMMENT ON TABLE education_levels IS 'Master data for education levels (SMA, D3, S1, S2, S3)';

-- ============================================
-- Table: experience_levels
-- Purpose: Master data for work experience requirements
-- ============================================
CREATE TABLE IF NOT EXISTS experience_levels (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    min_years INTEGER NOT NULL DEFAULT 0,
    max_years INTEGER,
    "order" INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_experience_levels_code ON experience_levels(code);
CREATE INDEX idx_experience_levels_min_years ON experience_levels(min_years);
CREATE INDEX idx_experience_levels_order ON experience_levels("order");

COMMENT ON TABLE experience_levels IS 'Master data for experience levels (Fresh Graduate, 1-3 years, 3-5 years, etc.)';
COMMENT ON COLUMN experience_levels.min_years IS 'Minimum years of experience';
COMMENT ON COLUMN experience_levels.max_years IS 'Maximum years of experience (NULL = unlimited)';

-- ============================================
-- Table: gender_preferences
-- Purpose: Master data for gender requirements in job postings
-- ============================================
CREATE TABLE IF NOT EXISTS gender_preferences (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    "order" INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_gender_preferences_code ON gender_preferences(code);
CREATE INDEX idx_gender_preferences_order ON gender_preferences("order");

COMMENT ON TABLE gender_preferences IS 'Master data for gender preferences (Male, Female, Any)';

-- ============================================
-- Grant permissions (if needed)
-- ============================================
-- GRANT SELECT ON ALL TABLES IN SCHEMA public TO app_user;
