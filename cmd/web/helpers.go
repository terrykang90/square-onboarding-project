package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/justinas/nosurf"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	sq_oauth_ctrl := NewSquareOauthController()
	return sq_oauth_ctrl.getSquareOAuthToken(r) != nil
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData) {
	td.CSRFToken = nosurf.Token(r)
	td.PageName = page
	td.IsAuthenticated = app.isAuthenticated(r)

	err := app.templates.ExecuteTemplate(w, page, td)
	if err != nil {
		panic(err)
	}
}

func setCookie(w http.ResponseWriter, key string, value string) {
	cookie := http.Cookie{
		Name:     key,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func clearCookie(w http.ResponseWriter, key string) {
	cookie := http.Cookie{
		Name:     key,
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}
