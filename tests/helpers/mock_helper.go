package helpers

import (
	"errors"

	"github.com/stretchr/testify/mock"
)

// MockSetup is a helper function to setup mock expectations
type MockSetup func()

// MockTeardown is a helper function to teardown mocks
type MockTeardown func()

// MockHelper provides utilities for working with mocks
type MockHelper struct {
	mocks []interface{}
}

// NewMockHelper creates a new mock helper
func NewMockHelper() *MockHelper {
	return &MockHelper{
		mocks: make([]interface{}, 0),
	}
}

// Register registers a mock for tracking
func (m *MockHelper) Register(mockObj interface{}) {
	m.mocks = append(m.mocks, mockObj)
}

// AssertExpectations asserts all registered mocks
func (m *MockHelper) AssertExpectations(t mock.TestingT) {
	for _, mockObj := range m.mocks {
		if mockWithExpectations, ok := mockObj.(interface{ AssertExpectations(mock.TestingT) bool }); ok {
			mockWithExpectations.AssertExpectations(t)
		}
	}
}

// ResetMocks resets all registered mocks
func (m *MockHelper) ResetMocks() {
	for _, mockObj := range m.mocks {
		if mockWithReset, ok := mockObj.(interface{ ExpectedCalls() []*mock.Call }); ok {
			_ = mockWithReset
		}
	}
}

// mockTestingT is a simple implementation of mock.TestingT for helper functions
type mockTestingT struct{}

func (t *mockTestingT) Logf(format string, args ...interface{})   {}
func (t *mockTestingT) Errorf(format string, args ...interface{}) {}
func (t *mockTestingT) FailNow()                                  {}

// CommonMockSetup provides common mock setup patterns
type CommonMockSetup struct{}

// NewCommonMockSetup creates new common mock setup
func NewCommonMockSetup() *CommonMockSetup {
	return &CommonMockSetup{}
}

// SetupSuccess sets up mock for successful operation
func (c *CommonMockSetup) SetupSuccess(mockObj *mock.Mock, method string, args []interface{}, returns []interface{}) {
	call := mockObj.On(method, args...)
	call.Return(returns...)
}

// SetupError sets up mock for error operation
func (c *CommonMockSetup) SetupError(mockObj *mock.Mock, method string, args []interface{}, err error) {
	call := mockObj.On(method, args...)
	returns := make([]interface{}, len(args))
	for i := range returns {
		returns[i] = nil
	}
	returns = append(returns, err)
	call.Return(returns...)
}

// MockAnything is a helper to match any argument
func MockAnything() interface{} {
	return mock.Anything
}

// MockAnythingOfType is a helper to match any argument of a specific type
func MockAnythingOfType(t string) interface{} {
	return mock.AnythingOfType(t)
}

// MockMatchedBy is a helper to match argument by custom function
func MockMatchedBy(fn interface{}) interface{} {
	return mock.MatchedBy(fn)
}

// MockContextMatcher matches any context
func MockContextMatcher() interface{} {
	return mock.AnythingOfType("*context.emptyCtx")
}

// MockInt64Matcher matches any int64
func MockInt64Matcher() interface{} {
	return mock.AnythingOfType("int64")
}

// MockStringMatcher matches any string
func MockStringMatcher() interface{} {
	return mock.AnythingOfType("string")
}

// ExpectationBuilder helps build mock expectations fluently
type ExpectationBuilder struct {
	mockObj *mock.Mock
	method  string
	args    []interface{}
	returns []interface{}
	times   int
	maybe   bool
}

// NewExpectation creates a new expectation builder
func NewExpectation(mockObj *mock.Mock, method string) *ExpectationBuilder {
	return &ExpectationBuilder{
		mockObj: mockObj,
		method:  method,
		args:    []interface{}{},
		returns: []interface{}{},
		times:   1,
		maybe:   false,
	}
}

// WithArgs sets the expected arguments
func (e *ExpectationBuilder) WithArgs(args ...interface{}) *ExpectationBuilder {
	e.args = args
	return e
}

// Returns sets the return values
func (e *ExpectationBuilder) Returns(returns ...interface{}) *ExpectationBuilder {
	e.returns = returns
	return e
}

// Build builds and registers the expectation
func (e *ExpectationBuilder) Build() *mock.Call {
	call := e.mockObj.On(e.method, e.args...)
	call.Return(e.returns...)
	if e.maybe {
		call.Maybe()
	} else if e.times > 0 {
		call.Times(e.times)
	}
	return call
}

// Common test errors
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrValidation   = errors.New("validation error")
	ErrDatabase     = errors.New("database error")
)
