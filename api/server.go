package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kinya-h/blog_engine/db"
)

type Server struct {
	db     *db.Queries
	router *chi.Mux
}

func NewServer(db *db.Queries) *Server {
	server := &Server{db: db}
	router := chi.NewRouter()

	server.router = router
	return server

}

func (server *Server) Start() {
	http.ListenAndServe(":3000", server.router)

}
