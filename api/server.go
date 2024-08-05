package api

import (
	"fmt"
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

func NewServer(config util.Config, db *db.Queries) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		db:         db,
		tokenMaker: tokenMaker,
		config:     config,
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	//users route
	router.Post("/users", server.createUser)
	router.Post("/users/login", server.loginUser)

	router.Post("/tokens/renew_access", server.renewAccessToken)

	//categories route
	router.Post("/categories", server.createCategory)
	router.Get("/categories", server.getCategories)
	router.Get("/categories/{id}", server.getCategory)
	router.Patch("/categories/{id}", server.updateCategory)
	router.Delete("/categories/{id}", server.deleteCategory)

	//posts route
	router.Post("/posts", server.createPost)
	router.Get("/posts", server.getPosts)
	router.Get("/posts/{id}", server.getPost)
	router.Patch("/posts/{id}", server.updatePost)
	router.Delete("/posts/{id}", server.deletePost)

	server.router = router
	return server, nil

}

func (server *Server) Start() {
	http.ListenAndServe(":3000", server.router)

}
