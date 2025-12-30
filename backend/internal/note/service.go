package note

import (
	"context"
	"time"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Note, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	n := &Note{
		Title:     req.Title,
		Content:   req.Content,
		Color:     req.Color,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.Create(ctx, n); err != nil {
		return nil, err
	}

	return n, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Note, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Note, error) {
	return s.store.List(ctx)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateRequest) (*Note, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	n, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	n.Title = req.Title
	n.Content = req.Content
	n.Color = req.Color
	n.UpdatedAt = time.Now()

	if err := s.store.Update(ctx, n); err != nil {
		return nil, err
	}

	return n, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
