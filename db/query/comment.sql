
-- name: CreateComment :execresult
INSERT INTO comments(post_id,user_id,content)
VALUES(?,?,?);

-- name: GetComments :one
SELECT p.title,c.content,u.username,u.email
FROM comments c
JOIN users u 
ON u.user_id = c.user_id
JOIN posts p
ON c.post_id = c.post_id;


-- name: GetComment :one
SELECT p.title,c.content,u.username,u.email
FROM comments c
JOIN users u 
ON u.user_id = c.user_id
JOIN posts p
ON c.post_id = c.post_id
WHERE c.comment_id = ?;

-- name: UpdateComment :exec
UPDATE comments 
SET content = ? WHERE comment_id = ?;

-- name: DeleteComment :exec
DELETE FROM comments WHERE comment_id = ?;
