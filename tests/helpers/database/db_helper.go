package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// TestDB holds test database configuration
type TestDB struct {
	DB     *sql.DB
	DBName string
}

// SetupTestDB creates and initializes a test database
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	testDBURL := getTestDBURL()

	db, err := sql.Open("postgres", testDBURL)
	require.NoError(t, err, "Failed to connect to test database")

	// Ping to verify connection
	err = db.Ping()
	require.NoError(t, err, "Failed to ping test database")

	testDB := &TestDB{
		DB:     db,
		DBName: "keerja_test",
	}

	// Cleanup on test completion
	t.Cleanup(func() {
		testDB.Close()
	})

	return testDB
}

// Close closes the database connection
func (tdb *TestDB) Close() {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
}

// CleanDatabase truncates all tables in the test database
func (tdb *TestDB) CleanDatabase(t *testing.T) {
	t.Helper()

	tables := []string{
		"application_notes",
		"application_documents",
		"application_stages",
		"job_applications",
		"device_tokens",
		"notifications",
		"saved_jobs",
		"job_skills",
		"jobs",
		"employer_users",
		"companies",
		"user_profiles",
		"users",
	}

	for _, table := range tables {
		_, err := tdb.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		require.NoError(t, err, "Failed to truncate table: %s", table)
	}
}

// BeginTransaction starts a new transaction
func (tdb *TestDB) BeginTransaction(t *testing.T) *sql.Tx {
	t.Helper()

	tx, err := tdb.DB.Begin()
	require.NoError(t, err, "Failed to begin transaction")

	return tx
}

// RollbackTransaction rolls back a transaction
func (tdb *TestDB) RollbackTransaction(t *testing.T, tx *sql.Tx) {
	t.Helper()

	err := tx.Rollback()
	require.NoError(t, err, "Failed to rollback transaction")
}

// CommitTransaction commits a transaction
func (tdb *TestDB) CommitTransaction(t *testing.T, tx *sql.Tx) {
	t.Helper()

	err := tx.Commit()
	require.NoError(t, err, "Failed to commit transaction")
}

// ExecQuery executes a query and returns the result
func (tdb *TestDB) ExecQuery(t *testing.T, query string, args ...interface{}) sql.Result {
	t.Helper()

	result, err := tdb.DB.Exec(query, args...)
	require.NoError(t, err, "Failed to execute query")

	return result
}

// QueryRow executes a query and returns a single row
func (tdb *TestDB) QueryRow(t *testing.T, query string, args ...interface{}) *sql.Row {
	t.Helper()
	return tdb.DB.QueryRow(query, args...)
}

// Query executes a query and returns multiple rows
func (tdb *TestDB) Query(t *testing.T, query string, args ...interface{}) *sql.Rows {
	t.Helper()

	rows, err := tdb.DB.Query(query, args...)
	require.NoError(t, err, "Failed to execute query")

	return rows
}

// SeedData inserts seed data for testing
func (tdb *TestDB) SeedData(t *testing.T, query string, args ...interface{}) {
	t.Helper()

	_, err := tdb.DB.Exec(query, args...)
	require.NoError(t, err, "Failed to seed data")
}

// CountRows returns the count of rows in a table
func (tdb *TestDB) CountRows(t *testing.T, tableName string) int {
	t.Helper()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := tdb.DB.QueryRow(query).Scan(&count)
	require.NoError(t, err, "Failed to count rows")

	return count
}

// TableExists checks if a table exists
func (tdb *TestDB) TableExists(t *testing.T, tableName string) bool {
	t.Helper()

	var exists bool
	query := `SELECT EXISTS (
SELECT FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name = $1
)`
	err := tdb.DB.QueryRow(query, tableName).Scan(&exists)
	require.NoError(t, err, "Failed to check table existence")

	return exists
}

// WithTransaction executes a function within a transaction
func (tdb *TestDB) WithTransaction(t *testing.T, fn func(*sql.Tx) error) {
	t.Helper()

	tx := tdb.BeginTransaction(t)
	defer tdb.RollbackTransaction(t, tx)

	err := fn(tx)
	require.NoError(t, err, "Transaction function failed")

	tdb.CommitTransaction(t, tx)
}

// WithContext executes a function with context
func (tdb *TestDB) WithContext(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

// getTestDBURL returns the test database URL
func getTestDBURL() string {
	// Check environment variable first
	if url := getEnv("TEST_DB_URL"); url != "" {
		return url
	}

	// Default test database URL
	return "postgres://postgres:postgres@localhost:5432/keerja_test?sslmode=disable"
}

// getEnv gets environment variable
func getEnv(key string) string {
	return os.Getenv(key)
}
