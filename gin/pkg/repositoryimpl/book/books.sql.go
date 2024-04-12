// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: books.sql

package book_repository_impl

import (
	"context"
)

const getAllByAuthorID = `-- name: GetAllByAuthorID :many
SELECT b.id, b.author_id, b.title, b.created_at, b.updated_at FROM books AS b LEFT JOIN authors AS a ON b.author_id = a.id WHERE a.id = ?
`

func (q *Queries) GetAllByAuthorID(ctx context.Context, id int64) ([]Book, error) {
	rows, err := q.db.QueryContext(ctx, getAllByAuthorID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Book
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.ID,
			&i.AuthorID,
			&i.Title,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getByID = `-- name: GetByID :one
SELECT id, author_id, title, created_at, updated_at FROM books WHERE id = ? LIMIT 1
`

func (q *Queries) GetByID(ctx context.Context, id int64) (Book, error) {
	row := q.db.QueryRowContext(ctx, getByID, id)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.AuthorID,
		&i.Title,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}