package application

import (
	"time"
)

// JobApplication represents a job application entity
type JobApplication struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	JobID            int64     `gorm:"column:job_id;not null;index:idx_application_unique,unique" json:"job_id" validate:"required"`
	UserID           int64     `gorm:"column:user_id;not null;index:idx_application_unique,unique" json:"user_id" validate:"required"`
	CompanyID        *int64    `gorm:"column:company_id;index" json:"company_id,omitempty"`
	AppliedAt        time.Time `gorm:"column:applied_at;default:now()" json:"applied_at"`
	Status           string    `gorm:"column:status;type:varchar(30);default:'applied'" json:"status" validate:"omitempty,oneof='applied' 'screening' 'shortlisted' 'interview' 'offered' 'hired' 'rejected' 'withdrawn'"`
	Source           string    `gorm:"column:source;type:varchar(50);default:'keerja_portal'" json:"source"`
	MatchScore       float64   `gorm:"column:match_score;type:numeric(5,2);default:0.00" json:"match_score"`
	NotesText        string    `gorm:"column:notes;type:text" json:"notes_text,omitempty"`
	ViewedByEmployer bool      `gorm:"column:viewed_by_employer;default:false" json:"viewed_by_employer"`
	IsBookmarked     bool      `gorm:"column:is_bookmarked;default:false" json:"is_bookmarked"`
	ResumeURL        string    `gorm:"column:resume_url;type:text" json:"resume_url,omitempty"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Stages           []JobApplicationStage `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:CASCADE" json:"stages,omitempty"`
	Documents        []ApplicationDocument `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:CASCADE" json:"documents,omitempty"`
	ApplicationNotes []ApplicationNote     `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:CASCADE" json:"application_notes,omitempty"`
	Interviews       []Interview           `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:CASCADE" json:"interviews,omitempty"`
}

// TableName specifies the table name for JobApplication
func (JobApplication) TableName() string {
	return "job_applications"
}

// IsApplied checks if application status is applied
func (ja *JobApplication) IsApplied() bool {
	return ja.Status == "applied"
}

// IsInProgress checks if application is in hiring process
func (ja *JobApplication) IsInProgress() bool {
	return ja.Status == "screening" || ja.Status == "shortlisted" || ja.Status == "interview" || ja.Status == "offered"
}

// IsCompleted checks if application has final status
func (ja *JobApplication) IsCompleted() bool {
	return ja.Status == "hired" || ja.Status == "rejected" || ja.Status == "withdrawn"
}

// IsHired checks if application resulted in hire
func (ja *JobApplication) IsHired() bool {
	return ja.Status == "hired"
}

// IsRejected checks if application was rejected
func (ja *JobApplication) IsRejected() bool {
	return ja.Status == "rejected"
}

// IsWithdrawn checks if application was withdrawn by user
func (ja *JobApplication) IsWithdrawn() bool {
	return ja.Status == "withdrawn"
}

// CanWithdraw checks if user can withdraw application
func (ja *JobApplication) CanWithdraw() bool {
	return !ja.IsCompleted()
}

// JobApplicationStage represents application stage tracking entity
type JobApplicationStage struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApplicationID int64      `gorm:"column:application_id;not null;index" json:"application_id" validate:"required"`
	StageName     string     `gorm:"column:stage_name;type:varchar(50);not null" json:"stage_name" validate:"required,oneof='applied' 'screening' 'shortlisted' 'interview' 'offered' 'hired' 'rejected' 'withdrawn'"`
	Description   string     `gorm:"column:description;type:text" json:"description,omitempty"`
	HandledBy     *int64     `gorm:"column:handled_by;index" json:"handled_by,omitempty"`
	StartedAt     time.Time  `gorm:"column:started_at;default:now()" json:"started_at"`
	CompletedAt   *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
	Duration      *string    `gorm:"column:duration;->;type:interval" json:"duration,omitempty"` // Generated column, read-only
	Notes         string     `gorm:"column:notes;type:text" json:"notes,omitempty"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Application *JobApplication   `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:CASCADE" json:"application,omitempty"`
	StageNotes  []ApplicationNote `gorm:"foreignKey:StageID;references:ID;constraint:OnDelete:SET NULL" json:"stage_notes,omitempty"`
	Interviews  []Interview       `gorm:"foreignKey:StageID;references:ID;constraint:OnDelete:SET NULL" json:"interviews,omitempty"`
}

// TableName specifies the table name for JobApplicationStage
func (JobApplicationStage) TableName() string {
	return "job_application_stages"
}

// IsCompleted checks if stage is completed
func (jas *JobApplicationStage) IsCompleted() bool {
	return jas.CompletedAt != nil
}

// IsInProgress checks if stage is in progress
func (jas *JobApplicationStage) IsInProgress() bool {
	return jas.CompletedAt == nil
}

// Complete marks the stage as completed
func (jas *JobApplicationStage) Complete() {
	now := time.Now()
	jas.CompletedAt = &now
}

// ApplicationDocument represents application document entity
type ApplicationDocument struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApplicationID int64      `gorm:"column:application_id;not null;index" json:"application_id" validate:"required"`
	UserID        int64      `gorm:"column:user_id;not null;index" json:"user_id" validate:"required"`
	DocumentType  string     `gorm:"column:document_type;type:varchar(50);default:'cv'" json:"document_type" validate:"omitempty,oneof='cv' 'cover_letter' 'portfolio' 'certificate' 'transcript' 'other'"`
	FileName      string     `gorm:"column:file_name;type:varchar(255)" json:"file_name,omitempty"`
	FileURL       string     `gorm:"column:file_url;type:text;not null" json:"file_url" validate:"required"`
	FileType      string     `gorm:"column:file_type;type:varchar(50)" json:"file_type,omitempty"`
	FileSize      int64      `gorm:"column:file_size" json:"file_size,omitempty"`
	UploadedAt    time.Time  `gorm:"column:uploaded_at;default:now()" json:"uploaded_at"`
	IsVerified    bool       `gorm:"column:is_verified;default:false;index" json:"is_verified"`
	VerifiedBy    *int64     `gorm:"column:verified_by" json:"verified_by,omitempty"`
	VerifiedAt    *time.Time `gorm:"column:verified_at" json:"verified_at,omitempty"`
	Notes         string     `gorm:"column:notes;type:text" json:"notes,omitempty"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Application *JobApplication `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:CASCADE" json:"application,omitempty"`
}

// TableName specifies the table name for ApplicationDocument
func (ApplicationDocument) TableName() string {
	return "application_documents"
}

// IsCV checks if document is a CV/resume
func (ad *ApplicationDocument) IsCV() bool {
	return ad.DocumentType == "cv"
}

// IsCoverLetter checks if document is a cover letter
func (ad *ApplicationDocument) IsCoverLetter() bool {
	return ad.DocumentType == "cover_letter"
}

// Verify marks the document as verified
func (ad *ApplicationDocument) Verify(verifierID int64) {
	ad.IsVerified = true
	ad.VerifiedBy = &verifierID
	now := time.Now()
	ad.VerifiedAt = &now
}

// ApplicationNote represents notes on application entity
type ApplicationNote struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApplicationID int64     `gorm:"column:application_id;not null;index" json:"application_id" validate:"required"`
	StageID       *int64    `gorm:"column:stage_id;index" json:"stage_id,omitempty"`
	AuthorID      int64     `gorm:"column:author_id;not null;index" json:"author_id" validate:"required"`
	NoteType      string    `gorm:"column:note_type;type:varchar(30);default:'internal'" json:"note_type" validate:"omitempty,oneof='evaluation' 'feedback' 'reminder' 'internal'"`
	NoteText      string    `gorm:"column:note_text;type:text;not null" json:"note_text" validate:"required"`
	Visibility    string    `gorm:"column:visibility;type:varchar(20);default:'internal'" json:"visibility" validate:"omitempty,oneof='internal' 'public'"`
	Sentiment     string    `gorm:"column:sentiment;type:varchar(20);default:'neutral'" json:"sentiment" validate:"omitempty,oneof='positive' 'neutral' 'negative'"`
	IsPinned      bool      `gorm:"column:is_pinned;default:false" json:"is_pinned"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Application *JobApplication      `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:CASCADE" json:"application,omitempty"`
	Stage       *JobApplicationStage `gorm:"foreignKey:StageID;references:ID;constraint:OnDelete:SET NULL" json:"stage,omitempty"`
}

// TableName specifies the table name for ApplicationNote
func (ApplicationNote) TableName() string {
	return "application_notes"
}

// IsInternal checks if note is internal only
func (an *ApplicationNote) IsInternal() bool {
	return an.Visibility == "internal"
}

// IsPublic checks if note is visible to candidate
func (an *ApplicationNote) IsPublic() bool {
	return an.Visibility == "public"
}

// IsEvaluation checks if note is evaluation type
func (an *ApplicationNote) IsEvaluation() bool {
	return an.NoteType == "evaluation"
}

// IsFeedback checks if note is feedback type
func (an *ApplicationNote) IsFeedback() bool {
	return an.NoteType == "feedback"
}

// IsPositive checks if note has positive sentiment
func (an *ApplicationNote) IsPositive() bool {
	return an.Sentiment == "positive"
}

// IsNegative checks if note has negative sentiment
func (an *ApplicationNote) IsNegative() bool {
	return an.Sentiment == "negative"
}

// Pin marks the note as pinned
func (an *ApplicationNote) Pin() {
	an.IsPinned = true
}

// Unpin removes pinned status
func (an *ApplicationNote) Unpin() {
	an.IsPinned = false
}

// Interview represents interview entity
type Interview struct {
	ID                 int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApplicationID      int64      `gorm:"column:application_id;not null;index" json:"application_id" validate:"required"`
	StageID            *int64     `gorm:"column:stage_id;index" json:"stage_id,omitempty"`
	InterviewerID      *int64     `gorm:"column:interviewer_id;index" json:"interviewer_id,omitempty"`
	ScheduledAt        time.Time  `gorm:"column:scheduled_at;not null;index" json:"scheduled_at" validate:"required"`
	EndedAt            *time.Time `gorm:"column:ended_at" json:"ended_at,omitempty"`
	InterviewType      string     `gorm:"column:interview_type;type:varchar(20);default:'online'" json:"interview_type" validate:"omitempty,oneof='online' 'onsite' 'hybrid'"`
	MeetingLink        string     `gorm:"column:meeting_link;type:text" json:"meeting_link,omitempty"`
	Location           string     `gorm:"column:location;type:text" json:"location,omitempty"`
	Status             string     `gorm:"column:status;type:varchar(20);default:'scheduled'" json:"status" validate:"omitempty,oneof='scheduled' 'completed' 'rescheduled' 'cancelled' 'no_show'"`
	OverallScore       *float64   `gorm:"column:overall_score;type:numeric(4,2)" json:"overall_score,omitempty" validate:"omitempty,min=0,max=100"`
	TechnicalScore     *float64   `gorm:"column:technical_score;type:numeric(4,2)" json:"technical_score,omitempty" validate:"omitempty,min=0,max=100"`
	CommunicationScore *float64   `gorm:"column:communication_score;type:numeric(4,2)" json:"communication_score,omitempty" validate:"omitempty,min=0,max=100"`
	PersonalityScore   *float64   `gorm:"column:personality_score;type:numeric(4,2)" json:"personality_score,omitempty" validate:"omitempty,min=0,max=100"`
	Remarks            string     `gorm:"column:remarks;type:text" json:"remarks,omitempty"`
	FeedbackSummary    string     `gorm:"column:feedback_summary;type:text" json:"feedback_summary,omitempty"`
	CreatedAt          time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relationships
	Application *JobApplication      `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:CASCADE" json:"application,omitempty"`
	Stage       *JobApplicationStage `gorm:"foreignKey:StageID;references:ID;constraint:OnDelete:SET NULL" json:"stage,omitempty"`
}

// TableName specifies the table name for Interview
func (Interview) TableName() string {
	return "interviews"
}

// IsScheduled checks if interview is scheduled
func (i *Interview) IsScheduled() bool {
	return i.Status == "scheduled"
}

// IsCompleted checks if interview is completed
func (i *Interview) IsCompleted() bool {
	return i.Status == "completed"
}

// IsCancelled checks if interview was cancelled
func (i *Interview) IsCancelled() bool {
	return i.Status == "cancelled"
}

// IsNoShow checks if candidate didn't show up
func (i *Interview) IsNoShow() bool {
	return i.Status == "no_show"
}

// IsOnline checks if interview is online
func (i *Interview) IsOnline() bool {
	return i.InterviewType == "online"
}

// IsOnsite checks if interview is onsite
func (i *Interview) IsOnsite() bool {
	return i.InterviewType == "onsite"
}

// Complete marks the interview as completed
func (i *Interview) Complete() {
	i.Status = "completed"
	now := time.Now()
	i.EndedAt = &now
}

// Cancel marks the interview as cancelled
func (i *Interview) Cancel() {
	i.Status = "cancelled"
}

// MarkNoShow marks the interview as no show
func (i *Interview) MarkNoShow() {
	i.Status = "no_show"
}

// HasScores checks if interview has evaluation scores
func (i *Interview) HasScores() bool {
	return i.OverallScore != nil || i.TechnicalScore != nil || i.CommunicationScore != nil || i.PersonalityScore != nil
}

// CalculateAverageScore calculates average from all scores
func (i *Interview) CalculateAverageScore() float64 {
	if !i.HasScores() {
		return 0.0
	}

	var total float64
	var count int

	if i.TechnicalScore != nil {
		total += *i.TechnicalScore
		count++
	}
	if i.CommunicationScore != nil {
		total += *i.CommunicationScore
		count++
	}
	if i.PersonalityScore != nil {
		total += *i.PersonalityScore
		count++
	}

	if count == 0 {
		return 0.0
	}

	return total / float64(count)
}
