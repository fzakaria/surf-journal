package main

import (
	"log"
	"net/http"

	surf_journal "github.com/fzakaria/surf-journal"
	"github.com/fzakaria/surf-journal/internal"
	"github.com/go-chi/chi/v5"
)

func main() {
	db := internal.ConnectDB()
	defer db.Close()

	r := chi.NewRouter()

	r.Get("/login", internal.LoginHandler)
	r.Post("/login", internal.LoginHandler)
	r.Post("/logout", internal.LogoutHandler)

	r.Group(func(auth chi.Router) {
		auth.Use(internal.AuthMiddleware)
		auth.Get("/", internal.HomeHandler)
	})

	staticServer := http.FileServer(http.FS(surf_journal.StaticFS))
	r.Handle("/static/*", staticServer)

	log.Println("Serving at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
