package sqlite

import (
	"context"
	"database/sql"

	"lofam/internal/domain"
)

type TaskRepository struct {
	db *DB
}

func NewTaskRepository(db *DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (title, description, status, priority, due_date)
		 VALUES (?, ?, ?, ?, ?)`,
		task.Title, task.Description, task.Status, task.Priority, task.DueDate,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = id

	return r.db.QueryRowContext(ctx,
		"SELECT created_at FROM tasks WHERE id = ?", id,
	).Scan(&task.CreatedAt)
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	var task domain.Task
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, status, priority, due_date, created_at
		 FROM tasks WHERE id = ?`, id,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
		&task.DueDate, &task.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound("task", id)
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) List(ctx context.Context) ([]domain.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, description, status, priority, due_date, created_at
		 FROM tasks ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.DueDate, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET title = ?, description = ?, status = ?, priority = ?, due_date = ?
		 WHERE id = ?`,
		task.Title, task.Description, task.Status, task.Priority, task.DueDate, task.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound("task", task.ID)
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound("task", id)
	}

	return nil
}
