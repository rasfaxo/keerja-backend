package response

import "time"

// ApplicationResponse represents application basic response
type ApplicationResponse struct {
	ID               int64     `json:"id"`
	JobID            int64     `json:"job_id"`
	JobTitle         string    `json:"job_title"`
	CompanyID        *int64    `json:"company_id,omitempty"`
	CompanyName      string    `json:"company_name"`
	CompanyLogoURL   string    `json:"company_logo_url,omitempty"`
	UserID           int64     `json:"user_id"`
	UserFullName     string    `json:"user_full_name"`
	UserEmail        string    `json:"user_email"`
	UserPhotoURL     string    `json:"user_photo_url,omitempty"`
	AppliedAt        time.Time `json:"applied_at"`
	Status           string    `json:"status"`
	Source           string    `json:"source"`
	MatchScore       float64   `json:"match_score"`
	ViewedByEmployer bool      `json:"viewed_by_employer"`
	IsBookmarked     bool      `json:"is_bookmarked"`
	ResumeURL        string    `json:"resume_url,omitempty"`
	NotesCount       int       `json:"notes_count"`
	InterviewsCount  int       `json:"interviews_count"`
	DocumentsCount   int       `json:"documents_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ApplicationDetailResponse represents detailed application response
type ApplicationDetailResponse struct {
	ID               int64                         `json:"id"`
	JobID            int64                         `json:"job_id"`
	JobTitle         string                        `json:"job_title"`
	JobSlug          string                        `json:"job_slug"`
	CompanyID        *int64                        `json:"company_id,omitempty"`
	CompanyName      string                        `json:"company_name"`
	CompanySlug      string                        `json:"company_slug"`
	CompanyLogoURL   string                        `json:"company_logo_url,omitempty"`
	UserID           int64                         `json:"user_id"`
	UserFullName     string                        `json:"user_full_name"`
	UserEmail        string                        `json:"user_email"`
	UserPhone        string                        `json:"user_phone,omitempty"`
	UserPhotoURL     string                        `json:"user_photo_url,omitempty"`
	AppliedAt        time.Time                     `json:"applied_at"`
	Status           string                        `json:"status"`
	Source           string                        `json:"source"`
	MatchScore       float64                       `json:"match_score"`
	NotesText        string                        `json:"notes_text,omitempty"`
	ViewedByEmployer bool                          `json:"viewed_by_employer"`
	IsBookmarked     bool                          `json:"is_bookmarked"`
	ResumeURL        string                        `json:"resume_url,omitempty"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
	Documents        []ApplicationDocumentResponse `json:"documents,omitempty"`
	Stages           []ApplicationStageResponse    `json:"stages,omitempty"`
	Notes            []ApplicationNoteResponse     `json:"notes,omitempty"`
	Interviews       []InterviewResponse           `json:"interviews,omitempty"`
}

// ApplicationAnswerResponse represents screening question answer response
type ApplicationAnswerResponse struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
}

// ApplicationDocumentResponse represents application document response
type ApplicationDocumentResponse struct {
	ID           int64      `json:"id"`
	DocumentType string     `json:"document_type"`
	FileName     string     `json:"file_name,omitempty"`
	FileURL      string     `json:"file_url"`
	FileType     string     `json:"file_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
	UploadedAt   time.Time  `json:"uploaded_at"`
	IsVerified   bool       `json:"is_verified"`
	VerifiedBy   *int64     `json:"verified_by,omitempty"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	Notes        string     `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ApplicationStageResponse represents application stage response
type ApplicationStageResponse struct {
	ID          int64      `json:"id"`
	StageName   string     `json:"stage_name"`
	Description string     `json:"description,omitempty"`
	HandledBy   *int64     `json:"handled_by,omitempty"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Duration    *string    `json:"duration,omitempty"`
	Notes       string     `json:"notes,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ApplicationNoteResponse represents application note response
type ApplicationNoteResponse struct {
	ID         int64     `json:"id"`
	StageID    *int64    `json:"stage_id,omitempty"`
	AuthorID   int64     `json:"author_id"`
	NoteType   string    `json:"note_type"`
	NoteText   string    `json:"note_text"`
	Visibility string    `json:"visibility"`
	Sentiment  string    `json:"sentiment"`
	IsPinned   bool      `json:"is_pinned"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// InterviewResponse represents interview response
type InterviewResponse struct {
	ID                 int64      `json:"id"`
	ApplicationID      int64      `json:"application_id"`
	StageID            *int64     `json:"stage_id,omitempty"`
	InterviewerID      *int64     `json:"interviewer_id,omitempty"`
	ScheduledAt        time.Time  `json:"scheduled_at"`
	EndedAt            *time.Time `json:"ended_at,omitempty"`
	InterviewType      string     `json:"interview_type"`
	MeetingLink        string     `json:"meeting_link,omitempty"`
	Location           string     `json:"location,omitempty"`
	Status             string     `json:"status"`
	OverallScore       *float64   `json:"overall_score,omitempty"`
	TechnicalScore     *float64   `json:"technical_score,omitempty"`
	CommunicationScore *float64   `json:"communication_score,omitempty"`
	PersonalityScore   *float64   `json:"personality_score,omitempty"`
	Remarks            string     `json:"remarks,omitempty"`
	FeedbackSummary    string     `json:"feedback_summary,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// ApplicationListResponse represents list of applications response
type ApplicationListResponse struct {
	Applications []ApplicationResponse `json:"applications"`
}

// ApplicationStatsResponse represents application statistics response
type ApplicationStatsResponse struct {
	TotalApplications       int64            `json:"total_applications"`
	PendingApplications     int64            `json:"pending_applications"`
	ScreeningApplications   int64            `json:"screening_applications"`
	ShortlistedApplications int64            `json:"shortlisted_applications"`
	InterviewApplications   int64            `json:"interview_applications"`
	OfferedApplications     int64            `json:"offered_applications"`
	RejectedApplications    int64            `json:"rejected_applications"`
	WithdrawnApplications   int64            `json:"withdrawn_applications"`
	AverageMatchScore       float64          `json:"average_match_score"`
	StatusBreakdown         map[string]int64 `json:"status_breakdown"`
	TimeToHire              *float64         `json:"time_to_hire,omitempty"`     // in days
	ApplicationRate         *float64         `json:"application_rate,omitempty"` // applications per day
}
