#!/bin/bash

# test_service.sh - Run service layer tests
# Usage: ./scripts/test_service.sh

set -e

echo "================================"
echo "Running Service Layer Tests"
echo "================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Set test environment
export APP_ENV=test
export TEST_DB_URL=${TEST_DB_URL:-"postgres://postgres:postgres@localhost:5432/keerja_test?sslmode=disable"}

echo -e "${YELLOW}Running service tests...${NC}"
echo ""

# Run service tests
go test -v -race -coverprofile=coverage_service.out ./internal/service/... ./tests/service/...

TEST_EXIT_CODE=$?

echo ""
echo "================================"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✓ Service tests passed!${NC}"
    
    # Show coverage
    if [ -f coverage_service.out ]; then
        COVERAGE=$(go tool cover -func=coverage_service.out | grep total | awk '{print $3}')
        echo "  Service Layer Coverage: $COVERAGE"
        
        # Generate HTML report
        go tool cover -html=coverage_service.out -o coverage_service.html
        echo -e "${GREEN}✓ Coverage report: coverage_service.html${NC}"
    fi
else
    echo -e "${RED}✗ Service tests failed!${NC}"
    exit $TEST_EXIT_CODE
fi

echo "================================"
