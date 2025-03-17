package handlers

import (
	"html/template"
	"net/http"

	surf_journal "github.com/fzakaria/surf-journal"
	"github.com/go-chi/chi/v5"
)

type SessionsResource struct{}

func (rs SessionsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)

	return r
}

func (rs SessionsResource) List(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	username, _ := session.Values["username"].(string)
	tmpls, err := template.ParseFS(surf_journal.TemplateFS,
		"templates/base.html.tmpl",
		"templates/flash.html.tmpl",
		"templates/nav.html.tmpl",
		"templates/sessions/index.html.tmpl")
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	data := map[string]interface{}{
		"UserName":    username,
		"CurrentPage": "sessions",
	}

	Render(tmpls, w, r, data)
}
