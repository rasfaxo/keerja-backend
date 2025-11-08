package mapper

import (
	"keerja-backend/internal/domain/admin"
	"keerja-backend/internal/dto/response"
)

// ToAdminCompanyListItemResponse converts admin.AdminCompanyListItem to response DTO
func ToAdminCompanyListItemResponse(item *admin.AdminCompanyListItem) response.AdminCompanyListItemResponse {
	return response.AdminCompanyListItemResponse{
		ID:                 item.ID,
		UUID:               item.UUID,
		CompanyName:        item.CompanyName,
		Slug:               item.Slug,
		LegalName:          item.LegalName,
		RegistrationNumber: item.RegistrationNumber,
		Industry:           item.Industry,
		CompanySize:        item.CompanySize,
		Location:           item.Location,
		Verified:           item.Verified,
		VerifiedAt:         item.VerifiedAt,
		IsActive:           item.IsActive,
		VerificationStatus: item.VerificationStatus,
		TotalJobs:          item.TotalJobs,
		ActiveJobs:         item.ActiveJobs,
		TotalApplications:  item.TotalApplications,
		CreatedAt:          item.CreatedAt,
		UpdatedAt:          item.UpdatedAt,
	}
}

// ToAdminCompanyDetailResponse converts admin.AdminCompanyDetailResponse to response DTO
func ToAdminCompanyDetailResponse(detail *admin.AdminCompanyDetailResponse) response.AdminCompanyDetailResponse {
	resp := response.AdminCompanyDetailResponse{
		ID:                 detail.ID,
		UUID:               detail.UUID,
		CompanyName:        detail.CompanyName,
		Slug:               detail.Slug,
		LegalName:          detail.LegalName,
		RegistrationNumber: detail.RegistrationNumber,
		IndustryID:         detail.IndustryID,
		CompanySizeID:      detail.CompanySizeID,
		WebsiteURL:         detail.WebsiteURL,
		EmailDomain:        detail.EmailDomain,
		Phone:              detail.Phone,
		FullAddress:        detail.FullAddress,
		Description:        detail.Description,
		About:              detail.About,
		Culture:            detail.Culture,
		Verified:           detail.Verified,
		VerifiedAt:         detail.VerifiedAt,
		IsActive:           detail.IsActive,
		CreatedAt:          detail.CreatedAt,
		UpdatedAt:          detail.UpdatedAt,
	}

	// Map verification detail
	if detail.VerificationDetail != nil {
		resp.VerificationDetail = &response.AdminCompanyVerificationDetailResponse{
			ID:                detail.VerificationDetail.ID,
			Status:            detail.VerificationDetail.Status,
			VerificationScore: detail.VerificationDetail.VerificationScore,
			VerificationNotes: detail.VerificationDetail.VerificationNotes,
			RejectionReason:   detail.VerificationDetail.RejectionReason,
			ReviewedBy:        detail.VerificationDetail.ReviewedBy,
			ReviewedByName:    detail.VerificationDetail.ReviewedByName,
			ReviewedAt:        detail.VerificationDetail.ReviewedAt,
			CreatedAt:         detail.VerificationDetail.CreatedAt,
		}
	}

	// Map documents
	if len(detail.Documents) > 0 {
		resp.Documents = make([]response.AdminCompanyDocumentResponse, 0, len(detail.Documents))
		for _, doc := range detail.Documents {
			resp.Documents = append(resp.Documents, response.AdminCompanyDocumentResponse{
				ID:              doc.ID,
				DocumentType:    doc.DocumentType,
				DocumentNumber:  doc.DocumentNumber,
				FilePath:        doc.FilePath,
				Status:          doc.Status,
				VerifiedBy:      doc.VerifiedBy,
				VerifiedAt:      doc.VerifiedAt,
				RejectionReason: doc.RejectionReason,
			})
		}
	}

	// Map owner info
	if detail.OwnerInfo != nil {
		resp.CreatorInfo = &response.AdminCompanyCreatorResponse{
			UserID:   detail.OwnerInfo.UserID,
			FullName: detail.OwnerInfo.FullName,
			Email:    detail.OwnerInfo.Email,
			Phone:    detail.OwnerInfo.Phone,
			Role:     detail.OwnerInfo.Role,
		}
	}

	// Map stats
	if detail.Stats != nil {
		resp.Stats = response.AdminCompanyStatsResponse{
			TotalJobs:           detail.Stats.TotalJobs,
			ActiveJobs:          detail.Stats.ActiveJobs,
			ClosedJobs:          detail.Stats.ClosedJobs,
			DraftJobs:           detail.Stats.DraftJobs,
			TotalApplications:   detail.Stats.TotalApplications,
			PendingApplications: detail.Stats.PendingApplications,
			TotalFollowers:      detail.Stats.TotalFollowers,
			TotalEmployees:      detail.Stats.TotalEmployees,
			TotalReviews:        detail.Stats.TotalReviews,
			AverageRating:       detail.Stats.AverageRating,
		}
	}

	return resp
}

// ToAdminCompanyStatsResponse converts admin stats to response DTO
func ToAdminCompanyStatsResponse(stats *admin.AdminCompanyStatsResponse) response.AdminCompanyStatsResponse {
	return response.AdminCompanyStatsResponse{
		TotalJobs:           stats.TotalJobs,
		ActiveJobs:          stats.ActiveJobs,
		ClosedJobs:          stats.ClosedJobs,
		DraftJobs:           stats.DraftJobs,
		TotalApplications:   stats.TotalApplications,
		PendingApplications: stats.PendingApplications,
		TotalFollowers:      stats.TotalFollowers,
		TotalEmployees:      stats.TotalEmployees,
		TotalReviews:        stats.TotalReviews,
		AverageRating:       stats.AverageRating,
	}
}

// ToAdminDashboardStatsResponse converts dashboard stats to response DTO
func ToAdminDashboardStatsResponse(stats *admin.AdminDashboardStatsResponse) response.AdminDashboardStatsResponse {
	return response.AdminDashboardStatsResponse{
		TotalCompanies:        stats.TotalCompanies,
		VerifiedCompanies:     stats.VerifiedCompanies,
		PendingVerification:   stats.PendingVerification,
		RejectedCompanies:     stats.RejectedCompanies,
		SuspendedCompanies:    stats.SuspendedCompanies,
		NewCompaniesThisMonth: stats.NewCompaniesThisMonth,
		NewCompaniesToday:     stats.NewCompaniesToday,
		TotalJobs:             stats.TotalJobs,
		ActiveJobs:            stats.ActiveJobs,
		TotalApplications:     stats.TotalApplications,
	}
}

// ToAuditLogListResponse converts audit logs to response DTO
func ToAuditLogListResponse(logs *admin.AuditLogListResponse) response.AuditLogListResponse {
	resp := response.AuditLogListResponse{
		Total:      logs.Total,
		Page:       logs.Page,
		Limit:      logs.Limit,
		TotalPages: logs.TotalPages,
	}

	if len(logs.Logs) > 0 {
		resp.Logs = make([]response.AuditLogEntry, 0, len(logs.Logs))
		for _, log := range logs.Logs {
			resp.Logs = append(resp.Logs, response.AuditLogEntry{
				ID:          log.ID,
				CompanyID:   log.CompanyID,
				CompanyName: log.CompanyName,
				AdminID:     log.AdminID,
				AdminName:   log.AdminName,
				Action:      log.Action,
				Description: log.Description,
				OldValue:    log.OldValue,
				NewValue:    log.NewValue,
				IPAddress:   log.IPAddress,
				UserAgent:   log.UserAgent,
				CreatedAt:   log.CreatedAt,
			})
		}
	}

	return resp
}
