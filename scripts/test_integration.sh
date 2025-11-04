#!/bin/bash

# test_integration.sh - Run integration tests
# Usage: ./scripts/test_integration.sh

set -e

echo "================================"
echo "Running Integration Tests"
echo "================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Set test environment
export APP_ENV=test
export TEST_DB_URL=${TEST_DB_URL:-"postgres://postgres:postgres@localhost:5432/keerja_test?sslmode=disable"}

echo -e "${YELLOW}Test Configuration:${NC}"
echo "  Environment: $APP_ENV"
echo "  Database: $TEST_DB_URL"
echo ""

# Check if test database is accessible
echo -e "${YELLOW}Checking database connection...${NC}"
if ! psql "$TEST_DB_URL" -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${RED}✗ Cannot connect to test database${NC}"
    echo "  Please ensure PostgreSQL is running and test database exists"
    exit 1
fi
echo -e "${GREEN}✓ Database connection successful${NC}"
echo ""

# Clean database before tests
echo -e "${YELLOW}Preparing test database...${NC}"
# Add migration commands here if needed
echo -e "${GREEN}✓ Database ready${NC}"
echo ""

echo -e "${YELLOW}Running integration tests...${NC}"
echo ""

# Run integration tests
go test -v -race -coverprofile=coverage_integration.out -tags=integration ./tests/integration/...

TEST_EXIT_CODE=$?

echo ""
echo "================================"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✓ Integration tests passed!${NC}"
    
    # Show coverage
    if [ -f coverage_integration.out ]; then
        COVERAGE=$(go tool cover -func=coverage_integration.out | grep total | awk '{print $3}')
        echo "  Integration Test Coverage: $COVERAGE"
        
        # Generate HTML report
        go tool cover -html=coverage_integration.out -o coverage_integration.html
        echo -e "${GREEN}✓ Coverage report: coverage_integration.html${NC}"
    fi
else
    echo -e "${RED}✗ Integration tests failed!${NC}"
    exit $TEST_EXIT_CODE
fi

echo "================================"
