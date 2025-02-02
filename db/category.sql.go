// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: category.sql

package db

import (
	"context"
	"database/sql"
)

const createCategory = `-- name: CreateCategory :execresult
INSERT INTO categories(name,description)
VALUES(?,?)
`

type CreateCategoryParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createCategory, arg.Name, arg.Description)
}

const deleteCategory = `-- name: DeleteCategory :exec
DELETE FROM categories
WHERE category_id = ?
`

func (q *Queries) DeleteCategory(ctx context.Context, categoryID int32) error {
	_, err := q.db.ExecContext(ctx, deleteCategory, categoryID)
	return err
}

const getCategories = `-- name: GetCategories :many
SELECT category_id, name, description FROM categories
`

func (q *Queries) GetCategories(ctx context.Context) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, getCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.CategoryID, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCategory = `-- name: GetCategory :one
SELECT category_id, name, description FROM categories 
WHERE category_id = ?
`

func (q *Queries) GetCategory(ctx context.Context, categoryID int32) (Category, error) {
	row := q.db.QueryRowContext(ctx, getCategory, categoryID)
	var i Category
	err := row.Scan(&i.CategoryID, &i.Name, &i.Description)
	return i, err
}

const updateCategory = `-- name: UpdateCategory :exec
UPDATE categories 
SET name = ?,
description = ?
WHERE category_id = ?
`

type UpdateCategoryParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  int32  `json:"category_id"`
}

func (q *Queries) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) error {
	_, err := q.db.ExecContext(ctx, updateCategory, arg.Name, arg.Description, arg.CategoryID)
	return err
}
