package wishlist

import "context"

type Store interface {
	Create(ctx context.Context, w *Wishlist) error
	GetByID(ctx context.Context, id int64) (*Wishlist, error)
	List(ctx context.Context) ([]Wishlist, error)
	Update(ctx context.Context, w *Wishlist) error
	Delete(ctx context.Context, id int64) error
}
