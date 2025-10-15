package response

import "time"

// ApplicationResponse represents application basic response
type ApplicationResponse struct {
	ID                 int64      `json:"id"`
	UUID               string     `json:"uuid"`
	JobID              int64      `json:"job_id"`
	JobTitle           string     `json:"job_title"`
	CompanyID          int64      `json:"company_id"`
	CompanyName        string     `json:"company_name"`
	CompanyLogoURL     string     `json:"company_logo_url,omitempty"`
	UserID             int64      `json:"user_id"`
	UserFullName       string     `json:"user_full_name"`
	UserEmail          string     `json:"user_email"`
	UserPhotoURL       string     `json:"user_photo_url,omitempty"`
	Status             string     `json:"status"`
	ResumeURL          string     `json:"resume_url"`
	CoverLetter        string     `json:"cover_letter,omitempty"`
	ExpectedSalary     *float64   `json:"expected_salary,omitempty"`
	Currency           string     `json:"currency,omitempty"`
	AvailableStartDate *time.Time `json:"available_start_date,omitempty"`
	AppliedAt          time.Time  `json:"applied_at"`
	ViewedAt           *time.Time `json:"viewed_at,omitempty"`
	IsViewed           bool       `json:"is_viewed"`
	MatchScore         *float64   `json:"match_score,omitempty"`
	NotesCount         int        `json:"notes_count"`
	InterviewsCount    int        `json:"interviews_count"`
	DocumentsCount     int        `json:"documents_count"`
}

// ApplicationDetailResponse represents detailed application response
type ApplicationDetailResponse struct {
	ID                 int64                         `json:"id"`
	UUID               string                        `json:"uuid"`
	JobID              int64                         `json:"job_id"`
	JobTitle           string                        `json:"job_title"`
	JobSlug            string                        `json:"job_slug"`
	CompanyID          int64                         `json:"company_id"`
	CompanyName        string                        `json:"company_name"`
	CompanySlug        string                        `json:"company_slug"`
	CompanyLogoURL     string                        `json:"company_logo_url,omitempty"`
	UserID             int64                         `json:"user_id"`
	UserFullName       string                        `json:"user_full_name"`
	UserEmail          string                        `json:"user_email"`
	UserPhone          string                        `json:"user_phone,omitempty"`
	UserPhotoURL       string                        `json:"user_photo_url,omitempty"`
	Status             string                        `json:"status"`
	ResumeURL          string                        `json:"resume_url"`
	CoverLetter        string                        `json:"cover_letter,omitempty"`
	ExpectedSalary     *float64                      `json:"expected_salary,omitempty"`
	Currency           string                        `json:"currency,omitempty"`
	AvailableStartDate *time.Time                    `json:"available_start_date,omitempty"`
	AppliedAt          time.Time                     `json:"applied_at"`
	ViewedAt           *time.Time                    `json:"viewed_at,omitempty"`
	ResponsedAt        *time.Time                    `json:"responsed_at,omitempty"`
	WithdrawnAt        *time.Time                    `json:"withdrawn_at,omitempty"`
	WithdrawReason     string                        `json:"withdraw_reason,omitempty"`
	RejectionReason    string                        `json:"rejection_reason,omitempty"`
	IsViewed           bool                          `json:"is_viewed"`
	MatchScore         *float64                      `json:"match_score,omitempty"`
	MatchReasons       []string                      `json:"match_reasons,omitempty"`
	CreatedAt          time.Time                     `json:"created_at"`
	UpdatedAt          time.Time                     `json:"updated_at"`
	Answers            []ApplicationAnswerResponse   `json:"answers,omitempty"`
	Documents          []ApplicationDocumentResponse `json:"documents,omitempty"`
	Stages             []ApplicationStageResponse    `json:"stages,omitempty"`
	Notes              []ApplicationNoteResponse     `json:"notes,omitempty"`
	Interviews         []InterviewResponse           `json:"interviews,omitempty"`
	Timeline           []ApplicationTimelineResponse `json:"timeline,omitempty"`
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
	ID           int64     `json:"id"`
	DocumentType string    `json:"document_type"`
	DocumentURL  string    `json:"document_url"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	Description  string    `json:"description,omitempty"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

// ApplicationStageResponse represents application stage response
type ApplicationStageResponse struct {
	ID          int64      `json:"id"`
	StageName   string     `json:"stage_name"`
	StageOrder  int16      `json:"stage_order"`
	Description string     `json:"description,omitempty"`
	EnteredAt   time.Time  `json:"entered_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Notes       string     `json:"notes,omitempty"`
	HandledBy   *int64     `json:"handled_by,omitempty"`
	HandlerName string     `json:"handler_name,omitempty"`
}

// ApplicationNoteResponse represents application note response
type ApplicationNoteResponse struct {
	ID          int64     `json:"id"`
	NoteText    string    `json:"note_text"`
	IsInternal  bool      `json:"is_internal"`
	CreatedBy   int64     `json:"created_by"`
	CreatorName string    `json:"creator_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InterviewResponse represents interview response
type InterviewResponse struct {
	ID             int64                       `json:"id"`
	UUID           string                      `json:"uuid"`
	ApplicationID  int64                       `json:"application_id"`
	InterviewType  string                      `json:"interview_type"`
	InterviewStage string                      `json:"interview_stage"`
	ScheduledAt    time.Time                   `json:"scheduled_at"`
	Duration       int16                       `json:"duration"`
	Location       string                      `json:"location,omitempty"`
	MeetingURL     string                      `json:"meeting_url,omitempty"`
	Status         string                      `json:"status"`
	Rating         *int16                      `json:"rating,omitempty"`
	Feedback       string                      `json:"feedback,omitempty"`
	Result         string                      `json:"result,omitempty"`
	Recommendation string                      `json:"recommendation,omitempty"`
	Notes          string                      `json:"notes,omitempty"`
	ConductedBy    *int64                      `json:"conducted_by,omitempty"`
	ConductorName  string                      `json:"conductor_name,omitempty"`
	CreatedAt      time.Time                   `json:"created_at"`
	UpdatedAt      time.Time                   `json:"updated_at"`
	CompletedAt    *time.Time                  `json:"completed_at,omitempty"`
	CancelledAt    *time.Time                  `json:"cancelled_at,omitempty"`
	Interviewers   []InterviewerResponse       `json:"interviewers,omitempty"`
	Feedbacks      []InterviewFeedbackResponse `json:"feedbacks,omitempty"`
}

// InterviewerResponse represents interviewer response
type InterviewerResponse struct {
	UserID      int64  `json:"user_id"`
	FullName    string `json:"full_name"`
	Position    string `json:"position,omitempty"`
	Email       string `json:"email"`
	PhotoURL    string `json:"photo_url,omitempty"`
	IsConfirmed bool   `json:"is_confirmed"`
}

// InterviewFeedbackResponse represents interview feedback response
type InterviewFeedbackResponse struct {
	ID              int64     `json:"id"`
	InterviewerID   int64     `json:"interviewer_id"`
	InterviewerName string    `json:"interviewer_name"`
	Rating          int16     `json:"rating"`
	Feedback        string    `json:"feedback"`
	Recommendation  string    `json:"recommendation"`
	CreatedAt       time.Time `json:"created_at"`
}

// ApplicationTimelineResponse represents application timeline event response
type ApplicationTimelineResponse struct {
	ID          int64     `json:"id"`
	EventType   string    `json:"event_type"`
	EventTitle  string    `json:"event_title"`
	Description string    `json:"description,omitempty"`
	ActorID     *int64    `json:"actor_id,omitempty"`
	ActorName   string    `json:"actor_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
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

// InterviewScheduleResponse represents interview schedule response
type InterviewScheduleResponse struct {
	Date       string              `json:"date"`
	Interviews []InterviewResponse `json:"interviews"`
}

// ApplicationExperienceResponse represents application experience rating response
type ApplicationExperienceResponse struct {
	ApplicationID int64     `json:"application_id"`
	Rating        int16     `json:"rating"`
	Feedback      string    `json:"feedback,omitempty"`
	RatedAt       time.Time `json:"rated_at"`
}
