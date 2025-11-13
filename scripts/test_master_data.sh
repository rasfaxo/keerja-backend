#!/bin/bash

# Master Data API Test Script
# Tests all 10 master data endpoints with success and error cases

BASE_URL="http://localhost:8080/api/v1/meta"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "================================================"
echo "Master Data API Testing"
echo "================================================"
echo ""

# Function to test endpoint
test_endpoint() {
    local name=$1
    local url=$2
    local expected_code=$3
    
    echo -n "Testing: $name ... "
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url")
    
    if [ "$response" == "$expected_code" ]; then
        echo -e "${GREEN}✓ PASS${NC} ($response)"
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (Expected: $expected_code, Got: $response)"
        return 1
    fi
}

# Counter for results
PASSED=0
FAILED=0

echo "================================================"
echo "INDUSTRIES - Testing 5 cases"
echo "================================================"

test_endpoint "Get all industries" "$BASE_URL/industries" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get active industries" "$BASE_URL/industries?active=true" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Search industries" "$BASE_URL/industries?active=true&search=tech" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get industry by ID (valid)" "$BASE_URL/industries/1" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get industry by ID (not found)" "$BASE_URL/industries/999" "404" && ((PASSED++)) || ((FAILED++))

echo ""
echo "================================================"
echo "COMPANY SIZES - Testing 4 cases"
echo "================================================"

test_endpoint "Get all company sizes" "$BASE_URL/company-sizes" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get active company sizes" "$BASE_URL/company-sizes?active=true" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get company size by ID (valid)" "$BASE_URL/company-sizes/3" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get company size by ID (not found)" "$BASE_URL/company-sizes/999" "404" && ((PASSED++)) || ((FAILED++))

echo ""
echo "================================================"
echo "PROVINCES - Testing 5 cases"
echo "================================================"

test_endpoint "Get all provinces" "$BASE_URL/locations/provinces" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get active provinces" "$BASE_URL/locations/provinces?active=true" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Search provinces" "$BASE_URL/locations/provinces?search=jawa" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get province by ID (valid)" "$BASE_URL/locations/provinces/32" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get province by ID (not found)" "$BASE_URL/locations/provinces/999" "404" && ((PASSED++)) || ((FAILED++))

echo ""
echo "================================================"
echo "CITIES - Testing 7 cases"
echo "================================================"

test_endpoint "Get cities by province" "$BASE_URL/locations/cities?province_id=32" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get active cities" "$BASE_URL/locations/cities?province_id=32&active=true" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Search cities" "$BASE_URL/locations/cities?province_id=32&search=bandung" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get city by ID (valid)" "$BASE_URL/locations/cities/150" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get city by ID (not found)" "$BASE_URL/locations/cities/999" "404" && ((PASSED++)) || ((FAILED++))
test_endpoint "Cities without province_id" "$BASE_URL/locations/cities" "400" && ((PASSED++)) || ((FAILED++))
test_endpoint "Cities with invalid province" "$BASE_URL/locations/cities?province_id=999" "404" && ((PASSED++)) || ((FAILED++))

echo ""
echo "================================================"
echo "DISTRICTS - Testing 7 cases"
echo "================================================"

test_endpoint "Get districts by city" "$BASE_URL/locations/districts?city_id=150" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get active districts" "$BASE_URL/locations/districts?city_id=150&active=true" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Search districts" "$BASE_URL/locations/districts?city_id=150&search=batu" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get district by ID (valid)" "$BASE_URL/locations/districts/2045" "200" && ((PASSED++)) || ((FAILED++))
test_endpoint "Get district by ID (not found)" "$BASE_URL/locations/districts/9999" "404" && ((PASSED++)) || ((FAILED++))
test_endpoint "Districts without city_id" "$BASE_URL/locations/districts" "400" && ((PASSED++)) || ((FAILED++))
test_endpoint "Districts with invalid city" "$BASE_URL/locations/districts?city_id=999" "404" && ((PASSED++)) || ((FAILED++))

echo ""
echo "================================================"
echo "TEST SUMMARY"
echo "================================================"
echo -e "${GREEN}✓ Passed:${NC} $PASSED"
echo -e "${RED}✗ Failed:${NC} $FAILED"
echo "Total: $((PASSED + FAILED))"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "\n${RED}❌ SOME TESTS FAILED${NC}"
    exit 1
fi
