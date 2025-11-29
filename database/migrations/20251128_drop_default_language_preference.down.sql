-- Revert change: set default back to 'id' and replace NULLs with 'id'
BEGIN;

-- Set existing NULL language_preference entries back to 'id'
UPDATE user_preferences SET language_preference = 'id' WHERE language_preference IS NULL;

-- Re-add default
ALTER TABLE user_preferences ALTER COLUMN language_preference SET DEFAULT 'id';

COMMIT;
-- Rollback: re-add DEFAULT 'id' to user_preferences.language_preference

ALTER TABLE IF EXISTS user_preferences
    ALTER COLUMN language_preference SET DEFAULT 'id';
