-- Migration: Create Job Categories and Subcategories
-- Created: 2025-11-07
-- Purpose: Create job_categories and job_subcategories tables which are referenced by job master tables

CREATE TABLE IF NOT EXISTS job_categories (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT,
    code VARCHAR(30) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_job_categories_parent FOREIGN KEY (parent_id) REFERENCES job_categories(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_job_categories_code ON job_categories(code);
CREATE INDEX IF NOT EXISTS idx_job_categories_is_active ON job_categories(is_active);

-- Subcategories
CREATE TABLE IF NOT EXISTS job_subcategories (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_job_subcategories_category FOREIGN KEY (category_id) REFERENCES job_categories(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_job_subcategories_category_id ON job_subcategories(category_id);
CREATE INDEX IF NOT EXISTS idx_job_subcategories_code ON job_subcategories(code);
CREATE INDEX IF NOT EXISTS idx_job_subcategories_is_active ON job_subcategories(is_active);

COMMENT ON TABLE job_categories IS 'Master table for job categories (Technology, Finance, etc)';
COMMENT ON TABLE job_subcategories IS 'Subcategories for job categories (e.g., Software Development under Technology)';
