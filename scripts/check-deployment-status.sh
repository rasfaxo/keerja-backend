#!/bin/bash

# check-deployment-status.sh
# Script to check the latest deployment status on Bump.sh

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' 

# Check for API token
if [ -z "$BUMP_TOKEN" ]; then
    echo -e "${RED}Error: BUMP_TOKEN environment variable is not set.${NC}"
    echo "Please set it with: export BUMP_TOKEN=your_token"
    exit 1
fi

# Configuration
DOC_ID="${BUMP_DOC_ID}" 
API_URL="https://bump.sh/api/v1"

echo -e "${YELLOW}Checking Bump.sh deployment status...${NC}"

# If DOC_ID is set, check specific documentation
if [ ! -z "$DOC_ID" ]; then
    echo "Checking documentation ID: $DOC_ID"
    RESPONSE=$(curl -s -H "Authorization: Token $BUMP_TOKEN" "$API_URL/docs/$DOC_ID")
else
    # List docs
    echo "Fetching all documentations..."
    RESPONSE=$(curl -s -H "Authorization: Token $BUMP_TOKEN" "$API_URL/docs")
fi

# Check if response contains error
if echo "$RESPONSE" | grep -q "error"; then
    echo -e "${RED}Error fetching status:${NC}"
    echo "$RESPONSE"
    exit 1
fi

echo -e "${GREEN}Successfully retrieved status!${NC}"
echo "----------------------------------------"
echo "$RESPONSE" # In a real script, we would parse JSON with jq
echo "----------------------------------------"
echo -e "${YELLOW}Tip: Use 'jq' to parse the JSON output for better readability.${NC}"
