package service

import (
	"context"

	"lofam/internal/domain"
	"lofam/internal/repository"
)

type TaskService struct {
	taskRepo    repository.TaskRepository
	projectRepo repository.ProjectRepository
}

func NewTaskService(taskRepo repository.TaskRepository, projectRepo repository.ProjectRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (s *TaskService) Create(ctx context.Context, req domain.CreateTaskRequest) (*domain.Task, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.projectRepo.GetByID(ctx, req.ProjectID); err != nil {
		return nil, err
	}

	priority := req.Priority
	if priority == "" {
		priority = domain.TaskPriorityMedium
	}

	task := &domain.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.TaskStatusTodo,
		Priority:    priority,
		ProjectID:   req.ProjectID,
		DueDate:     req.DueDate,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *TaskService) List(ctx context.Context) ([]domain.Task, error) {
	tasks, err := s.taskRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []domain.Task{}
	}
	return tasks, nil
}

func (s *TaskService) ListByProjectID(ctx context.Context, projectID int64) ([]domain.Task, error) {
	if _, err := s.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, err
	}

	tasks, err := s.taskRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []domain.Task{}
	}
	return tasks, nil
}

func (s *TaskService) Update(ctx context.Context, id int64, req domain.UpdateTaskRequest) (*domain.Task, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	return s.taskRepo.Delete(ctx, id)
}
