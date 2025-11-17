package upload

import "time"

// FileUpload represents an uploaded file
type FileUpload struct {
	ID           int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       int64      `json:"user_id" gorm:"not null;index"`
	FileName     string     `json:"file_name" gorm:"type:varchar(255);not null"`
	OriginalName string     `json:"original_name" gorm:"type:varchar(255);not null"`
	FileType     string     `json:"file_type" gorm:"type:varchar(100);not null"` // image/jpeg, application/pdf, etc.
	FileSize     int64      `json:"file_size" gorm:"not null"`                   // in bytes
	FileURL      string     `json:"file_url" gorm:"type:varchar(500);not null"`
	StoragePath  string     `json:"storage_path" gorm:"type:varchar(500);not null"`
	Provider     string     `json:"provider" gorm:"type:varchar(50);not null"` // local, s3, cloudinary, etc.
	Bucket       string     `json:"bucket" gorm:"type:varchar(100)"`
	Category     string     `json:"category" gorm:"type:varchar(50);not null;index"` // avatar, resume, document, company_logo, etc.
	IsPublic     bool       `json:"is_public" gorm:"default:false"`
	IsVerified   bool       `json:"is_verified" gorm:"default:false"`
	VerifiedAt   *time.Time `json:"verified_at"`
	VerifiedBy   *int64     `json:"verified_by"`
	Metadata     string     `json:"metadata" gorm:"type:json"` // Additional metadata as JSON
	DownloadURL  string     `json:"download_url" gorm:"type:varchar(500)"`
	ThumbnailURL string     `json:"thumbnail_url" gorm:"type:varchar(500)"`
	ExpiresAt    *time.Time `json:"expires_at"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:timestamp;default:now()"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:timestamp;default:now()"`
}

// TableName specifies the table name
func (FileUpload) TableName() string {
	return "file_uploads"
}

// IsImage checks if file is an image
func (f *FileUpload) IsImage() bool {
	return f.FileType == "image/jpeg" || f.FileType == "image/png" || f.FileType == "image/jpg" || f.FileType == "image/gif"
}

// IsPDF checks if file is a PDF
func (f *FileUpload) IsPDF() bool {
	return f.FileType == "application/pdf"
}

// IsDocument checks if file is a document
func (f *FileUpload) IsDocument() bool {
	return f.IsPDF() || f.FileType == "application/msword" || f.FileType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
}

// IsExpired checks if file has expired
func (f *FileUpload) IsExpired() bool {
	if f.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*f.ExpiresAt)
}

// CanDownload checks if file can be downloaded
func (f *FileUpload) CanDownload() bool {
	return !f.IsExpired()
}

// Verify marks file as verified
func (f *FileUpload) Verify(verifiedBy int64) {
	f.IsVerified = true
	now := time.Now()
	f.VerifiedAt = &now
	f.VerifiedBy = &verifiedBy
}

// GetSizeInKB returns file size in KB
func (f *FileUpload) GetSizeInKB() float64 {
	return float64(f.FileSize) / 1024
}

// GetSizeInMB returns file size in MB
func (f *FileUpload) GetSizeInMB() float64 {
	return float64(f.FileSize) / (1024 * 1024)
}
