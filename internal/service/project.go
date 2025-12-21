package service

import (
	"context"

	"lofam/internal/domain"
	"lofam/internal/repository"
)

type ProjectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, req domain.CreateProjectRequest) (*domain.Project, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	project := &domain.Project{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id int64) (*domain.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectService) List(ctx context.Context) ([]domain.Project, error) {
	projects, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	if projects == nil {
		projects = []domain.Project{}
	}
	return projects, nil
}

func (s *ProjectService) Update(ctx context.Context, id int64, req domain.UpdateProjectRequest) (*domain.Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}

	if err := s.repo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
