
-- name: CreatePost :execresult
INSERT INTO posts
(user_id,title,content)
VALUES(?,?,?);

-- name: GetPosts :many
SELECT post_id,title,content,created_at,updated_at
FROM posts;

-- name: GetPostsByCategory :many
SELECT p.* 
FROM posts p
JOIN post_categories pc ON p.post_id = pc.post_id
WHERE pc.category_id = ?;

-- name: GetPost :one
SELECT title,content,created_at,updated_at
FROM posts
WHERE post_id  = ?;

-- name: UpdatePost :exec
UPDATE posts
SET title = ?,
    content = ?
WHERE post_id = ? ;

-- name: DeletePost :exec
DELETE FROM posts
WHERE post_id = ? ;
