package handlers

import (
	"html/template"
	"net/http"

	surf_journal "github.com/fzakaria/surf-journal"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-very-secret-key"))

func AuthenticationRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/login", LoginHandler)
	r.Post("/login", LoginHandler)
	r.Post("/logout", LogoutHandler)
	return r
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	session.Values["authenticated"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")

	if r.Method == "POST" {
		username := r.FormValue("username")
		// password := r.FormValue("password")

		// TODO: Validate username/password via DB.

		session.Values["authenticated"] = true
		session.Values["username"] = username
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpls, err := template.ParseFS(surf_journal.TemplateFS,
		"templates/base.html.tmpl",
		"templates/nav.html.tmpl",
		"templates/login.html.tmpl")
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth-session")
		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
