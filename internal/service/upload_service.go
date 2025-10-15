package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UploadService defines the interface for file upload operations
type UploadService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, directory string) (string, error)
	DeleteFile(ctx context.Context, fileURL string) error
	GetFileURL(ctx context.Context, path string) string
	ValidateFile(file *multipart.FileHeader, allowedTypes []string, maxSize int64) error
	CalculateChecksum(file *multipart.FileHeader) (string, error)
}

// uploadService implements file upload functionality
type uploadService struct {
	storageProvider string
	uploadPath      string
	baseURL         string
}

// UploadServiceConfig holds configuration for upload service
type UploadServiceConfig struct {
	StorageProvider string // "local", "s3", "cloudinary"
	UploadPath      string
	BaseURL         string // Base URL for serving files
}

// NewUploadService creates a new upload service instance
func NewUploadService(config UploadServiceConfig) UploadService {
	return &uploadService{
		storageProvider: config.StorageProvider,
		uploadPath:      config.UploadPath,
		baseURL:         config.BaseURL,
	}
}

// UploadFile uploads a file to the configured storage
func (s *uploadService) UploadFile(ctx context.Context, file *multipart.FileHeader, directory string) (string, error) {
	// For now, only implement local storage
	// S3 and Cloudinary bisa ditambahkan nanti
	switch s.storageProvider {
	case "local":
		return s.uploadToLocal(ctx, file, directory)
	case "s3":
		// TODO: Implement S3 upload
		return "", fmt.Errorf("S3 storage not yet implemented")
	case "cloudinary":
		// TODO: Implement Cloudinary upload
		return "", fmt.Errorf("Cloudinary storage not yet implemented")
	default:
		return "", fmt.Errorf("unsupported storage provider: %s", s.storageProvider)
	}
}

// uploadToLocal uploads file to local filesystem
func (s *uploadService) uploadToLocal(ctx context.Context, file *multipart.FileHeader, directory string) (string, error) {
	// Create unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s", uuid.New().String(), time.Now().Format("20060102150405"), ext)

	// Create full path
	fullPath := filepath.Join(s.uploadPath, directory)

	// Ensure directory exists
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Full file path
	filePath := filepath.Join(fullPath, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return relative URL path
	relativePath := filepath.Join(directory, filename)
	return s.GetFileURL(ctx, relativePath), nil
}

// DeleteFile deletes a file from storage
func (s *uploadService) DeleteFile(ctx context.Context, fileURL string) error {
	if fileURL == "" {
		return nil
	}

	switch s.storageProvider {
	case "local":
		return s.deleteFromLocal(ctx, fileURL)
	case "s3":
		// TODO: Implement S3 deletion
		return fmt.Errorf("S3 storage not yet implemented")
	case "cloudinary":
		// TODO: Implement Cloudinary deletion
		return fmt.Errorf("Cloudinary storage not yet implemented")
	default:
		return fmt.Errorf("unsupported storage provider: %s", s.storageProvider)
	}
}

// deleteFromLocal deletes file from local filesystem
func (s *uploadService) deleteFromLocal(ctx context.Context, fileURL string) error {
	// Extract path from URL (remove base URL if present)
	path := strings.TrimPrefix(fileURL, s.baseURL)
	path = strings.TrimPrefix(path, "/")

	// Full file path
	filePath := filepath.Join(s.uploadPath, path)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File doesn't exist, not an error
		return nil
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetFileURL generates the full URL for a file path
func (s *uploadService) GetFileURL(ctx context.Context, path string) string {
	// For local storage, return URL relative to base URL
	if s.storageProvider == "local" {
		// Normalize path separators
		path = filepath.ToSlash(path)
		return fmt.Sprintf("%s/%s", strings.TrimSuffix(s.baseURL, "/"), strings.TrimPrefix(path, "/"))
	}

	// For S3/Cloudinary, return the full URL (to be implemented)
	return path
}

// ValidateFile validates file type and size
func (s *uploadService) ValidateFile(file *multipart.FileHeader, allowedTypes []string, maxSize int64) error {
	// Check file size
	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxSize)
	}

	// Check file type if allowedTypes is specified
	if len(allowedTypes) > 0 {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowed := false
		for _, allowedType := range allowedTypes {
			if ext == strings.ToLower(allowedType) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("file type %s is not allowed. Allowed types: %v", ext, allowedTypes)
		}
	}

	return nil
}

// CalculateChecksum calculates SHA256 checksum of a file
func (s *uploadService) CalculateChecksum(file *multipart.FileHeader) (string, error) {
	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Calculate SHA256 hash
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	// Return hex encoded hash
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// File type constants for validation
var (
	ImageTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg"}
	DocumentTypes = []string{".pdf", ".doc", ".docx", ".txt", ".rtf"}
	VideoTypes = []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm"}
	AllMediaTypes = append(append(ImageTypes, DocumentTypes...), VideoTypes...)
)

// File size constants (in bytes)
const (
	MaxAvatarSize    = 5 * 1024 * 1024   // 5 MB
	MaxCoverSize     = 10 * 1024 * 1024  // 10 MB
	MaxDocumentSize  = 20 * 1024 * 1024  // 20 MB
	MaxVideoSize     = 100 * 1024 * 1024 // 100 MB
)
