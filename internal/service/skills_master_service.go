package service

import (
	"context"
	"fmt"

	"keerja-backend/internal/domain/master"
)

// skillsMasterService implements master.SkillsMasterService
type skillsMasterService struct {
	repo master.SkillsMasterRepository
}

// NewSkillsMasterService creates a new instance of skills master service
func NewSkillsMasterService(repo master.SkillsMasterRepository) master.SkillsMasterService {
	return &skillsMasterService{
		repo: repo,
	}
}

// GetSkills retrieves skills with filtering and pagination
func (s *skillsMasterService) GetSkills(ctx context.Context, filter *master.SkillsFilter) (*master.SkillListResponse, error) {
	skills, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	// Convert to response
	skillResponses := make([]master.SkillResponse, len(skills))
	for i, skill := range skills {
		skillResponses[i] = s.toSkillResponse(&skill)
	}

	// Calculate total pages
	totalPages := 0
	if filter != nil && filter.PageSize > 0 {
		totalPages = int(total) / filter.PageSize
		if int(total)%filter.PageSize > 0 {
			totalPages++
		}
	}

	page := 1
	pageSize := len(skills)
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.PageSize > 0 {
			pageSize = filter.PageSize
		}
	}

	return &master.SkillListResponse{
		Skills:     skillResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchSkills searches skills by query string
func (s *skillsMasterService) SearchSkills(ctx context.Context, query string, page, pageSize int) (*master.SkillListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	skills, total, err := s.repo.SearchSkills(ctx, query, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to search skills: %w", err)
	}

	// Convert to response
	skillResponses := make([]master.SkillResponse, len(skills))
	for i, skill := range skills {
		skillResponses[i] = s.toSkillResponse(&skill)
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &master.SkillListResponse{
		Skills:     skillResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetSkillsByType retrieves skills by type
func (s *skillsMasterService) GetSkillsByType(ctx context.Context, skillType string) ([]master.SkillResponse, error) {
	skills, err := s.repo.ListByType(ctx, skillType)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills by type: %w", err)
	}

	responses := make([]master.SkillResponse, len(skills))
	for i, skill := range skills {
		responses[i] = s.toSkillResponse(&skill)
	}

	return responses, nil
}

// GetSkillsByIDs retrieves multiple skills by their IDs
func (s *skillsMasterService) GetSkillsByIDs(ctx context.Context, ids []int64) ([]master.SkillResponse, error) {
	if len(ids) == 0 {
		return []master.SkillResponse{}, nil
	}

	skills := make([]master.SkillResponse, 0, len(ids))

	// Fetch each skill by ID
	for _, id := range ids {
		skill, err := s.repo.FindByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get skill with ID %d: %w", id, err)
		}

		// Skip if skill not found (don't fail the entire request)
		if skill != nil {
			skills = append(skills, s.toSkillResponse(skill))
		}
	}

	return skills, nil
}

// GetSkillsByNames retrieves multiple skills by their names
func (s *skillsMasterService) GetSkillsByNames(ctx context.Context, names []string) ([]master.SkillResponse, error) {
	if len(names) == 0 {
		return []master.SkillResponse{}, nil
	}

	skills := make([]master.SkillResponse, 0, len(names))

	// Fetch each skill by name
	for _, name := range names {
		skill, err := s.repo.FindByName(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get skill '%s': %w", name, err)
		}

		// Skip if skill not found (don't fail the entire request)
		if skill != nil {
			skills = append(skills, s.toSkillResponse(skill))
		}
	}

	return skills, nil
}

// Helper methods

// toSkillResponse converts entity to response DTO
func (s *skillsMasterService) toSkillResponse(skill *master.SkillsMaster) master.SkillResponse {
	response := master.SkillResponse{
		ID:              skill.ID,
		Code:            skill.Code,
		Name:            skill.Name,
		NormalizedName:  skill.NormalizedName,
		CategoryID:      skill.CategoryID,
		Description:     skill.Description,
		SkillType:       skill.SkillType,
		DifficultyLevel: skill.DifficultyLevel,
		PopularityScore: skill.PopularityScore,
		Aliases:         skill.Aliases,
		ParentID:        skill.ParentID,
		IsActive:        skill.IsActive,
		ChildrenCount:   int64(len(skill.Children)),
	}

	if skill.Parent != nil {
		parent := s.toSkillResponse(skill.Parent)
		response.Parent = &parent
	}

	return response
}

// Stub implementations for interface compliance (not used by endpoints)

func (s *skillsMasterService) GetActiveSkills(ctx context.Context) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetSkill(ctx context.Context, id int64) (*master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetSkillByCode(ctx context.Context, code string) (*master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetSkillsByDifficulty(ctx context.Context, difficultyLevel string) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetMostPopular(ctx context.Context, limit int) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetSkillStats(ctx context.Context) (*master.SkillStatsResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) CreateSkill(ctx context.Context, req *master.CreateSkillRequest) (*master.SkillsMaster, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) UpdateSkill(ctx context.Context, id int64, req *master.UpdateSkillRequest) (*master.SkillsMaster, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) DeleteSkill(ctx context.Context, id int64) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetSkillsByCategory(ctx context.Context, categoryID int64) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetTypeStats(ctx context.Context) ([]master.TypeStatResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetDifficultyStats(ctx context.Context) ([]master.DifficultyStatResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetRootSkills(ctx context.Context) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetChildSkills(ctx context.Context, parentID int64) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetParentSkill(ctx context.Context, childID int64) (*master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetSkillTree(ctx context.Context, rootID int64) (*master.SkillTreeResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) SetParentSkill(ctx context.Context, skillID, parentID int64) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) RemoveParentSkill(ctx context.Context, skillID int64) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) UpdatePopularity(ctx context.Context, id int64, score float64) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) IncrementPopularity(ctx context.Context, id int64) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetPopularByType(ctx context.Context, skillType string, limit int) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) AddAlias(ctx context.Context, skillID int64, alias string) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) RemoveAlias(ctx context.Context, skillID int64, alias string) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) UpdateAliases(ctx context.Context, skillID int64, aliases []string) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) SearchByAlias(ctx context.Context, alias string) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) ActivateSkill(ctx context.Context, id int64) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) DeactivateSkill(ctx context.Context, id int64) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) BulkCreateSkills(ctx context.Context, skills []master.CreateSkillRequest) ([]master.SkillsMaster, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) BulkUpdatePopularity(ctx context.Context, updates map[int64]float64) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) ImportSkills(ctx context.Context, data []master.SkillImportData) (*master.ImportResult, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) ExportSkills(ctx context.Context, filter *master.SkillsFilter) ([]master.SkillExportData, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetRelatedSkills(ctx context.Context, skillID int64, limit int) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetComplementarySkills(ctx context.Context, skillIDs []int64, limit int) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetSkillSuggestions(ctx context.Context, userSkills []int64, limit int) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) GetTrendingSkills(ctx context.Context, limit int) ([]master.SkillResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) ValidateSkill(ctx context.Context, req *master.CreateSkillRequest) error {
	return fmt.Errorf("not implemented")
}

func (s *skillsMasterService) CheckDuplicateSkill(ctx context.Context, name, code string) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (s *skillsMasterService) ValidateHierarchy(ctx context.Context, skillID, parentID int64) error {
	return fmt.Errorf("not implemented")
}
