-- Pre-Migration Script: Transfer data from deprecated fields to FK fields
-- Version: 000031 PRE-MIGRATION
-- Run this BEFORE running 000031_refactor_jobs_table_add_requirements.up.sql
-- Description:
--   Transfer data from legacy string fields to new master data FK fields
--   This ensures no data loss when deprecated columns are dropped

-- ============================================================================
-- PREREQUISITE CHECK
-- ============================================================================

DO $$
BEGIN
    -- Check if master data tables exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'job_types') THEN
        RAISE EXCEPTION 'Master data table job_types does not exist. Run migration 000026 first.';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'experience_levels') THEN
        RAISE EXCEPTION 'Master data table experience_levels does not exist. Run migration 000026 first.';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'education_levels') THEN
        RAISE EXCEPTION 'Master data table education_levels does not exist. Run migration 000026 first.';
    END IF;

    RAISE NOTICE 'Prerequisite check passed. All master data tables exist.';
END $$;

-- ============================================================================
-- STEP 1: Analyze existing data
-- ============================================================================

DO $$
DECLARE
    total_jobs INTEGER;
    jobs_with_job_level INTEGER;
    jobs_with_employment_type INTEGER;
    jobs_with_education_level INTEGER;
BEGIN
    SELECT COUNT(*) INTO total_jobs FROM jobs;

    SELECT COUNT(*) INTO jobs_with_job_level
    FROM jobs WHERE job_level IS NOT NULL AND job_level != '';

    SELECT COUNT(*) INTO jobs_with_employment_type
    FROM jobs WHERE employment_type IS NOT NULL AND employment_type != '';

    SELECT COUNT(*) INTO jobs_with_education_level
    FROM jobs WHERE education_level IS NOT NULL AND education_level != '';

    RAISE NOTICE '========================================';
    RAISE NOTICE 'DATA ANALYSIS:';
    RAISE NOTICE 'Total jobs: %', total_jobs;
    RAISE NOTICE 'Jobs with job_level: %', jobs_with_job_level;
    RAISE NOTICE 'Jobs with employment_type: %', jobs_with_employment_type;
    RAISE NOTICE 'Jobs with education_level: %', jobs_with_education_level;
    RAISE NOTICE '========================================';
END $$;

-- ============================================================================
-- STEP 2: Migrate employment_type to job_type_id
-- ============================================================================

DO $$
DECLARE
    updated_count INTEGER := 0;
BEGIN
    RAISE NOTICE 'Migrating employment_type to job_type_id...';

    -- Full-Time
    UPDATE jobs
    SET job_type_id = (SELECT id FROM job_types WHERE LOWER(code) = 'full-time' LIMIT 1)
    WHERE (LOWER(employment_type) = 'full-time' OR LOWER(employment_type) = 'fulltime')
    AND job_type_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % Full-Time jobs', updated_count;

    -- Part-Time
    UPDATE jobs
    SET job_type_id = (SELECT id FROM job_types WHERE LOWER(code) = 'part-time' LIMIT 1)
    WHERE (LOWER(employment_type) = 'part-time' OR LOWER(employment_type) = 'parttime')
    AND job_type_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % Part-Time jobs', updated_count;

    -- Contract
    UPDATE jobs
    SET job_type_id = (SELECT id FROM job_types WHERE LOWER(code) = 'contract' LIMIT 1)
    WHERE LOWER(employment_type) = 'contract'
    AND job_type_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % Contract jobs', updated_count;

    -- Internship
    UPDATE jobs
    SET job_type_id = (SELECT id FROM job_types WHERE LOWER(code) = 'internship' LIMIT 1)
    WHERE (LOWER(employment_type) = 'internship' OR LOWER(employment_type) = 'intern')
    AND job_type_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % Internship jobs', updated_count;

    -- Freelance
    UPDATE jobs
    SET job_type_id = (SELECT id FROM job_types WHERE LOWER(code) = 'freelance' LIMIT 1)
    WHERE LOWER(employment_type) = 'freelance'
    AND job_type_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % Freelance jobs', updated_count;

    RAISE NOTICE 'employment_type migration completed.';
END $$;

-- ============================================================================
-- STEP 3: Migrate job_level to experience_level_id
-- ============================================================================

DO $$
DECLARE
    updated_count INTEGER := 0;
BEGIN
    RAISE NOTICE 'Migrating job_level to experience_level_id...';

    -- Internship / Fresh Graduate
    UPDATE jobs
    SET experience_level_id = (SELECT id FROM experience_levels WHERE LOWER(code) = 'fresh' LIMIT 1)
    WHERE (LOWER(job_level) LIKE '%intern%' OR LOWER(job_level) LIKE '%fresh%')
    AND experience_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % Fresh/Internship jobs', updated_count;

    -- Entry Level / Junior (1-2 years)
    UPDATE jobs
    SET experience_level_id = (SELECT id FROM experience_levels WHERE LOWER(code) = 'junior' LIMIT 1)
    WHERE (LOWER(job_level) LIKE '%entry%' OR LOWER(job_level) LIKE '%junior%')
    AND experience_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % Entry/Junior Level jobs', updated_count;

    -- Mid Level (3-5 years)
    UPDATE jobs
    SET experience_level_id = (SELECT id FROM experience_levels WHERE LOWER(code) = 'mid' LIMIT 1)
    WHERE (LOWER(job_level) LIKE '%mid%' OR LOWER(job_level) LIKE '%intermediate%')
    AND experience_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % Mid Level jobs', updated_count;

    -- Senior Level (6-10 years)
    UPDATE jobs
    SET experience_level_id = (SELECT id FROM experience_levels WHERE LOWER(code) = 'senior' LIMIT 1)
    WHERE LOWER(job_level) LIKE '%senior%'
    AND experience_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % Senior Level jobs', updated_count;

    -- Manager / Director (10+ years)
    UPDATE jobs
    SET experience_level_id = (SELECT id FROM experience_levels WHERE LOWER(code) = 'expert' LIMIT 1)
    WHERE (LOWER(job_level) LIKE '%manager%' OR LOWER(job_level) LIKE '%director%' OR LOWER(job_level) LIKE '%expert%')
    AND experience_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % Manager/Director/Expert jobs', updated_count;

    -- Any / Not Specified (default for remaining)
    UPDATE jobs
    SET experience_level_id = (SELECT id FROM experience_levels WHERE LOWER(code) = 'any' LIMIT 1)
    WHERE job_level IS NOT NULL
    AND job_level != ''
    AND experience_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % jobs to Any/Not Specified', updated_count;

    RAISE NOTICE 'job_level migration completed.';
END $$;

-- ============================================================================
-- STEP 4: Migrate education_level to education_level_id
-- ============================================================================

DO $$
DECLARE
    updated_count INTEGER := 0;
BEGIN
    RAISE NOTICE 'Migrating education_level to education_level_id...';

    -- SMA/SMK
    UPDATE jobs
    SET education_level_id = (SELECT id FROM education_levels WHERE LOWER(code) = 'sma' LIMIT 1)
    WHERE (LOWER(education_level) LIKE '%sma%' OR LOWER(education_level) LIKE '%smk%' OR LOWER(education_level) LIKE '%high school%')
    AND education_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % SMA/SMK jobs', updated_count;

    -- D3
    UPDATE jobs
    SET education_level_id = (SELECT id FROM education_levels WHERE LOWER(code) = 'd3' LIMIT 1)
    WHERE (LOWER(education_level) = 'd3' OR LOWER(education_level) LIKE '%diploma%')
    AND education_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % D3/Diploma jobs', updated_count;

    -- D4
    UPDATE jobs
    SET education_level_id = (SELECT id FROM education_levels WHERE LOWER(code) = 'd4' LIMIT 1)
    WHERE LOWER(education_level) = 'd4'
    AND education_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % D4 jobs', updated_count;

    -- S1
    UPDATE jobs
    SET education_level_id = (SELECT id FROM education_levels WHERE LOWER(code) = 's1' LIMIT 1)
    WHERE (LOWER(education_level) = 's1' OR LOWER(education_level) LIKE '%bachelor%' OR LOWER(education_level) LIKE '%sarjana%')
    AND education_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % S1/Bachelor jobs', updated_count;

    -- S2
    UPDATE jobs
    SET education_level_id = (SELECT id FROM education_levels WHERE LOWER(code) = 's2' LIMIT 1)
    WHERE (LOWER(education_level) = 's2' OR LOWER(education_level) LIKE '%master%' OR LOWER(education_level) LIKE '%magister%')
    AND education_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % S2/Master jobs', updated_count;

    -- S3
    UPDATE jobs
    SET education_level_id = (SELECT id FROM education_levels WHERE LOWER(code) = 's3' LIMIT 1)
    WHERE (LOWER(education_level) = 's3' OR LOWER(education_level) LIKE '%doctor%' OR LOWER(education_level) LIKE '%phd%')
    AND education_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % S3/PhD jobs', updated_count;

    -- Any / Not Specified (default for remaining)
    UPDATE jobs
    SET education_level_id = (SELECT id FROM education_levels WHERE LOWER(code) = 'any' LIMIT 1)
    WHERE education_level IS NOT NULL
    AND education_level != ''
    AND education_level_id IS NULL;
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE '  - Migrated % jobs to Any/Not Specified', updated_count;

    RAISE NOTICE 'education_level migration completed.';
END $$;

-- ============================================================================
-- STEP 5: Verify migration results
-- ============================================================================

DO $$
DECLARE
    unmigrated_employment INTEGER;
    unmigrated_job_level INTEGER;
    unmigrated_education INTEGER;
    migrated_employment INTEGER;
    migrated_job_level INTEGER;
    migrated_education INTEGER;
BEGIN
    RAISE NOTICE '========================================';
    RAISE NOTICE 'MIGRATION VERIFICATION:';

    -- Count unmigrated records
    SELECT COUNT(*) INTO unmigrated_employment
    FROM jobs
    WHERE employment_type IS NOT NULL
    AND employment_type != ''
    AND job_type_id IS NULL;

    SELECT COUNT(*) INTO unmigrated_job_level
    FROM jobs
    WHERE job_level IS NOT NULL
    AND job_level != ''
    AND experience_level_id IS NULL;

    SELECT COUNT(*) INTO unmigrated_education
    FROM jobs
    WHERE education_level IS NOT NULL
    AND education_level != ''
    AND education_level_id IS NULL;

    -- Count migrated records
    SELECT COUNT(*) INTO migrated_employment
    FROM jobs WHERE job_type_id IS NOT NULL;

    SELECT COUNT(*) INTO migrated_job_level
    FROM jobs WHERE experience_level_id IS NOT NULL;

    SELECT COUNT(*) INTO migrated_education
    FROM jobs WHERE education_level_id IS NOT NULL;

    RAISE NOTICE 'Migrated to job_type_id: %', migrated_employment;
    RAISE NOTICE 'Unmigrated employment_type: %', unmigrated_employment;
    RAISE NOTICE '';
    RAISE NOTICE 'Migrated to experience_level_id: %', migrated_job_level;
    RAISE NOTICE 'Unmigrated job_level: %', unmigrated_job_level;
    RAISE NOTICE '';
    RAISE NOTICE 'Migrated to education_level_id: %', migrated_education;
    RAISE NOTICE 'Unmigrated education_level: %', unmigrated_education;
    RAISE NOTICE '========================================';

    IF unmigrated_employment > 0 THEN
        RAISE WARNING '% jobs with employment_type could not be migrated. Review manually.', unmigrated_employment;
    END IF;

    IF unmigrated_job_level > 0 THEN
        RAISE WARNING '% jobs with job_level could not be migrated. Review manually.', unmigrated_job_level;
    END IF;

    IF unmigrated_education > 0 THEN
        RAISE WARNING '% jobs with education_level could not be migrated. Review manually.', unmigrated_education;
    END IF;
END $$;

-- ============================================================================
-- STEP 6: Show unmigrated records for manual review
-- ============================================================================

-- Show unmigrated employment_type values
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM jobs
        WHERE employment_type IS NOT NULL
        AND employment_type != ''
        AND job_type_id IS NULL
        LIMIT 1
    ) THEN
        RAISE NOTICE '';
        RAISE NOTICE '========================================';
        RAISE NOTICE 'UNMIGRATED EMPLOYMENT_TYPE VALUES:';
        RAISE NOTICE 'Please review these and migrate manually if needed:';
        RAISE NOTICE '========================================';
    END IF;
END $$;

SELECT DISTINCT employment_type, COUNT(*) as count
FROM jobs
WHERE employment_type IS NOT NULL
AND employment_type != ''
AND job_type_id IS NULL
GROUP BY employment_type
ORDER BY count DESC;

-- Show unmigrated job_level values
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM jobs
        WHERE job_level IS NOT NULL
        AND job_level != ''
        AND experience_level_id IS NULL
        LIMIT 1
    ) THEN
        RAISE NOTICE '';
        RAISE NOTICE '========================================';
        RAISE NOTICE 'UNMIGRATED JOB_LEVEL VALUES:';
        RAISE NOTICE 'Please review these and migrate manually if needed:';
        RAISE NOTICE '========================================';
    END IF;
END $$;

SELECT DISTINCT job_level, COUNT(*) as count
FROM jobs
WHERE job_level IS NOT NULL
AND job_level != ''
AND experience_level_id IS NULL
GROUP BY job_level
ORDER BY count DESC;

-- Show unmigrated education_level values
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM jobs
        WHERE education_level IS NOT NULL
        AND education_level != ''
        AND education_level_id IS NULL
        LIMIT 1
    ) THEN
        RAISE NOTICE '';
        RAISE NOTICE '========================================';
        RAISE NOTICE 'UNMIGRATED EDUCATION_LEVEL VALUES:';
        RAISE NOTICE 'Please review these and migrate manually if needed:';
        RAISE NOTICE '========================================';
    END IF;
END $$;

SELECT DISTINCT education_level, COUNT(*) as count
FROM jobs
WHERE education_level IS NOT NULL
AND education_level != ''
AND education_level_id IS NULL
GROUP BY education_level
ORDER BY count DESC;

-- ============================================================================
-- COMPLETED
-- ============================================================================

DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'PRE-MIGRATION DATA TRANSFER COMPLETED!';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Next steps:';
    RAISE NOTICE '1. Review any unmigrated records above';
    RAISE NOTICE '2. Manually fix unmigrated records if needed';
    RAISE NOTICE '3. Run backup: pg_dump -U user dbname > backup.sql';
    RAISE NOTICE '4. Run main migration: 000031_refactor_jobs_table_add_requirements.up.sql';
    RAISE NOTICE '========================================';
END $$;
