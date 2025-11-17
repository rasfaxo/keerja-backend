package http

// Error message constants for consistent error handling across handlers
const (
	// Generic errors
	ErrInvalidID        = "Invalid ID parameter"
	ErrInvalidRequest   = "Invalid request body"
	ErrInvalidBody      = "Invalid request body format"
	ErrInvalidJSON      = "Invalid JSON format"
	ErrValidationFailed = "Validation failed"
	ErrInternalServer   = "Internal server error"
	ErrNotFound         = "Resource not found"
	ErrUnauthorized     = "Unauthorized access"
	ErrForbidden        = "Access forbidden"
	ErrConflict         = "Resource already exists"

	// Query parameter errors
	ErrInvalidQueryParams = "Invalid query parameters"
	ErrInvalidPage        = "Invalid page parameter"
	ErrInvalidLimit       = "Invalid limit parameter"
	ErrInvalidSortBy      = "Invalid sort_by parameter"
	ErrInvalidSortOrder   = "Invalid sort_order parameter"
	ErrInvalidDateFormat  = "Invalid date format. Use RFC3339 format"
	ErrInvalidBoolParam   = "Invalid boolean parameter"

	// File upload errors
	ErrNoFileUploaded   = "No file uploaded"
	ErrFileTooLarge     = "File size exceeds maximum allowed"
	ErrInvalidFileType  = "File type not allowed"
	ErrInvalidFileExt   = "File extension not allowed"
	ErrFileUploadFailed = "Failed to upload file"

	// Authentication errors
	ErrMissingToken      = "Missing authentication token"
	ErrInvalidToken      = "Invalid or expired token"
	ErrTokenExpired      = "Token has expired"
	ErrInsufficientPerms = "Insufficient permissions"

	// User errors
	ErrUserNotFound       = "User not found"
	ErrEmailAlreadyExists = "Email already registered"
	ErrInvalidCredentials = "Invalid email or password"
	ErrAccountNotVerified = "Account not verified"
	ErrAccountSuspended   = "Account has been suspended"

	// Job errors
	ErrJobNotFound       = "Job not found"
	ErrJobExpired        = "Job posting has expired"
	ErrJobClosed         = "Job posting is closed"
	ErrNotJobOwner       = "You are not the owner of this job"
	ErrCannotApplyOwnJob = "Cannot apply to your own job posting"

	// Application errors
	ErrApplicationNotFound     = "Application not found"
	ErrAlreadyApplied          = "You have already applied to this job"
	ErrApplicationClosed       = "Application period has closed"
	ErrCannotWithdraw          = "Cannot withdraw application at this stage"
	ErrInvalidApplicationStage = "Invalid application stage"

	// Company errors
	ErrInvalidCompanyID   = "Invalid company ID"
	ErrCompanyNotFound    = "Company not found"
	ErrNotCompanyMember   = "You are not a member of this company"
	ErrCompanyNotVerified = "Company is not verified"
	ErrAlreadyFollowing   = "Already following this company"
	ErrNotFollowing       = "Not following this company"
	ErrFailedOperation    = "Operation failed. Please try again"

	// Review errors
	ErrReviewNotFound  = "Review not found"
	ErrAlreadyReviewed = "You have already reviewed this company"
	ErrCannotReviewOwn = "Cannot review your own company"
	ErrInvalidRating   = "Invalid rating value"

	// Interview errors
	ErrInterviewNotFound = "Interview not found"
	ErrInterviewPast     = "Interview date is in the past"
	ErrInterviewConflict = "Interview time conflicts with another interview"

	// Input validation errors
	ErrMissingRequiredField = "Missing required field"
	ErrInvalidEmail         = "Invalid email format"
	ErrInvalidPhone         = "Invalid phone number format"
	ErrInvalidURL           = "Invalid URL format"
	ErrPasswordTooShort     = "Password must be at least 8 characters"
	ErrPasswordTooWeak      = "Password must contain uppercase, lowercase, and numbers"

	// Business logic errors
	ErrExceedsMaxLimit    = "Exceeds maximum allowed limit"
	ErrBelowMinLimit      = "Below minimum required value"
	ErrInvalidDateRange   = "Invalid date range"
	ErrFutureDateRequired = "Date must be in the future"
	ErrPastDateRequired   = "Date must be in the past"

	// Rate limiting errors
	ErrRateLimitExceeded = "Rate limit exceeded. Please try again later"
	ErrTooManyRequests   = "Too many requests. Please slow down"

	// Security errors
	ErrSuspiciousActivity = "Suspicious activity detected"
	ErrInvalidInput       = "Input contains invalid characters"
	ErrPotentialXSS       = "Input contains potentially harmful content"
	ErrPotentialSQLi      = "Input contains potentially harmful SQL patterns"
)

// Success message constants
const (
	MsgCreatedSuccess    = "Resource created successfully"
	MsgUpdatedSuccess    = "Resource updated successfully"
	MsgDeletedSuccess    = "Resource deleted successfully"
	MsgFetchedSuccess    = "Resource retrieved successfully"
	MsgOperationSuccess  = "Operation completed successfully"
	MsgEmailSent         = "Email sent successfully"
	MsgUploadSuccess     = "File uploaded successfully"
	MsgApplicationSubmit = "Application submitted successfully"
	MsgStatusUpdated     = "Status updated successfully"
)

// GetValidationError returns user-friendly validation error message
func GetValidationError(field string, tag string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "url":
		return field + " must be a valid URL"
	case "min":
		return field + " is too short"
	case "max":
		return field + " is too long"
	case "gte":
		return field + " must be greater than or equal to specified value"
	case "lte":
		return field + " must be less than or equal to specified value"
	case "oneof":
		return field + " must be one of the allowed values"
	case "uuid":
		return field + " must be a valid UUID"
	case "numeric":
		return field + " must be numeric"
	case "alpha":
		return field + " must contain only letters"
	case "alphanum":
		return field + " must contain only letters and numbers"
	default:
		return field + " is invalid"
	}
}
