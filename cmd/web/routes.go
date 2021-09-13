package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(noSurf)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.homeView))
	mux.Get("/login", dynamicMiddleware.ThenFunc(app.loginView))
	mux.Get("/logout", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logout))
	mux.Get("/customers", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.listCustomersView))
	mux.Get("/giftcards", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.listGiftCardsView))
	mux.Get("/giftcards/:id/send", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.sendAGiftCardView))
	mux.Get("/giftcards/:id/activate", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.activateAGiftCardView))
	mux.Get("/giftcards/:id/deactivate", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deactivateAGiftCardView))
	mux.Get("/giftcards/:id/load", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.loadAGiftCardView))

	return standardMiddleware.Then(mux)
}