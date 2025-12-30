package note

import "context"

type Store interface {
	Create(ctx context.Context, n *Note) error
	GetByID(ctx context.Context, id int64) (*Note, error)
	List(ctx context.Context) ([]Note, error)
	Update(ctx context.Context, n *Note) error
	Delete(ctx context.Context, id int64) error
}
