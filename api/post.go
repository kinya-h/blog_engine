package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kinya-h/blog_engine/db"
)

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

	var req db.CreatePostParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	result, err := server.db.CreatePost(r.Context(), req)
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

func (server *Server) getPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	postId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, " id parameter Must be a number (id of the post)", http.StatusBadRequest)
		return
	}

	post, err := server.db.GetPost(r.Context(), int32(postId))
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)

}

func (server *Server) updatePost(w http.ResponseWriter, r *http.Request) {

	var req db.UpdatePostParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	result, err := server.db.UpdatePost(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusInternalServerError)
		return
	}

	updatedPostId, err := result.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Post with id %d updated successfully. ", updatedPostId)))

}

func (server *Server) deletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	postId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, " id parameter Must be a number (id of the post)", http.StatusBadRequest)
		return
	}

	errr := server.db.DeletePost(r.Context(), int32(postId))
	if errr != nil {
		http.Error(w, "Sorry!,Unknown error occured", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(fmt.Sprintf("Post with id %d deleted successfully. ", postId)))

}
