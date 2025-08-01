// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package author_repository_impl

import (
	"time"
)

type Author struct {
	ID        int64
	Name      string
	Bio       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Book struct {
	ID        int64
	AuthorID  int64
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
