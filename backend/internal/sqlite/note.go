package sqlite

import (
	"context"
	"database/sql"

	"github.com/stadtaev/lofam/backend/internal/note"
)

type NoteStore struct {
	db *DB
}

func NewNoteStore(db *DB) *NoteStore {
	return &NoteStore{db: db}
}

func (s *NoteStore) Create(ctx context.Context, n *note.Note) error {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO notes (title, content, color, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, n.Title, n.Content, n.Color, n.CreatedAt, n.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	n.ID = id
	return nil
}

func (s *NoteStore) GetByID(ctx context.Context, id int64) (*note.Note, error) {
	var n note.Note
	err := s.db.QueryRowContext(ctx, `
		SELECT id, title, content, color, created_at, updated_at
		FROM notes WHERE id = ?
	`, id).Scan(&n.ID, &n.Title, &n.Content, &n.Color, &n.CreatedAt, &n.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, note.ErrNotFound(id)
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (s *NoteStore) List(ctx context.Context) ([]note.Note, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, content, color, created_at, updated_at
		FROM notes ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []note.Note
	for rows.Next() {
		var n note.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Color, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	if notes == nil {
		notes = []note.Note{}
	}

	return notes, rows.Err()
}

func (s *NoteStore) Update(ctx context.Context, n *note.Note) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE notes SET title = ?, content = ?, color = ?, updated_at = ?
		WHERE id = ?
	`, n.Title, n.Content, n.Color, n.UpdatedAt, n.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return note.ErrNotFound(n.ID)
	}

	return nil
}

func (s *NoteStore) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM notes WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return note.ErrNotFound(id)
	}

	return nil
}
