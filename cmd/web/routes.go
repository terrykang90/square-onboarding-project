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

	mux.Get("/oauth/square/login", dynamicMiddleware.ThenFunc(app.squareOauthLogin))
	mux.Get("/oauth/square/callback", dynamicMiddleware.ThenFunc(app.squareOauthCallback))

	mux.Get("/giftcards/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showAGiftCardView))
	mux.Post("/giftcards", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createAGiftCard))

	mux.Post("/giftcards/:id/send", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.sendAGiftCard))
	mux.Post("/giftcards/:id/activate", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.activateAGiftCard))
	mux.Post("/giftcards/:id/deactivate", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deactivateAGiftCard))
	mux.Post("/giftcards/:id/load", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.loadAGiftCard))

	return standardMiddleware.Then(mux)
}
