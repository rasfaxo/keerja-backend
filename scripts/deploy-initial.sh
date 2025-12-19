#!/bin/bash
# ============================================================
# scripts/deploy-initial.sh
# ============================================================
# Initial Deployment Script for Keerja Backend on VPS
#
# PURPOSE:
#   First-time deployment of STAGING and DEMO environments
#   after VPS setup is complete.
#
# PREREQUISITES:
#   - setup-vps.sh has been run
#   - Repository cloned to /opt/keerja
#   - .env.staging and .env.demo configured
#   - Nginx configs installed
#
# USAGE:
#   cd /opt/keerja
#   chmod +x scripts/deploy-initial.sh
#   ./scripts/deploy-initial.sh [staging|demo|all]
#
# EXAMPLES:
#   ./scripts/deploy-initial.sh           # Deploy both
#   ./scripts/deploy-initial.sh staging   # Deploy staging only
#   ./scripts/deploy-initial.sh demo      # Deploy demo only
# ============================================================

set -e

# Configuration
DEPLOY_PATH="/opt/keerja"
COMPOSE_FILE="docker-compose.vps.yml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

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

log_step() {
    echo -e "\n${CYAN}━━━ $1 ━━━${NC}\n"
}

# Check prerequisites
check_prerequisites() {
    log_step "Checking Prerequisites"
    
    # Check if we're in the right directory
    if [ ! -f "$DEPLOY_PATH/$COMPOSE_FILE" ]; then
        log_error "docker-compose.vps.yml not found in $DEPLOY_PATH"
        log_error "Please ensure repository is cloned to $DEPLOY_PATH"
        exit 1
    fi
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Run setup-vps.sh first."
        exit 1
    fi
    
    # Check Docker Compose
    if ! docker compose version &> /dev/null; then
        log_error "Docker Compose is not installed. Run setup-vps.sh first."
        exit 1
    fi
    
    # Check Nginx
    if ! command -v nginx &> /dev/null; then
        log_error "Nginx is not installed. Run setup-vps.sh first."
        exit 1
    fi
    
    log_success "All prerequisites met"
}

# Check environment files
check_env_files() {
    log_step "Checking Environment Files"
    
    local missing_files=0
    
    if [ "$1" = "staging" ] || [ "$1" = "all" ]; then
        if [ ! -f "$DEPLOY_PATH/.env.staging" ]; then
            log_error ".env.staging not found!"
            log_info "Create it: cp .env.staging.example .env.staging"
            missing_files=1
        else
            log_success ".env.staging found"
        fi
    fi
    
    if [ "$1" = "demo" ] || [ "$1" = "all" ]; then
        if [ ! -f "$DEPLOY_PATH/.env.demo" ]; then
            log_error ".env.demo not found!"
            log_info "Create it: cp .env.demo.example .env.demo"
            missing_files=1
        else
            log_success ".env.demo found"
        fi
    fi
    
    if [ $missing_files -eq 1 ]; then
        log_error "Please create missing environment files and try again"
        exit 1
    fi
}

# Load environment variables
load_env() {
    log_info "Loading environment variables..."
    
    # Source common variables from .env.staging (for shared infra)
    if [ -f "$DEPLOY_PATH/.env.staging" ]; then
        set -a
        source "$DEPLOY_PATH/.env.staging"
        set +a
    fi
}

# Start infrastructure (PostgreSQL + Redis)
start_infrastructure() {
    log_step "Starting Infrastructure Services"
    
    cd "$DEPLOY_PATH"
    
    log_info "Starting PostgreSQL and Redis..."
    docker compose -f "$COMPOSE_FILE" up -d postgres redis
    
    log_info "Waiting for PostgreSQL to be ready..."
    sleep 15
    
    # Check PostgreSQL health
    for i in {1..30}; do
        if docker compose -f "$COMPOSE_FILE" exec -T postgres pg_isready -U postgres &>/dev/null; then
            log_success "PostgreSQL is ready"
            break
        fi
        echo -n "."
        sleep 2
    done
    
    # Check Redis health
    if docker compose -f "$COMPOSE_FILE" exec -T redis redis-cli ping &>/dev/null; then
        log_success "Redis is ready"
    else
        log_warning "Redis may not be ready yet"
    fi
}

# Build images
build_images() {
    log_step "Building Docker Images"
    
    cd "$DEPLOY_PATH"
    
    if [ "$1" = "staging" ] || [ "$1" = "all" ]; then
        log_info "Building STAGING image..."
        docker compose -f "$COMPOSE_FILE" build api-staging
        log_success "STAGING image built"
    fi
    
    if [ "$1" = "demo" ] || [ "$1" = "all" ]; then
        log_info "Building DEMO image..."
        docker compose -f "$COMPOSE_FILE" build api-demo
        log_success "DEMO image built"
    fi
}

# Run migrations
run_migrations() {
    log_step "Running Database Migrations"
    
    cd "$DEPLOY_PATH"
    
    if [ "$1" = "staging" ] || [ "$1" = "all" ]; then
        log_info "Running migrations for STAGING (keerja_staging)..."
        docker compose -f "$COMPOSE_FILE" --profile migrate-staging up --build
        log_success "STAGING migrations complete"
    fi
    
    if [ "$1" = "demo" ] || [ "$1" = "all" ]; then
        log_info "Running migrations for DEMO (keerja_demo)..."
        docker compose -f "$COMPOSE_FILE" --profile migrate-demo up --build
        log_success "DEMO migrations complete"
    fi
}

# Deploy applications
deploy_applications() {
    log_step "Deploying Applications"
    
    cd "$DEPLOY_PATH"
    
    if [ "$1" = "staging" ] || [ "$1" = "all" ]; then
        log_info "Deploying STAGING..."
        docker compose -f "$COMPOSE_FILE" --profile staging up -d api-staging
        log_success "STAGING deployed"
    fi
    
    if [ "$1" = "demo" ] || [ "$1" = "all" ]; then
        log_info "Deploying DEMO..."
        docker compose -f "$COMPOSE_FILE" --profile demo up -d api-demo
        log_success "DEMO deployed"
    fi
}

# Configure Nginx
configure_nginx() {
    log_step "Configuring Nginx"
    
    # Copy nginx configs if they exist in repo
    if [ -f "$DEPLOY_PATH/nginx/staging.conf" ]; then
        log_info "Installing STAGING nginx config..."
        cp "$DEPLOY_PATH/nginx/staging.conf" /etc/nginx/sites-available/keerja-staging.conf
        ln -sf /etc/nginx/sites-available/keerja-staging.conf /etc/nginx/sites-enabled/
        log_success "STAGING nginx config installed"
    fi
    
    if [ -f "$DEPLOY_PATH/nginx/demo.conf" ]; then
        log_info "Installing DEMO nginx config..."
        cp "$DEPLOY_PATH/nginx/demo.conf" /etc/nginx/sites-available/keerja-demo.conf
        ln -sf /etc/nginx/sites-available/keerja-demo.conf /etc/nginx/sites-enabled/
        log_success "DEMO nginx config installed"
    fi
    
    # Test nginx configuration
    log_info "Testing nginx configuration..."
    if nginx -t; then
        log_success "Nginx configuration is valid"
        systemctl reload nginx
        log_success "Nginx reloaded"
    else
        log_error "Nginx configuration is invalid!"
        exit 1
    fi
}

# Health check
health_check() {
    log_step "Running Health Checks"
    
    sleep 10  # Wait for services to start
    
    if [ "$1" = "staging" ] || [ "$1" = "all" ]; then
        log_info "Checking STAGING health..."
        for i in {1..10}; do
            if curl -sf http://localhost:8080/health/live &>/dev/null; then
                log_success "STAGING is healthy ✓"
                break
            fi
            echo -n "."
            sleep 3
        done
    fi
    
    if [ "$1" = "demo" ] || [ "$1" = "all" ]; then
        log_info "Checking DEMO health..."
        for i in {1..10}; do
            if curl -sf http://localhost:8081/health/live &>/dev/null; then
                log_success "DEMO is healthy ✓"
                break
            fi
            echo -n "."
            sleep 3
        done
    fi
}

# Print summary
print_summary() {
    local vps_ip=$(curl -s ifconfig.me 2>/dev/null || echo "145.79.8.227")
    
    echo ""
    echo "============================================================"
    echo -e "${GREEN}Initial Deployment Complete!${NC}"
    echo "============================================================"
    echo ""
    echo "Running containers:"
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""
    
    if [ "$1" = "staging" ] || [ "$1" = "all" ]; then
        echo "STAGING Environment:"
        echo "  URL: http://staging-api.${vps_ip}.nip.io"
        echo "  Direct: http://${vps_ip}:8080"
        echo "  Health: http://staging-api.${vps_ip}.nip.io/health/live"
        echo ""
    fi
    
    if [ "$1" = "demo" ] || [ "$1" = "all" ]; then
        echo "DEMO Environment:"
        echo "  URL: http://demo-api.${vps_ip}.nip.io"
        echo "  Direct: http://${vps_ip}:8081"
        echo "  Health: http://demo-api.${vps_ip}.nip.io/health/live"
        echo ""
    fi
    
    echo "Useful commands:"
    echo "  View logs:    docker compose -f $COMPOSE_FILE logs -f api-staging"
    echo "  Restart:      docker compose -f $COMPOSE_FILE --profile staging restart"
    echo "  Stop all:     docker compose -f $COMPOSE_FILE down"
    echo ""
}

# Main execution
main() {
    local deploy_target="${1:-all}"
    
    echo ""
    echo "============================================================"
    echo "Keerja Backend - Initial Deployment"
    echo "Target: $deploy_target"
    echo "============================================================"
    echo ""
    
    # Validate target
    if [[ ! "$deploy_target" =~ ^(staging|demo|all)$ ]]; then
        log_error "Invalid target: $deploy_target"
        log_info "Usage: $0 [staging|demo|all]"
        exit 1
    fi
    
    cd "$DEPLOY_PATH"
    
    check_prerequisites
    check_env_files "$deploy_target"
    load_env
    start_infrastructure
    build_images "$deploy_target"
    run_migrations "$deploy_target"
    deploy_applications "$deploy_target"
    configure_nginx
    health_check "$deploy_target"
    print_summary "$deploy_target"
}

# Run main function
main "$@"
