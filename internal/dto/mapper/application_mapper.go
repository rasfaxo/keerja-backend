package mapper

import (
	"keerja-backend/internal/domain/application"
	"keerja-backend/internal/dto/response"
)

// Application Entity to Response Mappers

// ToApplicationResponse maps JobApplication entity to ApplicationResponse DTO
func ToApplicationResponse(a *application.JobApplication) *response.ApplicationResponse {
	if a == nil {
		return nil
	}

	return &response.ApplicationResponse{
		ID:               a.ID,
		JobID:            a.JobID,
		CompanyID:        a.CompanyID,
		UserID:           a.UserID,
		AppliedAt:        a.AppliedAt,
		Status:           a.Status,
		Source:           a.Source,
		MatchScore:       a.MatchScore,
		ViewedByEmployer: a.ViewedByEmployer,
		IsBookmarked:     a.IsBookmarked,
		ResumeURL:        a.ResumeURL,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
		// Computed fields (set externally if needed)
		NotesCount:      len(a.ApplicationNotes),
		InterviewsCount: len(a.Interviews),
		DocumentsCount:  len(a.Documents),
	}
}

// ToApplicationDetailResponse maps JobApplication entity with relations to ApplicationDetailResponse DTO
func ToApplicationDetailResponse(a *application.JobApplication) *response.ApplicationDetailResponse {
	if a == nil {
		return nil
	}

	resp := &response.ApplicationDetailResponse{
		ID:               a.ID,
		JobID:            a.JobID,
		CompanyID:        a.CompanyID,
		UserID:           a.UserID,
		AppliedAt:        a.AppliedAt,
		Status:           a.Status,
		Source:           a.Source,
		MatchScore:       a.MatchScore,
		NotesText:        a.NotesText,
		ViewedByEmployer: a.ViewedByEmployer,
		IsBookmarked:     a.IsBookmarked,
		ResumeURL:        a.ResumeURL,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
	}

	// Map stages
	if len(a.Stages) > 0 {
		resp.Stages = make([]response.ApplicationStageResponse, len(a.Stages))
		for i, stage := range a.Stages {
			resp.Stages[i] = *ToApplicationStageResponse(&stage)
		}
	}

	// Map documents
	if len(a.Documents) > 0 {
		resp.Documents = make([]response.ApplicationDocumentResponse, len(a.Documents))
		for i, doc := range a.Documents {
			resp.Documents[i] = *ToApplicationDocumentResponse(&doc)
		}
	}

	// Map notes
	if len(a.ApplicationNotes) > 0 {
		resp.Notes = make([]response.ApplicationNoteResponse, len(a.ApplicationNotes))
		for i, note := range a.ApplicationNotes {
			resp.Notes[i] = *ToApplicationNoteResponse(&note)
		}
	}

	// Map interviews
	if len(a.Interviews) > 0 {
		resp.Interviews = make([]response.InterviewResponse, len(a.Interviews))
		for i, interview := range a.Interviews {
			resp.Interviews[i] = *ToInterviewResponse(&interview)
		}
	}

	return resp
}

// ToApplicationStageResponse maps JobApplicationStage entity to ApplicationStageResponse DTO
func ToApplicationStageResponse(s *application.JobApplicationStage) *response.ApplicationStageResponse {
	if s == nil {
		return nil
	}

	return &response.ApplicationStageResponse{
		ID:          s.ID,
		StageName:   s.StageName,
		Description: s.Description,
		HandledBy:   s.HandledBy,
		StartedAt:   s.StartedAt,
		CompletedAt: s.CompletedAt,
		Duration:    s.Duration,
		Notes:       s.Notes,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// ToApplicationDocumentResponse maps ApplicationDocument entity to ApplicationDocumentResponse DTO
func ToApplicationDocumentResponse(d *application.ApplicationDocument) *response.ApplicationDocumentResponse {
	if d == nil {
		return nil
	}

	return &response.ApplicationDocumentResponse{
		ID:           d.ID,
		DocumentType: d.DocumentType,
		FileName:     d.FileName,
		FileURL:      d.FileURL,
		FileType:     d.FileType,
		FileSize:     d.FileSize,
		UploadedAt:   d.UploadedAt,
		IsVerified:   d.IsVerified,
		VerifiedBy:   d.VerifiedBy,
		VerifiedAt:   d.VerifiedAt,
		Notes:        d.Notes,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

// ToApplicationNoteResponse maps ApplicationNote entity to ApplicationNoteResponse DTO
func ToApplicationNoteResponse(n *application.ApplicationNote) *response.ApplicationNoteResponse {
	if n == nil {
		return nil
	}

	return &response.ApplicationNoteResponse{
		ID:         n.ID,
		StageID:    n.StageID,
		AuthorID:   n.AuthorID,
		NoteType:   n.NoteType,
		NoteText:   n.NoteText,
		Visibility: n.Visibility,
		Sentiment:  n.Sentiment,
		IsPinned:   n.IsPinned,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
}

// ToInterviewResponse maps Interview entity to InterviewResponse DTO
func ToInterviewResponse(i *application.Interview) *response.InterviewResponse {
	if i == nil {
		return nil
	}

	return &response.InterviewResponse{
		ID:                 i.ID,
		ApplicationID:      i.ApplicationID,
		StageID:            i.StageID,
		InterviewerID:      i.InterviewerID,
		ScheduledAt:        i.ScheduledAt,
		EndedAt:            i.EndedAt,
		InterviewType:      i.InterviewType,
		MeetingLink:        i.MeetingLink,
		Location:           i.Location,
		Status:             i.Status,
		OverallScore:       i.OverallScore,
		TechnicalScore:     i.TechnicalScore,
		CommunicationScore: i.CommunicationScore,
		PersonalityScore:   i.PersonalityScore,
		Remarks:            i.Remarks,
		FeedbackSummary:    i.FeedbackSummary,
		CreatedAt:          i.CreatedAt,
		UpdatedAt:          i.UpdatedAt,
	}
}
