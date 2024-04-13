-- name: GetByID :one
SELECT * FROM books WHERE id = ? LIMIT 1;

-- name: GetAllByAuthorID :many
SELECT b.* FROM books AS b LEFT JOIN authors AS a ON b.author_id = a.id WHERE a.id = ?;

-- name: GetWithAuthorByID :one
SELECT b.*, a.id AS author_id, a.name AS author_name FROM books AS b LEFT JOIN authors AS a ON b.author_id = a.id WHERE b.id = ? LIMIT 1;

-- name: Create :execresult
INSERT INTO books (author_id, title) VALUES (?, ?);
