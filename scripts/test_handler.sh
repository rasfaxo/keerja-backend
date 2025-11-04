#!/bin/bash

# test_handler.sh - Run handler layer tests
# Usage: ./scripts/test_handler.sh

set -e

echo "================================"
echo "Running Handler Layer Tests"
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

echo -e "${YELLOW}Running handler tests...${NC}"
echo ""

# Run handler tests
go test -v -race -coverprofile=coverage_handler.out ./internal/handler/... ./tests/handler/...

TEST_EXIT_CODE=$?

echo ""
echo "================================"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✓ Handler tests passed!${NC}"
    
    # Show coverage
    if [ -f coverage_handler.out ]; then
        COVERAGE=$(go tool cover -func=coverage_handler.out | grep total | awk '{print $3}')
        echo "  Handler Layer Coverage: $COVERAGE"
        
        # Generate HTML report
        go tool cover -html=coverage_handler.out -o coverage_handler.html
        echo -e "${GREEN}✓ Coverage report: coverage_handler.html${NC}"
    fi
else
    echo -e "${RED}✗ Handler tests failed!${NC}"
    exit $TEST_EXIT_CODE
fi

echo "================================"
