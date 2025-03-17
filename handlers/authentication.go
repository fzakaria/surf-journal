package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
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
	r.Get("/register", RegistrationHandler)
	r.Post("/register", RegistrationHandler)
	return r
}

type contextKey string

const dbKey contextKey = "db"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	session.Values["authenticated"] = false
	if err := session.Save(r, w); err != nil {
		ErrorHandler(w, r, err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		db, ok := r.Context().Value(dbKey).(*sql.DB)
		if !ok {
			ErrorHandler(w, r, fmt.Errorf("missing database"))
			return
		}

		serialized, err := passwords.NewSerializedPassword(password)
		if err != nil {
			session.AddFlash("Could not hash password", "errors")
			if err := session.Save(r, w); err != nil {
				ErrorHandler(w, r, err)
				return
			}
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		err = database.AddPassword(db, username, serialized)
		if err != nil {
			session.AddFlash("Could not create user", "errors")
			if err := session.Save(r, w); err != nil {
				ErrorHandler(w, r, err)
				return
			}
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		session.Values["authenticated"] = true
		session.Values["username"] = username
		if err := session.Save(r, w); err != nil {
			ErrorHandler(w, r, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpls, err := template.ParseFS(surf_journal.TemplateFS,
		"templates/base.html.tmpl",
		"templates/flash.html.tmpl",
		"templates/nav.html.tmpl",
		"templates/register.html.tmpl")
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	Render(tmpls, w, r, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		db, ok := r.Context().Value(dbKey).(*sql.DB)
		if !ok {
			ErrorHandler(w, r, fmt.Errorf("missing database"))
			return
		}

		serialized, err := database.GetSerializedPassword(db, username)
		if len(serialized) == 0 || err != nil {
			session.AddFlash(fmt.Sprintf("Could not find user %s", username), "errors")
			if err := session.Save(r, w); err != nil {
				ErrorHandler(w, r, err)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err = passwords.CheckPassword(serialized, password)
		if err != nil {
			session.AddFlash("Password does not match", "errors")
			if err := session.Save(r, w); err != nil {
				ErrorHandler(w, r, err)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session.Values["authenticated"] = true
		session.Values["username"] = username
		if err := session.Save(r, w); err != nil {
			ErrorHandler(w, r, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpls, err := template.ParseFS(surf_journal.TemplateFS,
		"templates/base.html.tmpl",
		"templates/flash.html.tmpl",
		"templates/nav.html.tmpl",
		"templates/login.html.tmpl")
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	Render(tmpls, w, r, nil)
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
