#!/bin/bash
set -e

# Environment variables dengan default
MIGRATOR_USER="${POSTGRES_MIGRATOR_USER:-kustan}"
MIGRATOR_PASS="${POSTGRES_MIGRATOR_PASSWORD:-kustan_dev_pass_123}"
APP_USER="${POSTGRES_APP_USER:-bekerja}"
APP_PASS="${POSTGRES_APP_PASSWORD:-bekerja_dev_pass_123}"
DB="${POSTGRES_DB:-keerja}"

echo "Setting up database roles and permissions..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$DB" <<-EOSQL
-- Create migrator role (idempotent)
DO \$\$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '$MIGRATOR_USER') THEN
    CREATE ROLE "$MIGRATOR_USER" LOGIN PASSWORD '$MIGRATOR_PASS';
    GRANT CONNECT ON DATABASE "$DB" TO "$MIGRATOR_USER";
    GRANT USAGE, CREATE ON SCHEMA public TO "$MIGRATOR_USER";
  END IF;
END
\$\$;

-- Create app role (idempotent)
DO \$\$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '$APP_USER') THEN
    CREATE ROLE "$APP_USER" LOGIN PASSWORD '$APP_PASS';
    GRANT CONNECT ON DATABASE "$DB" TO "$APP_USER";
    GRANT USAGE ON SCHEMA public TO "$APP_USER";
  END IF;
END
\$\$;

-- Grant DML permissions to app on current tables
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO "$APP_USER";
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO "$APP_USER";

-- Set default privileges for future tables/sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO "$APP_USER";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO "$APP_USER";

-- Allow migrator to manage sequences
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO "$MIGRATOR_USER";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO "$MIGRATOR_USER";
EOSQL

echo "✓ Database roles and permissions setup completed"
