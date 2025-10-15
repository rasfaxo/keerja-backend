package upload

import (
	"context"
	"io"
)

// UploadService defines the interface for file upload operations
type UploadService interface {
	// UploadFile uploads a file
	UploadFile(ctx context.Context, req *UploadFileRequest) (*FileUpload, error)

	// UploadImage uploads an image file
	UploadImage(ctx context.Context, req *UploadFileRequest) (*FileUpload, error)

	// UploadAvatar uploads user avatar
	UploadAvatar(ctx context.Context, userID int64, file io.Reader, filename string) (*FileUpload, error)

	// UploadResume uploads resume/CV
	UploadResume(ctx context.Context, userID int64, file io.Reader, filename string) (*FileUpload, error)

	// UploadCompanyLogo uploads company logo
	UploadCompanyLogo(ctx context.Context, userID int64, file io.Reader, filename string) (*FileUpload, error)

	// UploadDocument uploads a document
	UploadDocument(ctx context.Context, userID int64, file io.Reader, filename, category string) (*FileUpload, error)

	// GetFile retrieves file by ID
	GetFile(ctx context.Context, id int64) (*FileUpload, error)

	// GetFileByURL retrieves file by URL
	GetFileByURL(ctx context.Context, url string) (*FileUpload, error)

	// GetUserFiles retrieves all files uploaded by user
	GetUserFiles(ctx context.Context, userID int64, category string, page, limit int) ([]FileUpload, int64, error)

	// DeleteFile deletes a file
	DeleteFile(ctx context.Context, id, userID int64) error

	// VerifyFile marks file as verified
	VerifyFile(ctx context.Context, id, verifiedBy int64) error

	// GetUnverifiedFiles retrieves unverified files
	GetUnverifiedFiles(ctx context.Context, category string, page, limit int) ([]FileUpload, int64, error)

	// GenerateThumbnail generates thumbnail for image
	GenerateThumbnail(ctx context.Context, fileID int64) error

	// GetDownloadURL generates temporary download URL
	GetDownloadURL(ctx context.Context, fileID int64, expiresIn int) (string, error)

	// ValidateFile validates file before upload
	ValidateFile(ctx context.Context, filename string, size int64, category string) error

	// GetStorageStats retrieves storage statistics for user
	GetStorageStats(ctx context.Context, userID int64) (*StorageStats, error)
}

// UploadFileRequest represents file upload request
type UploadFileRequest struct {
	UserID       int64
	File         io.Reader
	FileName     string
	OriginalName string
	FileType     string
	FileSize     int64
	Category     string
	IsPublic     bool
	Metadata     map[string]interface{}
}

// StorageStats represents storage statistics
type StorageStats struct {
	TotalFiles     int64   `json:"total_files"`
	TotalSize      int64   `json:"total_size"`       // in bytes
	TotalSizeInMB  float64 `json:"total_size_in_mb"` // in MB
	ImageCount     int64   `json:"image_count"`
	DocumentCount  int64   `json:"document_count"`
	ResumeCount    int64   `json:"resume_count"`
	StorageLimit   int64   `json:"storage_limit"`    // in bytes
	StorageUsedPct float64 `json:"storage_used_pct"` // percentage
}

// UploadRepository defines the interface for upload data operations
type UploadRepository interface {
	// Create creates a new file upload record
	Create(ctx context.Context, upload *FileUpload) error

	// FindByID finds file upload by ID
	FindByID(ctx context.Context, id int64) (*FileUpload, error)

	// FindByURL finds file upload by URL
	FindByURL(ctx context.Context, url string) (*FileUpload, error)

	// Update updates file upload record
	Update(ctx context.Context, upload *FileUpload) error

	// Delete deletes file upload record
	Delete(ctx context.Context, id int64) error

	// ListByUser lists files uploaded by user
	ListByUser(ctx context.Context, userID int64, category string, page, limit int) ([]FileUpload, int64, error)

	// ListByCategory lists files by category
	ListByCategory(ctx context.Context, category string, page, limit int) ([]FileUpload, int64, error)

	// GetUnverifiedFiles retrieves unverified files
	GetUnverifiedFiles(ctx context.Context, category string, page, limit int) ([]FileUpload, int64, error)

	// GetExpiredFiles retrieves expired files
	GetExpiredFiles(ctx context.Context, limit int) ([]FileUpload, error)

	// GetUserStorageStats retrieves user storage statistics
	GetUserStorageStats(ctx context.Context, userID int64) (*StorageStats, error)

	// CountByUser counts files by user
	CountByUser(ctx context.Context, userID int64) (int64, error)

	// GetTotalSizeByUser gets total file size by user
	GetTotalSizeByUser(ctx context.Context, userID int64) (int64, error)
}
