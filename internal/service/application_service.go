package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"keerja-backend/internal/domain/application"
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/email"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/domain/notification"
	"keerja-backend/internal/domain/user"
)

// applicationService implements application.ApplicationService interface
type applicationService struct {
	appRepo      application.ApplicationRepository
	jobRepo      job.JobRepository
	userRepo     user.UserRepository
	companyRepo  company.CompanyRepository
	emailService email.EmailService
	notifService notification.NotificationService
}

// NewApplicationService creates a new application service instance
func NewApplicationService(
	appRepo application.ApplicationRepository,
	jobRepo job.JobRepository,
	userRepo user.UserRepository,
	companyRepo company.CompanyRepository,
	emailService email.EmailService,
	notifService notification.NotificationService,
) application.ApplicationService {
	return &applicationService{
		appRepo:      appRepo,
		jobRepo:      jobRepo,
		userRepo:     userRepo,
		companyRepo:  companyRepo,
		emailService: emailService,
		notifService: notifService,
	}
}

// ===== Application Submission and Management (Job Seeker) =====

// ApplyForJob submits a job application
func (s *applicationService) ApplyForJob(ctx context.Context, req *application.ApplyJobRequest) (*application.JobApplication, error) {
	// Check if user can apply for this job
	if err := s.CanApplyForJob(ctx, req.JobID, req.UserID); err != nil {
		return nil, err
	}

	// Check if already applied
	existingApp, _ := s.appRepo.FindByJobAndUser(ctx, req.JobID, req.UserID)
	if existingApp != nil {
		return nil, errors.New("you have already applied for this job")
	}

	// Get job details
	j, err := s.jobRepo.FindByID(ctx, req.JobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Create application
	app := &application.JobApplication{
		JobID:     req.JobID,
		UserID:    req.UserID,
		CompanyID: &j.CompanyID,
		Status:    "applied",
		Source:    req.Source,
		ResumeURL: req.ResumeURL,
		NotesText: req.CoverLetter,
	}

	// Set default source if not provided
	if app.Source == "" {
		app.Source = "keerja_portal"
	}

	// Calculate match score (integrate with job matching service)
	// For now, set to 0 - should be calculated by job service
	app.MatchScore = 0.0

	// Validate application
	if err := s.ValidateApplication(ctx, app); err != nil {
		return nil, fmt.Errorf("application validation failed: %w", err)
	}

	// Create application
	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	// Create initial stage
	stage := &application.JobApplicationStage{
		ApplicationID: app.ID,
		StageName:     "applied",
		Description:   "Application submitted",
	}
	if err := s.appRepo.CreateStage(ctx, stage); err != nil {
		return nil, fmt.Errorf("failed to create initial stage: %w", err)
	}

	// Upload documents if provided
	for _, docReq := range req.Documents {
		docReq.ApplicationID = app.ID
		docReq.UserID = req.UserID
		if _, err := s.UploadApplicationDocument(ctx, &docReq); err != nil {
			// Log error but don't fail the application
			fmt.Printf("failed to upload document: %v\n", err)
		}
	}

	// Increment application count for job
	s.jobRepo.IncrementApplications(ctx, req.JobID)

	// Send notification (async)
	go s.NotifyApplicationReceived(ctx, app.ID)

	// Reload application with relationships
	return s.appRepo.FindByID(ctx, app.ID)
}

// WithdrawApplication withdraws a job application
func (s *applicationService) WithdrawApplication(ctx context.Context, applicationID, userID int64) error {
	// Check ownership
	if err := s.CheckApplicationOwnership(ctx, applicationID, userID); err != nil {
		return err
	}

	// Get application
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Check if can be withdrawn
	if !app.CanWithdraw() {
		return errors.New("application cannot be withdrawn in current status")
	}

	// Update status
	app.Status = "withdrawn"
	if err := s.appRepo.Update(ctx, app); err != nil {
		return fmt.Errorf("failed to withdraw application: %w", err)
	}

	// Create withdrawal stage
	stage := &application.JobApplicationStage{
		ApplicationID: applicationID,
		StageName:     "withdrawn",
		Description:   "Application withdrawn by applicant",
		Notes:         "Withdrawn by user",
	}
	stage.Complete()
	s.appRepo.CreateStage(ctx, stage)

	return nil
}

// GetMyApplications retrieves user's applications
func (s *applicationService) GetMyApplications(ctx context.Context, userID int64, filter application.ApplicationFilter, page, limit int) (*application.ApplicationListResponse, error) {
	// Set user ID in filter
	filter.UserID = userID

	// Get applications
	apps, total, err := s.appRepo.ListByUser(ctx, userID, filter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	// Build response
	return s.buildApplicationListResponse(ctx, apps, total, page, limit)
}

// GetApplicationDetail retrieves detailed application information
func (s *applicationService) GetApplicationDetail(ctx context.Context, applicationID, userID int64) (*application.ApplicationDetailResponse, error) {
	// Check ownership
	if err := s.CheckApplicationOwnership(ctx, applicationID, userID); err != nil {
		return nil, err
	}

	return s.buildApplicationDetailResponse(ctx, applicationID)
}

// GetMyApplicationStats retrieves user's application statistics
func (s *applicationService) GetMyApplicationStats(ctx context.Context, userID int64) (*application.UserApplicationStats, error) {
	return s.appRepo.GetUserApplicationStats(ctx, userID)
}

// ===== Application Review and Management (Employer) =====

// GetJobApplications retrieves applications for a job
func (s *applicationService) GetJobApplications(ctx context.Context, jobID int64, filter application.ApplicationFilter, page, limit int) (*application.ApplicationListResponse, error) {
	// Verify job exists
	_, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Set job ID in filter
	filter.JobID = jobID

	// Get applications
	apps, total, err := s.appRepo.ListByJob(ctx, jobID, filter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	// Build response
	return s.buildApplicationListResponse(ctx, apps, total, page, limit)
}

// GetCompanyApplications retrieves all applications for a company
func (s *applicationService) GetCompanyApplications(ctx context.Context, companyID int64, filter application.ApplicationFilter, page, limit int) (*application.ApplicationListResponse, error) {
	// Verify company exists
	_, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// Set company ID in filter
	filter.CompanyID = companyID

	// Get applications
	apps, total, err := s.appRepo.ListByCompany(ctx, companyID, filter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	// Build response
	return s.buildApplicationListResponse(ctx, apps, total, page, limit)
}

// GetApplicationForReview retrieves application for employer review
func (s *applicationService) GetApplicationForReview(ctx context.Context, applicationID, employerUserID int64) (*application.ApplicationDetailResponse, error) {
	// Check employer access
	if err := s.CheckEmployerAccess(ctx, applicationID, employerUserID); err != nil {
		return nil, err
	}

	// Mark as viewed if not already
	app, _ := s.appRepo.FindByID(ctx, applicationID)
	if app != nil && !app.ViewedByEmployer {
		s.MarkAsViewed(ctx, applicationID, employerUserID)
	}

	return s.buildApplicationDetailResponse(ctx, applicationID)
}

// MarkAsViewed marks application as viewed by employer
func (s *applicationService) MarkAsViewed(ctx context.Context, applicationID, employerUserID int64) error {
	// Check employer access
	if err := s.CheckEmployerAccess(ctx, applicationID, employerUserID); err != nil {
		return err
	}

	return s.appRepo.MarkAsViewed(ctx, applicationID)
}

// ToggleBookmark toggles application bookmark status
func (s *applicationService) ToggleBookmark(ctx context.Context, applicationID, employerUserID int64) error {
	// Check employer access
	if err := s.CheckEmployerAccess(ctx, applicationID, employerUserID); err != nil {
		return err
	}

	return s.appRepo.ToggleBookmark(ctx, applicationID)
}

// GetBookmarkedApplications retrieves bookmarked applications
func (s *applicationService) GetBookmarkedApplications(ctx context.Context, companyID int64, page, limit int) (*application.ApplicationListResponse, error) {
	apps, total, err := s.appRepo.GetBookmarkedApplications(ctx, companyID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarked applications: %w", err)
	}

	return s.buildApplicationListResponse(ctx, apps, total, page, limit)
}

// ===== Application Status Workflow (Employer) =====

// MoveToScreening moves application to screening stage
func (s *applicationService) MoveToScreening(ctx context.Context, applicationID, handledBy int64, notes string) error {
	return s.updateApplicationStage(ctx, applicationID, handledBy, "screening", "Moved to screening", notes)
}

// MoveToShortlist moves application to shortlist stage
func (s *applicationService) MoveToShortlist(ctx context.Context, applicationID, handledBy int64, notes string) error {
	return s.updateApplicationStage(ctx, applicationID, handledBy, "shortlisted", "Shortlisted for interview", notes)
}

// MoveToInterview moves application to interview stage
func (s *applicationService) MoveToInterview(ctx context.Context, applicationID, handledBy int64, notes string) error {
	return s.updateApplicationStage(ctx, applicationID, handledBy, "interview", "Interview scheduled", notes)
}

// MakeOffer makes job offer to applicant
func (s *applicationService) MakeOffer(ctx context.Context, applicationID, handledBy int64, notes string) error {
	return s.updateApplicationStage(ctx, applicationID, handledBy, "offered", "Job offer extended", notes)
}

// MarkAsHired marks applicant as hired
func (s *applicationService) MarkAsHired(ctx context.Context, applicationID, handledBy int64, notes string) error {
	return s.updateApplicationStage(ctx, applicationID, handledBy, "hired", "Applicant hired", notes)
}

// RejectApplication rejects an application
func (s *applicationService) RejectApplication(ctx context.Context, applicationID, handledBy int64, reason string) error {
	// Check employer access
	if err := s.CheckEmployerAccess(ctx, applicationID, handledBy); err != nil {
		return err
	}

	// Get application
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Complete current stage
	currentStage, _ := s.appRepo.GetCurrentStage(ctx, applicationID)
	if currentStage != nil && !currentStage.IsCompleted() {
		s.appRepo.CompleteStage(ctx, currentStage.ID, "Application rejected")
	}

	// Update status
	app.Status = "rejected"
	if err := s.appRepo.Update(ctx, app); err != nil {
		return fmt.Errorf("failed to reject application: %w", err)
	}

	// Create rejection stage
	stage := &application.JobApplicationStage{
		ApplicationID: applicationID,
		StageName:     "rejected",
		Description:   "Application rejected",
		HandledBy:     &handledBy,
		Notes:         reason,
	}
	stage.Complete()
	if err := s.appRepo.CreateStage(ctx, stage); err != nil {
		return fmt.Errorf("failed to create rejection stage: %w", err)
	}

	// Send notification
	go s.NotifyStatusUpdate(ctx, applicationID, "rejected")

	return nil
}

// BulkUpdateStatus updates status for multiple applications
func (s *applicationService) BulkUpdateStatus(ctx context.Context, applicationIDs []int64, status string, handledBy int64) error {
	for _, appID := range applicationIDs {
		// Check access for each application
		if err := s.CheckEmployerAccess(ctx, appID, handledBy); err != nil {
			continue // Skip unauthorized applications
		}

		// Update status based on the status value
		switch status {
		case "screening":
			s.MoveToScreening(ctx, appID, handledBy, "Bulk update")
		case "shortlisted":
			s.MoveToShortlist(ctx, appID, handledBy, "Bulk update")
		case "interview":
			s.MoveToInterview(ctx, appID, handledBy, "Bulk update")
		case "offered":
			s.MakeOffer(ctx, appID, handledBy, "Bulk update")
		case "hired":
			s.MarkAsHired(ctx, appID, handledBy, "Bulk update")
		case "rejected":
			s.RejectApplication(ctx, appID, handledBy, "Bulk rejection")
		}
	}

	return nil
}

// updateApplicationStage is a helper to update application stage
func (s *applicationService) updateApplicationStage(ctx context.Context, applicationID, handledBy int64, newStatus, description, notes string) error {
	// Check employer access
	if err := s.CheckEmployerAccess(ctx, applicationID, handledBy); err != nil {
		return err
	}

	// Get application
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Check if application is in valid state
	if app.IsCompleted() {
		return errors.New("cannot update completed application")
	}

	// Complete current stage if exists
	currentStage, _ := s.appRepo.GetCurrentStage(ctx, applicationID)
	if currentStage != nil && !currentStage.IsCompleted() {
		s.appRepo.CompleteStage(ctx, currentStage.ID, "Moved to next stage")
	}

	// Update application status
	app.Status = newStatus
	if err := s.appRepo.Update(ctx, app); err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	// Create new stage
	stage := &application.JobApplicationStage{
		ApplicationID: applicationID,
		StageName:     newStatus,
		Description:   description,
		HandledBy:     &handledBy,
		Notes:         notes,
	}
	if err := s.appRepo.CreateStage(ctx, stage); err != nil {
		return fmt.Errorf("failed to create stage: %w", err)
	}

	// Send notification
	go s.NotifyStatusUpdate(ctx, applicationID, newStatus)

	return nil
}

// ===== Stage Management =====

// GetApplicationStages retrieves all stages for an application
func (s *applicationService) GetApplicationStages(ctx context.Context, applicationID int64) ([]application.JobApplicationStage, error) {
	return s.appRepo.ListStagesByApplication(ctx, applicationID)
}

// GetCurrentStage retrieves current active stage
func (s *applicationService) GetCurrentStage(ctx context.Context, applicationID int64) (*application.JobApplicationStage, error) {
	return s.appRepo.GetCurrentStage(ctx, applicationID)
}

// GetStageHistory retrieves stage history
func (s *applicationService) GetStageHistory(ctx context.Context, applicationID int64) ([]application.JobApplicationStage, error) {
	return s.appRepo.GetStageHistory(ctx, applicationID)
}

// CompleteStage marks a stage as completed
func (s *applicationService) CompleteStage(ctx context.Context, stageID, handledBy int64, notes string) error {
	// Get stage
	stage, err := s.appRepo.FindStageByID(ctx, stageID)
	if err != nil {
		return fmt.Errorf("stage not found: %w", err)
	}

	// Check employer access
	if err := s.CheckEmployerAccess(ctx, stage.ApplicationID, handledBy); err != nil {
		return err
	}

	// Complete stage
	return s.appRepo.CompleteStage(ctx, stageID, notes)
}

// ===== Document Management =====

// UploadApplicationDocument uploads a document for application
func (s *applicationService) UploadApplicationDocument(ctx context.Context, req *application.UploadDocumentRequest) (*application.ApplicationDocument, error) {
	// Check ownership
	if err := s.CheckApplicationOwnership(ctx, req.ApplicationID, req.UserID); err != nil {
		return nil, err
	}

	// Create document
	doc := &application.ApplicationDocument{
		ApplicationID: req.ApplicationID,
		UserID:        req.UserID,
		DocumentType:  req.DocumentType,
		FileName:      req.FileName,
		FileURL:       req.FileURL,
		FileType:      req.FileType,
		FileSize:      req.FileSize,
		Notes:         req.Notes,
	}

	if err := s.appRepo.CreateDocument(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to upload document: %w", err)
	}

	return doc, nil
}

// UpdateDocument updates document information
func (s *applicationService) UpdateDocument(ctx context.Context, documentID int64, req *application.UpdateDocumentRequest) (*application.ApplicationDocument, error) {
	// Get document
	doc, err := s.appRepo.FindDocumentByID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// Update fields
	if req.FileName != "" {
		doc.FileName = req.FileName
	}
	if req.FileURL != "" {
		doc.FileURL = req.FileURL
	}
	if req.Notes != "" {
		doc.Notes = req.Notes
	}

	if err := s.appRepo.UpdateDocument(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	return doc, nil
}

// DeleteDocument deletes an application document
func (s *applicationService) DeleteDocument(ctx context.Context, documentID, userID int64) error {
	// Get document
	doc, err := s.appRepo.FindDocumentByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Check ownership
	if err := s.CheckApplicationOwnership(ctx, doc.ApplicationID, userID); err != nil {
		return err
	}

	return s.appRepo.DeleteDocument(ctx, documentID)
}

// GetApplicationDocuments retrieves all documents for an application
func (s *applicationService) GetApplicationDocuments(ctx context.Context, applicationID int64) ([]application.ApplicationDocument, error) {
	return s.appRepo.ListDocumentsByApplication(ctx, applicationID)
}

// GetDocumentsByType retrieves documents by type
func (s *applicationService) GetDocumentsByType(ctx context.Context, applicationID int64, docType string) ([]application.ApplicationDocument, error) {
	return s.appRepo.GetDocumentsByType(ctx, applicationID, docType)
}

// VerifyDocument marks a document as verified
func (s *applicationService) VerifyDocument(ctx context.Context, documentID, verifiedBy int64, notes string) error {
	// Get document
	doc, err := s.appRepo.FindDocumentByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Check employer access
	if err := s.CheckEmployerAccess(ctx, doc.ApplicationID, verifiedBy); err != nil {
		return err
	}

	return s.appRepo.VerifyDocument(ctx, documentID, verifiedBy)
}

// GetUnverifiedDocuments retrieves unverified documents
func (s *applicationService) GetUnverifiedDocuments(ctx context.Context, page, limit int) ([]application.ApplicationDocument, int64, error) {
	return s.appRepo.GetUnverifiedDocuments(ctx, page, limit)
}

// ===== Notes Management (Employer) =====

// AddNote adds a note to application
func (s *applicationService) AddNote(ctx context.Context, req *application.AddNoteRequest) (*application.ApplicationNote, error) {
	// Check employer access
	if err := s.CheckEmployerAccess(ctx, req.ApplicationID, req.AuthorID); err != nil {
		return nil, err
	}

	// Create note
	note := &application.ApplicationNote{
		ApplicationID: req.ApplicationID,
		StageID:       req.StageID,
		AuthorID:      req.AuthorID,
		NoteType:      req.NoteType,
		NoteText:      req.NoteText,
		Visibility:    req.Visibility,
		Sentiment:     req.Sentiment,
		IsPinned:      req.IsPinned,
	}

	// Set defaults
	if note.NoteType == "" {
		note.NoteType = "internal"
	}
	if note.Visibility == "" {
		note.Visibility = "internal"
	}
	if note.Sentiment == "" {
		note.Sentiment = "neutral"
	}

	if err := s.appRepo.CreateNote(ctx, note); err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	return note, nil
}

// UpdateNote updates an application note
func (s *applicationService) UpdateNote(ctx context.Context, noteID int64, req *application.UpdateNoteRequest) (*application.ApplicationNote, error) {
	// Get note
	note, err := s.appRepo.FindNoteByID(ctx, noteID)
	if err != nil {
		return nil, fmt.Errorf("note not found: %w", err)
	}

	// Update fields
	if req.NoteText != "" {
		note.NoteText = req.NoteText
	}
	if req.Visibility != "" {
		note.Visibility = req.Visibility
	}
	if req.Sentiment != "" {
		note.Sentiment = req.Sentiment
	}
	if req.IsPinned != nil {
		note.IsPinned = *req.IsPinned
	}

	if err := s.appRepo.UpdateNote(ctx, note); err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return note, nil
}

// DeleteNote deletes an application note
func (s *applicationService) DeleteNote(ctx context.Context, noteID, authorID int64) error {
	// Get note
	note, err := s.appRepo.FindNoteByID(ctx, noteID)
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	// Check if author can delete
	if note.AuthorID != authorID {
		return errors.New("only the author can delete this note")
	}

	return s.appRepo.DeleteNote(ctx, noteID)
}

// GetApplicationNotes retrieves application notes
func (s *applicationService) GetApplicationNotes(ctx context.Context, applicationID int64, visibility string) ([]application.ApplicationNote, error) {
	notes, err := s.appRepo.ListNotesByApplication(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	// Filter by visibility if specified
	if visibility != "" {
		filtered := make([]application.ApplicationNote, 0)
		for _, note := range notes {
			if note.Visibility == visibility {
				filtered = append(filtered, note)
			}
		}
		return filtered, nil
	}

	return notes, nil
}

// GetStageNotes retrieves notes for a specific stage
func (s *applicationService) GetStageNotes(ctx context.Context, stageID int64) ([]application.ApplicationNote, error) {
	return s.appRepo.ListNotesByStage(ctx, stageID)
}

// PinNote pins a note
func (s *applicationService) PinNote(ctx context.Context, noteID, employerUserID int64) error {
	// Get note
	note, err := s.appRepo.FindNoteByID(ctx, noteID)
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	// Check employer access
	if err := s.CheckEmployerAccess(ctx, note.ApplicationID, employerUserID); err != nil {
		return err
	}

	return s.appRepo.PinNote(ctx, noteID)
}

// UnpinNote unpins a note
func (s *applicationService) UnpinNote(ctx context.Context, noteID, employerUserID int64) error {
	// Get note
	note, err := s.appRepo.FindNoteByID(ctx, noteID)
	if err != nil {
		return fmt.Errorf("note not found: %w", err)
	}

	// Check employer access
	if err := s.CheckEmployerAccess(ctx, note.ApplicationID, employerUserID); err != nil {
		return err
	}

	return s.appRepo.UnpinNote(ctx, noteID)
}

// GetPinnedNotes retrieves pinned notes for an application
func (s *applicationService) GetPinnedNotes(ctx context.Context, applicationID int64) ([]application.ApplicationNote, error) {
	return s.appRepo.GetPinnedNotes(ctx, applicationID)
}

// ===== Interview Scheduling and Management =====

// ScheduleInterview schedules an interview
func (s *applicationService) ScheduleInterview(ctx context.Context, req *application.ScheduleInterviewRequest) (*application.Interview, error) {
	// Check employer access
	app, err := s.appRepo.FindByID(ctx, req.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// Ensure application is in interview stage or later
	if app.Status != "interview" && app.Status != "offered" {
		// Auto-move to interview stage if not already
		if err := s.MoveToInterview(ctx, req.ApplicationID, *req.InterviewerID, "Interview scheduled"); err != nil {
			return nil, err
		}
	}

	// Create interview
	interview := &application.Interview{
		ApplicationID: req.ApplicationID,
		StageID:       req.StageID,
		InterviewerID: req.InterviewerID,
		ScheduledAt:   req.ScheduledAt,
		InterviewType: req.InterviewType,
		MeetingLink:   req.MeetingLink,
		Location:      req.Location,
		Status:        "scheduled",
	}

	// Set default interview type
	if interview.InterviewType == "" {
		interview.InterviewType = "online"
	}

	if err := s.appRepo.CreateInterview(ctx, interview); err != nil {
		return nil, fmt.Errorf("failed to schedule interview: %w", err)
	}

	// Send notification
	go s.NotifyInterviewScheduled(ctx, interview.ID)

	return interview, nil
}

// RescheduleInterview reschedules an interview
func (s *applicationService) RescheduleInterview(ctx context.Context, interviewID int64, req *application.RescheduleInterviewRequest) (*application.Interview, error) {
	// Get interview
	interview, err := s.appRepo.FindInterviewByID(ctx, interviewID)
	if err != nil {
		return nil, fmt.Errorf("interview not found: %w", err)
	}

	// Update interview
	interview.ScheduledAt = req.ScheduledAt
	interview.Status = "rescheduled"
	if req.MeetingLink != "" {
		interview.MeetingLink = req.MeetingLink
	}
	if req.Location != "" {
		interview.Location = req.Location
	}

	if err := s.appRepo.UpdateInterview(ctx, interview); err != nil {
		return nil, fmt.Errorf("failed to reschedule interview: %w", err)
	}

	// Add note about rescheduling
	if req.Reason != "" {
		noteReq := &application.AddNoteRequest{
			ApplicationID: interview.ApplicationID,
			StageID:       interview.StageID,
			AuthorID:      *interview.InterviewerID,
			NoteType:      "reminder",
			NoteText:      fmt.Sprintf("Interview rescheduled: %s", req.Reason),
			Visibility:    "internal",
		}
		s.AddNote(ctx, noteReq)
	}

	// Send notification
	go s.NotifyInterviewScheduled(ctx, interviewID)

	return interview, nil
}

// CancelInterview cancels an interview
func (s *applicationService) CancelInterview(ctx context.Context, interviewID int64, cancelledBy int64, reason string) error {
	// Get interview
	interview, err := s.appRepo.FindInterviewByID(ctx, interviewID)
	if err != nil {
		return fmt.Errorf("interview not found: %w", err)
	}

	// Update status
	interview.Status = "cancelled"
	if err := s.appRepo.UpdateInterview(ctx, interview); err != nil {
		return fmt.Errorf("failed to cancel interview: %w", err)
	}

	// Add note about cancellation
	if reason != "" {
		noteReq := &application.AddNoteRequest{
			ApplicationID: interview.ApplicationID,
			StageID:       interview.StageID,
			AuthorID:      cancelledBy,
			NoteType:      "internal",
			NoteText:      fmt.Sprintf("Interview cancelled: %s", reason),
			Visibility:    "internal",
		}
		s.AddNote(ctx, noteReq)
	}

	return nil
}

// CompleteInterview marks interview as completed with evaluation
func (s *applicationService) CompleteInterview(ctx context.Context, interviewID int64, req *application.CompleteInterviewRequest) (*application.Interview, error) {
	// Get interview
	interview, err := s.appRepo.FindInterviewByID(ctx, interviewID)
	if err != nil {
		return nil, fmt.Errorf("interview not found: %w", err)
	}

	// Update interview with scores
	interview.Status = "completed"
	now := time.Now()
	interview.EndedAt = &now
	interview.OverallScore = req.OverallScore
	interview.TechnicalScore = req.TechnicalScore
	interview.CommunicationScore = req.CommunicationScore
	interview.PersonalityScore = req.PersonalityScore
	interview.Remarks = req.Remarks
	interview.FeedbackSummary = req.FeedbackSummary

	if err := s.appRepo.UpdateInterview(ctx, interview); err != nil {
		return nil, fmt.Errorf("failed to complete interview: %w", err)
	}

	// Add feedback note
	if req.FeedbackSummary != "" {
		noteReq := &application.AddNoteRequest{
			ApplicationID: interview.ApplicationID,
			StageID:       interview.StageID,
			AuthorID:      req.CompletedBy,
			NoteType:      "feedback",
			NoteText:      req.FeedbackSummary,
			Visibility:    "internal",
			Sentiment:     "neutral",
		}
		s.AddNote(ctx, noteReq)
	}

	return interview, nil
}

// MarkInterviewNoShow marks interview as no-show
func (s *applicationService) MarkInterviewNoShow(ctx context.Context, interviewID int64, markedBy int64) error {
	// Get interview
	interview, err := s.appRepo.FindInterviewByID(ctx, interviewID)
	if err != nil {
		return fmt.Errorf("interview not found: %w", err)
	}

	// Update status
	interview.Status = "no_show"
	if err := s.appRepo.UpdateInterview(ctx, interview); err != nil {
		return fmt.Errorf("failed to mark as no-show: %w", err)
	}

	// Add note
	noteReq := &application.AddNoteRequest{
		ApplicationID: interview.ApplicationID,
		StageID:       interview.StageID,
		AuthorID:      markedBy,
		NoteType:      "internal",
		NoteText:      "Candidate did not attend scheduled interview",
		Visibility:    "internal",
		Sentiment:     "negative",
	}
	s.AddNote(ctx, noteReq)

	return nil
}

// GetApplicationInterviews retrieves all interviews for an application
func (s *applicationService) GetApplicationInterviews(ctx context.Context, applicationID int64) ([]application.Interview, error) {
	return s.appRepo.ListInterviewsByApplication(ctx, applicationID)
}

// GetInterviewDetail retrieves interview details
func (s *applicationService) GetInterviewDetail(ctx context.Context, interviewID int64) (*application.Interview, error) {
	return s.appRepo.FindInterviewByID(ctx, interviewID)
}

// GetUpcomingInterviews retrieves upcoming interviews for employer
func (s *applicationService) GetUpcomingInterviews(ctx context.Context, employerUserID int64, days int) ([]application.Interview, error) {
	endDate := time.Now().AddDate(0, 0, days)
	return s.appRepo.GetUpcomingInterviews(ctx, endDate, 100)
}

// GetInterviewsByDateRange retrieves interviews within date range
func (s *applicationService) GetInterviewsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]application.Interview, error) {
	return s.appRepo.GetInterviewsByDateRange(ctx, startDate, endDate)
}

// SendInterviewReminder sends interview reminder
func (s *applicationService) SendInterviewReminder(ctx context.Context, interviewID int64) error {
	// Get interview
	interview, err := s.appRepo.FindInterviewByID(ctx, interviewID)
	if err != nil {
		return fmt.Errorf("interview not found: %w", err)
	}

	// Check if interview is scheduled
	if interview.Status != "scheduled" && interview.Status != "rescheduled" {
		return errors.New("interview is not in scheduled status")
	}

	// Send notification
	return s.NotifyInterviewReminder(ctx, interviewID)
}

// ===== Search and Filtering =====

// SearchApplications performs advanced application search
func (s *applicationService) SearchApplications(ctx context.Context, filter application.ApplicationSearchFilter, page, limit int) (*application.ApplicationListResponse, error) {
	apps, total, err := s.appRepo.SearchApplications(ctx, filter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search applications: %w", err)
	}

	return s.buildApplicationListResponse(ctx, apps, total, page, limit)
}

// GetHighScoreApplications retrieves high-score applications
func (s *applicationService) GetHighScoreApplications(ctx context.Context, companyID int64, minScore float64, limit int) ([]application.JobApplication, error) {
	// Get applications with high scores
	apps, err := s.appRepo.GetApplicationsWithHighScore(ctx, minScore, limit)
	if err != nil {
		return nil, err
	}

	// Filter by company
	filtered := make([]application.JobApplication, 0)
	for _, app := range apps {
		if app.CompanyID != nil && *app.CompanyID == companyID {
			filtered = append(filtered, app)
		}
	}

	return filtered, nil
}

// GetRecentApplications retrieves recent applications
func (s *applicationService) GetRecentApplications(ctx context.Context, companyID int64, hours int, limit int) ([]application.JobApplication, error) {
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)

	filter := application.ApplicationFilter{
		CompanyID:    companyID,
		AppliedAfter: &cutoffTime,
		SortBy:       "latest",
	}

	apps, _, err := s.appRepo.ListByCompany(ctx, companyID, filter, 1, limit)
	return apps, err
}

// ===== Analytics and Reporting =====

// GetApplicationAnalytics retrieves detailed application analytics
func (s *applicationService) GetApplicationAnalytics(ctx context.Context, applicationID int64) (*application.ApplicationAnalytics, error) {
	// Get application
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// Get stages
	stages, _ := s.appRepo.ListStagesByApplication(ctx, applicationID)

	// Get documents
	docs, _ := s.appRepo.ListDocumentsByApplication(ctx, applicationID)

	// Get interviews
	interviews, _ := s.appRepo.ListInterviewsByApplication(ctx, applicationID)

	// Build analytics
	analytics := &application.ApplicationAnalytics{
		ApplicationID:  applicationID,
		Timeline:       s.buildTimeline(app, stages, interviews),
		StageProgress:  s.buildStageProgress(stages),
		DocumentStats:  s.buildDocumentStats(docs),
		InterviewStats: s.buildInterviewStats(interviews),
		MatchAnalysis: application.MatchAnalysis{
			OverallScore: app.MatchScore,
			// Other match details would come from job matching service
		},
	}

	return analytics, nil
}

// GetJobApplicationAnalytics retrieves job application analytics
func (s *applicationService) GetJobApplicationAnalytics(ctx context.Context, jobID int64, startDate, endDate time.Time) (*application.JobApplicationAnalytics, error) {
	// Get job
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Get job stats
	stats, err := s.appRepo.GetJobApplicationStats(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Build analytics
	analytics := &application.JobApplicationAnalytics{
		JobID:             jobID,
		JobTitle:          j.Title,
		Period:            fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalApplications: stats.TotalApplications,
		StatusBreakdown: map[string]int64{
			"applied":     stats.AppliedCount,
			"screening":   stats.ScreeningCount,
			"shortlisted": stats.ShortlistedCount,
			"interview":   stats.InterviewCount,
			"offered":     stats.OfferedCount,
			"hired":       stats.HiredCount,
			"rejected":    stats.RejectedCount,
		},
		AverageMatchScore: stats.AverageMatchScore,
		AverageTimeToHire: stats.AverageTimeToHire,
	}

	// Get conversion funnel
	funnel, _ := s.appRepo.GetConversionFunnel(ctx, jobID)
	analytics.ConversionFunnel = funnel

	return analytics, nil
}

// GetCompanyApplicationAnalytics retrieves company application analytics
func (s *applicationService) GetCompanyApplicationAnalytics(ctx context.Context, companyID int64, startDate, endDate time.Time) (*application.CompanyApplicationAnalytics, error) {
	// Get company
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// Get company stats
	stats, err := s.appRepo.GetCompanyApplicationStats(ctx, companyID)
	if err != nil {
		return nil, err
	}

	// Build analytics
	analytics := &application.CompanyApplicationAnalytics{
		CompanyID:         companyID,
		CompanyName:       comp.CompanyName,
		Period:            fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalApplications: stats.TotalApplications,
		TotalHires:        stats.TotalHires,
		ConversionRate:    stats.ConversionRate,
		AverageTimeToHire: stats.AverageTimeToHire,
	}

	// Get stage time analysis
	stageTime, _ := s.appRepo.GetAverageTimePerStage(ctx, companyID)
	analytics.StageTimeAnalysis = stageTime

	return analytics, nil
}

// GetConversionFunnel retrieves hiring funnel metrics
func (s *applicationService) GetConversionFunnel(ctx context.Context, jobID int64) (*application.ConversionFunnel, error) {
	return s.appRepo.GetConversionFunnel(ctx, jobID)
}

// GetApplicationTrends retrieves application trends
func (s *applicationService) GetApplicationTrends(ctx context.Context, companyID int64, startDate, endDate time.Time) ([]application.ApplicationTrend, error) {
	return s.appRepo.GetApplicationTrends(ctx, startDate, endDate)
}

// GetAverageTimePerStage retrieves average time per stage
func (s *applicationService) GetAverageTimePerStage(ctx context.Context, companyID int64) ([]application.StageTimeStats, error) {
	return s.appRepo.GetAverageTimePerStage(ctx, companyID)
}

// GetTopApplicants retrieves top applicants for a job
func (s *applicationService) GetTopApplicants(ctx context.Context, jobID int64, limit int) ([]application.JobApplication, error) {
	return s.appRepo.GetTopApplicants(ctx, jobID, limit)
}

// GetApplicationSourceAnalytics retrieves application source analytics
func (s *applicationService) GetApplicationSourceAnalytics(ctx context.Context, companyID int64) ([]application.SourceStats, error) {
	return s.appRepo.GetApplicationSourceStats(ctx, companyID)
}

// ===== Notifications =====

// NotifyApplicationReceived sends notification for new application
func (s *applicationService) NotifyApplicationReceived(ctx context.Context, applicationID int64) error {
	// Get application details
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Get job details
	j, err := s.jobRepo.FindByID(ctx, app.JobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Get user details
	user, err := s.userRepo.FindByID(ctx, app.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Get company details
	var companyName string
	if app.CompanyID != nil {
		comp, _ := s.companyRepo.FindByID(ctx, *app.CompanyID)
		if comp != nil {
			companyName = comp.CompanyName
		}
	}

	// Send notification to user
	if s.notifService != nil {
		if err := s.notifService.NotifyJobApplication(ctx, app.UserID, app.JobID, applicationID); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to send notification: %v\n", err)
		}
	}

	// Send email confirmation to user
	if s.emailService != nil {
		if err := s.emailService.SendJobApplicationEmail(ctx, user.Email, j.Title, companyName); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to send email: %v\n", err)
		}
	}

	return nil
}

// NotifyStatusUpdate sends notification for status change
func (s *applicationService) NotifyStatusUpdate(ctx context.Context, applicationID int64, newStatus string) error {
	// Get application details
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Get job details
	j, err := s.jobRepo.FindByID(ctx, app.JobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Get user details
	user, err := s.userRepo.FindByID(ctx, app.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Send notification to user
	if s.notifService != nil {
		if err := s.notifService.NotifyStatusUpdate(ctx, app.UserID, applicationID, app.Status, newStatus); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to send notification: %v\n", err)
		}
	}

	// Send email notification to user
	if s.emailService != nil {
		if err := s.emailService.SendJobStatusUpdateEmail(ctx, user.Email, j.Title, newStatus); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to send email: %v\n", err)
		}
	}

	return nil
}

// NotifyInterviewScheduled sends notification for scheduled interview
func (s *applicationService) NotifyInterviewScheduled(ctx context.Context, interviewID int64) error {
	// Get interview details
	interview, err := s.appRepo.FindInterviewByID(ctx, interviewID)
	if err != nil {
		return fmt.Errorf("interview not found: %w", err)
	}

	// Get application details
	app, err := s.appRepo.FindByID(ctx, interview.ApplicationID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Get job details
	j, err := s.jobRepo.FindByID(ctx, app.JobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Get user details
	user, err := s.userRepo.FindByID(ctx, app.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Send notification to user
	if s.notifService != nil {
		if err := s.notifService.NotifyInterviewScheduled(ctx, app.UserID, interviewID, interview.ScheduledAt); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to send notification: %v\n", err)
		}
	}

	// Send email invitation to user
	if s.emailService != nil {
		interviewDateStr := interview.ScheduledAt.Format("2006-01-02 15:04 MST")
		if err := s.emailService.SendInterviewInvitationEmail(ctx, user.Email, j.Title, interviewDateStr); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to send email: %v\n", err)
		}
	}

	return nil
}

// NotifyInterviewReminder sends interview reminder
func (s *applicationService) NotifyInterviewReminder(ctx context.Context, interviewID int64) error {
	// Get interview details
	interview, err := s.appRepo.FindInterviewByID(ctx, interviewID)
	if err != nil {
		return fmt.Errorf("interview not found: %w", err)
	}

	// Get application details
	app, err := s.appRepo.FindByID(ctx, interview.ApplicationID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Get job details
	j, err := s.jobRepo.FindByID(ctx, app.JobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Get user details
	user, err := s.userRepo.FindByID(ctx, app.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Send reminder notification to user
	if s.notifService != nil {
		// Note: NotificationService doesn't have NotifyInterviewReminder method
		// We can reuse NotifyInterviewScheduled for reminders
		if err := s.notifService.NotifyInterviewScheduled(ctx, app.UserID, interviewID, interview.ScheduledAt); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to send reminder notification: %v\n", err)
		}
	}

	// Send reminder email to user (using same template as interview invitation)
	if s.emailService != nil {
		interviewDateStr := interview.ScheduledAt.Format("2006-01-02 15:04 MST")
		if err := s.emailService.SendInterviewInvitationEmail(ctx, user.Email, j.Title, interviewDateStr); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("failed to send reminder email: %v\n", err)
		}
	}

	return nil
}

// ===== Validation and Permissions =====

// ValidateApplication validates application data
func (s *applicationService) ValidateApplication(ctx context.Context, app *application.JobApplication) error {
	// Check required fields
	if app.JobID == 0 {
		return errors.New("job ID is required")
	}
	if app.UserID == 0 {
		return errors.New("user ID is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"applied": true, "screening": true, "shortlisted": true,
		"interview": true, "offered": true, "hired": true,
		"rejected": true, "withdrawn": true,
	}
	if !validStatuses[app.Status] {
		return fmt.Errorf("invalid status: %s", app.Status)
	}

	return nil
}

// CheckApplicationOwnership verifies user owns the application
func (s *applicationService) CheckApplicationOwnership(ctx context.Context, applicationID, userID int64) error {
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	if app.UserID != userID {
		return errors.New("you do not own this application")
	}

	return nil
}

// CheckEmployerAccess verifies employer has access to application
func (s *applicationService) CheckEmployerAccess(ctx context.Context, applicationID, employerUserID int64) error {
	// Get application
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Check if employer belongs to company
	if app.CompanyID == nil {
		return errors.New("application has no company association")
	}

	// Check employer user permission
	employerUser, err := s.companyRepo.FindEmployerUserByUserAndCompany(ctx, employerUserID, *app.CompanyID)
	if err != nil || employerUser == nil {
		return errors.New("you do not have access to this application")
	}

	// Check if role has permission (viewer and above can view applications)
	if employerUser.Role != "viewer" && employerUser.Role != "recruiter" && employerUser.Role != "admin" && employerUser.Role != "owner" {
		return errors.New("insufficient permissions")
	}

	return nil
}

// CanApplyForJob checks if user can apply for job
func (s *applicationService) CanApplyForJob(ctx context.Context, jobID, userID int64) error {
	// Get job
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Check if job is active
	if !j.CanApply() {
		return errors.New("this job is not accepting applications")
	}

	// Check if already applied
	existingApp, _ := s.appRepo.FindByJobAndUser(ctx, jobID, userID)
	if existingApp != nil {
		return errors.New("you have already applied for this job")
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check if user is active
	if !user.IsActive() {
		return errors.New("your account is not active")
	}

	// Check if user is jobseeker
	if !user.IsJobseeker() {
		return errors.New("only job seekers can apply for jobs")
	}

	return nil
}

// ===== Bulk Operations =====

// BulkRejectApplications rejects multiple applications
func (s *applicationService) BulkRejectApplications(ctx context.Context, applicationIDs []int64, rejectedBy int64, reason string) error {
	for _, appID := range applicationIDs {
		if err := s.RejectApplication(ctx, appID, rejectedBy, reason); err != nil {
			// Log error but continue with others
			fmt.Printf("failed to reject application %d: %v\n", appID, err)
		}
	}
	return nil
}

// BulkMoveToStage moves multiple applications to a stage
func (s *applicationService) BulkMoveToStage(ctx context.Context, applicationIDs []int64, stage string, handledBy int64) error {
	return s.BulkUpdateStatus(ctx, applicationIDs, stage, handledBy)
}

// ExportApplications exports applications to CSV/Excel
func (s *applicationService) ExportApplications(ctx context.Context, companyID int64, filter application.ApplicationFilter) ([]byte, error) {
	// TODO: Implement export logic
	// This would generate CSV or Excel file with application data
	return nil, errors.New("export not yet implemented")
}

// ===== Helper Methods =====

// buildApplicationListResponse builds application list response
func (s *applicationService) buildApplicationListResponse(ctx context.Context, apps []application.JobApplication, total int64, page, limit int) (*application.ApplicationListResponse, error) {
	summaries := make([]application.ApplicationSummary, 0, len(apps))

	for _, app := range apps {
		// Get job details
		j, _ := s.jobRepo.FindByID(ctx, app.JobID)
		jobTitle := ""
		companyName := ""
		if j != nil {
			jobTitle = j.Title
		}

		// Get company details
		if app.CompanyID != nil {
			comp, _ := s.companyRepo.FindByID(ctx, *app.CompanyID)
			if comp != nil {
				companyName = comp.CompanyName
			}
		}

		// Get user details
		user, _ := s.userRepo.FindByID(ctx, app.UserID)
		userName := ""
		if user != nil {
			userName = user.FullName
		}

		// Get current stage
		currentStage, _ := s.appRepo.GetCurrentStage(ctx, app.ID)
		currentStageName := app.Status
		if currentStage != nil {
			currentStageName = currentStage.StageName
		}

		// Calculate days since applied
		daysSince := int(time.Since(app.AppliedAt).Hours() / 24)

		summaries = append(summaries, application.ApplicationSummary{
			ID:               app.ID,
			JobID:            app.JobID,
			JobTitle:         jobTitle,
			CompanyName:      companyName,
			UserID:           app.UserID,
			UserName:         userName,
			Status:           app.Status,
			MatchScore:       app.MatchScore,
			AppliedAt:        app.AppliedAt,
			ViewedByEmployer: app.ViewedByEmployer,
			IsBookmarked:     app.IsBookmarked,
			CurrentStage:     currentStageName,
			DaysSinceApplied: daysSince,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &application.ApplicationListResponse{
		Applications: summaries,
		Total:        total,
		Page:         page,
		Limit:        limit,
		TotalPages:   totalPages,
	}, nil
}

// buildApplicationDetailResponse builds detailed application response
func (s *applicationService) buildApplicationDetailResponse(ctx context.Context, applicationID int64) (*application.ApplicationDetailResponse, error) {
	// Get application
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// Get job details
	j, _ := s.jobRepo.FindByID(ctx, app.JobID)
	jobDetail := application.JobDetail{}
	if j != nil {
		jobDetail.ID = j.ID
		jobDetail.Title = j.Title
		jobDetail.CompanyID = j.CompanyID
		jobDetail.Location = j.Location
		// Use ExperienceLevel master data instead of deprecated JobLevel
		if j.ExperienceLevelM != nil {
			jobDetail.JobLevel = j.ExperienceLevelM.Name
		} else {
			jobDetail.JobLevel = ""
		}
		jobDetail.Status = j.Status

		// Get company name
		comp, _ := s.companyRepo.FindByID(ctx, j.CompanyID)
		if comp != nil {
			jobDetail.CompanyName = comp.CompanyName
		}
	}

	// Get applicant profile
	user, _ := s.userRepo.FindByID(ctx, app.UserID)
	applicantProfile := application.ApplicantProfile{}
	if user != nil {
		applicantProfile.UserID = user.ID
		applicantProfile.FullName = user.FullName
		applicantProfile.Email = user.Email
		if user.Phone != nil {
			applicantProfile.Phone = *user.Phone
		}
		applicantProfile.ResumeURL = app.ResumeURL
	}

	// Get stages
	stages, _ := s.appRepo.ListStagesByApplication(ctx, applicationID)

	// Get documents
	documents, _ := s.appRepo.ListDocumentsByApplication(ctx, applicationID)

	// Get notes
	notes, _ := s.appRepo.ListNotesByApplication(ctx, applicationID)

	// Get interviews
	interviews, _ := s.appRepo.ListInterviewsByApplication(ctx, applicationID)

	// Get stats
	stats, _ := s.appRepo.GetApplicationStats(ctx, applicationID)

	return &application.ApplicationDetailResponse{
		Application: *app,
		Job:         jobDetail,
		Applicant:   applicantProfile,
		Stages:      stages,
		Documents:   documents,
		Notes:       notes,
		Interviews:  interviews,
		Stats:       stats,
	}, nil
}

// buildTimeline builds timeline from application events
func (s *applicationService) buildTimeline(app *application.JobApplication, stages []application.JobApplicationStage, interviews []application.Interview) []application.TimelineEvent {
	timeline := make([]application.TimelineEvent, 0)

	// Add application submitted event
	timeline = append(timeline, application.TimelineEvent{
		Date:        app.AppliedAt,
		EventType:   "application_submitted",
		Description: "Application submitted",
		Actor:       "Applicant",
	})

	// Add stage events
	for _, stage := range stages {
		timeline = append(timeline, application.TimelineEvent{
			Date:        stage.StartedAt,
			EventType:   "stage_change",
			Description: fmt.Sprintf("Moved to %s stage", stage.StageName),
			Actor:       "Recruiter",
		})
	}

	// Add interview events
	for _, interview := range interviews {
		timeline = append(timeline, application.TimelineEvent{
			Date:        interview.ScheduledAt,
			EventType:   "interview_scheduled",
			Description: fmt.Sprintf("%s interview scheduled", interview.InterviewType),
			Actor:       "Recruiter",
		})
	}

	return timeline
}

// buildStageProgress builds stage progress information
func (s *applicationService) buildStageProgress(stages []application.JobApplicationStage) []application.StageProgress {
	progress := make([]application.StageProgress, 0, len(stages))

	for _, stage := range stages {
		duration := ""
		if stage.CompletedAt != nil {
			d := stage.CompletedAt.Sub(stage.StartedAt)
			duration = fmt.Sprintf("%.1f hours", d.Hours())
		}

		status := "in_progress"
		if stage.IsCompleted() {
			status = "completed"
		}

		progress = append(progress, application.StageProgress{
			StageName:   stage.StageName,
			StartedAt:   stage.StartedAt,
			CompletedAt: stage.CompletedAt,
			Duration:    duration,
			Status:      status,
		})
	}

	return progress
}

// buildDocumentStats builds document statistics
func (s *applicationService) buildDocumentStats(docs []application.ApplicationDocument) application.DocumentStats {
	stats := application.DocumentStats{
		TotalDocuments: int64(len(docs)),
		DocumentTypes:  make(map[string]int64),
	}

	for _, doc := range docs {
		if doc.IsVerified {
			stats.VerifiedDocuments++
		}
		stats.DocumentTypes[doc.DocumentType]++
	}

	return stats
}

// buildInterviewStats builds interview statistics
func (s *applicationService) buildInterviewStats(interviews []application.Interview) application.InterviewStats {
	stats := application.InterviewStats{
		TotalInterviews: int64(len(interviews)),
	}

	var totalScore float64
	var scoreCount int

	for _, interview := range interviews {
		if interview.IsCompleted() {
			stats.CompletedInterviews++
		}

		if interview.OverallScore != nil {
			totalScore += *interview.OverallScore
			scoreCount++

			if *interview.OverallScore > stats.HighestScore {
				stats.HighestScore = *interview.OverallScore
			}
		}
	}

	if scoreCount > 0 {
		stats.AverageScore = totalScore / float64(scoreCount)
	}

	return stats
}
