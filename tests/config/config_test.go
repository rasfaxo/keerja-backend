package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	cfgpkg "keerja-backend/internal/config"
)

func TestLoadConfig_AllowsMobileRedirectsAndRedisURL(t *testing.T) {
	// Set temporary env vars
	// Ensure required DB_PASSWORD is present to avoid LoadConfig validation failure
	os.Setenv("DB_PASSWORD", "notempty")
	os.Setenv("REDIS_URL", "redis://:password@localhost:6379/0")
	os.Setenv("ALLOWED_MOBILE_REDIRECT_URIS", "myapp://oauth-callback, myapp://prod-callback")

	c := cfgpkg.LoadConfig()

	require.Equal(t, "redis://:password@localhost:6379/0", c.RedisURL)
	require.Contains(t, c.AllowedMobileRedirectURIs, "myapp://oauth-callback")
	require.Contains(t, c.AllowedMobileRedirectURIs, "myapp://prod-callback")

	// cleanup
	os.Unsetenv("REDIS_URL")
	os.Unsetenv("ALLOWED_MOBILE_REDIRECT_URIS")
	os.Unsetenv("DB_PASSWORD")
}
