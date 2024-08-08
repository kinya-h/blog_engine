
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


