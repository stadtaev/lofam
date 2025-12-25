package task

import "context"

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Task, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	priority := req.Priority
	if priority == "" {
		priority = PriorityMedium
	}

	t := &Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      StatusTodo,
		Priority:    priority,
		DueDate:     req.DueDate,
	}

	if err := s.store.Create(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Task, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Task, error) {
	tasks, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, nil
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateRequest) (*Task, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	t, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		t.Title = *req.Title
	}
	if req.Description != nil {
		t.Description = *req.Description
	}
	if req.Status != nil {
		t.Status = *req.Status
	}
	if req.Priority != nil {
		t.Priority = *req.Priority
	}
	if req.DueDate != nil {
		t.DueDate = req.DueDate
	}

	if err := s.store.Update(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
