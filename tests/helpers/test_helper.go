package helpers

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig holds common test configuration
type TestConfig struct {
	TestDBURL      string
	TestTimeout    time.Duration
	TestDataPath   string
	MocksEnabled   bool
	CleanupEnabled bool
}

// DefaultTestConfig returns default test configuration
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		TestDBURL:      getEnvOrDefault("TEST_DB_URL", "postgres://postgres:postgres@localhost:5432/keerja_test?sslmode=disable"),
		TestTimeout:    30 * time.Second,
		TestDataPath:   "./testdata",
		MocksEnabled:   true,
		CleanupEnabled: true,
	}
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupTest initializes test environment
func SetupTest(t *testing.T) context.Context {
	t.Helper()

	// Set test environment
	os.Setenv("APP_ENV", "test")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// Cleanup on test completion
	t.Cleanup(func() {
		cancel()
	})

	return ctx
}

// TeardownTest cleans up test environment
func TeardownTest(t *testing.T) {
	t.Helper()
	// Additional cleanup if needed
}

// AssertNoError is a helper to assert no error
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.NoError(t, err, msgAndArgs...)
}

// AssertError is a helper to assert error exists
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.Error(t, err, msgAndArgs...)
}

// AssertEqual is a helper for equality assertion
func AssertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, expected, actual, msgAndArgs...)
}

// AssertNotNil is a helper to assert value is not nil
func AssertNotNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	assert.NotNil(t, object, msgAndArgs...)
}

// AssertContains is a helper to assert string contains substring
func AssertContains(t *testing.T, s, contains string, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Contains(t, s, contains, msgAndArgs...)
}

// RandomString generates random string for testing
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[int(time.Now().UnixNano())%len(charset)]
	}
	return string(b)
}

// RandomEmail generates random email for testing
func RandomEmail() string {
	return RandomString(10) + "@test.com"
}

// RandomPhone generates random phone number for testing
func RandomPhone() string {
	return "+628" + RandomString(9)
}
