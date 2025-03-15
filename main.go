package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/fzakaria/surf-journal/internal"
	"github.com/go-chi/chi/v5"
)

//go:embed static/output.css
var staticFS embed.FS

//go:generate tailwindcss -i ./static/input.css -o ./static/output.css

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

	// Static Files; go:embed creates an embed.FS that mirrors the
	// paths on disk, so "output.css" would be available at
	// static/output.css. But the Handler below registers the embed.FS
	// at static/, so all the paths would be static/static. fs.Sub
	// re-roots the tree.
	staticRootFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(fmt.Sprintf("could not create sub-fileystem rooted at static: %+v", err))
	}
	staticServer := http.FileServer(http.FS(staticRootFS))
	r.Handle("/static/*", http.StripPrefix("/static/", staticServer))

	log.Println("Serving at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
