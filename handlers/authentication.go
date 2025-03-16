package handlers

import (
	"context"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	surf_journal "github.com/fzakaria/surf-journal"
	"github.com/fzakaria/surf-journal/database"
	"github.com/fzakaria/surf-journal/passwords"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-very-secret-key"))

func AuthenticationRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/login", LoginHandler)
	r.Post("/login", LoginHandler)
	r.Get("/logout", LogoutHandler)
	r.Post("/logout", LogoutHandler)
	return r
}

type contextKey string

const dbKey contextKey = "db"

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
		password := r.FormValue("password")

		db, ok := r.Context().Value(dbKey).(*sql.DB)
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print("login/ missing database")
			return
		}

		serialized, err := database.GetSerializedPassword(db, username)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Printf("login/ GetSerializedPassword: %+v", err)
			return
		}

		if len(serialized) == 0 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			log.Printf("login/ unknown user %s", username)
			return
		}

		err = passwords.CheckPassword(serialized, password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			log.Printf("login/ CheckPassword: %+v", err)
			return
		}

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

func DatabaseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := database.Connect()
		ctx := context.WithValue(r.Context(), dbKey, db)
		context.AfterFunc(ctx, func() {
			db.Close()
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
