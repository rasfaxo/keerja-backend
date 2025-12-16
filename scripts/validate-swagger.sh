#!/bin/bash

# validate-swagger.sh
# Script to validate Swagger/OpenAPI documentation
# Part of the Keerja Backend CI/CD pipeline

set -e

SWAGGER_FILE="docs/swagger.yaml"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting Swagger Validation...${NC}"

# Check if swagger file exists
if [ ! -f "$SWAGGER_FILE" ]; then
    echo -e "${RED}Error: Swagger file not found at $SWAGGER_FILE${NC}"
    exit 1
fi

echo -e "Found swagger file at: $SWAGGER_FILE"

# Check for swagger-cli
if ! command -v swagger-cli &> /dev/null; then
    echo -e "${YELLOW}swagger-cli not found. Installing via npx...${NC}"
    # Try using npx
    if command -v npx &> /dev/null; then
         echo -e "Validating with @apidevtools/swagger-cli..."
         npx @apidevtools/swagger-cli validate "$SWAGGER_FILE"
    else
         echo -e "${RED}Error: npx is not installed. Cannot run swagger-cli.${NC}"
         echo -e "Please install Node.js and npm to run full validation."
         # Fallback: Basic content check
         echo -e "${YELLOW}Performing basic structural check...${NC}"
         if grep -q "swagger: \"2.0\"" "$SWAGGER_FILE"; then
             echo -e "${GREEN}Basic check passed: File contains 'swagger: \"2.0\"'${NC}"
         else
             echo -e "${RED}Basic check failed: File does not look like OpenAPI 2.0${NC}"
             exit 1
         fi
         exit 0
    fi
else
    echo -e "Validating with installed swagger-cli..."
    swagger-cli validate "$SWAGGER_FILE"
fi

echo -e "${GREEN}✅ Swagger validation passed successfully!${NC}"
exit 0
