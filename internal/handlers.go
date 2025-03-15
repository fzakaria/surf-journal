package internal

import (
	"html/template"
	"log"
	"net/http"

	surf_journal "github.com/fzakaria/surf-journal"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-very-secret-key"))

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
		log.Fatal(err)
	}

	err = tmpls.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println("Template execution error:", err)
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
		log.Fatal(err)
		return
	}

	data := struct {
		UserName string
	}{
		UserName: username,
	}

	tmpls.ExecuteTemplate(w, "base", data)
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
