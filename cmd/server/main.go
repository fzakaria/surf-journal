package main

import (
	"log"
	"net/http"

	surf_journal "github.com/fzakaria/surf-journal"
	"github.com/fzakaria/surf-journal/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handlers.DatabaseMiddleware)

	r.Mount("/", handlers.AuthenticationRouter())

	r.Group(func(auth chi.Router) {
		auth.Use(handlers.AuthMiddleware)
		auth.Get("/", handlers.HomeHandler)
		auth.Mount("/sessions", handlers.SessionsResource{}.Routes())
	})

	r.NotFound(handlers.NotFoundHandler)

	staticServer := http.FileServer(http.FS(surf_journal.StaticFS))
	r.Handle("/static/*", staticServer)

	log.Println("Serving at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
