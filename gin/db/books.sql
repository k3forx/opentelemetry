-- name: GetByID :one
SELECT * FROM books WHERE id = ? LIMIT 1;

-- name: GetAllByAuthorID :many
SELECT b.* FROM books AS b LEFT JOIN authors AS a ON b.author_id = a.id WHERE a.id = ?;
