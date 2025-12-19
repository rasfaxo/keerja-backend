# Keerja Backend - Deployment Guide

This guide provides step-by-step instructions for deploying the Keerja Backend to a VPS with multi-environment setup (STAGING + DEMO).

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [VPS Initial Setup](#vps-initial-setup)
4. [Repository Setup](#repository-setup)
5. [Environment Configuration](#environment-configuration)
6. [Nginx Configuration](#nginx-configuration)
7. [Initial Deployment](#initial-deployment)
8. [GitHub Actions Setup](#github-actions-setup)
9. [Bump.sh Documentation Setup](#bumpsh-documentation-setup)
10. [Maintenance Commands](#maintenance-commands)
11. [Troubleshooting](#troubleshooting)
12. [Rollback Procedures](#rollback-procedures)

---

## Overview

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        NGINX                                 │
│   staging-api.145.79.8.227.nip.io → localhost:8080          │
│   demo-api.145.79.8.227.nip.io    → localhost:8081          │
└─────────────────────────────────────────────────────────────┘
                   │                    │
           ┌───────┴───────┐    ┌───────┴───────┐
           │   STAGING     │    │     DEMO      │
           │   api:8080    │    │   api:8081    │
           └───────┬───────┘    └───────┬───────┘
                   │                    │
           ┌───────┴────────────────────┴───────┐
           │           PostgreSQL               │
           │   keerja_staging | keerja_demo     │
           └────────────────────────────────────┘
                           │
                   ┌───────┴───────┐
                   │     Redis     │
                   └───────────────┘
```

### Environment URLs

| Environment | API URL                                | Documentation                          |
| ----------- | -------------------------------------- | -------------------------------------- |
| STAGING     | http://staging-api.145.79.8.227.nip.io | https://bump.sh/doc/keerja-api-staging |
| DEMO        | http://demo-api.145.79.8.227.nip.io    | https://bump.sh/doc/keerja-api-demo    |

### Deployment Triggers

| Environment | Trigger               | Description                                 |
| ----------- | --------------------- | ------------------------------------------- |
| STAGING     | Push to `main` branch | Automatic deployment on every merge to main |
| DEMO        | Git tag (v*.*.\*)     | Deploy stable releases for client demos     |

---

## Prerequisites

### VPS Requirements

- **IP:** 145.79.8.227
- **OS:** Ubuntu 22.04 or Debian 11+
- **RAM:** 4GB minimum
- **Storage:** 50GB SSD
- **SSH access** with key-based authentication

### Required Tools (Local)

- Git
- SSH client
- rsync (for manual deployments)

### GitHub Repository Secrets

You'll need to add these secrets to your GitHub repository:

| Secret                | Description                    | Example                                  |
| --------------------- | ------------------------------ | ---------------------------------------- |
| `VPS_HOST`            | VPS IP address                 | `145.79.8.227`                           |
| `VPS_USER`            | SSH username                   | `root` or `deploy`                       |
| `SSH_PRIVATE_KEY`     | SSH private key for VPS access | `-----BEGIN OPENSSH PRIVATE KEY-----...` |
| `POSTGRES_PASSWORD`   | PostgreSQL admin password      | `secure_postgres_password_123`           |
| `REDIS_PASSWORD`      | Redis password                 | `secure_redis_password_123`              |
| `STAGING_DB_PASSWORD` | STAGING database password      | `staging_db_pass_123`                    |
| `DEMO_DB_PASSWORD`    | DEMO database password         | `demo_db_pass_123`                       |
| `JWT_SECRET_STAGING`  | JWT secret for STAGING         | `your_staging_jwt_secret_min_32_chars`   |
| `JWT_SECRET_DEMO`     | JWT secret for DEMO            | `your_demo_jwt_secret_min_32_chars`      |
| `BUMP_TOKEN`          | Bump.sh API token              | `your_bump_api_token`                    |
| `STAGING_BUMP_DOC_ID` | Bump.sh doc ID for STAGING     | `keerja-api-staging`                     |
| `DEMO_BUMP_DOC_ID`    | Bump.sh doc ID for DEMO        | `keerja-api-demo`                        |

---

## VPS Initial Setup

### Step 1: SSH into VPS

```bash
ssh root@145.79.8.227
```

### Step 2: Download and Run Setup Script

Option A: Download from GitHub (if repo is public):

```bash
curl -sSL https://raw.githubusercontent.com/rasfaxo/keerja-backend/main/scripts/setup-vps.sh | sudo bash
```

Option B: Clone repository first:

```bash
cd /opt
git clone https://github.com/rasfaxo/keerja-backend.git keerja
cd keerja
chmod +x scripts/setup-vps.sh
sudo ./scripts/setup-vps.sh
```

The setup script will install:

- Docker & Docker Compose
- Nginx
- Node.js & Bump CLI
- UFW Firewall
- Create necessary directories
- Configure swap (2GB)

### Step 3: Verify Installation

```bash
docker --version
docker compose version
nginx -v
node --version
bump --version
```

---

## Repository Setup

### Clone Repository (if not already done)

```bash
cd /opt/keerja
git clone https://github.com/rasfaxo/keerja-backend.git .
```

### Verify Structure

```bash
ls -la
# Should see: docker-compose.vps.yml, .env.staging.example, .env.demo.example, nginx/, scripts/
```

---

## Environment Configuration

### Step 1: Create Environment Files

```bash
cd /opt/keerja

# Create STAGING environment file
cp .env.staging.example .env.staging

# Create DEMO environment file
cp .env.demo.example .env.demo
```

### Step 2: Edit STAGING Configuration

```bash
nano .env.staging
```

**Critical variables to update:**

```env
# Database
DB_PASSWORD=your_secure_staging_password
STAGING_DB_PASSWORD=your_secure_staging_password
POSTGRES_PASSWORD=your_postgres_admin_password

# Redis
REDIS_PASSWORD=your_redis_password

# JWT
JWT_SECRET=your_staging_jwt_secret_minimum_32_characters_long

# CORS (add your frontend domains)
ALLOWED_ORIGINS=http://localhost:3000,http://staging.yourdomain.com

# Email (optional, for testing)
MAIL_HOST=smtp.gmail.com
MAIL_USERNAME=your-email@gmail.com
MAIL_PASSWORD=your-app-password
```

### Step 3: Edit DEMO Configuration

```bash
nano .env.demo
```

**Use different passwords from STAGING for security!**

---

## Nginx Configuration

### Step 1: Install Nginx Configs

```bash
# Copy configs
cp nginx/staging.conf /etc/nginx/sites-available/keerja-staging.conf
cp nginx/demo.conf /etc/nginx/sites-available/keerja-demo.conf

# Create symlinks
ln -sf /etc/nginx/sites-available/keerja-staging.conf /etc/nginx/sites-enabled/
ln -sf /etc/nginx/sites-available/keerja-demo.conf /etc/nginx/sites-enabled/

# Remove default site (optional)
rm -f /etc/nginx/sites-enabled/default
```

### Step 2: Test and Reload Nginx

```bash
# Test configuration
nginx -t

# Reload Nginx
systemctl reload nginx
```

---

## Initial Deployment

### Option 1: Using Deploy Script (Recommended)

```bash
cd /opt/keerja
chmod +x scripts/deploy-initial.sh

# Deploy both environments
./scripts/deploy-initial.sh all

# Or deploy individually
./scripts/deploy-initial.sh staging
./scripts/deploy-initial.sh demo
```

### Option 2: Manual Deployment

```bash
cd /opt/keerja

# Step 1: Start infrastructure
docker compose -f docker-compose.vps.yml up -d postgres redis

# Wait for PostgreSQL
sleep 15
docker compose -f docker-compose.vps.yml exec postgres pg_isready -U postgres

# Step 2: Run STAGING migrations
docker compose -f docker-compose.vps.yml --profile migrate-staging up --build

# Step 3: Deploy STAGING
docker compose -f docker-compose.vps.yml --profile staging up -d api-staging

# Step 4: Run DEMO migrations
docker compose -f docker-compose.vps.yml --profile migrate-demo up --build

# Step 5: Deploy DEMO
docker compose -f docker-compose.vps.yml --profile demo up -d api-demo
```

### Step 3: Verify Deployment

```bash
# Check containers
docker ps

# Health checks
curl http://localhost:8080/health/live
curl http://localhost:8081/health/live

# Via nginx
curl http://staging-api.145.79.8.227.nip.io/health/live
curl http://demo-api.145.79.8.227.nip.io/health/live
```

---

## GitHub Actions Setup

### Step 1: Add Repository Secrets

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Add each secret from the [Prerequisites](#required-tools-local) section

### Step 2: Create GitHub Environments (Optional but Recommended)

1. Go to **Settings** → **Environments**
2. Create `staging` environment
3. Create `demo` environment
4. Add environment-specific secrets if needed

### Step 3: Test Workflows

**STAGING Auto-Deploy:**

```bash
# Push to main branch
git checkout main
git push origin main
# → Triggers deploy-staging.yml
```

**DEMO Deploy:**

```bash
# Create and push a tag
git tag v1.0.0
git push origin v1.0.0
# → Triggers deploy-demo.yml
```

---

## Bump.sh Documentation Setup

### Step 1: Create Bump.sh Account

1. Go to https://bump.sh
2. Create an account or log in
3. Create an organization (optional)

### Step 2: Create Documentation Projects

1. Click **"Create new documentation"**
2. Create **"Keerja API - Staging"**
   - Note the documentation ID (slug)
3. Create **"Keerja API - Demo"**
   - Note the documentation ID (slug)

### Step 3: Get API Token

1. Go to **Settings** → **API tokens**
2. Create a new token with deploy permissions
3. Save the token securely

### Step 4: Add to GitHub Secrets

```
BUMP_TOKEN=your_api_token
STAGING_BUMP_DOC_ID=keerja-api-staging
DEMO_BUMP_DOC_ID=keerja-api-demo
```

### Step 5: Manual Documentation Update

```bash
# Export credentials
export BUMP_TOKEN=your_token
export STAGING_BUMP_DOC_ID=keerja-api-staging

# Update staging docs
./scripts/update-bump-docs.sh staging
```

---

## Maintenance Commands

### Using Makefile (from local machine)

```bash
# SSH into VPS
make vps-ssh

# Check status
make vps-status

# View logs
make vps-logs-staging
make vps-logs-demo

# Restart services
make vps-restart-staging
make vps-restart-demo

# Run migrations
make vps-migrate-staging
make vps-migrate-demo

# Backup databases
make vps-backup-staging
make vps-backup-demo

# Deploy manually
make vps-deploy-staging
make vps-deploy-demo
```

### Direct SSH Commands

```bash
# SSH into VPS
ssh root@145.79.8.227

# Navigate to project
cd /opt/keerja

# View all containers
docker compose -f docker-compose.vps.yml ps

# View logs
docker compose -f docker-compose.vps.yml logs -f api-staging

# Restart container
docker compose -f docker-compose.vps.yml --profile staging restart api-staging

# Stop all
docker compose -f docker-compose.vps.yml --profile staging --profile demo down

# Run health check
./scripts/health-check.sh
```

---

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker compose -f docker-compose.vps.yml logs api-staging

# Check if port is in use
netstat -tlnp | grep 8080

# Check environment file
cat .env.staging | grep -v PASSWORD
```

### Database Connection Failed

```bash
# Check PostgreSQL is running
docker compose -f docker-compose.vps.yml exec postgres pg_isready -U postgres

# Check databases exist
docker compose -f docker-compose.vps.yml exec postgres psql -U postgres -c "\l"

# Manually create database if missing
docker compose -f docker-compose.vps.yml exec postgres psql -U postgres -c "CREATE DATABASE keerja_staging;"
```

### Nginx 502 Bad Gateway

```bash
# Check if backend is running
curl http://localhost:8080/health/live

# Check nginx error logs
tail -f /var/log/nginx/keerja-staging-error.log

# Verify nginx config
nginx -t
```

### CORS Issues

1. Check `ALLOWED_ORIGINS` in `.env.staging`
2. Verify nginx is adding CORS headers:
   ```bash
   curl -I -X OPTIONS http://staging-api.145.79.8.227.nip.io/api/v1/health
   ```
3. For mobile apps, ensure null origin is handled

### Out of Memory

```bash
# Check memory usage
free -h
docker stats --no-stream

# Add/increase swap
fallocate -l 4G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile

# Reduce container memory limits in docker-compose.vps.yml
```

---

## Rollback Procedures

### Rollback STAGING

```bash
# SSH into VPS
ssh root@145.79.8.227
cd /opt/keerja

# List available backups
ls -la backups/

# Restore database
docker exec -i keerja-vps-postgres psql -U postgres -d keerja_staging < backups/staging_20241219_120000.sql

# Rollback to previous code (if using git)
git checkout HEAD~1

# Rebuild and restart
docker compose -f docker-compose.vps.yml build api-staging
docker compose -f docker-compose.vps.yml --profile staging up -d api-staging
```

### Rollback DEMO (to specific tag)

```bash
# Check available tags
git tag -l

# Checkout previous tag
git checkout v0.9.0

# Rebuild and restart
docker compose -f docker-compose.vps.yml build api-demo
docker compose -f docker-compose.vps.yml --profile demo up -d api-demo
```

### Emergency: Stop Everything

```bash
# Stop all containers
docker compose -f docker-compose.vps.yml --profile staging --profile demo down

# Or stop specific environment
docker compose -f docker-compose.vps.yml stop api-staging
```

---

## Quick Reference

### URLs

| Environment | API                                    | Health Check | Docs                                   |
| ----------- | -------------------------------------- | ------------ | -------------------------------------- |
| STAGING     | http://staging-api.145.79.8.227.nip.io | /health/live | https://bump.sh/doc/keerja-api-staging |
| DEMO        | http://demo-api.145.79.8.227.nip.io    | /health/live | https://bump.sh/doc/keerja-api-demo    |

### Important Files

| File                      | Purpose                       |
| ------------------------- | ----------------------------- |
| `docker-compose.vps.yml`  | VPS Docker orchestration      |
| `.env.staging`            | STAGING environment variables |
| `.env.demo`               | DEMO environment variables    |
| `nginx/staging.conf`      | STAGING nginx config          |
| `nginx/demo.conf`         | DEMO nginx config             |
| `scripts/health-check.sh` | Health check script           |

### Support

For issues, contact: support@keerja.com
