package handlers

import (
	"html/template"
	"log"
	"net/http"

	surf_journal "github.com/fzakaria/surf-journal"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpls, err := template.ParseFS(surf_journal.TemplateFS,
		"templates/base.html.tmpl",
		"templates/nav.html.tmpl",
		"templates/404.html.tmpl")
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	err = tmpls.ExecuteTemplate(w, "base", nil)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	username, _ := session.Values["username"].(string)
	tmpls, err := template.ParseFS(surf_journal.TemplateFS,
		"templates/base.html.tmpl",
		"templates/nav.html.tmpl",
		"templates/home.html.tmpl")
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	data := struct {
		UserName    string
		CurrentPage string
	}{
		UserName:    username,
		CurrentPage: "home",
	}

	err = tmpls.ExecuteTemplate(w, "base", data)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpls, err := template.ParseFS(surf_journal.TemplateFS,
		"templates/base.html.tmpl",
		"templates/nav.html.tmpl",
		"templates/505.html.tmpl")
	// Render a friendly HTML error page
	err = tmpls.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Panic(err)
	}
}
