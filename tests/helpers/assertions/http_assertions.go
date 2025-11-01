package assertions

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HTTPAssertion provides HTTP-specific assertions
type HTTPAssertion struct {
	T *testing.T
}

// NewHTTPAssertion creates a new HTTP assertion helper
func NewHTTPAssertion(t *testing.T) *HTTPAssertion {
	return &HTTPAssertion{T: t}
}

// AssertStatusCode asserts the HTTP status code
func (ha *HTTPAssertion) AssertStatusCode(resp *http.Response, expectedStatus int, msgAndArgs ...interface{}) {
	ha.T.Helper()
	assert.Equal(ha.T, expectedStatus, resp.StatusCode, msgAndArgs...)
}

// AssertStatusOK asserts status code is 200
func (ha *HTTPAssertion) AssertStatusOK(resp *http.Response, msgAndArgs ...interface{}) {
	ha.AssertStatusCode(resp, http.StatusOK, msgAndArgs...)
}

// AssertStatusCreated asserts status code is 201
func (ha *HTTPAssertion) AssertStatusCreated(resp *http.Response, msgAndArgs ...interface{}) {
	ha.AssertStatusCode(resp, http.StatusCreated, msgAndArgs...)
}

// AssertStatusBadRequest asserts status code is 400
func (ha *HTTPAssertion) AssertStatusBadRequest(resp *http.Response, msgAndArgs ...interface{}) {
	ha.AssertStatusCode(resp, http.StatusBadRequest, msgAndArgs...)
}

// AssertStatusUnauthorized asserts status code is 401
func (ha *HTTPAssertion) AssertStatusUnauthorized(resp *http.Response, msgAndArgs ...interface{}) {
	ha.AssertStatusCode(resp, http.StatusUnauthorized, msgAndArgs...)
}

// AssertStatusForbidden asserts status code is 403
func (ha *HTTPAssertion) AssertStatusForbidden(resp *http.Response, msgAndArgs ...interface{}) {
	ha.AssertStatusCode(resp, http.StatusForbidden, msgAndArgs...)
}

// AssertStatusNotFound asserts status code is 404
func (ha *HTTPAssertion) AssertStatusNotFound(resp *http.Response, msgAndArgs ...interface{}) {
	ha.AssertStatusCode(resp, http.StatusNotFound, msgAndArgs...)
}

// AssertStatusConflict asserts status code is 409
func (ha *HTTPAssertion) AssertStatusConflict(resp *http.Response, msgAndArgs ...interface{}) {
	ha.AssertStatusCode(resp, http.StatusConflict, msgAndArgs...)
}

// AssertStatusInternalServerError asserts status code is 500
func (ha *HTTPAssertion) AssertStatusInternalServerError(resp *http.Response, msgAndArgs ...interface{}) {
	ha.AssertStatusCode(resp, http.StatusInternalServerError, msgAndArgs...)
}

// AssertHeader asserts a response header value
func (ha *HTTPAssertion) AssertHeader(resp *http.Response, key, expectedValue string, msgAndArgs ...interface{}) {
	ha.T.Helper()
	actualValue := resp.Header.Get(key)
	assert.Equal(ha.T, expectedValue, actualValue, msgAndArgs...)
}

// AssertHeaderExists asserts a response header exists
func (ha *HTTPAssertion) AssertHeaderExists(resp *http.Response, key string, msgAndArgs ...interface{}) {
	ha.T.Helper()
	_, exists := resp.Header[key]
	assert.True(ha.T, exists, msgAndArgs...)
}

// AssertContentType asserts the Content-Type header
func (ha *HTTPAssertion) AssertContentType(resp *http.Response, expectedContentType string, msgAndArgs ...interface{}) {
	ha.AssertHeader(resp, "Content-Type", expectedContentType, msgAndArgs...)
}

// AssertContentTypeJSON asserts Content-Type is application/json
func (ha *HTTPAssertion) AssertContentTypeJSON(resp *http.Response, msgAndArgs ...interface{}) {
	ha.T.Helper()
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(ha.T, contentType, "application/json", msgAndArgs...)
}

// AssertJSONResponse asserts the response body is valid JSON
func (ha *HTTPAssertion) AssertJSONResponse(resp *http.Response, msgAndArgs ...interface{}) {
	ha.T.Helper()
	body, err := io.ReadAll(resp.Body)
	require.NoError(ha.T, err, "Failed to read response body")

	var js interface{}
	err = json.Unmarshal(body, &js)
	assert.NoError(ha.T, err, msgAndArgs...)
}

// AssertJSONField asserts a specific field in JSON response
func (ha *HTTPAssertion) AssertJSONField(resp *http.Response, fieldPath string, expectedValue interface{}, msgAndArgs ...interface{}) {
	ha.T.Helper()
	body, err := io.ReadAll(resp.Body)
	require.NoError(ha.T, err, "Failed to read response body")

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	require.NoError(ha.T, err, "Failed to parse JSON response")

	actualValue := data[fieldPath]
	assert.Equal(ha.T, expectedValue, actualValue, msgAndArgs...)
}

// AssertJSONFieldExists asserts a field exists in JSON response
func (ha *HTTPAssertion) AssertJSONFieldExists(resp *http.Response, fieldPath string, msgAndArgs ...interface{}) {
	ha.T.Helper()
	body, err := io.ReadAll(resp.Body)
	require.NoError(ha.T, err, "Failed to read response body")

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	require.NoError(ha.T, err, "Failed to parse JSON response")

	_, exists := data[fieldPath]
	assert.True(ha.T, exists, msgAndArgs...)
}

// AssertResponseBody asserts the response body content
func (ha *HTTPAssertion) AssertResponseBody(resp *http.Response, expectedBody string, msgAndArgs ...interface{}) {
	ha.T.Helper()
	body, err := io.ReadAll(resp.Body)
	require.NoError(ha.T, err, "Failed to read response body")

	assert.Equal(ha.T, expectedBody, string(body), msgAndArgs...)
}

// AssertResponseBodyContains asserts the response body contains a string
func (ha *HTTPAssertion) AssertResponseBodyContains(resp *http.Response, substring string, msgAndArgs ...interface{}) {
	ha.T.Helper()
	body, err := io.ReadAll(resp.Body)
	require.NoError(ha.T, err, "Failed to read response body")

	assert.Contains(ha.T, string(body), substring, msgAndArgs...)
}

// AssertEmptyBody asserts the response body is empty
func (ha *HTTPAssertion) AssertEmptyBody(resp *http.Response, msgAndArgs ...interface{}) {
	ha.T.Helper()
	body, err := io.ReadAll(resp.Body)
	require.NoError(ha.T, err, "Failed to read response body")

	assert.Empty(ha.T, body, msgAndArgs...)
}

// ParseJSONResponse parses JSON response into a struct
func (ha *HTTPAssertion) ParseJSONResponse(resp *http.Response, v interface{}) {
	ha.T.Helper()
	body, err := io.ReadAll(resp.Body)
	require.NoError(ha.T, err, "Failed to read response body")

	err = json.Unmarshal(body, v)
	require.NoError(ha.T, err, "Failed to parse JSON response")
}

// AssertSuccessResponse asserts a successful API response structure
func (ha *HTTPAssertion) AssertSuccessResponse(resp *http.Response, msgAndArgs ...interface{}) {
	ha.T.Helper()
	ha.AssertStatusOK(resp)
	ha.AssertContentTypeJSON(resp)
	ha.AssertJSONFieldExists(resp, "success")
	ha.AssertJSONField(resp, "success", true, msgAndArgs...)
}

// AssertErrorResponse asserts an error API response structure
func (ha *HTTPAssertion) AssertErrorResponse(resp *http.Response, expectedStatus int, msgAndArgs ...interface{}) {
	ha.T.Helper()
	ha.AssertStatusCode(resp, expectedStatus)
	ha.AssertContentTypeJSON(resp)
	ha.AssertJSONFieldExists(resp, "error")
}

// AssertPaginatedResponse asserts a paginated response structure
func (ha *HTTPAssertion) AssertPaginatedResponse(resp *http.Response, msgAndArgs ...interface{}) {
	ha.T.Helper()
	ha.AssertStatusOK(resp)
	ha.AssertContentTypeJSON(resp)
	ha.AssertJSONFieldExists(resp, "data")
	ha.AssertJSONFieldExists(resp, "pagination")
}

// AssertValidationError asserts a validation error response
func (ha *HTTPAssertion) AssertValidationError(resp *http.Response, fieldName string, msgAndArgs ...interface{}) {
	ha.T.Helper()
	ha.AssertStatusBadRequest(resp)
	ha.AssertContentTypeJSON(resp)

	body, err := io.ReadAll(resp.Body)
	require.NoError(ha.T, err, "Failed to read response body")

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	require.NoError(ha.T, err, "Failed to parse JSON response")

	errors, ok := data["errors"].(map[string]interface{})
	require.True(ha.T, ok, "Response should contain errors field")

	_, exists := errors[fieldName]
	assert.True(ha.T, exists, "Validation error should exist for field: %s", fieldName)
}

// Response structure helpers
type APIResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    interface{}            `json:"data"`
	Error   string                 `json:"error"`
	Errors  map[string]interface{} `json:"errors"`
}

type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
}

// AssertAPIResponse parses and asserts API response structure
func (ha *HTTPAssertion) AssertAPIResponse(resp *http.Response, expectedSuccess bool) *APIResponse {
	ha.T.Helper()

	var apiResp APIResponse
	ha.ParseJSONResponse(resp, &apiResp)

	assert.Equal(ha.T, expectedSuccess, apiResp.Success)

	return &apiResp
}

// Quick assertion functions without creating instance

// AssertStatus asserts HTTP status code
func AssertStatus(t *testing.T, resp *http.Response, expectedStatus int, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, expectedStatus, resp.StatusCode, msgAndArgs...)
}

// AssertJSONContentType asserts Content-Type is JSON
func AssertJSONContentType(t *testing.T, resp *http.Response, msgAndArgs ...interface{}) {
	t.Helper()
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json", msgAndArgs...)
}

// AssertValidJSON asserts response body is valid JSON
func AssertValidJSON(t *testing.T, resp *http.Response, msgAndArgs ...interface{}) {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	var js interface{}
	err = json.Unmarshal(body, &js)
	assert.NoError(t, err, msgAndArgs...)
}

// AssertBodyContains asserts response body contains substring
func AssertBodyContains(t *testing.T, resp *http.Response, substring string, msgAndArgs ...interface{}) {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	assert.Contains(t, string(body), substring, msgAndArgs...)
}
