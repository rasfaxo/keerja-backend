#!/bin/bash

# test.sh - Run all tests with coverage
# Usage: ./scripts/test.sh

set -e

echo "================================"
echo "Running Keerja Backend Tests"
echo "================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

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
    echo "  Database URL: $TEST_DB_URL"
    exit 1
fi
echo -e "${GREEN}✓ Database connection successful${NC}"
echo ""

# Run tests
echo -e "${YELLOW}Running tests...${NC}"
echo ""

# Test with race detector and coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./... 2>&1 | tee test_output.log

# Check test result
TEST_EXIT_CODE=${PIPESTATUS[0]}

echo ""
echo "================================"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
else
    echo -e "${RED}✗ Tests failed!${NC}"
    exit $TEST_EXIT_CODE
fi

echo "================================"
echo ""

# Generate coverage report
if [ -f coverage.out ]; then
    echo -e "${YELLOW}Generating coverage report...${NC}"
    
    # Calculate coverage percentage
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    
    echo "  Total Coverage: $COVERAGE"
    
    # Generate HTML report
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"
    
    # Check coverage threshold
    COVERAGE_NUM=$(echo $COVERAGE | sed 's/%//')
    THRESHOLD=80
    
    if (( $(echo "$COVERAGE_NUM >= $THRESHOLD" | bc -l) )); then
        echo -e "${GREEN}✓ Coverage meets threshold (>= ${THRESHOLD}%)${NC}"
    else
        echo -e "${YELLOW}⚠ Coverage below threshold (${THRESHOLD}%): ${COVERAGE}${NC}"
    fi
else
    echo -e "${YELLOW}⚠ No coverage report generated${NC}"
fi

echo ""
echo "================================"
echo -e "${GREEN}Test run completed successfully!${NC}"
echo "================================"
