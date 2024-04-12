-- name: GetByID :one
SELECT * FROM authors WHERE id = ? LIMIT 1;

-- name: Create :execresult
INSERT INTO authors (name, bio) VALUES (?, ?);
