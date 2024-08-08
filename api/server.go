package api

import (
	"database/sql"
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

	// Group the user routes
	router.Route("/users", func(r chi.Router) {
		r.Post("/", server.createUser)
		r.Post("/login", server.loginUser)
	})

	// Group the token routes
	router.Post("/tokens/renew_access", server.renewAccessToken)

	// Group the category routes
	router.Route("/categories", func(r chi.Router) {
		r.Post("/", server.createCategory)
		r.Get("/", server.getCategories)
		r.Get("/{id}", server.getCategory)
		r.Patch("/{id}", server.updateCategory)
		r.Delete("/{id}", server.deleteCategory)
	})

	// Group the post routes
	router.Route("/posts", func(r chi.Router) {
		r.Post("/", server.createPost)
		r.Get("/", server.getPosts)
		r.Get("/{id}", server.getPost)
		r.Patch("/{id}", server.updatePost)
		r.Delete("/{id}", server.deletePost)
		//Get posts by category
		r.Get("/category/{id}", server.getPostsByCategory)
	})

	// Group the comment routes
	router.Route("/comments", func(r chi.Router) {
		r.Post("/", server.createComment)
		r.Get("/", server.getComments)
		r.Get("/{id}", server.getComment)
		r.Patch("/{id}", server.updateComment)
		r.Delete("/{id}", server.deleteComment)
	})

	server.router = router
	return server, nil

}

func (server *Server) Start() {
	http.ListenAndServe(":3000", server.router)

}

func (s *Server) GetDBConn() (*sql.DB, error) {
	dbtx := s.db.GetDBTX()
	conn, ok := dbtx.(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("failed to assert DBTX to *sql.DB")
	}
	return conn, nil
}
