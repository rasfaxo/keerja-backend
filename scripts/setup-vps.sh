#!/bin/bash
# ============================================================
# scripts/setup-vps.sh
# ============================================================
# VPS Initial Setup Script for Keerja Backend
#
# PURPOSE:
#   Install and configure all dependencies on a fresh VPS:
#   - Docker & Docker Compose
#   - Nginx
#   - Git
#   - Bump CLI (for documentation)
#   - UFW Firewall
#
# USAGE:
#   chmod +x scripts/setup-vps.sh
#   sudo ./scripts/setup-vps.sh
#
# VPS SPECIFICATIONS:
#   - IP: 145.79.8.227
#   - OS: Ubuntu 22.04 / Debian 11+
#   - RAM: 4GB
#   - Storage: 50GB SSD
#
# AFTER RUNNING THIS SCRIPT:
#   1. Clone repository to /opt/keerja
#   2. Copy .env.staging.example to .env.staging
#   3. Copy .env.demo.example to .env.demo
#   4. Configure nginx sites
#   5. Run deploy-initial.sh
# ============================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

# Check if running as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        log_error "Please run as root (sudo ./setup-vps.sh)"
        exit 1
    fi
}

# Update system packages
update_system() {
    log_info "Updating system packages..."
    apt-get update
    apt-get upgrade -y
    log_success "System packages updated"
}

# Install essential tools
install_essentials() {
    log_info "Installing essential tools..."
    apt-get install -y \
        curl \
        wget \
        git \
        vim \
        htop \
        unzip \
        apt-transport-https \
        ca-certificates \
        gnupg \
        lsb-release \
        software-properties-common
    log_success "Essential tools installed"
}

# Install Docker
install_docker() {
    log_info "Installing Docker..."
    
    # Remove old versions
    apt-get remove -y docker docker-engine docker.io containerd runc 2>/dev/null || true
    
    # Add Docker's official GPG key
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg
    
    # Set up repository
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # Install Docker Engine
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    
    # Start and enable Docker
    systemctl start docker
    systemctl enable docker
    
    # Verify installation
    docker --version
    docker compose version
    
    log_success "Docker installed successfully"
}

# Install Nginx
install_nginx() {
    log_info "Installing Nginx..."
    apt-get install -y nginx
    
    # Start and enable Nginx
    systemctl start nginx
    systemctl enable nginx
    
    # Verify installation
    nginx -v
    
    log_success "Nginx installed successfully"
}

# Install Node.js (for Bump CLI and other tools)
install_nodejs() {
    log_info "Installing Node.js..."
    
    # Install Node.js 20.x LTS
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
    apt-get install -y nodejs
    
    # Verify installation
    node --version
    npm --version
    
    log_success "Node.js installed successfully"
}

# Install Bump CLI
install_bump_cli() {
    log_info "Installing Bump CLI..."
    npm install -g bump-cli
    
    # Verify installation
    bump --version
    
    log_success "Bump CLI installed successfully"
}

# Configure UFW Firewall
configure_firewall() {
    log_info "Configuring UFW Firewall..."
    
    # Install UFW if not present
    apt-get install -y ufw
    
    # Default policies
    ufw default deny incoming
    ufw default allow outgoing
    
    # Allow SSH
    ufw allow ssh
    ufw allow 22/tcp
    
    # Allow HTTP and HTTPS
    ufw allow 80/tcp
    ufw allow 443/tcp
    
    # Allow internal Docker ports (from localhost only)
    # These are not exposed externally, Nginx proxies to them
    
    # Enable firewall
    echo "y" | ufw enable
    
    # Show status
    ufw status verbose
    
    log_success "Firewall configured successfully"
}

# Create application directory structure
create_directories() {
    log_info "Creating application directories..."
    
    # Main application directory
    mkdir -p /opt/keerja
    mkdir -p /opt/keerja/logs
    mkdir -p /opt/keerja/backups
    mkdir -p /opt/keerja/uploads
    mkdir -p /opt/keerja/config
    
    # Nginx directories
    mkdir -p /etc/nginx/sites-available
    mkdir -p /etc/nginx/sites-enabled
    
    # Set permissions
    chown -R root:root /opt/keerja
    chmod -R 755 /opt/keerja
    
    log_success "Directories created successfully"
}

# Configure Nginx
configure_nginx() {
    log_info "Configuring Nginx..."
    
    # Backup default config
    if [ -f /etc/nginx/sites-enabled/default ]; then
        mv /etc/nginx/sites-enabled/default /etc/nginx/sites-available/default.backup
    fi
    
    # Create main nginx config if needed
    cat > /etc/nginx/nginx.conf << 'NGINX_CONF'
user www-data;
worker_processes auto;
pid /run/nginx.pid;
include /etc/nginx/modules-enabled/*.conf;

events {
    worker_connections 1024;
    multi_accept on;
}

http {
    ##
    # Basic Settings
    ##
    sendfile on;
    tcp_nopush on;
    types_hash_max_size 2048;
    server_tokens off;

    # Increase buffer sizes
    client_body_buffer_size 10K;
    client_header_buffer_size 1k;
    client_max_body_size 50M;
    large_client_header_buffers 4 32k;

    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    ##
    # SSL Settings
    ##
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;

    ##
    # Logging Settings
    ##
    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;

    ##
    # Gzip Settings
    ##
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml application/json application/javascript application/rss+xml application/atom+xml image/svg+xml;

    ##
    # Virtual Host Configs
    ##
    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*;
}
NGINX_CONF
    
    log_success "Nginx configured successfully"
}

# Configure swap (for 4GB RAM VPS)
configure_swap() {
    log_info "Configuring swap space..."
    
    # Check if swap already exists
    if [ -f /swapfile ]; then
        log_warning "Swap file already exists, skipping..."
        return
    fi
    
    # Create 2GB swap file
    fallocate -l 2G /swapfile
    chmod 600 /swapfile
    mkswap /swapfile
    swapon /swapfile
    
    # Make swap permanent
    echo '/swapfile none swap sw 0 0' >> /etc/fstab
    
    # Configure swappiness
    sysctl vm.swappiness=10
    echo 'vm.swappiness=10' >> /etc/sysctl.conf
    
    log_success "Swap configured (2GB)"
}

# Configure Docker daemon
configure_docker() {
    log_info "Configuring Docker daemon..."
    
    # Create Docker daemon config for resource management
    mkdir -p /etc/docker
    cat > /etc/docker/daemon.json << 'DOCKER_CONF'
{
    "log-driver": "json-file",
    "log-opts": {
        "max-size": "50m",
        "max-file": "5"
    },
    "storage-driver": "overlay2",
    "live-restore": true,
    "default-ulimits": {
        "nofile": {
            "Name": "nofile",
            "Hard": 65536,
            "Soft": 65536
        }
    }
}
DOCKER_CONF
    
    # Restart Docker to apply changes
    systemctl restart docker
    
    log_success "Docker daemon configured"
}

# Create deploy user (optional, for better security)
create_deploy_user() {
    log_info "Creating deploy user..."
    
    # Check if user already exists
    if id "deploy" &>/dev/null; then
        log_warning "Deploy user already exists, skipping..."
        return
    fi
    
    # Create user
    useradd -m -s /bin/bash deploy
    
    # Add to docker group
    usermod -aG docker deploy
    
    # Set permissions on /opt/keerja
    chown -R deploy:deploy /opt/keerja
    
    log_success "Deploy user created"
    log_info "To use deploy user for SSH, add your SSH key to /home/deploy/.ssh/authorized_keys"
}

# Print summary
print_summary() {
    echo ""
    echo "============================================================"
    echo -e "${GREEN}VPS Setup Complete!${NC}"
    echo "============================================================"
    echo ""
    echo "Installed components:"
    echo "  ✅ Docker $(docker --version | cut -d' ' -f3)"
    echo "  ✅ Docker Compose $(docker compose version | cut -d' ' -f4)"
    echo "  ✅ Nginx $(nginx -v 2>&1 | cut -d'/' -f2)"
    echo "  ✅ Node.js $(node --version)"
    echo "  ✅ Bump CLI $(bump --version 2>/dev/null || echo 'installed')"
    echo "  ✅ UFW Firewall (enabled)"
    echo "  ✅ Swap (2GB)"
    echo ""
    echo "Next steps:"
    echo "  1. Clone repository:"
    echo "     cd /opt/keerja"
    echo "     git clone https://github.com/rasfaxo/keerja-backend.git ."
    echo ""
    echo "  2. Configure environment files:"
    echo "     cp .env.staging.example .env.staging"
    echo "     cp .env.demo.example .env.demo"
    echo "     # Edit both files with actual values"
    echo ""
    echo "  3. Copy nginx configs:"
    echo "     cp nginx/staging.conf /etc/nginx/sites-available/keerja-staging.conf"
    echo "     cp nginx/demo.conf /etc/nginx/sites-available/keerja-demo.conf"
    echo "     ln -s /etc/nginx/sites-available/keerja-staging.conf /etc/nginx/sites-enabled/"
    echo "     ln -s /etc/nginx/sites-available/keerja-demo.conf /etc/nginx/sites-enabled/"
    echo "     nginx -t && systemctl reload nginx"
    echo ""
    echo "  4. Run initial deployment:"
    echo "     chmod +x scripts/deploy-initial.sh"
    echo "     ./scripts/deploy-initial.sh"
    echo ""
    echo "VPS IP: $(curl -s ifconfig.me 2>/dev/null || echo '145.79.8.227')"
    echo ""
    echo "Environment URLs (after deployment):"
    echo "  STAGING: http://staging-api.145.79.8.227.nip.io"
    echo "  DEMO:    http://demo-api.145.79.8.227.nip.io"
    echo ""
}

# Main execution
main() {
    echo ""
    echo "============================================================"
    echo "Keerja Backend - VPS Setup Script"
    echo "============================================================"
    echo ""
    
    check_root
    update_system
    install_essentials
    configure_swap
    install_docker
    configure_docker
    install_nginx
    configure_nginx
    install_nodejs
    install_bump_cli
    configure_firewall
    create_directories
    create_deploy_user
    print_summary
}

# Run main function
main "$@"
