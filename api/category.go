package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kinya-h/blog_engine/db"
)

type categoryResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func newCategoryResponse(category db.Category) categoryResponse {
	return categoryResponse{
		ID:          category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
	}
}

func (server *Server) createCategory(w http.ResponseWriter, r *http.Request) {
	var req db.CreateCategoryParams

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	result, err := server.db.CreateCategory(context.Background(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf(" Error %s", err), http.StatusInternalServerError)
		return
	}

	userId, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("An Error Occured Creating Category : %s", err.Error())
	}

	rsp := newCategoryResponse(db.Category{CategoryID: int32(userId), Name: req.Name, Description: req.Description})
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rsp)
}

func (server *Server) getCategories(w http.ResponseWriter, r *http.Request) {

	categories, err := server.db.GetCategories(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func (server *Server) getCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	categoryId, err := strconv.Atoi(id)

	if err != nil {

		http.Error(w, " id parameter Must be a number (id of the category)", http.StatusBadRequest)

	}
	category, err := server.db.GetCategory(r.Context(), int32(categoryId))

	if err != nil {

		if err == sql.ErrNoRows {
			emptyCategory := make([]db.Category, 0) // Return [] instead of a struct with empty fields
			json.NewEncoder(w).Encode(emptyCategory)

			return
		}
		http.Error(w, fmt.Sprintf("An Error Occured %s", err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)

}

func (server *Server) updateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	categoryId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "id parameter Must be a number (id of the category)", http.StatusBadRequest)
		return
	}

	var req db.UpdateCategoryParams
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("An error occured %s", decodeErr), http.StatusBadRequest)
		return
	}

	arg := db.UpdateCategoryParams{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  int32(categoryId),
	}

	result, err := server.db.UpdateCategory(r.Context(), arg)
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
	w.Write([]byte(fmt.Sprintf("Category with id %d updated successfully. ", updatedPostId)))

}

func (server *Server) deleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	categoryId, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "id parameter Must be a number (id of the category)", http.StatusBadRequest)
		return
	}

	errr := server.db.DeleteCategory(r.Context(), int32(categoryId))
	if errr != nil {
		http.Error(w, "Sorry!,Unknown error occured", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(fmt.Sprintf("Category with id %d deleted successfully. ", categoryId)))

}
