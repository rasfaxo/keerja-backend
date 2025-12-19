# Bump.sh API Documentation Setup Guide

This guide explains how to set up Bump.sh for automatic API documentation hosting with multi-environment support (STAGING + DEMO).

## Table of Contents

1. [What is Bump.sh?](#what-is-bumpsh)
2. [Creating Documentation Projects](#creating-documentation-projects)
3. [Getting API Credentials](#getting-api-credentials)
4. [GitHub Integration](#github-integration)
5. [Manual Updates](#manual-updates)
6. [Configuration Files](#configuration-files)
7. [Troubleshooting](#troubleshooting)

---

## What is Bump.sh?

Bump.sh is a hosted API documentation platform that:

- Automatically generates beautiful documentation from OpenAPI/Swagger files
- Provides API diff on pull requests
- Supports versioning and change history
- Offers a developer-friendly interface

### Why Two Documentation Projects?

We maintain separate documentation for each environment:

| Environment | Purpose              | URL                             | Content                                    |
| ----------- | -------------------- | ------------------------------- | ------------------------------------------ |
| STAGING     | Development testing  | staging-api.145.79.8.227.nip.io | Latest features, may have breaking changes |
| DEMO        | Client presentations | demo-api.145.79.8.227.nip.io    | Stable releases only                       |

---

## Creating Documentation Projects

### Step 1: Create Bump.sh Account

1. Go to [https://bump.sh](https://bump.sh)
2. Click **"Sign up"** or **"Get started for free"**
3. Create an account using GitHub, GitLab, or email

### Step 2: Create Organization (Optional)

1. Click your profile → **"Create organization"**
2. Name it `keerja` or your company name
3. This helps organize multiple documentation projects

### Step 3: Create STAGING Documentation

1. Click **"+ New documentation"**
2. Fill in the details:
   - **Name:** `Keerja API - Staging`
   - **Slug:** `keerja-api-staging` (this becomes your DOC_ID)
   - **Visibility:** Private or Public (your choice)
3. Click **"Create"**
4. **Save the slug** - this is your `STAGING_BUMP_DOC_ID`

### Step 4: Create DEMO Documentation

1. Click **"+ New documentation"** again
2. Fill in the details:
   - **Name:** `Keerja API - Demo`
   - **Slug:** `keerja-api-demo`
   - **Visibility:** Public (recommended for client access)
3. Click **"Create"**
4. **Save the slug** - this is your `DEMO_BUMP_DOC_ID`

### Step 5: Initial Upload

For each documentation:

1. Click **"Deploy new version"**
2. Upload your `docs/swagger.yaml` file
3. Verify the documentation renders correctly

---

## Getting API Credentials

### Step 1: Create API Token

1. Click your profile icon → **"Account settings"**
2. Go to **"API tokens"** tab
3. Click **"Create new token"**
4. Configure the token:
   - **Name:** `GitHub Actions Deploy`
   - **Permissions:** `Deploy` (minimum required)
   - **Scope:** Select your organization or personal account
5. Click **"Create"**
6. **Copy the token immediately** - it won't be shown again!

### Step 2: Save Credentials

Store these securely:

```
BUMP_TOKEN=bmp_xxxxxxxxxxxxxxxxxxxxxxxxxxxx
STAGING_BUMP_DOC_ID=keerja-api-staging
DEMO_BUMP_DOC_ID=keerja-api-demo
```

---

## GitHub Integration

### Step 1: Add Secrets to GitHub Repository

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **"New repository secret"** for each:

| Secret Name           | Value                | Description                |
| --------------------- | -------------------- | -------------------------- |
| `BUMP_TOKEN`          | `bmp_xxx...`         | Your Bump.sh API token     |
| `STAGING_BUMP_DOC_ID` | `keerja-api-staging` | STAGING documentation slug |
| `DEMO_BUMP_DOC_ID`    | `keerja-api-demo`    | DEMO documentation slug    |

### Step 2: Verify Workflow Configuration

The `.github/workflows/deploy-api-docs.yml` workflow is already configured to:

- **On push to main:** Update STAGING documentation
- **On release/tag:** Update DEMO documentation
- **On pull request:** Show API diff as PR comment

### Step 3: Test the Integration

**Test STAGING docs update:**

```bash
# Make a change to swagger.yaml
echo "# Test comment" >> docs/swagger.yaml
git add docs/swagger.yaml
git commit -m "docs: test bump.sh integration"
git push origin main
```

Check GitHub Actions to see the workflow run.

**Test DEMO docs update:**

```bash
git tag v1.0.0-test
git push origin v1.0.0-test
```

---

## Manual Updates

### Using the Script

```bash
# Set environment variables
export BUMP_TOKEN="your_bump_token"
export STAGING_BUMP_DOC_ID="keerja-api-staging"
export DEMO_BUMP_DOC_ID="keerja-api-demo"

# Update STAGING documentation
./scripts/update-bump-docs.sh staging

# Update DEMO documentation
./scripts/update-bump-docs.sh demo

# Update both
./scripts/update-bump-docs.sh all

# Preview changes (dry run)
./scripts/update-bump-docs.sh preview-staging
```

### Using Bump CLI Directly

```bash
# Install Bump CLI
npm install -g bump-cli

# Deploy to STAGING
bump deploy docs/swagger.yaml \
  --doc keerja-api-staging \
  --token $BUMP_TOKEN

# Deploy to DEMO
bump deploy docs/swagger.yaml \
  --doc keerja-api-demo \
  --token $BUMP_TOKEN

# Preview diff
bump diff docs/swagger.yaml \
  --doc keerja-api-staging \
  --token $BUMP_TOKEN
```

---

## Configuration Files

### .bump.yml (Repository Root)

Create a `.bump.yml` file for local development:

```yaml
# Bump.sh Configuration
# This file is used for local development and CI/CD

# Documentation settings
version: "1.0"

# Default documentation (used when not specified)
documentation:
  # The documentation slug on Bump.sh
  # Override with --doc flag or BUMP_DOC_ID env var
  id: keerja-api-staging

  # Path to your OpenAPI specification
  specification: docs/swagger.yaml

# Environment-specific overrides (for reference)
# These are handled by CI/CD workflows
environments:
  staging:
    id: keerja-api-staging
    host: staging-api.145.79.8.227.nip.io
  demo:
    id: keerja-api-demo
    host: demo-api.145.79.8.227.nip.io
```

### Swagger Host Configuration

The GitHub Actions workflow automatically updates the `host` value in swagger.yaml for each environment. If doing manual updates, ensure the host is correct:

**For STAGING:**

```yaml
servers:
  - url: http://staging-api.145.79.8.227.nip.io
    description: Staging server
```

**For DEMO:**

```yaml
servers:
  - url: http://demo-api.145.79.8.227.nip.io
    description: Demo server
```

---

## Troubleshooting

### "Authentication failed" Error

```bash
# Verify token is valid
curl -H "Authorization: Token $BUMP_TOKEN" https://bump.sh/api/v1/ping

# Check token permissions
# Go to Bump.sh → Account settings → API tokens
# Ensure token has "Deploy" permission
```

### "Documentation not found" Error

```bash
# Verify doc ID exists
# Go to Bump.sh dashboard and check the documentation slug
# It should match exactly (case-sensitive)

# Common issues:
# - Typo in doc ID
# - Using name instead of slug
# - Wrong organization scope on token
```

### Swagger Validation Errors

```bash
# Validate locally before deploying
npm install -g @apidevtools/swagger-cli
swagger-cli validate docs/swagger.yaml

# Common issues:
# - Invalid $ref paths
# - Missing required fields
# - Duplicate operation IDs
```

### Diff Not Showing in PR

1. Ensure `GITHUB_TOKEN` has `pull-requests: write` permission
2. Check that the workflow has correct permissions:
   ```yaml
   permissions:
     pull-requests: write
     contents: read
   ```
3. Verify the PR modifies `docs/swagger.yaml`

### Documentation Not Updating

```bash
# Force a new deployment
bump deploy docs/swagger.yaml \
  --doc $STAGING_BUMP_DOC_ID \
  --token $BUMP_TOKEN \
  --auto-create

# Check deployment history on Bump.sh dashboard
```

---

## Documentation URLs

After setup, your documentation will be available at:

| Environment | Documentation URL                      |
| ----------- | -------------------------------------- |
| STAGING     | https://bump.sh/doc/keerja-api-staging |
| DEMO        | https://bump.sh/doc/keerja-api-demo    |

You can also use custom domains if you have a Bump.sh paid plan.

---

## Best Practices

### 1. Keep Swagger File Updated

Always update `docs/swagger.yaml` when adding/modifying endpoints:

```bash
# After changing an endpoint
vim docs/swagger.yaml
git add docs/swagger.yaml
git commit -m "docs: add new user profile endpoint"
git push origin main
# → Automatically updates STAGING docs
```

### 2. Use Semantic Versioning for DEMO

Only create tags for stable releases:

```bash
# v1.0.0 - First stable release
# v1.1.0 - New features (backward compatible)
# v1.1.1 - Bug fixes
# v2.0.0 - Breaking changes

git tag v1.1.0
git push origin v1.1.0
# → Updates DEMO docs
```

### 3. Review Diffs in PRs

The API diff comment helps reviewers understand API changes:

- Breaking changes are highlighted
- New endpoints are listed
- Removed fields are flagged

### 4. Notify Frontend Teams

When documentation is updated, notify teams:

```markdown
🔄 API Documentation Updated

Environment: STAGING
Changes: Added /api/v1/users/me/profile endpoint
Docs: https://bump.sh/doc/keerja-api-staging

Please review and update your integrations.
```

---

## Quick Reference

### Environment Variables

```bash
export BUMP_TOKEN="bmp_xxxxxxxxxxxxxxxxxxxx"
export STAGING_BUMP_DOC_ID="keerja-api-staging"
export DEMO_BUMP_DOC_ID="keerja-api-demo"
export VPS_IP="145.79.8.227"
```

### Common Commands

```bash
# Update staging docs
./scripts/update-bump-docs.sh staging

# Update demo docs
./scripts/update-bump-docs.sh demo

# Validate swagger
swagger-cli validate docs/swagger.yaml

# Preview changes
bump diff docs/swagger.yaml --doc $STAGING_BUMP_DOC_ID --token $BUMP_TOKEN
```

### GitHub Secrets Required

| Secret                | Required For            |
| --------------------- | ----------------------- |
| `BUMP_TOKEN`          | All Bump.sh operations  |
| `STAGING_BUMP_DOC_ID` | STAGING docs deployment |
| `DEMO_BUMP_DOC_ID`    | DEMO docs deployment    |
