#!/bin/bash

# Configuration
API_URL="http://localhost:8080/api/v1"
ADMIN_EMAIL="admin@keerja.com"  # Replace with your admin user
ADMIN_PASSWORD="admin123"       # Replace with actual password

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Function to get auth token
get_auth_token() {
    local response=$(curl -s -X POST "${API_URL}/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"${ADMIN_EMAIL}\",\"password\":\"${ADMIN_PASSWORD}\"}")
    echo $response | jq -r '.data.token'
}

# Function to test endpoints
test_endpoint() {
    local method=$1
    local endpoint=$2
    local payload=$3
    local token=$4
    local description=$5

    echo -e "\n${GREEN}Testing ${description}${NC}"
    echo "Method: ${method}"
    echo "Endpoint: ${endpoint}"
    
    if [ ! -z "$payload" ]; then
        echo "Payload: ${payload}"
    fi

    local response
    if [ -z "$payload" ]; then
        response=$(curl -s -X "${method}" \
            -H "Authorization: Bearer ${token}" \
            "${API_URL}${endpoint}")
    else
        response=$(curl -s -X "${method}" \
            -H "Authorization: Bearer ${token}" \
            -H "Content-Type: application/json" \
            -d "${payload}" \
            "${API_URL}${endpoint}")
    fi

    echo "Response:"
    echo "${response}" | jq '.'
    echo -e "${GREEN}Test completed${NC}\n"
}

# Get auth token
echo "Getting authentication token..."
TOKEN=$(get_auth_token)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get authentication token${NC}"
    exit 1
fi

echo -e "${GREEN}Successfully got authentication token${NC}"

# Test Province endpoints
echo -e "\n${GREEN}Testing Province Endpoints${NC}"
test_endpoint "POST" "/admin/master/provinces" \
    '{"code":"TEST","name":"Test Province"}' \
    "$TOKEN" \
    "Create Province"

test_endpoint "GET" "/admin/master/provinces" "" "$TOKEN" \
    "List Provinces"

# Test City endpoints
echo -e "\n${GREEN}Testing City Endpoints${NC}"
test_endpoint "POST" "/admin/master/cities" \
    '{"name":"Test City","provinceId":1}' \
    "$TOKEN" \
    "Create City"

test_endpoint "GET" "/admin/master/cities?province_id=1" "" "$TOKEN" \
    "List Cities"

# Test District endpoints
echo -e "\n${GREEN}Testing District Endpoints${NC}"
test_endpoint "POST" "/admin/master/districts" \
    '{"name":"Test District","cityId":1}' \
    "$TOKEN" \
    "Create District"

test_endpoint "GET" "/admin/master/districts?city_id=1" "" "$TOKEN" \
    "List Districts"

# Test Industry endpoints
echo -e "\n${GREEN}Testing Industry Endpoints${NC}"
test_endpoint "POST" "/admin/master/industries" \
    '{"name":"Test Industry"}' \
    "$TOKEN" \
    "Create Industry"

test_endpoint "GET" "/admin/master/industries" "" "$TOKEN" \
    "List Industries"

# Test Job Type endpoints
echo -e "\n${GREEN}Testing Job Type Endpoints${NC}"
test_endpoint "POST" "/admin/master/job-types" \
    '{"code":"TEST","name":"Test Job Type","order":1}' \
    "$TOKEN" \
    "Create Job Type"

test_endpoint "GET" "/admin/master/job-types" "" "$TOKEN" \
    "List Job Types"

# Test Company Size endpoints
echo -e "\n${GREEN}Testing Company Size Endpoints${NC}"
test_endpoint "POST" "/admin/meta/company-sizes" \
    '{"label":"Test Size","minEmployees":1,"maxEmployees":10}' \
    "$TOKEN" \
    "Create Company Size"

test_endpoint "GET" "/admin/meta/company-sizes" "" "$TOKEN" \
    "List Company Sizes"

echo -e "\n${GREEN}All tests completed${NC}"