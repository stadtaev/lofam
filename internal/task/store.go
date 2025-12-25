package task

import "context"

// Store defines the interface for task persistence.
// Defined here (consumer) rather than in sqlite package (implementer) - idiomatic Go.
type Store interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id int64) (*Task, error)
	List(ctx context.Context) ([]Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id int64) error
}
