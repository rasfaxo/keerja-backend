#!/bin/bash
# ============================================================
# scripts/health-check.sh
# ============================================================
# Health Check Script for Keerja Backend VPS
#
# PURPOSE:
#   Check the health status of all services running on VPS:
#   - PostgreSQL
#   - Redis
#   - STAGING API
#   - DEMO API
#   - Nginx
#
# USAGE:
#   chmod +x scripts/health-check.sh
#   ./scripts/health-check.sh
#
# EXIT CODES:
#   0 - All services healthy
#   1 - One or more services unhealthy
# ============================================================

set -e

# Configuration
DEPLOY_PATH="/opt/keerja"
COMPOSE_FILE="docker-compose.vps.yml"
VPS_IP="${VPS_IP:-145.79.8.227}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# Functions
check_service() {
    local name="$1"
    local check_cmd="$2"
    local expected="$3"
    
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    
    printf "%-30s" "$name"
    
    if eval "$check_cmd" &>/dev/null; then
        echo -e "${GREEN}✓ OK${NC}"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

check_http() {
    local name="$1"
    local url="$2"
    local expected_code="${3:-200}"
    
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    
    printf "%-30s" "$name"
    
    local response=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 5 "$url" 2>/dev/null || echo "000")
    
    if [ "$response" = "$expected_code" ]; then
        echo -e "${GREEN}✓ OK${NC} (HTTP $response)"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    else
        echo -e "${RED}✗ FAILED${NC} (HTTP $response, expected $expected_code)"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

check_container() {
    local name="$1"
    local container="$2"
    
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    
    printf "%-30s" "$name"
    
    local status=$(docker inspect -f '{{.State.Status}}' "$container" 2>/dev/null || echo "not found")
    local health=$(docker inspect -f '{{.State.Health.Status}}' "$container" 2>/dev/null || echo "none")
    
    if [ "$status" = "running" ]; then
        if [ "$health" = "healthy" ] || [ "$health" = "none" ]; then
            echo -e "${GREEN}✓ OK${NC} (running, health: $health)"
            PASSED_CHECKS=$((PASSED_CHECKS + 1))
            return 0
        else
            echo -e "${YELLOW}⚠ WARNING${NC} (running, health: $health)"
            PASSED_CHECKS=$((PASSED_CHECKS + 1))
            return 0
        fi
    else
        echo -e "${RED}✗ FAILED${NC} (status: $status)"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

# Main checks
main() {
    echo ""
    echo "============================================================"
    echo "Keerja Backend - Health Check"
    echo "Timestamp: $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
    echo "============================================================"
    echo ""
    
    # Docker checks
    echo -e "${BLUE}Docker Services${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    check_service "Docker Engine" "docker info"
    check_container "PostgreSQL" "keerja-vps-postgres"
    check_container "Redis" "keerja-vps-redis"
    check_container "API Staging" "keerja-api-staging"
    check_container "API Demo" "keerja-api-demo"
    echo ""
    
    # Database checks
    echo -e "${BLUE}Database Connectivity${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    check_service "PostgreSQL Ready" "docker exec keerja-vps-postgres pg_isready -U postgres"
    check_service "Redis Ping" "docker exec keerja-vps-redis redis-cli ping"
    check_service "Staging DB Exists" "docker exec keerja-vps-postgres psql -U postgres -d keerja_staging -c 'SELECT 1'"
    check_service "Demo DB Exists" "docker exec keerja-vps-postgres psql -U postgres -d keerja_demo -c 'SELECT 1'"
    echo ""
    
    # HTTP endpoints (internal)
    echo -e "${BLUE}Internal HTTP Endpoints${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    check_http "Staging Health (internal)" "http://localhost:8080/health/live"
    check_http "Demo Health (internal)" "http://localhost:8081/health/live"
    echo ""
    
    # HTTP endpoints (via Nginx)
    echo -e "${BLUE}External HTTP Endpoints (via Nginx)${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    check_http "Staging (nip.io)" "http://staging-api.${VPS_IP}.nip.io/health/live"
    check_http "Demo (nip.io)" "http://demo-api.${VPS_IP}.nip.io/health/live"
    echo ""
    
    # Nginx checks
    echo -e "${BLUE}Nginx Status${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    check_service "Nginx Running" "systemctl is-active nginx"
    check_service "Nginx Config Valid" "nginx -t"
    echo ""
    
    # Disk and memory
    echo -e "${BLUE}System Resources${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Disk Usage:"
    df -h / | tail -1 | awk '{print "  Root: " $5 " used (" $3 "/" $2 ")"}'
    df -h /var/lib/docker 2>/dev/null | tail -1 | awk '{print "  Docker: " $5 " used (" $3 "/" $2 ")"}' || true
    echo ""
    echo "Memory Usage:"
    free -h | awk '/^Mem:/ {print "  RAM: " $3 "/" $2 " used"}'
    free -h | awk '/^Swap:/ {print "  Swap: " $3 "/" $2 " used"}'
    echo ""
    
    # Container resource usage
    echo -e "${BLUE}Container Resource Usage${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" 2>/dev/null | head -10 || echo "Unable to get container stats"
    echo ""
    
    # Summary
    echo "============================================================"
    echo -e "Health Check Summary"
    echo "============================================================"
    echo ""
    echo -e "Total checks: $TOTAL_CHECKS"
    echo -e "Passed: ${GREEN}$PASSED_CHECKS${NC}"
    echo -e "Failed: ${RED}$FAILED_CHECKS${NC}"
    echo ""
    
    if [ $FAILED_CHECKS -eq 0 ]; then
        echo -e "${GREEN}All systems operational! ✓${NC}"
        exit 0
    else
        echo -e "${RED}Some checks failed. Please investigate.${NC}"
        exit 1
    fi
}

# Run main
main "$@"
