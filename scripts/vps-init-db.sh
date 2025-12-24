#!/bin/bash
# ============================================================
# scripts/vps-init-db.sh
# ============================================================
# PostgreSQL initialization script for VPS multi-environment
#
# PURPOSE:
#   Creates both keerja_staging and keerja_demo databases
#   in a single PostgreSQL container
#
# USAGE:
#   Automatically executed by PostgreSQL on first start
#   Mounted in docker-compose.vps.yml as:
#   ./scripts/vps-init-db.sh:/docker-entrypoint-initdb.d/00-init-databases.sh
#
# INTEGRATION:
#   - Works with docker-compose.vps.yml
#   - Creates isolated databases for STAGING and DEMO
#   - Does NOT run migrations (handled by migrate-staging/migrate-demo services)
# ============================================================

set -e

echo "=== Keerja VPS PostgreSQL Initialization ==="
echo "Timestamp: $(date)"

# Function to create database if not exists
create_database() {
    local database=$1
    echo "Creating database: $database"
    
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        -- Check if database exists
        SELECT 'CREATE DATABASE $database ENCODING ''UTF8'' LC_COLLATE ''C'' LC_CTYPE ''C'' TEMPLATE template0'
        WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$database')\gexec
        
        -- Grant privileges to postgres user
        GRANT ALL PRIVILEGES ON DATABASE $database TO $POSTGRES_USER;
EOSQL

    echo "  - Activating extensions for $database..."
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$database" <<-EOSQL
        CREATE EXTENSION IF NOT EXISTS pgcrypto;
EOSQL
    
    echo "✓ Database $database created/verified"
}

# Function to create application user if not exists
create_app_user() {
    local username=$1
    local password=$2
    echo "Creating application user: $username"
    
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        -- Create user if not exists
        DO \$\$
        BEGIN
            IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$username') THEN
                CREATE USER $username WITH PASSWORD '$password';
                RAISE NOTICE 'User $username created';
            ELSE
                -- Update password if user exists
                ALTER USER $username WITH PASSWORD '$password';
                RAISE NOTICE 'User $username password updated';
            END IF;
        END
        \$\$;
EOSQL
    
    echo "✓ User $username created/verified"
}

# Function to grant privileges on database
grant_privileges() {
    local database=$1
    local username=$2
    echo "Granting privileges on $database to $username"
    
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$database" <<-EOSQL
        -- Grant connect privilege
        GRANT CONNECT ON DATABASE $database TO $username;
        
        -- Grant usage on public schema
        GRANT ALL ON SCHEMA public TO $username;
        
        -- Grant all privileges on all tables (current and future)
        GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO $username;
        GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO $username;
        GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO $username;
        
        -- Default privileges for future objects
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $username;
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO $username;
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO $username;
EOSQL
    
    echo "✓ Privileges granted on $database to $username"
}

# ==========================================
# Main execution
# ==========================================

echo ""
echo "Step 1: Creating databases..."
echo "--------------------------------------------"
create_database "keerja_staging"
create_database "keerja_demo"

echo ""
echo "Step 2: Creating application user..."
echo "--------------------------------------------"
# Use environment variables or defaults
APP_USER="${DB_APP_USER:-bekerja}"
APP_PASSWORD="${DB_APP_PASSWORD:-KeerjaStaging2025!}"
create_app_user "$APP_USER" "$APP_PASSWORD"

echo ""
echo "Step 3: Granting privileges..."
echo "--------------------------------------------"
grant_privileges "keerja_staging" "$APP_USER"
grant_privileges "keerja_demo" "$APP_USER"

echo ""
echo "=== VPS PostgreSQL Initialization Complete ==="
echo ""
echo "Databases created:"
echo "  - keerja_staging (for STAGING environment)"
echo "  - keerja_demo (for DEMO environment)"
echo ""
echo "Application user: $APP_USER"
echo ""
echo "Next steps:"
echo "  1. Run migrations for staging: docker-compose -f docker-compose.vps.yml --profile migrate-staging up"
echo "  2. Run migrations for demo: docker-compose -f docker-compose.vps.yml --profile migrate-demo up"
echo ""
