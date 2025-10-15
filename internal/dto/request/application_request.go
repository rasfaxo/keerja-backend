package request

import (
	"keerja-backend/internal/dto"
	"time"
)

// ApplyJobRequest represents apply for job request
type ApplyJobRequest struct {
	JobID               int64                      `json:"job_id" validate:"required"`
	ResumeURL           string                     `json:"resume_url" validate:"required,url"`
	CoverLetter         string                     `json:"cover_letter" validate:"max=2000"`
	ExpectedSalary      *float64                   `json:"expected_salary" validate:"omitempty,gt=0"`
	Currency            string                     `json:"currency" validate:"required_with=ExpectedSalary,len=3"`
	AvailableStartDate  *time.Time                 `json:"available_start_date" validate:"omitempty"`
	Answers             []ApplicationAnswerRequest `json:"answers" validate:"omitempty,dive"`
	AdditionalDocuments []string                   `json:"additional_documents" validate:"omitempty,dive,url"`
}

// ApplicationAnswerRequest represents answer to screening question
type ApplicationAnswerRequest struct {
	QuestionID int64  `json:"question_id" validate:"required"`
	Answer     string `json:"answer" validate:"required,max=1000"`
}

// WithdrawApplicationRequest represents withdraw application request
type WithdrawApplicationRequest struct {
	Reason string `json:"reason" validate:"required,min=10,max=500"`
}

// UpdateApplicationStatusRequest represents update application status request
type UpdateApplicationStatusRequest struct {
	Status          string `json:"status" validate:"required,oneof=pending screening shortlisted interview offered rejected withdrawn"`
	Notes           string `json:"notes" validate:"max=1000"`
	RejectionReason string `json:"rejection_reason" validate:"omitempty,max=500"`
	HandledBy       *int64 `json:"handled_by" validate:"omitempty"`
}

// UpdateApplicationStageRequest represents update application stage request
type UpdateApplicationStageRequest struct {
	StageName   string `json:"stage_name" validate:"required,min=2,max=100"`
	StageOrder  int16  `json:"stage_order" validate:"required,min=1"`
	Description string `json:"description" validate:"max=500"`
	Notes       string `json:"notes" validate:"max=1000"`
	HandledBy   *int64 `json:"handled_by" validate:"omitempty"`
}

// AddApplicationNoteRequest represents add note to application request
type AddApplicationNoteRequest struct {
	NoteText   string `json:"note_text" validate:"required,min=5,max=2000"`
	IsInternal bool   `json:"is_internal"` // Internal notes visible only to company
}

// UpdateApplicationNoteRequest represents update note request
type UpdateApplicationNoteRequest struct {
	NoteText   string `json:"note_text" validate:"required,min=5,max=2000"`
	IsInternal bool   `json:"is_internal"`
}

// ScheduleInterviewRequest represents schedule interview request
type ScheduleInterviewRequest struct {
	ApplicationID   int64     `json:"application_id" validate:"required"`
	InterviewType   string    `json:"interview_type" validate:"required,oneof=phone video onsite technical hr final"`
	InterviewStage  string    `json:"interview_stage" validate:"required,max=100"`
	ScheduledAt     time.Time `json:"scheduled_at" validate:"required,gtfield=Now"`
	Duration        int16     `json:"duration" validate:"required,min=15,max=480"` // in minutes
	Location        string    `json:"location" validate:"omitempty,max=200"`
	MeetingURL      string    `json:"meeting_url" validate:"omitempty,url"`
	Interviewers    []int64   `json:"interviewers" validate:"omitempty,dive,gt=0"`
	Notes           string    `json:"notes" validate:"max=1000"`
	NotifyCandidate bool      `json:"notify_candidate"`
}

// UpdateInterviewRequest represents update interview request
type UpdateInterviewRequest struct {
	InterviewType   *string    `json:"interview_type" validate:"omitempty,oneof=phone video onsite technical hr final"`
	InterviewStage  *string    `json:"interview_stage" validate:"omitempty,max=100"`
	ScheduledAt     *time.Time `json:"scheduled_at" validate:"omitempty"`
	Duration        *int16     `json:"duration" validate:"omitempty,min=15,max=480"`
	Location        *string    `json:"location" validate:"omitempty,max=200"`
	MeetingURL      *string    `json:"meeting_url" validate:"omitempty,url"`
	Status          *string    `json:"status" validate:"omitempty,oneof=scheduled rescheduled completed cancelled no_show"`
	Notes           *string    `json:"notes" validate:"omitempty,max=1000"`
	NotifyCandidate *bool      `json:"notify_candidate"`
}

// CompleteInterviewRequest represents complete interview with feedback request
type CompleteInterviewRequest struct {
	Status         string `json:"status" validate:"required,oneof=completed cancelled no_show"`
	Rating         *int16 `json:"rating" validate:"omitempty,min=1,max=5"`
	Feedback       string `json:"feedback" validate:"max=2000"`
	Result         string `json:"result" validate:"required,oneof=pass fail pending"`
	Recommendation string `json:"recommendation" validate:"max=1000"`
	ConductedBy    *int64 `json:"conducted_by" validate:"omitempty"`
}

// RescheduleInterviewRequest represents reschedule interview request
type RescheduleInterviewRequest struct {
	ScheduledAt     time.Time `json:"scheduled_at" validate:"required,gtfield=Now"`
	Duration        int16     `json:"duration" validate:"required,min=15,max=480"`
	Reason          string    `json:"reason" validate:"required,min=10,max=500"`
	Location        string    `json:"location" validate:"omitempty,max=200"`
	MeetingURL      string    `json:"meeting_url" validate:"omitempty,url"`
	NotifyCandidate bool      `json:"notify_candidate"`
}

// BulkUpdateApplicationsRequest represents bulk update applications request
type BulkUpdateApplicationsRequest struct {
	ApplicationIDs []int64 `json:"application_ids" validate:"required,min=1,dive,gt=0"`
	Status         string  `json:"status" validate:"required,oneof=pending screening shortlisted interview offered rejected withdrawn"`
	Notes          string  `json:"notes" validate:"max=1000"`
}

// ApplicationSearchRequest represents search applications request
type ApplicationSearchRequest struct {
	Query           string     `json:"query" validate:"omitempty,max=100"`
	JobID           *int64     `json:"job_id" validate:"omitempty"`
	CompanyID       *int64     `json:"company_id" validate:"omitempty"`
	UserID          *int64     `json:"user_id" validate:"omitempty"`
	Status          string     `json:"status" validate:"omitempty,oneof=pending screening shortlisted interview offered rejected withdrawn"`
	AppliedFrom     *time.Time `json:"applied_from" validate:"omitempty"`
	AppliedTo       *time.Time `json:"applied_to" validate:"omitempty,gtefield=AppliedFrom"`
	MinSalary       *float64   `json:"min_salary" validate:"omitempty,gt=0"`
	MaxSalary       *float64   `json:"max_salary" validate:"omitempty,gtefield=MinSalary"`
	HasInterview    *bool      `json:"has_interview"`
	InterviewStatus string     `json:"interview_status" validate:"omitempty,oneof=scheduled rescheduled completed cancelled no_show"`
	SortBy          string     `json:"sort_by" validate:"omitempty,oneof=created_at updated_at applied_date status"`
	SortOrder       string     `json:"sort_order" validate:"omitempty,oneof=asc desc"`
	dto.PaginationRequest
}

// ApplicationFilterRequest represents filter applications for employer
type ApplicationFilterRequest struct {
	JobID         *int64   `json:"job_id" validate:"omitempty"`
	CompanyID     *int64   `json:"company_id" validate:"omitempty"`
	Status        string   `json:"status" validate:"omitempty,oneof=pending screening shortlisted interview offered rejected withdrawn"`
	IsViewed      *bool    `json:"is_viewed"`
	HasNotes      *bool    `json:"has_notes"`
	HasInterview  *bool    `json:"has_interview"`
	MinMatchScore *float64 `json:"min_match_score" validate:"omitempty,min=0,max=100"`
	SortBy        string   `json:"sort_by" validate:"omitempty,oneof=created_at match_score status"`
	SortOrder     string   `json:"sort_order" validate:"omitempty,oneof=asc desc"`
	dto.PaginationRequest
}

// UploadApplicationDocumentRequest represents upload document request
type UploadApplicationDocumentRequest struct {
	DocumentType string `json:"document_type" validate:"required,oneof=resume cover_letter certificate portfolio other"`
	DocumentURL  string `json:"document_url" validate:"required,url"`
	FileName     string `json:"file_name" validate:"required,max=255"`
	FileSize     int64  `json:"file_size" validate:"required,gt=0"`
	Description  string `json:"description" validate:"max=500"`
}

// RateApplicationExperienceRequest represents rate application experience request
type RateApplicationExperienceRequest struct {
	Rating   int16  `json:"rating" validate:"required,min=1,max=5"`
	Feedback string `json:"feedback" validate:"max=1000"`
}
