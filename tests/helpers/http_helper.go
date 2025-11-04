package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// HTTPTestHelper provides utilities for HTTP testing
type HTTPTestHelper struct {
	App *fiber.App
	T   *testing.T
}

// NewHTTPTestHelper creates a new HTTP test helper
func NewHTTPTestHelper(t *testing.T, app *fiber.App) *HTTPTestHelper {
	return &HTTPTestHelper{
		App: app,
		T:   t,
	}
}

// Request represents an HTTP request for testing
type Request struct {
	Method      string
	URL         string
	Body        interface{}
	Headers     map[string]string
	QueryParams map[string]string
	FormData    map[string]string
	Files       map[string]string
}

// Response represents an HTTP response for testing
type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

// MakeRequest makes an HTTP request and returns the response
func (h *HTTPTestHelper) MakeRequest(req Request) *Response {
	h.T.Helper()

	var body io.Reader
	contentType := "application/json"

	// Handle different request body types
	if len(req.Files) > 0 || len(req.FormData) > 0 {
		// Multipart form data
		bodyBuf, ct, err := h.createMultipartBody(req.FormData, req.Files)
		require.NoError(h.T, err, "Failed to create multipart body")
		body = bodyBuf
		contentType = ct
	} else if req.Body != nil {
		// JSON body
		jsonBody, err := json.Marshal(req.Body)
		require.NoError(h.T, err, "Failed to marshal request body")
		body = bytes.NewBuffer(jsonBody)
	}

	// Create HTTP request
	httpReq := httptest.NewRequest(req.Method, req.URL, body)

	// Set content type
	httpReq.Header.Set("Content-Type", contentType)

	// Set custom headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Set query parameters
	if len(req.QueryParams) > 0 {
		q := httpReq.URL.Query()
		for key, value := range req.QueryParams {
			q.Add(key, value)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	// Execute request
	resp, err := h.App.Test(httpReq, -1)
	require.NoError(h.T, err, "Failed to execute request")

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(h.T, err, "Failed to read response body")
	resp.Body.Close()

	// Extract response headers
	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    headers,
	}
}

// createMultipartBody creates a multipart form body with files
func (h *HTTPTestHelper) createMultipartBody(formData map[string]string, files map[string]string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range formData {
		err := writer.WriteField(key, value)
		if err != nil {
			return nil, "", err
		}
	}

	// Add files
	for fieldName, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, "", err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
		if err != nil {
			return nil, "", err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, "", err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

// ParseResponseJSON parses JSON response body into a struct
func (h *HTTPTestHelper) ParseResponseJSON(resp *Response, v interface{}) {
	h.T.Helper()
	err := json.Unmarshal(resp.Body, v)
	require.NoError(h.T, err, "Failed to parse response JSON")
}

// MakeJSONRequest makes a JSON request
func (h *HTTPTestHelper) MakeJSONRequest(method, url string, body interface{}, headers map[string]string) *Response {
	return h.MakeRequest(Request{
		Method:  method,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

// MakeAuthenticatedRequest makes an authenticated request with JWT token
func (h *HTTPTestHelper) MakeAuthenticatedRequest(method, url string, token string, body interface{}) *Response {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	return h.MakeJSONRequest(method, url, body, headers)
}

// MakeMultipartRequest makes a multipart form request with files
func (h *HTTPTestHelper) MakeMultipartRequest(method, url string, formData map[string]string, files map[string]string, headers map[string]string) *Response {
	return h.MakeRequest(Request{
		Method:   method,
		URL:      url,
		FormData: formData,
		Files:    files,
		Headers:  headers,
	})
}

// RequestBuilder provides a fluent interface for building requests
type RequestBuilder struct {
	request Request
	helper  *HTTPTestHelper
}

// NewRequestBuilder creates a new request builder
func NewRequestBuilder(helper *HTTPTestHelper) *RequestBuilder {
	return &RequestBuilder{
		request: Request{
			Headers:     make(map[string]string),
			QueryParams: make(map[string]string),
			FormData:    make(map[string]string),
			Files:       make(map[string]string),
		},
		helper: helper,
	}
}

// Method sets the HTTP method
func (rb *RequestBuilder) Method(method string) *RequestBuilder {
	rb.request.Method = method
	return rb
}

// GET sets the method to GET
func (rb *RequestBuilder) GET() *RequestBuilder {
	return rb.Method("GET")
}

// POST sets the method to POST
func (rb *RequestBuilder) POST() *RequestBuilder {
	return rb.Method("POST")
}

// PUT sets the method to PUT
func (rb *RequestBuilder) PUT() *RequestBuilder {
	return rb.Method("PUT")
}

// DELETE sets the method to DELETE
func (rb *RequestBuilder) DELETE() *RequestBuilder {
	return rb.Method("DELETE")
}

// PATCH sets the method to PATCH
func (rb *RequestBuilder) PATCH() *RequestBuilder {
	return rb.Method("PATCH")
}

// URL sets the request URL
func (rb *RequestBuilder) URL(url string) *RequestBuilder {
	rb.request.URL = url
	return rb
}

// Body sets the request body
func (rb *RequestBuilder) Body(body interface{}) *RequestBuilder {
	rb.request.Body = body
	return rb
}

// Header sets a request header
func (rb *RequestBuilder) Header(key, value string) *RequestBuilder {
	rb.request.Headers[key] = value
	return rb
}

// WithAuth sets the Authorization header
func (rb *RequestBuilder) WithAuth(token string) *RequestBuilder {
	return rb.Header("Authorization", "Bearer "+token)
}

// QueryParam sets a query parameter
func (rb *RequestBuilder) QueryParam(key, value string) *RequestBuilder {
	rb.request.QueryParams[key] = value
	return rb
}

// FormField sets a form field
func (rb *RequestBuilder) FormField(key, value string) *RequestBuilder {
	rb.request.FormData[key] = value
	return rb
}

// File sets a file for upload
func (rb *RequestBuilder) File(fieldName, filePath string) *RequestBuilder {
	rb.request.Files[fieldName] = filePath
	return rb
}

// Send executes the request
func (rb *RequestBuilder) Send() *Response {
	return rb.helper.MakeRequest(rb.request)
}

// SendAndParse executes the request and parses the JSON response
func (rb *RequestBuilder) SendAndParse(v interface{}) *Response {
	resp := rb.Send()
	rb.helper.ParseResponseJSON(resp, v)
	return resp
}

// CreateTestFile creates a temporary file for testing
func CreateTestFile(t *testing.T, content []byte, filename string) string {
	t.Helper()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, filename)

	err := os.WriteFile(filePath, content, 0644)
	require.NoError(t, err, "Failed to create test file")

	return filePath
}

// CreateTestImageFile creates a temporary image file for testing
func CreateTestImageFile(t *testing.T) string {
	t.Helper()

	// Create a simple 1x1 pixel PNG
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}

	return CreateTestFile(t, pngData, "test_image.png")
}

// CreateTestPDFFile creates a temporary PDF file for testing
func CreateTestPDFFile(t *testing.T) string {
	t.Helper()

	// Minimal PDF content
	pdfData := []byte(`%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj
2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj
3 0 obj
<<
/Type /Page
/Parent 2 0 R
/MediaBox [0 0 612 792]
>>
endobj
xref
0 4
0000000000 65535 f 
0000000009 00000 n 
0000000058 00000 n 
0000000115 00000 n 
trailer
<<
/Size 4
/Root 1 0 R
>>
startxref
190
%%EOF`)

	return CreateTestFile(t, pdfData, "test_document.pdf")
}
