package shopping

import "context"

type Store interface {
	Create(ctx context.Context, item *Item) error
	List(ctx context.Context) ([]Item, error)
	Delete(ctx context.Context, id int64) error
}
