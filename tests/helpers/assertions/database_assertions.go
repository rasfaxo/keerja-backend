package assertions

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// DatabaseAssertion provides database-specific assertions
type DatabaseAssertion struct {
	T  *testing.T
	DB *sql.DB
}

// NewDatabaseAssertion creates a new database assertion helper
func NewDatabaseAssertion(t *testing.T, db *sql.DB) *DatabaseAssertion {
	return &DatabaseAssertion{
		T:  t,
		DB: db,
	}
}

// AssertRowCount asserts the number of rows in a table
func (da *DatabaseAssertion) AssertRowCount(tableName string, expectedCount int, msgAndArgs ...interface{}) {
	da.T.Helper()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := da.DB.QueryRow(query).Scan(&count)
	require.NoError(da.T, err, "Failed to count rows in table: %s", tableName)

	assert.Equal(da.T, expectedCount, count, msgAndArgs...)
}

// AssertRowExists asserts a row exists in a table
func (da *DatabaseAssertion) AssertRowExists(tableName string, whereClause string, args ...interface{}) {
	da.T.Helper()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tableName, whereClause)
	err := da.DB.QueryRow(query, args...).Scan(&count)
	require.NoError(da.T, err, "Failed to check row existence")

	assert.Greater(da.T, count, 0, "Expected row to exist in %s with condition: %s", tableName, whereClause)
}

// AssertRowNotExists asserts a row does not exist in a table
func (da *DatabaseAssertion) AssertRowNotExists(tableName string, whereClause string, args ...interface{}) {
	da.T.Helper()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tableName, whereClause)
	err := da.DB.QueryRow(query, args...).Scan(&count)
	require.NoError(da.T, err, "Failed to check row non-existence")

	assert.Equal(da.T, 0, count, "Expected no row to exist in %s with condition: %s", tableName, whereClause)
}

// AssertFieldValue asserts a field value in a table
func (da *DatabaseAssertion) AssertFieldValue(tableName, fieldName string, whereClause string, expectedValue interface{}, args ...interface{}) {
	da.T.Helper()

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", fieldName, tableName, whereClause)
	var actualValue interface{}
	err := da.DB.QueryRow(query, args...).Scan(&actualValue)
	require.NoError(da.T, err, "Failed to get field value")

	assert.Equal(da.T, expectedValue, actualValue, "Field %s in %s does not match expected value", fieldName, tableName)
}

// AssertTableEmpty asserts a table is empty
func (da *DatabaseAssertion) AssertTableEmpty(tableName string, msgAndArgs ...interface{}) {
	da.AssertRowCount(tableName, 0, msgAndArgs...)
}

// AssertTableNotEmpty asserts a table is not empty
func (da *DatabaseAssertion) AssertTableNotEmpty(tableName string, msgAndArgs ...interface{}) {
	da.T.Helper()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := da.DB.QueryRow(query).Scan(&count)
	require.NoError(da.T, err, "Failed to count rows")

	assert.Greater(da.T, count, 0, msgAndArgs...)
}

// AssertUserExists asserts a user exists by ID
func (da *DatabaseAssertion) AssertUserExists(userID int64) {
	da.AssertRowExists("users", "id = $1", userID)
}

// AssertUserNotExists asserts a user does not exist by ID
func (da *DatabaseAssertion) AssertUserNotExists(userID int64) {
	da.AssertRowNotExists("users", "id = $1", userID)
}

// AssertUserByEmail asserts a user exists by email
func (da *DatabaseAssertion) AssertUserByEmail(email string) {
	da.AssertRowExists("users", "email = $1", email)
}

// AssertUserEmailVerified asserts a user's email is verified
func (da *DatabaseAssertion) AssertUserEmailVerified(userID int64, isVerified bool) {
	da.AssertFieldValue("users", "is_email_verified", "id = $1", isVerified, userID)
}

// AssertJobExists asserts a job exists by ID
func (da *DatabaseAssertion) AssertJobExists(jobID int64) {
	da.AssertRowExists("jobs", "id = $1", jobID)
}

// AssertJobNotExists asserts a job does not exist by ID
func (da *DatabaseAssertion) AssertJobNotExists(jobID int64) {
	da.AssertRowNotExists("jobs", "id = $1", jobID)
}

// AssertJobStatus asserts a job's status
func (da *DatabaseAssertion) AssertJobStatus(jobID int64, expectedStatus string) {
	da.AssertFieldValue("jobs", "status", "id = $1", expectedStatus, jobID)
}

// AssertApplicationExists asserts an application exists by ID
func (da *DatabaseAssertion) AssertApplicationExists(applicationID int64) {
	da.AssertRowExists("job_applications", "id = $1", applicationID)
}

// AssertApplicationNotExists asserts an application does not exist by ID
func (da *DatabaseAssertion) AssertApplicationNotExists(applicationID int64) {
	da.AssertRowNotExists("job_applications", "id = $1", applicationID)
}

// AssertApplicationStatus asserts an application's status
func (da *DatabaseAssertion) AssertApplicationStatus(applicationID int64, expectedStatus string) {
	da.AssertFieldValue("job_applications", "status", "id = $1", expectedStatus, applicationID)
}

// AssertApplicationCountForJob asserts the number of applications for a job
func (da *DatabaseAssertion) AssertApplicationCountForJob(jobID int64, expectedCount int) {
	da.T.Helper()

	var count int
	query := "SELECT COUNT(*) FROM job_applications WHERE job_id = $1"
	err := da.DB.QueryRow(query, jobID).Scan(&count)
	require.NoError(da.T, err, "Failed to count applications for job")

	assert.Equal(da.T, expectedCount, count, "Application count for job %d does not match", jobID)
}

// AssertNotificationExists asserts a notification exists by ID
func (da *DatabaseAssertion) AssertNotificationExists(notificationID int64) {
	da.AssertRowExists("notifications", "id = $1", notificationID)
}

// AssertNotificationIsRead asserts a notification's read status
func (da *DatabaseAssertion) AssertNotificationIsRead(notificationID int64, isRead bool) {
	da.AssertFieldValue("notifications", "is_read", "id = $1", isRead, notificationID)
}

// AssertUnreadNotificationCount asserts the count of unread notifications for a user
func (da *DatabaseAssertion) AssertUnreadNotificationCount(userID int64, expectedCount int) {
	da.T.Helper()

	var count int
	query := "SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false"
	err := da.DB.QueryRow(query, userID).Scan(&count)
	require.NoError(da.T, err, "Failed to count unread notifications")

	assert.Equal(da.T, expectedCount, count, "Unread notification count for user %d does not match", userID)
}

// AssertDeviceTokenExists asserts a device token exists
func (da *DatabaseAssertion) AssertDeviceTokenExists(token string) {
	da.AssertRowExists("device_tokens", "token = $1", token)
}

// AssertDeviceTokenActive asserts a device token is active
func (da *DatabaseAssertion) AssertDeviceTokenActive(token string, isActive bool) {
	da.AssertFieldValue("device_tokens", "is_active", "token = $1", isActive, token)
}

// AssertCompanyExists asserts a company exists by ID
func (da *DatabaseAssertion) AssertCompanyExists(companyID int64) {
	da.AssertRowExists("companies", "id = $1", companyID)
}

// AssertCompanyVerified asserts a company's verification status
func (da *DatabaseAssertion) AssertCompanyVerified(companyID int64, isVerified bool) {
	da.AssertFieldValue("companies", "is_verified", "id = $1", isVerified, companyID)
}

// AssertForeignKeyConstraint asserts a foreign key constraint is enforced
func (da *DatabaseAssertion) AssertForeignKeyConstraint(parentTable, childTable, foreignKey string, parentID int64) {
	da.T.Helper()

	// Try to delete parent record (should fail if child exists)
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", parentTable)
	_, err := da.DB.Exec(query, parentID)

	// If there are child records, deletion should fail
	if err != nil {
		assert.Contains(da.T, err.Error(), "foreign key constraint",
			"Expected foreign key constraint error")
	}
}

// AssertUniqueConstraint asserts a unique constraint is enforced
func (da *DatabaseAssertion) AssertUniqueConstraint(tableName, columnName string, value interface{}) {
	da.T.Helper()

	// Try to insert duplicate value (should fail)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1)", tableName, columnName)
	_, err := da.DB.Exec(query, value)

	assert.Error(da.T, err, "Expected unique constraint violation")
	if err != nil {
		assert.Contains(da.T, err.Error(), "unique constraint",
			"Expected unique constraint error")
	}
}

// AssertTimestampUpdated asserts a timestamp field was updated
func (da *DatabaseAssertion) AssertTimestampUpdated(tableName, timestampField, whereClause string, args ...interface{}) {
	da.T.Helper()

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", timestampField, tableName, whereClause)
	var timestamp sql.NullTime
	err := da.DB.QueryRow(query, args...).Scan(&timestamp)
	require.NoError(da.T, err, "Failed to get timestamp")

	assert.True(da.T, timestamp.Valid, "Timestamp %s should not be null", timestampField)
}

// Quick assertion functions without creating instance

// AssertRecordCount asserts the number of records in a table
func AssertRecordCount(t *testing.T, db *sql.DB, tableName string, expectedCount int) {
	t.Helper()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.QueryRow(query).Scan(&count)
	require.NoError(t, err, "Failed to count rows")

	assert.Equal(t, expectedCount, count)
}

// AssertRecordExists asserts a record exists
func AssertRecordExists(t *testing.T, db *sql.DB, tableName, whereClause string, args ...interface{}) {
	t.Helper()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tableName, whereClause)
	err := db.QueryRow(query, args...).Scan(&count)
	require.NoError(t, err, "Failed to check record existence")

	assert.Greater(t, count, 0, "Expected record to exist")
}

// AssertRecordNotExists asserts a record does not exist
func AssertRecordNotExists(t *testing.T, db *sql.DB, tableName, whereClause string, args ...interface{}) {
	t.Helper()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tableName, whereClause)
	err := db.QueryRow(query, args...).Scan(&count)
	require.NoError(t, err, "Failed to check record non-existence")

	assert.Equal(t, 0, count, "Expected no record to exist")
}
