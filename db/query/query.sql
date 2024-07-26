
-- name: FetchUsers :many
SELECT username, email,
       last_login,created_at,role, updated_at
       FROM users; 