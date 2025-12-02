# Config Directory

This directory contains configuration files for the Keerja Backend application.

## Directory Structure

```
config/
├── README.md                          # This file
└── firebase-service-account.json      # Firebase service account credentials (GITIGNORED)
```

## Firebase Service Account

### File: `firebase-service-account.json`

**Description**: Firebase Admin SDK service account credentials for Firebase Cloud Messaging (FCM).

**Security Level**: **CRITICAL - TOP SECRET**

**How to Get**:

1. Visit [Firebase Console](https://console.firebase.google.com/)
2. Select your project
3. Go to: ⚙️ Project Settings → Service Accounts tab
4. Click: "Generate new private key"
5. Save the downloaded file as: `config/firebase-service-account.json`

**File Format**:

```json
{
  "type": "service_account",
  "project_id": "your-project-id",
  "private_key_id": "...",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk-xxxxx@your-project.iam.gserviceaccount.com",
  "client_id": "...",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "..."
}
```

### Security Best Practices

1. **NEVER commit this file to Git**

   - Already added to `.gitignore`
   - Verify: `git status` should NOT show this file

2. **File Permissions** (Linux/Mac only)

   ```bash
   chmod 600 config/firebase-service-account.json
   ```

3. **Storage**

   - Development: Store locally in `config/` folder
   - Production: Use environment variables or secret management
   - CI/CD: Store as encrypted secret in GitHub Actions/GitLab CI

4. **Rotation**

   - Rotate key every 3-6 months
   - Rotate immediately if compromised
   - Delete old keys from Firebase Console

5. **Access Control**
   - Only team leads and DevOps should have access
   - Use separate keys for dev/staging/prod environments

### What to Do If Key is Compromised

If the service account key is accidentally exposed:

```bash
# 1. Immediately rotate the key
# Go to Firebase Console → Project Settings → Service Accounts
# Generate new key and replace the file

# 2. Delete the exposed key
# Firebase Console → Service Accounts → Click ⋮ → Delete key

# 3. If committed to Git, remove from history
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch config/firebase-service-account.json" \
  --prune-empty --tag-name-filter cat -- --all

# 4. Force push (WARNING: Coordinate with team!)
git push origin --force --all

# 5. Notify security team
# Email: security@keerja.com
```

## Environment Variables

Instead of storing the file path in code, use environment variables:

```env
# .env file
FCM_CREDENTIALS_FILE=config/firebase-service-account.json
```

For production, consider using:

- **AWS Secrets Manager**
- **Google Secret Manager**
- **HashiCorp Vault**
- **Environment-specific encrypted files**

## Adding New Configuration Files

When adding new config files:

1. **Update `.gitignore`** if file contains secrets
2. **Document** the file format in this README
3. **Provide** `.example` version for reference
4. **Add** to verification script if critical

## Verification

To verify your Firebase credentials:

```bash
# Windows PowerShell
.\scripts\verify-fcm-setup.ps1

# Linux/Mac
bash scripts/verify-fcm-setup.sh
```

## Related Documentation

- [Firebase Setup Guide](../docs/FIREBASE_SETUP_GUIDE.md) - Complete setup instructions
- [Firebase Quick Start](../docs/FIREBASE_QUICK_START.md) - Quick reference card
- [FCM Instruction Prompt](../.github/prompts/fcm_intruction.prompt.md) - Implementation phases

## Troubleshooting

### Error: "FCM credentials file not found"

```bash
# Check if file exists
ls -la config/firebase-service-account.json

# Check path in .env
cat .env | grep FCM_CREDENTIALS_FILE
```

### Error: "Invalid service account"

```bash
# Validate JSON format
jq empty config/firebase-service-account.json

# Check project_id
jq -r '.project_id' config/firebase-service-account.json
```

### Error: "Permission denied"

```bash
# Fix file permissions (Linux/Mac)
chmod 600 config/firebase-service-account.json

# Check file ownership
ls -la config/firebase-service-account.json
```

---

**⚠️ REMEMBER**: This directory contains sensitive credentials. Always practice security best practices!
