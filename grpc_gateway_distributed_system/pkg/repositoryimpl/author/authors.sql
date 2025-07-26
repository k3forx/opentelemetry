-- name: GetByID :one
SELECT id, name, bio, created_at FROM authors WHERE id = ?; 
