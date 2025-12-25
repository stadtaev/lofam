package sqlite

import (
	"context"
	"database/sql"

	"github.com/stadtaev/lofam/internal/task"
)

type TaskStore struct {
	db *DB
}

func NewTaskStore(db *DB) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) Create(ctx context.Context, t *task.Task) error {
	result, err := s.db.ExecContext(ctx,
		`INSERT INTO tasks (title, description, status, priority, due_date)
		 VALUES (?, ?, ?, ?, ?)`,
		t.Title, t.Description, t.Status, t.Priority, t.DueDate,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = id

	return s.db.QueryRowContext(ctx,
		"SELECT created_at FROM tasks WHERE id = ?", id,
	).Scan(&t.CreatedAt)
}

func (s *TaskStore) GetByID(ctx context.Context, id int64) (*task.Task, error) {
	var t task.Task
	err := s.db.QueryRowContext(ctx,
		`SELECT id, title, description, status, priority, due_date, created_at
		 FROM tasks WHERE id = ?`, id,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
		&t.DueDate, &t.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, task.ErrNotFound(id)
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *TaskStore) List(ctx context.Context) ([]task.Task, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, title, description, status, priority, due_date, created_at
		 FROM tasks ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.DueDate, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

func (s *TaskStore) Update(ctx context.Context, t *task.Task) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE tasks SET title = ?, description = ?, status = ?, priority = ?, due_date = ?
		 WHERE id = ?`,
		t.Title, t.Description, t.Status, t.Priority, t.DueDate, t.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return task.ErrNotFound(t.ID)
	}

	return nil
}

func (s *TaskStore) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return task.ErrNotFound(id)
	}

	return nil
}
