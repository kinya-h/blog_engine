

-- name: CreatePostCategory :execresult
INSERT INTO post_categories
(post_id,category_id)
VALUES (?,?);


-- name: GetPostCategory :one
SELECT * FROM post_categories 
WHERE post_id =?;


-- name: UpdatePostCategory :exec
UPDATE post_categories
SET category_id = ?
WHERE post_id = ?;