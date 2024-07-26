package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kinya-h/blog_engine/db"
	"github.com/kinya-h/blog_engine/token"
	"github.com/kinya-h/blog_engine/util"
)

type Server struct {
	config     util.Config
	db         *db.Queries
	tokenMaker token.Maker
	router     *chi.Mux
}

func NewServer(db *db.Queries) *Server {
	server := &Server{db: db}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/api/users", server.createUser)

	server.router = router
	return server

}

func (server *Server) Start() {
	http.ListenAndServe(":3000", server.router)

}
