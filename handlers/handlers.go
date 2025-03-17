package handlers

import (
	"html/template"
	"log"
	"net/http"

	surf_journal "github.com/fzakaria/surf-journal"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	username, _ := session.Values["username"].(string)

	tmpls, err := template.ParseFS(surf_journal.TemplateFS,
		"templates/base.html.tmpl",
		"templates/flash.html.tmpl",
		"templates/nav.html.tmpl",
		"templates/404.html.tmpl")
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	data := map[string]interface{}{
		"UserName":    username,
		"CurrentPage": "home",
	}

	Render(tmpls, w, r, data)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	username, _ := session.Values["username"].(string)
	tmpls, err := template.ParseFS(surf_journal.TemplateFS,
		"templates/base.html.tmpl",
		"templates/flash.html.tmpl",
		"templates/nav.html.tmpl",
		"templates/home.html.tmpl")
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	data := map[string]interface{}{
		"UserName":    username,
		"CurrentPage": "home",
	}

	Render(tmpls, w, r, data)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpls, err := template.ParseFS(surf_journal.TemplateFS,
		"templates/base.html.tmpl",
		"templates/flash.html.tmpl",
		"templates/nav.html.tmpl",
		"templates/505.html.tmpl")
	// Render a friendly HTML error page
	err = tmpls.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Panic(err)
	}
}

func Render(tmpls *template.Template, w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	session, _ := store.Get(r, "auth-session")

	if data == nil {
		data = make(map[string]interface{})
	}

	data["errors"] = session.Flashes("errors")

	if err := session.Save(r, w); err != nil {
		ErrorHandler(w, r, err)
		return
	}

	if err := tmpls.ExecuteTemplate(w, "base", data); err != nil {
		ErrorHandler(w, nil, err)
		return
	}
}
