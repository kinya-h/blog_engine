package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kinya-h/blog_engine/db"
)

type commentResponse struct {
	CommentID int32     `json:"comment_id"`
	PostID    int32     `json:"post_id"`
	UserID    int32     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func newCommentResponse(comment db.Comment) commentResponse {
	return commentResponse{
		CommentID: comment.CommentID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
}

func (server *Server) createComment(w http.ResponseWriter, r *http.Request) {
	var req db.CreateCommentParams

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	result, err := server.db.CreateComment(context.Background(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(" Error %s", err), http.StatusInternalServerError)
		return
	}

	commentId, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("An Error Occured Creating Comment : %s", err.Error())
	}

	rsp := newCommentResponse(db.Comment{CommentID: int32(commentId), PostID: req.PostID, UserID: req.UserID, Content: req.Content})
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rsp)
}

func (server *Server) getComments(w http.ResponseWriter, r *http.Request) {

	categories, err := server.db.GetComments(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func (server *Server) getComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	commentId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, " id parameter Must be a number (id of the comment)", http.StatusBadRequest)
	}

	comment, err := server.db.GetComment(r.Context(), int32(commentId))

	if err != nil {

		if err == sql.ErrNoRows {
			emptyComment := make([]db.Comment, 0) // Return [] instead of a struct with empty fields
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(emptyComment)
			return
		}

		http.Error(w, fmt.Sprintf("An Error Occured %s", err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)

}

func (server *Server) updateComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	commentId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "id parameter Must be a number (id of the comment)", http.StatusBadRequest)
		return
	}

	var req db.UpdateCommentParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	arg := db.UpdateCommentParams{
		CommentID: req.CommentID,
		Content:   req.Content,
	}

	// Will not error out in case the id is not found (default mysql behaviour).
	if err := server.db.UpdateComment(r.Context(), arg); err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("comment with id %d updated successfully.", commentId)))

}

func (server *Server) deleteComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	commentId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "id parameter Must be a number (id of the Comment)", http.StatusBadRequest)
		return
	}

	if err := server.db.DeleteComment(r.Context(), int32(commentId)); err != nil {
		http.Error(w, "Sorry!,Unknown error occured", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(fmt.Sprintf("Comment with id %d deleted successfully. ", commentId)))

}
