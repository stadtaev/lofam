package shopping

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

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Item, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	item := &Item{
		Title:     req.Title,
		CreatedAt: time.Now(),
	}

	if err := s.store.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) List(ctx context.Context) ([]Item, error) {
	return s.store.List(ctx)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
