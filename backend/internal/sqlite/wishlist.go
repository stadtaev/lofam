package sqlite

import (
	"context"
	"database/sql"

	"github.com/stadtaev/lofam/backend/internal/wishlist"
)

type WishlistStore struct {
	db *DB
}

func NewWishlistStore(db *DB) *WishlistStore {
	return &WishlistStore{db: db}
}

func (s *WishlistStore) Create(ctx context.Context, w *wishlist.Wishlist) error {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO wishlists (title, content, color, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, w.Title, w.Content, w.Color, w.CreatedAt, w.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	w.ID = id
	return nil
}

func (s *WishlistStore) GetByID(ctx context.Context, id int64) (*wishlist.Wishlist, error) {
	var w wishlist.Wishlist
	err := s.db.QueryRowContext(ctx, `
		SELECT id, title, content, color, created_at, updated_at
		FROM wishlists WHERE id = ?
	`, id).Scan(&w.ID, &w.Title, &w.Content, &w.Color, &w.CreatedAt, &w.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, wishlist.ErrNotFound(id)
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (s *WishlistStore) List(ctx context.Context) ([]wishlist.Wishlist, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, content, color, created_at, updated_at
		FROM wishlists ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []wishlist.Wishlist
	for rows.Next() {
		var w wishlist.Wishlist
		if err := rows.Scan(&w.ID, &w.Title, &w.Content, &w.Color, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, w)
	}

	if items == nil {
		items = []wishlist.Wishlist{}
	}

	return items, rows.Err()
}

func (s *WishlistStore) Update(ctx context.Context, w *wishlist.Wishlist) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE wishlists SET title = ?, content = ?, color = ?, updated_at = ?
		WHERE id = ?
	`, w.Title, w.Content, w.Color, w.UpdatedAt, w.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return wishlist.ErrNotFound(w.ID)
	}

	return nil
}

func (s *WishlistStore) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM wishlists WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return wishlist.ErrNotFound(id)
	}

	return nil
}
