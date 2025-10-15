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

	var matchScore *float64
	if a.MatchScore > 0 {
		matchScore = &a.MatchScore
	}

	return &response.ApplicationResponse{
		ID:         a.ID,
		JobID:      a.JobID,
		UserID:     a.UserID,
		Status:     a.Status,
		ResumeURL:  a.ResumeURL,
		AppliedAt:  a.AppliedAt,
		IsViewed:   a.ViewedByEmployer,
		MatchScore: matchScore,
	}
}

// ToApplicationDetailResponse maps JobApplication entity with relations to ApplicationDetailResponse DTO
func ToApplicationDetailResponse(a *application.JobApplication) *response.ApplicationDetailResponse {
	if a == nil {
		return nil
	}

	var matchScore *float64
	if a.MatchScore > 0 {
		matchScore = &a.MatchScore
	}

	resp := &response.ApplicationDetailResponse{
		ID:         a.ID,
		JobID:      a.JobID,
		UserID:     a.UserID,
		Status:     a.Status,
		ResumeURL:  a.ResumeURL,
		AppliedAt:  a.AppliedAt,
		IsViewed:   a.ViewedByEmployer,
		MatchScore: matchScore,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
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
		StageOrder:  0, // Calculate based on stage name
		Description: s.Description,
		EnteredAt:   s.StartedAt,
		CompletedAt: s.CompletedAt,
		Notes:       s.Notes,
		HandledBy:   s.HandledBy,
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
		DocumentURL:  d.FileURL,
		FileName:     d.FileName,
		FileSize:     d.FileSize,
		Description:  d.Notes,
		UploadedAt:   d.UploadedAt,
	}
}

// ToApplicationNoteResponse maps ApplicationNote entity to ApplicationNoteResponse DTO
func ToApplicationNoteResponse(n *application.ApplicationNote) *response.ApplicationNoteResponse {
	if n == nil {
		return nil
	}

	return &response.ApplicationNoteResponse{
		ID:         n.ID,
		NoteText:   n.NoteText,
		IsInternal: n.IsInternal(),
		CreatedBy:  n.AuthorID,
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
		ID:             i.ID,
		ApplicationID:  i.ApplicationID,
		InterviewType:  i.InterviewType,
		InterviewStage: "", // Not in entity, handler should set
		ScheduledAt:    i.ScheduledAt,
		Duration:       0, // Not in entity, calculate from ScheduledAt and EndedAt
		Location:       i.Location,
		MeetingURL:     i.MeetingLink,
		Status:         i.Status,
		Rating:         nil, // Use OverallScore if needed
		Feedback:       i.FeedbackSummary,
		Result:         "", // Not in entity, handler should set
		Notes:          i.Remarks,
		ConductedBy:    i.InterviewerID,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
		CompletedAt:    i.EndedAt,
		CancelledAt:    nil, // Not directly in entity
	}
}
