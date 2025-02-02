// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: posts.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createPost = `-- name: CreatePost :execresult
INSERT INTO posts
(user_id,title,content)
VALUES(?,?,?)
`

type CreatePostParams struct {
	UserID  int32  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createPost, arg.UserID, arg.Title, arg.Content)
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts
WHERE post_id = ?
`

func (q *Queries) DeletePost(ctx context.Context, postID int32) error {
	_, err := q.db.ExecContext(ctx, deletePost, postID)
	return err
}

const getPost = `-- name: GetPost :one
SELECT title,content,created_at,updated_at
FROM posts
WHERE post_id  = ?
`

type GetPostRow struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) GetPost(ctx context.Context, postID int32) (GetPostRow, error) {
	row := q.db.QueryRowContext(ctx, getPost, postID)
	var i GetPostRow
	err := row.Scan(
		&i.Title,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPosts = `-- name: GetPosts :many
SELECT post_id,title,content,created_at,updated_at
FROM posts
`

type GetPostsRow struct {
	PostID    int32     `json:"post_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) GetPosts(ctx context.Context) ([]GetPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPostsRow{}
	for rows.Next() {
		var i GetPostsRow
		if err := rows.Scan(
			&i.PostID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
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

const getPostsByCategory = `-- name: GetPostsByCategory :many
SELECT p.post_id, p.user_id, p.title, p.content, p.created_at, p.updated_at 
FROM posts p
JOIN post_categories pc ON p.post_id = pc.post_id
WHERE pc.category_id = ?
`

func (q *Queries) GetPostsByCategory(ctx context.Context, categoryID int32) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByCategory, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.PostID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
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

const updatePost = `-- name: UpdatePost :exec
UPDATE posts
SET title = ?,
    content = ?
WHERE post_id = ?
`

type UpdatePostParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	PostID  int32  `json:"post_id"`
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) error {
	_, err := q.db.ExecContext(ctx, updatePost, arg.Title, arg.Content, arg.PostID)
	return err
}
