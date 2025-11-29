-- Drop default for language_preference on user_preferences and set existing 'id' values to NULL
BEGIN;

-- Ensure we don't error if column already has no default
ALTER TABLE user_preferences ALTER COLUMN language_preference DROP DEFAULT;

-- Convert legacy 'id' default entries to NULL so semantics are consistent
UPDATE user_preferences SET language_preference = NULL WHERE language_preference = 'id';

COMMIT;
-- Migration: remove DEFAULT on user_preferences.language_preference
-- This will make new rows default to NULL when language_preference is not provided

ALTER TABLE IF EXISTS user_preferences
    ALTER COLUMN language_preference DROP DEFAULT;
