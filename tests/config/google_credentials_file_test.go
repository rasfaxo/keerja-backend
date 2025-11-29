package config_test

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/require"

    cfgpkg "keerja-backend/internal/config"
)

func TestLoadConfig_ReadsGoogleCredentialsFile(t *testing.T) {
    // prepare a temporary credentials JSON file
    tmpDir := t.TempDir()
    filePath := filepath.Join(tmpDir, "client_secret_test.json")

    content := `{"web":{"client_id":"abc123","project_id":"keerja-backend","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","client_secret":"secretXYZ","redirect_uris":["http://localhost:8080/api/v1/auth/oauth/google/callback"]}}`
    require.NoError(t, os.WriteFile(filePath, []byte(content), 0600))

    // Set envs needed for LoadConfig validation
    os.Setenv("DB_PASSWORD", "notempty")
    os.Setenv("GOOGLE_CREDENTIALS_FILE", filePath)

    cfg := cfgpkg.LoadConfig()

    require.Equal(t, "abc123", cfg.GoogleClientID)
    require.Equal(t, "secretXYZ", cfg.GoogleClientSecret)
    require.Equal(t, "http://localhost:8080/api/v1/auth/oauth/google/callback", cfg.GoogleRedirectURI)

    // cleanup
    os.Unsetenv("GOOGLE_CREDENTIALS_FILE")
    os.Unsetenv("DB_PASSWORD")
}
