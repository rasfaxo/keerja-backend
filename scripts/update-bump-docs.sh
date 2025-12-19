#!/bin/bash
# ============================================================
# scripts/update-bump-docs.sh
# ============================================================
# Manual Bump.sh Documentation Update Script
#
# PURPOSE:
#   Manually update Bump.sh API documentation for
#   STAGING and/or DEMO environments.
#
# PREREQUISITES:
#   - Bump CLI installed (npm install -g bump-cli)
#   - BUMP_TOKEN environment variable set
#   - STAGING_BUMP_DOC_ID or DEMO_BUMP_DOC_ID set
#
# USAGE:
#   # Update staging docs
#   BUMP_TOKEN=xxx STAGING_BUMP_DOC_ID=xxx ./scripts/update-bump-docs.sh staging
#   
#   # Update demo docs
#   BUMP_TOKEN=xxx DEMO_BUMP_DOC_ID=xxx ./scripts/update-bump-docs.sh demo
#   
#   # Update both
#   BUMP_TOKEN=xxx STAGING_BUMP_DOC_ID=xxx DEMO_BUMP_DOC_ID=xxx ./scripts/update-bump-docs.sh all
#
# ENVIRONMENT VARIABLES:
#   BUMP_TOKEN           - Bump.sh API token (required)
#   STAGING_BUMP_DOC_ID  - Documentation ID for staging
#   DEMO_BUMP_DOC_ID     - Documentation ID for demo
#   VPS_IP               - VPS IP address (default: 145.79.8.227)
# ============================================================

set -e

# Configuration
VPS_IP="${VPS_IP:-145.79.8.227}"
SWAGGER_FILE="docs/swagger.yaml"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check bump CLI
    if ! command -v bump &> /dev/null; then
        log_error "Bump CLI not found. Install with: npm install -g bump-cli"
        exit 1
    fi
    
    # Check BUMP_TOKEN
    if [ -z "$BUMP_TOKEN" ]; then
        log_error "BUMP_TOKEN environment variable is not set"
        log_info "Get your token from: https://bump.sh/docs/bump-cli"
        exit 1
    fi
    
    # Check swagger file
    if [ ! -f "$PROJECT_ROOT/$SWAGGER_FILE" ]; then
        log_error "Swagger file not found: $PROJECT_ROOT/$SWAGGER_FILE"
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

# Validate swagger file
validate_swagger() {
    log_info "Validating swagger file..."
    
    if command -v swagger-cli &> /dev/null; then
        swagger-cli validate "$PROJECT_ROOT/$SWAGGER_FILE"
        log_success "Swagger validation passed"
    else
        log_warning "swagger-cli not installed, skipping validation"
    fi
}

# Prepare swagger file for environment
prepare_swagger() {
    local env="$1"
    local output_file="$2"
    
    log_info "Preparing swagger file for $env..."
    
    cp "$PROJECT_ROOT/$SWAGGER_FILE" "$output_file"
    
    case "$env" in
        staging)
            sed -i "s|url: .*|url: http://staging-api.${VPS_IP}.nip.io|g" "$output_file"
            sed -i "s|title: .*|title: Keerja API (Staging)|g" "$output_file"
            ;;
        demo)
            sed -i "s|url: .*|url: http://demo-api.${VPS_IP}.nip.io|g" "$output_file"
            sed -i "s|title: .*|title: Keerja API (Demo)|g" "$output_file"
            ;;
    esac
    
    log_success "Swagger file prepared for $env"
}

# Deploy to Bump.sh
deploy_to_bump() {
    local env="$1"
    local swagger_file="$2"
    local doc_id="$3"
    
    if [ -z "$doc_id" ]; then
        log_error "Documentation ID not set for $env"
        log_info "Set ${env^^}_BUMP_DOC_ID environment variable"
        return 1
    fi
    
    log_info "Deploying $env documentation to Bump.sh..."
    
    bump deploy "$swagger_file" \
        --doc "$doc_id" \
        --token "$BUMP_TOKEN"
    
    log_success "$env documentation deployed to Bump.sh"
    log_info "View at: https://bump.sh/doc/$doc_id"
}

# Update staging docs
update_staging() {
    local temp_file=$(mktemp)
    
    log_info "Updating STAGING documentation..."
    
    prepare_swagger "staging" "$temp_file"
    deploy_to_bump "staging" "$temp_file" "$STAGING_BUMP_DOC_ID"
    
    rm -f "$temp_file"
}

# Update demo docs
update_demo() {
    local temp_file=$(mktemp)
    
    log_info "Updating DEMO documentation..."
    
    prepare_swagger "demo" "$temp_file"
    deploy_to_bump "demo" "$temp_file" "$DEMO_BUMP_DOC_ID"
    
    rm -f "$temp_file"
}

# Preview changes (dry run)
preview_changes() {
    local env="$1"
    local doc_id="$2"
    local temp_file=$(mktemp)
    
    log_info "Previewing changes for $env..."
    
    prepare_swagger "$env" "$temp_file"
    
    bump diff "$temp_file" \
        --doc "$doc_id" \
        --token "$BUMP_TOKEN" || true
    
    rm -f "$temp_file"
}

# Print usage
print_usage() {
    echo "Usage: $0 [staging|demo|all|preview-staging|preview-demo]"
    echo ""
    echo "Commands:"
    echo "  staging          Update STAGING documentation"
    echo "  demo             Update DEMO documentation"
    echo "  all              Update both STAGING and DEMO"
    echo "  preview-staging  Preview changes for STAGING (dry run)"
    echo "  preview-demo     Preview changes for DEMO (dry run)"
    echo ""
    echo "Environment Variables:"
    echo "  BUMP_TOKEN           Bump.sh API token (required)"
    echo "  STAGING_BUMP_DOC_ID  Documentation ID for staging"
    echo "  DEMO_BUMP_DOC_ID     Documentation ID for demo"
    echo "  VPS_IP               VPS IP address (default: 145.79.8.227)"
}

# Main execution
main() {
    local command="${1:-}"
    
    echo ""
    echo "============================================================"
    echo "Keerja API - Bump.sh Documentation Update"
    echo "============================================================"
    echo ""
    
    if [ -z "$command" ]; then
        print_usage
        exit 1
    fi
    
    check_prerequisites
    validate_swagger
    
    case "$command" in
        staging)
            update_staging
            ;;
        demo)
            update_demo
            ;;
        all)
            update_staging
            echo ""
            update_demo
            ;;
        preview-staging)
            preview_changes "staging" "$STAGING_BUMP_DOC_ID"
            ;;
        preview-demo)
            preview_changes "demo" "$DEMO_BUMP_DOC_ID"
            ;;
        *)
            log_error "Unknown command: $command"
            print_usage
            exit 1
            ;;
    esac
    
    echo ""
    log_success "Done!"
}

# Run main
main "$@"
