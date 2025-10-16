#!/bin/bash
set -e

# Create additional users if needed
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create kustan user if referenced in schema
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'kustan') THEN
            CREATE USER kustan;
        END IF;
    END
    \$\$;
    
    -- Grant necessary permissions
    GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_DB TO bekerja;
    GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_DB TO kustan;
EOSQL

echo "Additional users created successfully"
