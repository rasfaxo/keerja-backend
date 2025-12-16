# Bump.sh Integration Setup Guide

This guide describes how to set up, configure, and maintain the integration between **Keerja Backend** and **Bump.sh** for API documentation.

## Part 1: Prerequisites

Before you begin, ensure you have:
1.  **Admin access** to the [Keerja Backend GitHub Repository](https://github.com/rasfaxo/keerja-backend).
2.  A **Bump.sh account** (Sign up at [bump.sh](https://bump.sh)).
3.  `docs/swagger.yaml` file present in the repository.

## Part 2: Secure Token Generation

**⚠️ SECURITY WARNING: Never commit tokens to the repository!**

1.  Log in to your Bump.sh account.
2.  Create a new documentation or select an existing one.
3.  Go to **Settings** > **CI deployment**.
4.  Copy the **Doc ID** (slug) and the **Token**.
    *   *Note: Treat this token like a password.*

## Part 3: GitHub Secrets Setup

To secure your tokens, we use GitHub Secrets.

1.  Go to your GitHub Repository.
2.  Navigate to **Settings** > **Secrets and variables** > **Actions**.
3.  Click **New repository secret**.
4.  Add the following secrets:

    | Name | Value | Description |
    | :--- | :--- | :--- |
    | `BUMP_TOKEN` | `<your-bump-sh-token>` | The secret token from Bump.sh settings. |
    | `BUMP_DOC_ID` | `<your-doc-slug>` | The ID/Slug of your documentation (e.g., `keerja-api`). |

## Part 4: Deployment Setup

The integration is handled by GitHub Actions.

### Files Overview
*   **Workflow**: `.github/workflows/deploy-api-docs.yml` - Orchestrates validation and deployment.
*   **Config**: `.bump.yml` - Configures Bump.sh behavior (validation rules, preview settings).
*   **Validation**: `scripts/validate-swagger.sh` - Local validation script.

### How it works
1.  **On Push to Main**: The workflow validates the swagger file and deploys it to Bump.sh.
2.  **On Pull Request**: The workflow validates the swagger file and posts a comment with the **API Diff** (changes preview) to the PR.
3.  **Daily Schedule**: Runs a health check on the documentation.

## Part 5: Maintenance

### Updating Documentation
1.  Modify `docs/swagger.yaml` in a new branch.
2.  Run local validation:
    ```bash
    ./scripts/validate-swagger.sh
    ```
3.  Commit and push to create a Pull Request.
4.  Check the PR comment for the API Diff to ensure changes are correct.
5.  Merge to main to deploy.

### Monitoring
*   Check the **Actions** tab in GitHub to see deployment history.
*   Run the status check script locally:
    ```bash
    export BUMP_TOKEN=your_token
    ./scripts/check-deployment-status.sh
    ```

## Part 6: Troubleshooting

### Common Issues

**1. "Validation Failed" in CI**
*   **Cause**: Invalid Swagger/OpenAPI syntax.
*   **Fix**: Check the CI logs for specific line numbers. Run `swagger-cli validate docs/swagger.yaml` locally.

**2. "Unauthorized" or "Token Invalid"**
*   **Cause**: Incorrect or expired `BUMP_TOKEN` in GitHub Secrets.
*   **Fix**: Regenerate the token in Bump.sh and update the GitHub Secret.

**3. API Diff comment not appearing**
*   **Cause**: GitHub Actions permissions issue.
*   **Fix**: Ensure the workflow has `pull-requests: write` permission (already configured in the provided YAML).

### Contact
For support, contact the dev team at dev@keerja.com.
