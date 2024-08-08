package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kinya-h/blog_engine/db"
)

type CreatePostRequest struct {
	CategoryID int32 `json:"category_id"`
	*db.CreatePostParams
}
type postResponse struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newPostResponse(post db.Post) postResponse {
	return postResponse{
		ID:        post.PostID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

func (server *Server) createPost(w http.ResponseWriter, r *http.Request) {

	var req CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}
	dbConn, err := server.GetDBConn()
	if err != nil {

		fmt.Println("Error retrieving DB connection:", err)
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusInternalServerError)
		return
	}

	tx, err := dbConn.Begin()

	if err != nil {
		fmt.Println("Error creating a transaction :", err)
		http.Error(w, fmt.Sprintf("An error occured. %s", err.Error()), http.StatusInternalServerError)
		return
	}

	defer tx.Rollback()

	qtx := server.db.WithTx(tx)
	arg := db.CreatePostParams{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	}
	result, err := qtx.CreatePost(r.Context(), arg)
	if err != nil {
		http.Error(w, fmt.Sprintf(" Error %s", err), http.StatusInternalServerError)
		return
	}

	postId, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("An Error Occured Creating Post : %s", err.Error())
		http.Error(w, fmt.Sprintf(" Error %s", err), http.StatusInternalServerError)
		return
	}

	if _, err := qtx.CreatePostCategory(r.Context(), db.CreatePostCategoryParams{PostID: int32(postId), CategoryID: req.CategoryID}); err != nil {
		http.Error(w, fmt.Sprintf(" Error %s", err), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("An Error Occured Creating Post : %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	rsp := newPostResponse(db.Post{PostID: int32(postId), Title: req.Title, Content: req.Content})
	json.NewEncoder(w).Encode(rsp)
}

func (server *Server) getPosts(w http.ResponseWriter, r *http.Request) {

	posts, err := server.db.GetPosts(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)

}

func (server *Server) getPostsByCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	categoryId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, " id parameter Must be a number (id of the category)", http.StatusBadRequest)

		return
	}
	posts, err := server.db.GetPostsByCategory(r.Context(), int32(categoryId))
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)

}

func (server *Server) getPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	postId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, " id parameter Must be a number (id of the post)", http.StatusBadRequest)

		return
	}

	post, err := server.db.GetPost(r.Context(), int32(postId))
	if err != nil {

		if err == sql.ErrNoRows {
			emptyPost := make([]db.Post, 0) // Return [] instead of a struct with empty fields
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(emptyPost)

			return
		}

		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)

}

func (server *Server) updatePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	postId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "id parameter Must be a number (id of the post)", http.StatusBadRequest)
		return
	}

	var req db.UpdatePostParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	arg := db.UpdatePostParams{
		Title:   req.Title,
		Content: req.Content,
		PostID:  int32(postId),
	}

	// Will not error out in case the id (primary key) is not found (default mysql behaviour).
	if err := server.db.UpdatePost(r.Context(), arg); err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Post with id %d updated successfully. ", postId)))

}

func (server *Server) deletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	postId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, " id parameter Must be a number (id of the post)", http.StatusBadRequest)
		return
	}

	if err := server.db.DeletePost(r.Context(), int32(postId)); err != nil {
		http.Error(w, "Sorry!,Unknown error occured", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(fmt.Sprintf("Post with id %d deleted successfully. ", postId)))
}
