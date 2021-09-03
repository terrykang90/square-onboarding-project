package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			// setCookie(w, OAUTH_TOKEN_COOKIE_KEY, "eyJhY2Nlc3NfdG9rZW4iOiJFQUFBRVBQUkUxc181VURHSjhZUzFXN192TExvVDZPdDNoZGw3ejRLaXFQWUVsVE1Ub2ZHUjBWZzZEbEZhaEI1IiwidG9rZW5fdHlwZSI6ImJlYXJlciIsInJlZnJlc2hfdG9rZW4iOiJFUUFBRUZjc1E4T2VtZmdORlpEOVlhNnB5OWZkV3pSYmpJRDg2ejQ2dlBkbk9EV3ZPQkRHa3MxUThBT2FqR1dVIiwiZXhwaXJ5IjoiMDAwMS0wMS0wMVQwMDowMDowMFoifQ==")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
