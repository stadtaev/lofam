package sqlite

import (
	"context"
	"database/sql"

	"lofam/internal/domain"
)

type ProjectRepository struct {
	db *DB
}

func NewProjectRepository(db *DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO projects (name, description) VALUES (?, ?)",
		project.Name, project.Description,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	project.ID = id

	return r.db.QueryRowContext(ctx,
		"SELECT created_at FROM projects WHERE id = ?", id,
	).Scan(&project.CreatedAt)
}

func (r *ProjectRepository) GetByID(ctx context.Context, id int64) (*domain.Project, error) {
	var project domain.Project
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, description, created_at FROM projects WHERE id = ?", id,
	).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound("project", id)
	}
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepository) List(ctx context.Context) ([]domain.Project, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name, description, created_at FROM projects ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, rows.Err()
}

func (r *ProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE projects SET name = ?, description = ? WHERE id = ?",
		project.Name, project.Description, project.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound("project", project.ID)
	}

	return nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound("project", id)
	}

	return nil
}
