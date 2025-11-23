-- name: GetUserByUsername :one
SELECT id, username, password
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT id, username
FROM users
WHERE id = $1
LIMIT 1;
