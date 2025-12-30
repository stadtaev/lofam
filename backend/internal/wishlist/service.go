package wishlist

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

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Wishlist, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	w := &Wishlist{
		Title:     req.Title,
		Content:   req.Content,
		Color:     req.Color,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.Create(ctx, w); err != nil {
		return nil, err
	}

	return w, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Wishlist, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Wishlist, error) {
	return s.store.List(ctx)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateRequest) (*Wishlist, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	w, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	w.Title = req.Title
	w.Content = req.Content
	w.Color = req.Color
	w.UpdatedAt = time.Now()

	if err := s.store.Update(ctx, w); err != nil {
		return nil, err
	}

	return w, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
