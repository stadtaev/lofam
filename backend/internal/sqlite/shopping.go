package sqlite

import (
	"context"

	"github.com/stadtaev/lofam/backend/internal/shopping"
)

type ShoppingStore struct {
	db *DB
}

func NewShoppingStore(db *DB) *ShoppingStore {
	return &ShoppingStore{db: db}
}

func (s *ShoppingStore) Create(ctx context.Context, item *shopping.Item) error {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO shopping_items (title, created_at)
		VALUES (?, ?)
	`, item.Title, item.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = id
	return nil
}

func (s *ShoppingStore) List(ctx context.Context) ([]shopping.Item, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, created_at
		FROM shopping_items ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []shopping.Item
	for rows.Next() {
		var item shopping.Item
		if err := rows.Scan(&item.ID, &item.Title, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if items == nil {
		items = []shopping.Item{}
	}

	return items, rows.Err()
}

func (s *ShoppingStore) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM shopping_items WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return shopping.ErrNotFound(id)
	}

	return nil
}
