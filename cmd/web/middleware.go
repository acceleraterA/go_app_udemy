package main

import (
	"net/http"

	"github.com/acceleraterA/go_app_udemy/internal/helpers"
	"github.com/justinas/nosurf"
)

/*
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("HIT THE PAGE")
		next.ServeHTTP(w, r)
	})
}
*/
// Nosurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// has access to request
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "login first!")
			http.Redirect(w, r, "user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
