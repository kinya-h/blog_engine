
-- name: CreateCategory :execresult
INSERT INTO categories(name,description)
VALUES(?,?);

-- name: GetCategory :one
SELECT * FROM categories 
WHERE category_id = ?;

-- name: GetCategories :many
SELECT * FROM categories; 

-- name: UpdateCategory :exec
UPDATE categories 
SET name = ?,
description = ?
WHERE category_id = ?;


-- name: DeleteCategory :exec
DELETE FROM categories
WHERE category_id = ?;


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

