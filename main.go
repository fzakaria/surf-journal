package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/fzakaria/surf-journal/internal"
	"github.com/go-chi/chi/v5"
)

//go:embed static/*
var staticFS embed.FS

//go:embed templates/*
var templatesFS embed.FS

//go:generate tailwindcss -i ./static/css/input.css -o ./static/css/output.css

func main() {
	fmt.Print("Starting Surf Journal\n")
	db := internal.ConnectDB()
	defer db.Close()

	r := chi.NewRouter()

	r.Get("/login", internal.LoginHandler)
	r.Post("/login", internal.LoginHandler)
	r.Post("/logout", internal.LogoutHandler)

	r.Group(func(auth chi.Router) {
		auth.Use(internal.AuthMiddleware)
		auth.Get("/", internal.IndexHandler)
	})

	// Static Files
	staticServer := http.FileServer(http.FS(staticFS))
	r.Handle("/static/*", http.StripPrefix("/static/", staticServer))

	log.Println("Serving at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
