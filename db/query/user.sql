-- name: CreateUser :execresult
INSERT INTO users(username, email, password_hash)
VALUES (?,?,?);

-- name: GetUser :one
SELECT * FROM users
WHERE username = ?;