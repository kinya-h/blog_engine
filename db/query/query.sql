
-- name: FetchUsers :many
SELECT username, email,
       last_login,created_at, updated_at
       FROM users; 